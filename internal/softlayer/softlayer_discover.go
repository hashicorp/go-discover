// Package softlayer provides node discovery for Softlayer.
package softlayer

import (
	"fmt"
	"log"

	"github.com/softlayer/softlayer-go/filter"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

func Discover(m map[string]string, l *log.Logger) ([]string, error) {
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
