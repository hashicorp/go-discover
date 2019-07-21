// Package digitalocean provides node discovery for DigitalOcean.
package digitalocean

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

type Provider struct {
	userAgent string
}

func (p *Provider) SetUserAgent(s string) {
	p.userAgent = s
}

func (p *Provider) Help() string {
	return `DigitalOcean:

    provider:  "digitalocean"
    region:    The DigitalOcean region to filter on
    tag_names:  The tag name to filter on
    api_token: The DigitalOcean API token to use
`
}

type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func tagExist(tag string, tagList []string) bool {
	for _, check := range tagList {
		if check == tag {
			return true
		}
	}

	return false
}

func listDropletsByTags(c *godo.Client, tagNames string) ([]godo.Droplet, error) {
	tagList := strings.Split(tagNames, ",")
	first, tagList := tagList[0], tagList[1:]

	if dropletList, err := listDropletsByTag(c, first); err != nil {
		return nil, err
	} else {
		for _, tag := range tagList {
			for _, droplet := range dropletList {
				if !tagExist(tag, droplet.Tags) {
					droplet = godo.Droplet{}
				}
			}
		}

		return dropletList, nil
	}
}

func listDropletsByTag(c *godo.Client, tagName string) ([]godo.Droplet, error) {
	dropletList := []godo.Droplet{}
	pageOpt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	for {
		droplets, resp, err := c.Droplets.ListByTag(context.TODO(), tagName, pageOpt)
		if err != nil {
			return nil, err
		}

		for _, d := range droplets {
			dropletList = append(dropletList, d)
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}

		pageOpt.Page = page + 1
	}

	return dropletList, nil
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "digitalocean" {
		return nil, fmt.Errorf("discover-digitalocean: invalid provider " + args["provider"])
	}

	if l == nil {
		l = log.New(ioutil.Discard, "", 0)
	}

	region := args["region"]
	tagNames := args["tag_names"]
	apiToken := args["api_token"]
	l.Printf("[DEBUG] discover-digitalocean: Using region=%s tag_names=%s", region, tagNames)

	tokenSource := &TokenSource{
		AccessToken: apiToken,
	}

	oauthClient := oauth2.NewClient(context.TODO(), tokenSource)
	client := godo.NewClient(oauthClient)
	if p.userAgent != "" {
		client.UserAgent = p.userAgent
	}

	droplets, err := listDropletsByTags(client, tagNames)
	if err != nil {
		return nil, fmt.Errorf("discover-digitalocean: %s", err)
	}

	var addrs []string
	for _, d := range droplets {
		if d.Region.Slug == region || region == "" {
			privateIP, err := d.PrivateIPv4()
			if err != nil {
				return nil, fmt.Errorf("discover-digitalocean: %s", err)
			}

			if privateIP != "" {
				l.Printf("[INFO] discover-digitalocean: Found instance %s (%d) with private IP: %s", d.Name, d.ID, privateIP)
				addrs = append(addrs, privateIP)
			}
		}
	}

	l.Printf("[DEBUG] discover-digitalocean: Found ip addresses: %v", addrs)
	return addrs, nil
}
