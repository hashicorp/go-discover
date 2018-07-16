package packet

import (
	"fmt"
	"log"
	"os"

	"github.com/packethost/packngo"
)

// Provider struct
type Provider struct{}

const baseURL = "https://api.packet.net/"

// Help function
func (p *Provider) Help() string {
	return `Packet:
	provider:		     "packet"
	// packet_organization: UUID of packet organization
	packet_project: 	 UUID of packet project
	packet_url: 		 Packet REST URL
	packet_auth_token:   Packet authentication token
	`
}

// Addrs function
func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	authToken := argsOrEnv(args, "packet_auth_token", "PACKET_AUTH_TOKEN")
	// organizationID := argsOrEnv(args, "packet_organization", "PACKET_ORGANIZATION")
	projectID := argsOrEnv(args, "packet_project", "PACKET_PROJECT")
	packetURL := argsOrEnv(args, "packet_url", "PACKET_URL")
	fmt.Println("ProjectID", projectID)
	c, err := client(packetURL, authToken)
	if err != nil {
		return nil, fmt.Errorf("discover-packet: Initializing Packet client failed: %s", err)
	}

	var devices []packngo.Device

	if projectID == "" {
		return nil, fmt.Errorf("discover-packet: 'packet_project' parameter must be provider")
	}

	devices, _, err = c.Devices.List(projectID, nil)
	if err != nil {
		return nil, fmt.Errorf("discover-packet: Fetching Packet devices failed: %s", err)
	}
	var addrs []string
	for _, d := range devices {
		for _, n := range d.Network {
			addrs = append(addrs, n.Address)
		}
	}
	return addrs, nil
}

func client(url, token string) (*packngo.Client, error) {
	if url == "" {
		url = baseURL
	}

	return packngo.NewClientWithBaseURL("packet go-discover", token, nil, url)
}
func argsOrEnv(args map[string]string, key, env string) string {
	if value := args[key]; value != "" {
		return value
	}
	return os.Getenv(env)
}
