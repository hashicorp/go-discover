package upcloud

import (
	"fmt"
	"log"
	"os"
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
    username:          The UpCloud username (alternative env var: UPCLOUD_API_USERNAME)
	password:          The UpCloud password (alternative env var: UPCLOUD_API_PASSWORD)
`
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	username := args["username"]
	password := args["password"]
	discoverTag := args["tag"]

	if username == "" {
		username = os.Getenv("UPCLOUD_API_USERNAME")
	}
	if password == "" {
		password = os.Getenv("UPCLOUD_API_PASSWORD")
	}
	if discoverTag == "" {
		return nil, fmt.Errorf("must provider a search tag")
	}

	c := client.New(username, password)
	// Set some longer timeout
	c.SetTimeout(time.Second * 30)

	svc := service.New(c)
	servers, err := svc.GetServers()
	if err != nil {
		return nil, fmt.Errorf("getting servers from upcloud: %x", err)
	}

	var ipAddrs = make([]string, 0)
	// Iterate over servers and check for matching tags and get the IP address
	for _, server := range servers.Servers {
		for _, tag := range server.Tags {
			if tag == discoverTag {
				details, err := svc.GetServerDetails(&request.GetServerDetailsRequest{
					UUID: server.UUID,
				})
				if err != nil {
					return nil, fmt.Errorf("getting server details from upcloud with UUID: %s: %x", server.UUID, err)
				}
				// Primitive, but let's just take the first IP address
				ipAddrs = append(ipAddrs, details.IPAddresses[0].Address)
			}
		}
	}
	l.Printf("[DEBUG] discover-upcloud: Found ip addresses: %v", ipAddrs)
	return ipAddrs, nil
}
