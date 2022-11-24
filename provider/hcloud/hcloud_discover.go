// Package hcloud provides node discovery for Hetzner Cloud.
package hcloud

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

type Provider struct{}

func (p *Provider) Help() string {
	return `Hetzner Cloud:

	provider:       "hcloud"
	location:       The Hetzner Cloud datacenter location to filter by (eg. "fsn1")
	label_selector: The label selector to filter by
	address_type:   "private_v4", "public_v4" or "public_v6", defaults to "private_v4". In the case of private networks, the first one will be used
	api_token:      The Hetzner Cloud API token to use, can also be provided by environment variable: HCLOUD_TOKEN
`
}

// serverIP returns the IP address of the specified type for the hcloud server.
func serverIP(s *hcloud.Server, addrType string, l *log.Logger) string {
	switch addrType {
	case "public_v4":
		if !s.PublicNet.IPv4.Blocked {
			l.Printf("[INFO] discover-hcloud: instance %s (%d) has public IP %s", s.Name, s.ID, s.PublicNet.IPv4.IP.String())
			return s.PublicNet.IPv4.IP.String()
		} else if len(s.PublicNet.FloatingIPs) != 0 {
			l.Printf("[INFO] discover-hcloud: public IPv4 for instance %s (%d) is blocked, checking associated floating IPs", s.Name, s.ID)
			for _, floatingIP := range s.PublicNet.FloatingIPs {
				if floatingIP.Type == hcloud.FloatingIPTypeIPv4 && !floatingIP.Blocked {
					l.Printf("[INFO] discover-hcloud: instance %s (%d) has floating IP %s", s.Name, s.ID, floatingIP.IP.String())
					return floatingIP.IP.String()
				}
			}
		}
	case "public_v6":
		if !s.PublicNet.IPv6.Blocked {
			l.Printf("[INFO] discover-hcloud: instance %s (%d) has public IP %s", s.Name, s.ID, s.PublicNet.IPv6.IP.String())
			return s.PublicNet.IPv6.IP.String()
		} else if len(s.PublicNet.FloatingIPs) != 0 {
			l.Printf("[INFO] discover-hcloud: public IPv6 for instance %s (%d) is blocked, checking associated floating IPs", s.Name, s.ID)
			for _, floatingIP := range s.PublicNet.FloatingIPs {
				if floatingIP.Type == hcloud.FloatingIPTypeIPv6 && !floatingIP.Blocked {
					l.Printf("[INFO] discover-hcloud: instance %s (%d) has floating IP %s", s.Name, s.ID, floatingIP.IP.String())
					return floatingIP.IP.String()
				}
			}
		}
	case "private_v4":
		if len(s.PrivateNet) == 0 {
			l.Printf("[INFO] discover-hcloud: instance %s (%d) has no private IP", s.Name, s.ID)
		} else {
			l.Printf("[INFO] discover-hcloud: instance %s (%d) has private IP %s", s.Name, s.ID, s.PrivateNet[0].IP.String())
			return s.PrivateNet[0].IP.String()
		}
	default:
	}

	l.Printf("[DEBUG] discover-hcloud: instance %s (%d) has no valid associated IP address", s.Name, s.ID)
	return ""
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "hcloud" {
		return nil, fmt.Errorf("discover-hcloud: invalid provider %s", args["provider"])
	}

	if l == nil {
		l = log.New(ioutil.Discard, "", 0)
	}

	labelSelector := args["label_selector"]
	addrType := args["address_type"]
	location := args["location"]
	apiToken := args["api_token"]

	if apiToken == "" {
		l.Printf("[INFO] no API token specified, checking environment variable HCLOUD_TOKEN")
		apiToken = os.Getenv("HCLOUD_TOKEN")
		if apiToken == "" {
			return nil, fmt.Errorf("discover-hcloud: no api_token specified")
		}
	}

	if addrType == "" {
		l.Printf("[INFO] discover-hcloud: address type not provided, using 'private_v4'")
		addrType = "private_v4"
	}

	if addrType != "private_v4" && addrType != "public_v4" && addrType != "public_v6" {
		l.Printf("[INFO] discover-hcloud: address_type %s is invalid, falling back to 'private_v4'. valid values are: private_v4, public_v4, public_v6", addrType)
		addrType = "private_v4"
	}

	l.Printf("[DEBUG] discover-hcloud: using address_type=%s label_selector=%s location=%s", addrType, labelSelector, location)

	client := hcloud.NewClient(hcloud.WithToken(apiToken))

	options := hcloud.ServerListOpts{
		ListOpts: hcloud.ListOpts{
			LabelSelector: labelSelector,
		},
		Status: []hcloud.ServerStatus{hcloud.ServerStatusRunning},
	}

	servers, err := client.Server.AllWithOpts(context.Background(), options)
	if err != nil {
		return nil, fmt.Errorf("discover-hcloud: %s", err)
	}

	var addrs []string
	for _, s := range servers {
		if location == "" || location == s.Datacenter.Location.Name {
			if serverIP := serverIP(s, addrType, l); serverIP != "" {
				addrs = append(addrs, serverIP)
			}
		}
	}

	log.Printf("[DEBUG] discover-hcloud: found IP addresses: %v", addrs)
	return addrs, nil
}
