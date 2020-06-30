// Package exoscale provides node discovery for Exoscale.
package exoscale

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/exoscale/egoscale"
)

const (
	defaultAPIEndpoint = "https://api.exoscale.com/v1"
)

type Provider struct{}

func (p *Provider) Help() string {
	return `Exoscale:

    provider:       "exoscale"
    api_key:        The Exoscale API key (required)
    api_secret:     The Exoscale API secret key (required)
    api_endpoint:   The Exoscale API endpoint
    tag_key:        The tag key to filter on
    tag_value:      The tag value to filter on
    zone:           The Exoscale Zone
    addr_type:      "public_v4" or "public_v6". (default: "public_v4")
`
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "exoscale" {
		return nil, fmt.Errorf("discover-exoscale: invalid provider " + args["provider"])
	}

	if l == nil {
		l = log.New(ioutil.Discard, "", 0)
	}

	client, err := buildClient(args)
	if err != nil {
		return nil, err
	}

	tagKey := args["tag_key"]
	tagValue := args["tag_value"]
	zone := args["zone"]
	addrType := args["addr_type"]

	if addrType == "" {
		addrType = "public_v4"
	}

	if addrType != "public_v4" && addrType != "public_v6" {
		l.Printf("[INFO] discover-exoscale: Address type %s is not supported. Valid values are public_v4 or public_v6. Falling back to 'public_v4'", addrType)
		addrType = "public_v4"
	}

	l.Printf(
		"[DEBUG] discover-exoscale: Using zone=%s tag_key=%s, tag_value=%s, addr_type=%s",
		zone,
		tagKey,
		tagValue,
		addrType,
	)

	var zoneID *egoscale.UUID
	if zone != "" {
		resp, err := client.Get(egoscale.Zone{Name: zone})
		if err != nil {
			return nil, err
		}
		z := resp.(*egoscale.Zone)
		zoneID = z.ID
	}

	var tags []egoscale.ResourceTag
	if tagKey != "" && tagValue != "" {
		tags = []egoscale.ResourceTag{
			{
				Key:   tagKey,
				Value: tagValue,
			},
		}
	}

	resp, err := client.RequestWithContext(context.Background(), egoscale.ListVirtualMachines{
		ZoneID: zoneID,
		Tags:   tags,
	})
	if err != nil {
		return nil, err
	}

	instances := resp.(*egoscale.ListVirtualMachinesResponse).VirtualMachine

	var addrs []string
	for _, instance := range instances {
		ip := instance.IP().String()
		if addrType == "public_v6" {
			nic := instance.DefaultNic()
			if nic != nil {
				if nic.IP6Address.String() == "<nil>" {
					continue
				}

				ip = nic.IP6Address.String()
			}
		}

		l.Printf("[DEBUG] discover-exoscale: Found instance %s - with %s IP: %s",
			instance.Name,
			addrType,
			ip,
		)

		addrs = append(addrs, ip)
	}

	l.Printf("[DEBUG] discover-exoscale: Found ip addresses: %v", addrs)
	return addrs, nil
}

func buildClient(args map[string]string) (*egoscale.Client, error) {
	apiKey := args["api_key"]
	apiSecret := args["api_secret"]
	apiEndpoint := args["api_endpoint"]

	apiKeyEnv := os.Getenv("EXOSCALE_API_KEY")
	apiSecretEnv := os.Getenv("EXOSCALE_API_SECRET")
	apiEndpointEnv := os.Getenv("EXOSCALE_API_ENDPOINT")

	if apiKeyEnv != "" {
		apiKey = apiKeyEnv
	}
	if apiSecretEnv != "" {
		apiSecret = apiSecretEnv
	}
	if apiEndpointEnv != "" {
		apiEndpoint = apiEndpointEnv
	}

	if apiKey == "" || apiSecret == "" {
		return nil, errors.New("incomplete or missing for API credentials")
	}

	if apiEndpoint == "" {
		apiEndpoint = defaultAPIEndpoint
	}

	return egoscale.NewClient(apiEndpoint, apiKey, apiSecret), nil
}
