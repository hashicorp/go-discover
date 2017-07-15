// Package azure provides node discovery for Microsoft Azure.
package azure

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/arm/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

func Discover(m map[string]string, l *log.Logger) ([]string, error) {
	tenantID := m["tenant_id"]
	clientID := m["client_id"]
	subscriptionID := m["subscription_id"]
	secretKey := m["secret_access_key"]
	tagName := m["tag_name"]
	tagValue := m["tag_value"]

	// Only works for the Azure PublicCLoud for now; no ability to test other Environment
	oauthConfig, err := azure.PublicCloud.OAuthConfigForTenant(tenantID)
	if err != nil {
		return nil, fmt.Errorf("discover-azure: %s", err)
	}

	// Get the ServicePrincipalToken for use searching the NetworkInterfaces
	sbt, err := azure.NewServicePrincipalToken(*oauthConfig, clientID, secretKey, azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		return nil, fmt.Errorf("discover-azure: %s", err)
	}

	// Setup the client using autorest; followed the structure from Terraform
	vmnet := network.NewInterfacesClient(subscriptionID)
	vmnet.Client.UserAgent = "Hashicorp-Consul"
	vmnet.Sender = autorest.CreateSender(autorest.WithLogging(l))
	vmnet.Authorizer = sbt

	// Get all network interfaces across resource groups
	// unless there is a compelling reason to restrict
	netres, err := vmnet.ListAll()
	if err != nil {
		return nil, fmt.Errorf("discover-azure: %s", err)
	}

	if netres.Value == nil {
		return nil, fmt.Errorf("discover-azure: no interfaces")
	}

	// Choose any PrivateIPAddress with the matching tag
	var addrs []string
	for _, v := range *netres.Value {
		if v.Tags == nil {
			continue
		}
		tv := (*v.Tags)[tagName] // *string
		if tv == nil || *tv != tagValue {
			continue
		}
		if v.IPConfigurations == nil {
			continue
		}
		for _, x := range *v.IPConfigurations {
			if x.PrivateIPAddress == nil {
				continue
			}
			addrs = append(addrs, *x.PrivateIPAddress)
		}
	}
	return addrs, nil
}
