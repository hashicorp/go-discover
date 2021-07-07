package equinixmetal

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/packethost/packngo"
)

const baseURL = "https://api.equinix.com/metal/v1/"

var (
	addressTypes       = []string{packngo.PublicIPv4, packngo.PrivateIPv4, packngo.PublicIPv6}
	defaultAddressType = packngo.PrivateIPv4
)

// Provider struct
type Provider struct {
	userAgent string
}

// SetUserAgent setter
func (p *Provider) SetUserAgent(s string) {
	p.userAgent = s
}

// Help function
func (p *Provider) Help() string {
	return fmt.Sprintf(`Equinix Metal:
	provider:       "equinixmetal"
	project:        UUID of metal project. Required
	auth_token:     Equinix Metal authentication token. Required
	url:            Equinix Metal REST URL. Optional
	address_type:   Address type, one of: %s. Defaults to %q. Optional
	facility:       Filter for specific facility (Examples: "sv15,ny5")
	metro:          Filter for specific metro (Examples: "sv,ny")
	tag:            Filter by tag (Examples: "tag1,tag2")
	
	Variables can also be provided by environmental variables:
	export METAL_PROJECT for project
	export METAL_URL for url
	export METAL_AUTH_TOKEN for auth_token
`, strings.Join(addressTypes, ", "), defaultAddressType)
}

// Addrs function
func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	authToken := argsOrEnv(args, "auth_token", "METAL_AUTH_TOKEN")
	projectID := argsOrEnv(args, "project", "METAL_PROJECT")
	metalURL := argsOrEnv(args, "url", "METAL_URL")
	addressType := args["address_type"]
	metalFacilities := args["facility"]
	metalMetros := args["metro"]
	metalTags := args["tag"]

	if !Include(addressTypes, addressType) {
		l.Printf("[INFO] discover-metal: Address type %s is not supported. Valid values are {%s}. Falling back to '%s'", addressType, strings.Join(addressTypes, ","), defaultAddressType)
		addressType = defaultAddressType
	}

	includeMetros := includeArgs(metalMetros)
	includeFacilities := includeArgs(metalFacilities)
	includeTags := includeArgs(metalTags)

	c, err := client(p.userAgent, metalURL, authToken)
	if err != nil {
		return nil, fmt.Errorf("discover-metal: Initializing Equinix Metal client failed: %s", err)
	}

	var devices []packngo.Device

	if projectID == "" {
		return nil, fmt.Errorf("discover-metal: 'project' parameter must be provided")
	}

	getOpts := &packngo.GetOptions{}
	getOpts.Including("facility", "metro", "ip_addresses")
	getOpts.Excluding("ssh_keys", "project")
	devices, _, err = c.Devices.List(projectID, getOpts)
	if err != nil {
		return nil, fmt.Errorf("discover-metal: Fetching Equinix Metal devices failed: %s", err)
	}

	var addrs []string
	for _, d := range devices {

		if len(includeFacilities) > 0 && !Include(includeFacilities, d.Facility.Code) {
			continue
		}

		if len(includeMetros) > 0 && !Include(includeMetros, d.Metro.Code) {
			continue
		}

		if len(includeTags) > 0 && !Any(d.Tags, func(v string) bool { return Include(includeTags, v) }) {
			continue
		}

		for _, n := range d.Network {
			if ipMatchesType(addressType, n) {
				addrs = append(addrs, n.Address)
			}
		}
	}
	return addrs, nil
}

func ipMatchesType(addressType string, n *packngo.IPAddressAssignment) bool {
	switch addressType {
	case packngo.PublicIPv4:
		return n.Public && n.AddressFamily == 4
	case packngo.PublicIPv6:
		return n.Public && n.AddressFamily == 6
	case packngo.PrivateIPv4:
		return !n.Public && n.AddressFamily == 4
	default:
		return false
	}
}

func client(useragent, url, token string) (*packngo.Client, error) {
	if url == "" {
		url = baseURL
	}

	client, err := packngo.NewClientWithBaseURL(useragent, token, nil, url)
	if err == nil {
		client.UserAgent = fmt.Sprintf("%s %s", useragent, client.UserAgent)
	}
	return client, err
}

func argsOrEnv(args map[string]string, key, env string) string {
	if value := args[key]; value != "" {
		return value
	}
	return os.Getenv(env)
}

func includeArgs(s string) []string {
	var include []string
	for _, localstring := range strings.Split(s, ",") {
		if len(localstring) > 0 {
			include = append(include, localstring)
		}
	}
	return include
}

// Index returns the first index of the target string t, or -1 if no match is found.
func Index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

// Include returns true if the target string t is in the slice.
func Include(vs []string, t string) bool {
	return Index(vs, t) >= 0
}

//Any returns true if one of the strings in the slice satisfies the predicate f.
func Any(vs []string, f func(string) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}
