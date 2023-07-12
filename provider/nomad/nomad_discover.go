// Package nomad provides allocation discovery for Nomad.
package nomad

import (
	"fmt"
	"log"

	"github.com/hashicorp/nomad/api"
)

type Provider struct{
	client *api.Client
}

func (p *Provider) Help() string {
	return `Nomad (Nomad):

    provider:         "nomad"
    namespace:        Namespace to search for allocations (defaults to "default")
    service_name:     Nomad service to discover allocations for
		region:     			Nomad region to discover allocations for (optional)
`
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "nomad" {
		return nil, fmt.Errorf("discover-nomad: invalid provider " + args["provider"])
	}

	namespace := args["namespace"]
	if namespace == "" {
		namespace = "default"
	}
	p.client.SetNamespace(namespace)

	serviceName := args["service_name"]
	if serviceName == "" {
		return nil, fmt.Errorf("discover-nomad: must provider a service_name")
	}

	region := args["region"]
	if region == "" {
		p.client.SetRegion(region)
	}

	queryOpts := &api.QueryOptions{}
	services, _, err := p.client.Services().Get(serviceName, queryOpts)
	if err != nil {
		return nil, fmt.Errorf("discover-nomad: error retrieving services: %v", err)
	}

	addrs := []string{}
	for _, service := range services {
		addr := service.Address
		port := service.Port
		addrs = append(addrs, fmt.Sprintf("%s:%d", addr, port)
	}

	return addrs, nil
}
