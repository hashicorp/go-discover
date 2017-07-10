// Package softlayer provides node discovery for Softlayer.
package softlayer

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-discover/config"
	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

// Discover returns the private ip addresses of all Softlayer
// instances in a datacenter with a certain tag value.
//
// cfg contains the configuration in "key=val key=val ..." format. The
// values are URL encoded.
//
// The supported keys are:
//
//   datacenter: The SoftLayer datacenter to filter on
//   tag_value:  The tag value to filter on
//   username:   The SoftLayer username to use
//   api_key:    The SoftLayer api key to use
//
// Example:
//
//  provider=softlayer datacenter=dal06 tag_value=consul username=... api_key=...
//
func Discover(cfg string, l *log.Logger) ([]string, error) {
	m, err := config.Parse(cfg)
	if err != nil {
		return nil, fmt.Errorf("discover-softlayer: %s", err)
	}

	datacenter := m["datacenter"]
	tagValue := m["tag_value"]
	username := m["username"]
	apiKey := m["api_key"]

	l.Printf("[INFO] discover-softlayer: Datacenter is %q", datacenter)

	// Create a session and get a service
	sess := session.New(username, apiKey)
	service := services.GetAccountService(sess)

	// Compose the filter
	mask := "id,hostname,domain,tagReferences[tag[name]],primaryBackendIpAddress,datacenter"
	filterVMs := filter.Build(
		filter.Path("virtualGuests.datacenter.name").Eq(datacenter),
		filter.Path("virtualGuests.tagReferences.tag.name").Eq(tagValue),
	)

	// Get the virtual machines that match the filter
	vms, err := service.Mask(mask).Filter(filterVMs).GetVirtualGuests()
	if err != nil {
		return nil, fmt.Errorf("discover-softlayer: %s", err)
	}

	var addrs []string
	for _, vm := range vms {
		l.Printf("[INFO] discover-softlayer: Found instance (%d) %s.%s with private IP: %s",
			*vm.Id, *vm.Hostname, *vm.Domain, *vm.PrimaryBackendIpAddress)
		addrs = append(addrs, *vm.PrimaryBackendIpAddress)
	}
	return addrs, nil
}
