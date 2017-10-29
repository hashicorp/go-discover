// Package scaleway provides node discovery for Scaleway.
package scaleway

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/nicolai86/scaleway-sdk/api"
)

type Provider struct{}

func (p *Provider) Help() string {
	return `Scaleway:
	provider: "scaleway"
	organization: The Scaleway organization access key
	tag_name: The tag name to filter on
	api_key: The Scaleway api access token
	region: The Scalway region
`
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "scaleway" {
		return nil, fmt.Errorf("discover-scaleway: invalid provider " + args["provider"])
	}

	if l == nil {
		l = log.New(ioutil.Discard, "", 0)
	}

	organization := args["organization"]
	tagName := args["tag_name"]
	apiKey := args["api_key"]
	region := args["region"]

	l.Printf("[INFO] discover-scaleway: Organization is %q", organization)
	l.Printf("[INFO] discover-scaleway: Region is %q", region)

	// Create a new API client
	api, err := api.New(
		organization,
		apiKey,
		region,
	)
	if err != nil {
		return nil, fmt.Errorf("discover-scaleway: %s", err)
	}

	// Currently fetching all servers since the API doesn't support
	// filter options
	servers, err := api.GetServers(true, 0)
	if err != nil {
		return nil, fmt.Errorf("discover-scaleway: %s", err)
	}

	// Get a list of private ips that match the tag name
	return filterServersForTagName(servers, tagName), nil
}

func filterServersForTagName(servers *[]api.ScalewayServer, tagName string) []string {
	var serverAddrs []string
	for _, server := range *servers {
		if stringInSlice(tagName, server.Tags) {
			l.Printf("[INFO] discover-scaleway: Found server (%d) - %s with private IP: %s",
				server.Name, server.Hostname, server.PrivateIP)
			serverAddrs = append(serverAddrs, server.PrivateIP)
		}
	}

	return serverAddrs
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
