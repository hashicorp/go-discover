// Copyright IBM Corp. 2017, 2026
// SPDX-License-Identifier: MPL-2.0

// Package azure provides node discovery for Microsoft Azure.
package azure

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
)

type Provider struct {
	userAgent string
}

func (p *Provider) SetUserAgent(s string) {
	p.userAgent = s
}

func (p *Provider) Help() string {
	return `Microsoft Azure:

   provider:          "azure"
   tenant_id:         The id of the tenant
   client_id:         The id of the client
   subscription_id:   The id of the subscription
   secret_access_key: The authentication credential
   msft_telemetry_opt_in: Optional boolean string to opt in to sending telemetry to Microsoft

    **NOTE** The secret_access_key value often may have an equals sign in it's value,
    especially if generated from the Azure Portal. So is important to wrap in single quotes
    eg. secret_acccess_key='fpOfcHQJAQBczjAxiVpeyLmX1M0M0KPBST+GU2GvEN4='

   Variables can also be provided by environmental variables:
    export ARM_SUBSCRIPTION_ID for subscription
    export ARM_TENANT_ID for tenant
    export ARM_CLIENT_ID for client
    export ARM_CLIENT_SECRET for secret access key

   Set the following environment variables to enable AzureSDK client's log and telemetry:
	export AZURE_SDK_GO_LOGGING=all
	export OPT_IN_MSFT_TELEMETRY=true


   If none of those options are given, the Azure SDK is using the default  environment based authentication outlined
   here https://docs.microsoft.com/en-us/go/azure/azure-sdk-go-authorization#use-environment-based-authentication
   This will fallback to MSI if nothing is explicitly specified.

   Use these configuration parameters when using tags:

   tag_name:          The name of the tag to filter on
   tag_value:         The value of the tag to filter on

   Use these configuration parameters when using Virtual Machine Scale Sets:

   resource_group:    The name of the resource group to filter on
   vm_scale_set:      The name of the virtual machine scale set to filter on

   When using tags the only permission needed is Microsoft.Network/networkInterfaces/*

   When using Virtual Machine Scale Sets the only role action needed is Microsoft.Compute/virtualMachineScaleSets/*/read.
   The Azure provider only supports Virtual Machine Scale Sets deployed in [Uniform mode](https://learn.microsoft.com/en-us/azure/virtual-machine-scale-sets/virtual-machine-scale-sets-orchestration-modes#scale-sets-with-uniform-orchestration).
   As of 2023 VMSS deploys using Flexible mode by default.

   It is recommended you make a dedicated key used only for auto-joining.
`
}

// argsOrEnv allows you to pick an environmental variable for a setting if the arg is not set
func argsOrEnv(args map[string]string, key, env string) string {
	if value, ok := args[key]; ok {
		return value
	}
	return os.Getenv(env)
}

// Prepare helper functions for overriding UserAgent behavior in policy.ClientOptions

// policyFunc implements the azcore Policy interface and is used to
// set a custom user agent in the azure client configuration
type policyFunc func(*policy.Request) (*http.Response, error)

// Do implements the Policy interface on policyFunc
func (pf policyFunc) Do(req *policy.Request) (*http.Response, error) {
	return pf(req)
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "azure" {
		return nil, fmt.Errorf("discover-azure: invalid provider %s", args["provider"])
	}

	if l == nil {
		l = log.New(io.Discard, "", 0)
	}

	// check for environmental variables, and use if the argument hasn't been set in config
	tenantID := argsOrEnv(args, "tenant_id", "ARM_TENANT_ID")
	clientID := argsOrEnv(args, "client_id", "ARM_CLIENT_ID")
	subscriptionID := argsOrEnv(args, "subscription_id", "ARM_SUBSCRIPTION_ID")
	secretKey := argsOrEnv(args, "secret_access_key", "ARM_CLIENT_SECRET")

	var clientPolicies []policy.Policy
	// AzureSDK clients create their own connection pipelines and inherit
	// from the credential config's policy.ClientOptions object
	// We will:
	// -Relay any user provided userAgent configuration
	// -Turn off MSFT telemetry unless go-discover users opt-in.

	if p.userAgent != "" {
		policyFunc := policyFunc(func(req *policy.Request) (*http.Response, error) {
			req.Raw().Header.Set("User-Agent", p.userAgent)
			return req.Next()
		})
		clientPolicies = append(clientPolicies, policyFunc)
	}

	telemetryOptIn, err := strconv.ParseBool(argsOrEnv(args, "msft_telemetry_opt_in", "OPT_IN_MSFT_TELEMETRY"))
	if err != nil {
		telemetryOptIn = false
	}
	clientOpts := policy.ClientOptions{
		PerCallPolicies: clientPolicies,
		Telemetry: policy.TelemetryOptions{
			// Toggle the telemetryOptIn bool to align our opt-in config with their opt-out config
			// OPT_IN_MSFT_TELEMETRY = false | Telemetry.Disabled = true (our default)
			// OPT_IN_MSFT_TELEMETRY = true | Telemetry.Disabled = false
			Disabled: !telemetryOptIn,
		},
	}

	// Try to use the argument and environment provided arguments first, if this fails fall back to the SDK's
	// DefaultCredential which attempts to find default config for Azure envars, AzureWorkloadIdentityCredentials,
	// Azure ManagedIdentityCredentials, as well as local credentials
	// https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#DefaultAzureCredential

	var (
		clientSecretCred *azidentity.ClientSecretCredential
		defaultCred      *azidentity.DefaultAzureCredential
	)
	if tenantID != "" && clientID != "" && secretKey != "" {
		var err error
		clientSecretCred, err = azidentity.NewClientSecretCredential(tenantID, clientID, secretKey,
			&azidentity.ClientSecretCredentialOptions{ClientOptions: clientOpts})

		if err != nil {
			return nil, fmt.Errorf("discover-azure (ClientCredentials): %w", err)
		}
	} else {
		var err error
		defaultCred, err = azidentity.NewDefaultAzureCredential(&azidentity.DefaultAzureCredentialOptions{ClientOptions: clientOpts})
		if err != nil {
			return nil, fmt.Errorf("discover-azure (EnvironmentCredentials): %w", err)
		}
	}

	// Use tags if using network interfaces
	tagName := args["tag_name"]
	tagValue := args["tag_value"]

	// Use resourceGroup and vmScaleSet if using vm scale sets
	resourceGroup := args["resource_group"]
	vmScaleSet := args["vm_scale_set"]

	if subscriptionID == "" {
		return nil, fmt.Errorf("discover-azure (Credentials): subscription_id not provided as argument or environment variable")
	}
	// Create NetworkInterfaceClient with the appropriate credential and no additional configuration
	// as the client will inherit telemetry config from the credentials
	var vmnet *armnetwork.InterfacesClient
	if clientSecretCred != nil {
		vmnet, err = armnetwork.NewInterfacesClient(subscriptionID, clientSecretCred, nil)
		if err != nil {
			return nil, fmt.Errorf("discover-azure (Azure Client): %w", err)
		}
	} else {
		vmnet, err = armnetwork.NewInterfacesClient(subscriptionID, defaultCred, nil)
		if err != nil {
			return nil, fmt.Errorf("discover-azure (Azure Client): %w", err)
		}
	}

	if tagName != "" && tagValue != "" && resourceGroup == "" && vmScaleSet == "" {
		l.Printf("[DEBUG] discover-azure: using tag method. tag_name: %s, tag_value: %s", tagName, tagValue)
		return fetchAddrsWithTags(tagName, tagValue, *vmnet, l)
	} else if resourceGroup != "" && vmScaleSet != "" && tagName == "" && tagValue == "" {
		l.Printf("[DEBUG] discover-azure: using vm scale set method. resource_group: %s, vm_scale_set: %s", resourceGroup, vmScaleSet)
		return fetchAddrsWithVmScaleSet(resourceGroup, vmScaleSet, *vmnet, l)
	} else {
		l.Printf("[ERROR] discover-azure: tag_name: %s, tag_value: %s", tagName, tagValue)
		l.Printf("[ERROR] discover-azure: resource_group %s, vm_scale_set %s", resourceGroup, vmScaleSet)
		return nil, fmt.Errorf("discover-azure: unclear configuration. use (tag name and value) or (resouce_group and vm_scale_set)")
	}
}

func fetchAddrsWithTags(tagName string, tagValue string, vmnet armnetwork.InterfacesClient, l *log.Logger) ([]string, error) {
	// Get all network interfaces across resource groups
	// unless there is a compelling reason to restrict

	ctx := context.Background()
	pager := vmnet.NewListAllPager(nil)
	var addrs []string
	for pager.More() {

		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("discover-azure: %w", err)
		}
		if len(page.Value) == 0 {
			return nil, fmt.Errorf("discover-azure: no interfaces")
		}
		// Collect any PrivateIPAddress with the matching tag on each page of results
		for _, v := range page.Value {
			var id string
			if v.ID != nil {
				id = *v.ID
			} else {
				id = "unknown_interface_id"
			}
			if v.Tags == nil {
				l.Printf("[DEBUG] discover-azure: Interface %s has no tags", id)
				continue
			}
			tv := v.Tags[tagName] // *string
			if tv == nil {
				l.Printf("[DEBUG] discover-azure: Interface %s did not have tag: %s", id, tagName)
				continue
			}
			if *tv != tagValue {
				l.Printf("[DEBUG] discover-azure: Interface %s tag value was: %s which did not match: %s", id, *tv, tagValue)
				continue
			}
			if v.Properties == nil {
				l.Printf("[DEBUG] discover-azure: Interface %s had no properties", id)
				continue
			}
			for _, x := range v.Properties.IPConfigurations {
				if x.Properties.PrivateIPAddress == nil {
					l.Printf("[DEBUG] discover-azure: Interface %s had no private ip", id)
					continue
				}
				iAddr := *x.Properties.PrivateIPAddress
				l.Printf("[DEBUG] discover-azure: Interface %s has private ip: %s", id, iAddr)
				addrs = append(addrs, iAddr)
			}
		}
		l.Printf("[DEBUG] discover-azure: Found ip addresses: %v", addrs)
	}

	return addrs, nil
}

func fetchAddrsWithVmScaleSet(resourceGroup string, vmScaleSet string, vmnet armnetwork.InterfacesClient, l *log.Logger) ([]string, error) {
	// Get all network interfaces for a specific virtual machine scale set
	ctx := context.Background()
	pager := vmnet.NewListVirtualMachineScaleSetNetworkInterfacesPager(resourceGroup, vmScaleSet, nil)
	var addrs []string

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("discover-azure: %w", err)
		}
		if len(page.Value) == 0 {
			return nil, fmt.Errorf("discover-azure: no interfaces")
		}
		// Collect all of PrivateIPAddresses we can
		for _, v := range page.Value {
			var id string
			if v.ID != nil {
				id = *v.ID
			} else {
				id = "unknown_interface_id"
			}
			if v.Properties == nil {
				l.Printf("[DEBUG] discover-azure: Interface %s had properties", id)
				continue
			}

			for _, x := range v.Properties.IPConfigurations {
				if x.Properties.PrivateIPAddress == nil {
					l.Printf("[DEBUG] discover-azure: Interface %s had no private ip", id)
					continue
				}
				iAddr := *x.Properties.PrivateIPAddress
				l.Printf("[DEBUG] discover-azure: Interface %s has private ip: %s", id, iAddr)
				addrs = append(addrs, iAddr)
			}
		}
		l.Printf("[DEBUG] discover-azure: Found ip addresses: %v", addrs)
	}
	return addrs, nil
}
