// Package oci provides node discovery for Oracle Cloud Infrastructure.
package oci

import (
	"fmt"
	"io/ioutil"
	"log"
	"context"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/resourcesearch"
	"github.com/oracle/oci-go-sdk/core"
)

type Provider struct{}

type Config struct{
	tenancyOcid    string
	userOcid       string
	keyFingerprint string
	region         string
}

func (p *Provider) Help() string {
	return `Oracle Cloud Infrastructure:

    provider:        "oci"
    tag_namespace:   The namespace of the tag
    tag_key:         The tag key to filter on
    tag_value:       The tag value to filter on
    addr_type:       "private" or "public". Defaults to "private".
    tenancy_ocid:    The Tenancy OCID of the OCI account.
    user_ocid:       The OCID of the user to use.
    key_fingerprint: The fingerprint of the key associated with the user.
    region:          The OCI region. Default to region of instance.
    
    Values for tenancy_ocid, user_ocid, key_fingerprint, and region can be omitted if these
    are supplied in ~/.oci/config as described at
    https://docs.cloud.oracle.com/iaas/Content/API/Concepts/sdkconfig.htm#FileNameandLocation.		
`
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "oci" {
		return nil, fmt.Errorf("discover-oci: invalid provider " + args["provider"])
	}

	if l == nil {
		l = log.New(ioutil.Discard, "", 0)
	}

	var tagNamespace, tagKey, tagValue string
	
	if args["tag_namespace"] != "" {
		tagNamespace = fmt.Sprintf("definedTags.namespace = '%s' && ", args["tag_namespace"])
	}
	
	if args["tag_key"] != "" {
		if tagNamespace != "" {
			tagKey = fmt.Sprintf("definedTags.key = '%s'", args["tag_key"])
		} else {
			tagKey = fmt.Sprintf("freeformTags.key = '%s'", args["tag_key"])
		}		
	}
	
	if args["tag_value"] != "" {
		if tagNamespace != "" {
			tagValue = fmt.Sprintf(" && definedTags.value = '%s'", args["tag_value"])
		} else {
			tagValue = fmt.Sprintf(" && freeformTags.value = '%s'", args["tag_value"])
		}		
	}

	query := fmt.Sprintf("query instance resources where (lifecycleState = 'RUNNING' && %s%s%s)", tagNamespace, tagKey, tagValue)

	addrType := args["addr_type"]

	var config common.ConfigurationProvider
	// if args["tenancy_ocid"] == "" || args["user_ocid"] == "" || args["key_fingerprint"] == "" || args["region"] == "" {
		// log.Printf("[DEBUG] discover-oci: Incomplete static configuration provided")
		l.Printf("[DEBUG] discover-oci: Using default values from ~/.oci/config")
		config = common.DefaultConfigProvider()
	// } else {
	// 	log.Printf("[DEBUG] discover-oci: Static configuration provided")
	// 	config = Config{args["tenancy_ocid"], args["user_ocid"], args["key_fingerprint"], args["region"]}
	// }
	
	if addrType != "private" && addrType != "public" {
		l.Printf("[INFO] discover-oci: Address type: %s. Falling back to 'private'", addrType)
		addrType = "private"
	}

	l.Printf("[DEBUG] discover-oci: Query is: %s", query)
	
	l.Printf("[DEBUG] discover-oci: Creating session...")
	search, err := resourcesearch.NewResourceSearchClientWithConfigurationProvider(config)
	if err != nil {
		return nil, fmt.Errorf("discover-oci: NewResourceSearchClientWithConfigurationProvider failed: %s", err)
	}

	l.Printf("[INFO] discover-oci: Filter instances with %s=%s", args["tag_key"], args["tag_value"])
	instances, err := search.SearchResources(
		context.Background(),
		resourcesearch.SearchResourcesRequest{
			SearchDetails: resourcesearch.StructuredSearchDetails{Query: &query},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("discover-oci: SearchResources failed: %s", err)
	}

	l.Printf("[DEBUG] discover-oci: Found %d resources", len(instances.Items))
	var addrs []string
	for _, inst := range instances.Items {
		l.Printf("[DEBUG] discover-oci: Found instance %s", *inst.Identifier)
		compute, err := core.NewComputeClientWithConfigurationProvider(config)
		if err != nil {
			return nil, fmt.Errorf("discover-oci: NewComputeClientWithConfigurationProvider failed: %s", err)
		}

		attachments, err := compute.ListVnicAttachments(
			context.Background(),
			core.ListVnicAttachmentsRequest{
				CompartmentId: inst.CompartmentId,
				InstanceId: inst.Identifier,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("discover-oci: ListVnicAttachments failed: %s", err)
		}

		l.Printf("[DEBUG] discover-oci: Instance %s has %d vnics", *inst.Identifier, len(attachments.Items))

		for _, attachment := range attachments.Items {
			l.Printf("[DEBUG] discover-oci: Checking vnic %s on Instance %s", *attachment.VnicId, *inst.Identifier)
			vnet, err := core.NewVirtualNetworkClientWithConfigurationProvider(config)
			if err != nil {
				return nil, fmt.Errorf("discover-oci: NewVirtualNetworkClientWithConfigurationProvider failed: %s", err)
			}

			vnic, err := vnet.GetVnic(
				context.Background(),
				core.GetVnicRequest{
					VnicId: attachment.VnicId,
				},
			)
			if err != nil {
				return nil, fmt.Errorf("discover-oci: GetVnic failed: %s", err)
			}

			switch addrType {
			case "private":
				l.Printf("[INFO] discover-oci: Vnic %s on instance %s has private ip %s", *vnic.Id, *inst.Identifier, *vnic.PrivateIp)
				addrs = append(addrs, *vnic.PrivateIp)
			case "public":
				if vnic.PublicIp != nil {
					l.Printf("[INFO] discover-oci: Vnic %s on instance %s has public ip %s", *vnic.Id, *inst.Identifier, *vnic.PublicIp)
					addrs = append(addrs, *vnic.PublicIp)
				} else {
					l.Printf("[INFO] discover-oci: Vnic %s on instance %s has no public ip. Skipping.", *vnic.Id, *inst.Identifier)
				}
			}			
		}
	}

	l.Printf("[DEBUG] discover-oci: Found ip addresses: %v", addrs)
	return addrs, nil
}

func (c Config) TenancyOCID() (string, error) {
	return c.tenancyOcid, nil
}

func (c Config) UserOCID() (string, error) {
	return c.userOcid, nil
}

func (c Config) KeyFingerprint() (string, error) {
	return c.keyFingerprint, nil
}

func (c Config) Region() (string, error) {
	return c.region, nil
}
