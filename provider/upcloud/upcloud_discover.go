package upcloud

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/upcloud/client"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/request"
	"github.com/UpCloudLtd/upcloud-go-api/upcloud/service"
)

// Provider defines a provider for UpCloud to discover server IP addresses
type Provider struct{}

func (p *Provider) Help() string {
	return `UpCloud:

    provider:          "upcloud"
    tag:               The tag to filter servers on
    title_match:       A regular expression to filter server titles by
    username:          The UpCloud username (alternative env var: UPCLOUD_API_USERNAME)
    password:          The UpCloud password (alternative env var: UPCLOUD_API_PASSWORD)
`
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	username := args["username"]
	password := args["password"]
	discoverTag := args["tag"]
	titleMatch := args["title_match"]

	if username == "" {
		username = os.Getenv("UPCLOUD_API_USERNAME")
	}
	if password == "" {
		password = os.Getenv("UPCLOUD_API_PASSWORD")
	}
	if discoverTag == "" && titleMatch == "" {
		return nil, fmt.Errorf("must provider either a search tag or a title_match regular expression")
	}

	c := client.New(username, password)
	// Set some longer timeout
	c.SetTimeout(time.Second * 30)

	svc := service.New(c)
	servers, err := svc.GetServers()
	if err != nil {
		return nil, fmt.Errorf("getting servers from upcloud: %w", err)
	}

	var ipAddrs = make([]string, 0)
	titleMatcher, err := regexp.Compile(titleMatch)
	if err != nil {
		return nil, fmt.Errorf("invalid title_match regular expression: %w", err)
	}

	// Iterate over servers and check for matching tags and get the IP address
	for _, server := range servers.Servers {
		var hasTag bool
		for _, tag := range server.Tags {
			if discoverTag != "" && tag == discoverTag {
				hasTag = true
			}
		}
		// Skip a tag was given to search for, but the server doesn't have that tag, then skip
		if discoverTag != "" && !hasTag {
			continue
		}
		// Skip server if no titleMatch given and the given titleMatch does not match
		if titleMatch != "" && !titleMatcher.MatchString(server.Title) {
			continue
		}
		details, err := svc.GetServerDetails(&request.GetServerDetailsRequest{
			UUID: server.UUID,
		})
		if err != nil {
			return nil, fmt.Errorf("getting server details from upcloud with UUID: %s: %w", server.UUID, err)
		}
		// Primitive, but let's just take the first IP address
		ipAddrs = append(ipAddrs, details.IPAddresses[0].Address)
	}
	l.Printf("[DEBUG] discover-upcloud: Found ip addresses: %v", ipAddrs)
	return ipAddrs, nil
}
