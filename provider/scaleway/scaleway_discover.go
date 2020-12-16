// Package scaleway provides node discovery for Scaleway.
package scaleway

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type Provider struct{}

func (p *Provider) Help() string {
	return `Scaleway:

    provider:           "scaleway"
    name:               The name of the instance to filter on
    organization:       The Scaleway organization access key
    tag_name:           The tag names to filter on (use "," as a separator)
    commercial_type:    The commercial type to filter on
    private_network_id: The private network id to filter on
    token:              The Scaleway API secret key
    zone:               The Scaleway zone (default fr-par-1)
    addr_type:          "private_v4", "public_v4" or "public_v6". Defaults to "private_v4".

	The following location are looked for in that order:
	1. Configuration file
	2. Environment variables
	3. Command arguments
`
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "scaleway" {
		return nil, fmt.Errorf("discover-scaleway: invalid provider " + args["provider"])
	}

	if l == nil {
		l = log.New(ioutil.Discard, "scaleway", 0)
	}

	addrType := args["addr_type"]
	tagName := args["tag_name"]
	commercialType := args["commercial_type"]
	privateNetwork := args["private_network_id"]
	name := args["name"]

	// Create a new API client
	scwClient, err := newClient(args, l)
	if err != nil {
		return nil, fmt.Errorf("discover-scaleway: %s", err)
	}

	instanceAPI := instance.NewAPI(scwClient)
	req := &instance.ListServersRequest{}

	if commercialType != "" {
		req.CommercialType = &commercialType
	}

	if privateNetwork != "" {
		req.PrivateNetwork = &privateNetwork
	}

	if name != "" {
		req.Name = &name
	}

	if tagName != "" {
		req.Tags = strings.Split(tagName, ",")
	}

	if addrType != "private_v4" && addrType != "public_v4" && addrType != "public_v6" {
		l.Printf("[INFO] discover-scaleway: Address type %q is not supported. Valid values for addr_type are {private_v4,public_v4,public_v6}. Falling back to 'private_v4'", addrType)
		addrType = "private_v4"
	}

	servers, err := instanceAPI.ListServers(req, scw.WithAllPages())
	if err != nil {
		return nil, fmt.Errorf("discover-scaleway: %s", err)
	}

	var addrs []string
	for _, server := range servers.Servers {
		a := ""
		switch addrType {

		case "public_v4":
			if server.PublicIP != nil {
				a = server.PublicIP.Address.String()
			}
		case "public_v6":
			if server.IPv6 != nil {
				a = server.IPv6.Address.String()
			}

		case "private_v4":
			a = *server.PrivateIP
		}

		if a != "" {
			addrs = append(addrs, a)
		}
	}

	l.Printf("discover-scaleway: Found ip addresses: %v", addrs)
	return addrs, nil
}

func newClient(args map[string]string, l *log.Logger) (*scw.Client, error) {
	// By default we set default zone and region to fr-par
	defaultZoneProfile := &scw.Profile{
		DefaultRegion: scw.StringPtr(scw.RegionFrPar.String()),
		DefaultZone:   scw.StringPtr(scw.ZoneFrPar1.String()),
	}

	// config file loading
	config, err := scw.LoadConfig()
	// If the config file do not exist, don't return an error as we may find config in ENV or flags.
	if _, isNotFoundError := err.(*scw.ConfigFileNotFoundError); isNotFoundError {
		config = &scw.Config{}
	} else if err != nil {
		return nil, err
	}
	configProfile := &config.Profile

	// Environment config loading
	envProfile := scw.LoadEnvProfile()

	// Command line argument loading
	argsProfile := &scw.Profile{}

	projectID := args["project-id"]
	if projectID != "" {
		argsProfile.DefaultProjectID = scw.StringPtr(projectID)
	}

	secretKey := args["secret-key"]
	if secretKey != "" {
		argsProfile.SecretKey = scw.StringPtr(secretKey)
	}

	zone := args["zone"]
	if zone != "" {
		z, err := scw.ParseZone(zone)
		if err != nil {
			return nil, fmt.Errorf("error during zone parsing: %s", err)
		}
		argsProfile.DefaultZone = scw.StringPtr(z.String())
	}

	organization := args["organization"]
	if organization != "" {
		argsProfile.DefaultOrganizationID = scw.StringPtr(organization)
	}

	profile := scw.MergeProfiles(defaultZoneProfile, configProfile, envProfile, argsProfile)

	l.Printf("[INFO] discover-scaleway: Organization is %q", *profile.DefaultOrganizationID)
	l.Printf("[INFO] discover-scaleway: Project is %q", *profile.DefaultProjectID)
	l.Printf("[INFO] discover-scaleway: Secret is %q", hideSecretKey(*profile.SecretKey))
	l.Printf("[INFO] discover-scaleway: Zone is %q", *profile.DefaultZone)

	opts := []scw.ClientOption{
		scw.WithUserAgent(fmt.Sprintf("go-discover")),
		scw.WithProfile(profile),
	}

	scwClient, err := scw.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	return scwClient, nil
}

func hideSecretKey(k string) string {
	switch {
	case len(k) == 0:
		return ""
	case len(k) > 8:
		return k[0:8] + "-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	default:
		return "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	}
}
