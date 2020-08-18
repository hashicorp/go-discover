// Package hcloud provides node discovery for Hetzner Cloud.
package hcloud

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

type Provider struct{}

func (p *Provider) Help() string {
	return `Hetzner Cloud:

    provider:       "hcloud"
    api_token:      The Hetzner CLoud API token to use (required)
    label_selector: The label selector to filter servers by
    addr_type:      "private_v4", "public_v4" or "public_v6". Defaults to "private_v4". If multiple private networks are defined, the first will be used.
`
}

func listServersByLabel(c *hcloud.Client, labelSelector string) ([]hcloud.Server, error) {
	var serverList []hcloud.Server

	opts := hcloud.ServerListOpts{
		ListOpts: hcloud.ListOpts{
			LabelSelector: labelSelector,
		},
		Status: []hcloud.ServerStatus{hcloud.ServerStatusRunning},
	}

	servers, err := c.Server.AllWithOpts(context.Background(), opts)

	if err != nil {
		return nil, err
	}

	for _, s := range servers {
		serverList = append(serverList, *s)
	}

	return serverList, nil
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "hcloud" {
		return nil, fmt.Errorf("discover-hcloud: invalid provider " + args["provider"])
	}

	if l == nil {
		l = log.New(ioutil.Discard, "", 0)
	}

	labelSelector := args["label_selector"]
	apiToken := args["api_token"]
	addrType := args["addr_type"]

	if addrType == "" {
		l.Printf("[DEBUG] discover-hcloud: Address type not provided. Using 'private_v4'")
		addrType = "private_v4"
	}

	if addrType != "private_v4" && addrType != "public_v4" && addrType != "public_v6" {
		l.Printf("[INFO] discover-hcloud: Address type %s is not supported. Valid values are {private_v4,public_v4,public_v6}. Falling back to 'private_v4'", addrType)
		addrType = "private_v4"
	}

	l.Printf("[DEBUG] discover-hcloud: Using label_selector=%s", labelSelector)

	client := hcloud.NewClient(hcloud.WithToken(apiToken))
	servers, err := listServersByLabel(client, labelSelector)

	if err != nil {
		return nil, fmt.Errorf("discover-hcloud: %s", err)
	}

	var addrs []string

	for _, s := range servers {
		l.Printf("[INFO] discover-hcloud: Found instance %s (%d)", s.Name, s.ID)

		switch addrType {
		case "public_v4":
			l.Printf("[INFO] discover-hcloud: Instance %s (%d) has public ip %s", s.Name, s.ID, s.PublicNet.IPv4.IP.String())
			addrs = append(addrs, s.PublicNet.IPv4.IP.String())
		case "private_v4":
			if len(s.PrivateNet) == 0 {
				l.Printf("[INFO] discover-hcloud: Instance %s (%d) has no private ip", s.Name, s.ID)
			} else {
				l.Printf("[INFO] discover-hcloud: Instance %s (%d) has private ip %s", s.Name, s.ID, s.PrivateNet[0].IP.String())
				addrs = append(addrs, s.PrivateNet[0].IP.String())
			}
		case "public_v6":
			l.Printf("[INFO] discover-hcloud: Instance %s (%d) has public ip %s", s.Name, s.ID, s.PublicNet.IPv6.IP.String())
			addrs = append(addrs, s.PublicNet.IPv6.IP.String())
		default:
		}
	}

	l.Printf("[DEBUG] discover-hcloud: Found ip addresses: %v", addrs)
	return addrs, nil
}
