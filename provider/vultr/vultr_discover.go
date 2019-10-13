// Package digitalocean provides node discovery for DigitalOcean.
package vultr

import (
        "fmt"
        "io/ioutil"
        "log"
	"context"

        "github.com/vultr/govultr"
)

type Provider struct {

	userAgent string

}

func (p *Provider) SetUserAgent(s string) {
	p.userAgent = s
}

func (p *Provider) Help() string {
	return `Vultr:

    provider: "vultr"
    region: The Vultr region to filter on
    tag_name: The tag name to filter on
    api_token: The Vultr API Token to use
`
}

func listServersByTag(c *govultr.Client, tagName string) ([]govultr.Server, error) {
	servers, err := c.Server.ListByTag(context.Background(), tagName)

	if err != nil {
		return nil, err
	}

	return servers, nil
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error){
	if args["provider"] != "vultr" {
		return nil, fmt.Errorf("discover-vultr: invalid provider " + args["privder"])
	}

	if l == nil {
		l = log.New(ioutil.Discard, "", 0)
	}

	region := args["region"]
	tagName := args["tag_name"]
	apiToken := args["api_token"]
	l.Printf("[DEBUG] discover-vultr: Using region=%s tag_name=%s", region, tagName)

	client := govultr.NewClient(nil, apiToken)

	servers, err := listServersByTag(client, tagName)
	if err != nil{
		return nil, fmt.Errorf("discover-vultr: %s", err)
	}

	var addrs []string
	var privateIP string

	for _, s := range servers {
		if s.RegionID == region || region == "" {
			privateIP = s.InternalIP
			if privateIP != "" {
				l.Printf("[INFO] discover-vultr: Found instance %s (%s) with private IP: %s", s.Label, s.InstanceID, privateIP)
				addrs = append(addrs, privateIP)
			}
		}
	}
	l.Printf("[DEBUG] discover-vultr: Found ip addresses: %v", addrs)
	return addrs, nil
}

