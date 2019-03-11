// Package linode provides node discovery for Linode.
package linode

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/linode/linodego"
	"golang.org/x/oauth2"
)

type Filter struct {
	Region string `json:"region,omitempty"`
	Tag    string `json:"tags,omitempty"`
}

type Provider struct {
	userAgent string
}

func (p *Provider) SetUserAgent(s string) {
	p.userAgent = s
}

func (p *Provider) Help() string {
	return `Linode:
    provider:     "linode"
    api_token:    The Linode API token to use
    region:       The Linode region to filter on
    tag_name:     The tag name to filter on
    address_type: "private_v4", "public_v4", "private_v6" or "public_v6". (default: "private_v4")

    Variables can also be provided by environment variables:
    export LINODE_TOKEN for api_token
`
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "linode" {
		return nil, fmt.Errorf("discover-linode: invalid provider " + args["provider"])
	}

	if l == nil {
		l = log.New(ioutil.Discard, "", 0)
	}

	addressType := args["address_type"]
	region := args["region"]
	tagName := args["tag_name"]
	apiToken := argsOrEnv(args, "api_token", "LINODE_TOKEN")
	l.Printf("[DEBUG] discover-linode: Using address_type=%s region=%s tag_name=%s", addressType, region, tagName)

	client := getLinodeClient(p.userAgent, apiToken)

	filters := Filter{
		Region: "",
		Tag:    "",
	}

	if region != "" {
		filters.Region = region
	}
	if tagName != "" {
		filters.Tag = tagName
	}

	jsonFilters, _ := json.Marshal(filters)
	filterOpt := linodego.ListOptions{Filter: string(jsonFilters)}

	linodes, err := client.ListInstances(context.Background(), &filterOpt)
	if err != nil {
		return nil, fmt.Errorf("discover-linode: Fetching Linode instances failed: %s", err)
	}

	var addrs []string
	for _, linode := range linodes {
		if linodeAddrs, err := getLinodeAddresses(client, addressType, linode); err == nil {
			addrs = append(addrs, linodeAddrs...)
		} else {
			return nil, err
		}
	}

	return addrs, nil
}

func getLinodeAddresses(client linodego.Client, addressType string, linode linodego.Instance) (addresses []string, err error) {
	switch addressType {
	case "public_v4":
		for _, ip := range linode.IPv4 {
			if !privateIP(*ip) {
				addresses = append(addresses, ip.String())
			}
		}
	case "public_v6":
		v6Addr := strings.SplitN(linode.IPv6, "/", 2)
		if len(v6Addr) > 0 {
			addresses = append(addresses, v6Addr[0])
		}
	case "private_v6":
		var addr *linodego.InstanceIPAddressResponse
		if addr, err = client.GetInstanceIPAddresses(context.Background(), linode.ID); err == nil {
			if addr.IPv6.LinkLocal.Address != "" {
				addresses = append(addresses, addr.IPv6.LinkLocal.Address)
			}
		} else {
			err = fmt.Errorf("discover-linode: Fetching Linode IP address for instance %v failed: %s", linode.ID, err)
		}
	// private_v4
	default:
		for _, ip := range linode.IPv4 {
			if privateIP(*ip) {
				addresses = append(addresses, ip.String())
			}
		}
	}
	return addresses, err
}

func getLinodeClient(userAgent, apiToken string) linodego.Client {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: apiToken})

	oauth2Client := &http.Client{
		Transport: &oauth2.Transport{
			Source: tokenSource,
		},
	}

	client := linodego.NewClient(oauth2Client)
	if userAgent != "" {
		client.SetUserAgent(userAgent)
	}

	return client
}

func argsOrEnv(args map[string]string, key, env string) string {
	if value := args[key]; value != "" {
		return value
	}
	return os.Getenv(env)
}

// privateIP determines if an IP is for private use (RFC1918)
// https://stackoverflow.com/a/41273687
func privateIP(ip net.IP) bool {
	return ipInCIDR(ip, "10.0.0.0/8") || ipInCIDR(ip, "172.16.0.0/12") || ipInCIDR(ip, "192.168.0.0/16")
}

func ipInCIDR(ip net.IP, CIDR string) bool {
	_, ipNet, err := net.ParseCIDR(CIDR)
	if err != nil {
		return false
	}
	return ipNet.Contains(ip)
}
