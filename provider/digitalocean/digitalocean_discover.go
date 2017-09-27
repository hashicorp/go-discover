// Package digitalocean provides node discovery for DigitalOcean.
package digitalocean

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/oauth2"

	"github.com/digitalocean/godo"
)

type Provider struct{}

func (p *Provider) Help() string {
	return `DigitalOcean:

    provider:   "digitalocean"
    region:     The DigitalOcean region to filter on
    tag_value:  The tag value to filter on
    addr_type:  "private_v4", "public_v4" or "public_v6". Defaults to "public_v4".
    api_key:    The DigitalOcean api key to use
`
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "digitalocean" {
		return nil, fmt.Errorf("discover-digitalocean: invalid provider " + args["provider"])
	}

	if l == nil {
		l = log.New(ioutil.Discard, "", 0)
	}

	region := args["region"]
	tagValue := args["tag_value"]
	apiKey := args["api_key"]
	addrType := args["addr_type"]

	l.Printf("[INFO] discover-digitalocean: Region is %q", region)

	if addrType != "private_v4" && addrType != "public_v4" && addrType != "public_v6" {
		l.Printf("[INFO] discover-digitalocean: Address type %s is not supported. Valid values are {private_v4,public_v4,public_v6}. Falling back to 'public_v4'", addrType)
		addrType = "public_v4"
	}

	if addrType == "" {
		l.Printf("[DEBUG] discover-digitalocean: Address type not provided. Using 'public_v4'")
		addrType = "public_v4"
	}

	client := newClientFromToken(apiKey)

	// Get the droplets for that tag
	droplets, _, err := client.Droplets.ListByTag(context.Background(), tagValue, nil)
	if err != nil {
		return nil, fmt.Errorf("discover-digitalocean: %s", err)
	}

	var addrs []string
	for _, droplet := range droplets {
		// Ignore droplets not in the requested region
		if droplet.Region.Slug != region {
			continue
		}

		var (
			addr string
			err error
		)
		switch addrType {
		case "public_v4":
			addr, err = droplet.PublicIPv4()
			if err != nil {
				fmt.Errorf("discover-digitalocean: %s", err)
				continue
			}
		case "private_v4":
			addr, err = droplet.PrivateIPv4()
			if err != nil {
				fmt.Errorf("discover-digitalocean: %s", err)
				continue
			}
		case "public_v6":
			addr, err = droplet.PublicIPv6()
			if err != nil {
				fmt.Errorf("discover-digitalocean: %s", err)
				continue
			}
		}

		l.Printf("[INFO] discover-digitalocean: Found instance (%d) %s with IP: %s",
			droplet.ID, droplet.Name, addr)
		addrs = append(addrs, addr)
	}
	return addrs, nil
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

func newClientFromToken(pat string) *godo.Client {
	tokenSource := &TokenSource{
		AccessToken: pat,
	}
	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	return godo.NewClient(oauthClient)
}
