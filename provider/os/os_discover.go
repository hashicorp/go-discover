// Package os provides node discovery for Openstack.
package os

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/gophercloud/utils/terraform/auth"
	"github.com/gophercloud/utils/terraform/mutexkv"
)

type Provider struct {
	userAgent string
}

func (p *Provider) SetUserAgent(s string) {
	p.userAgent = s
}

func (p *Provider) Help() string {
	return `Openstack:

    provider:   "os"
    cloud:	An entry in a clouds.yaml file to use. 
    auth_url:   The endpoint of OS identity
    project_id: The id of the project (tenant id)
    tag_key:    The tag key to filter on
    tag_value:  The tag value to filter on
    user_name:  The user used to authenticate
    password:   The password of the provided user
    token:      The token to use
    insecure:   Sets if the api certificate shouldn't be check. Any value means true
    application_credential_id: ID of the application credentials to be used.
    application_credential_name: Name of the application credentials to be used
    application_credential_secret: Secret for the given application credential ID/Name

    Variables can also be provided by environmental variables.
`
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "os" {
		return nil, fmt.Errorf("discover-os: invalid provider " + args["provider"])
	}

	if l == nil {
		l = log.New(ioutil.Discard, "", 0)
	}

	client, config, err := newClient(args, l)
	if err != nil {
		l.Printf("[DEBUG] Client Creation Error - %s", err)
		return nil, err
	}

	projectID := config.TenantID
	tagKey := args["tag_key"]
	tagValue := args["tag_value"]

	if projectID == "" && config.Cloud == "" { // Use the one on the instance if not provided either by parameter or env
		l.Printf("[INFO] discover-os: ProjectID not provided. Looking up in metadata...")
		projectID, err = getProjectID()
		if err != nil {
			return nil, err
		}
		l.Printf("[INFO] discover-os: ProjectID is %s", projectID)
		args["project_id"] = projectID
	}

	log.Printf("[DEBUG] discover-os: Using project_id=%s tag_key=%s tag_value=%s", projectID, tagKey, tagValue)

	pager := servers.List(client, ListOpts{ListOpts: servers.ListOpts{Status: "ACTIVE"}, ProjectID: projectID})
	if err := pager.Err; err != nil {
		return nil, fmt.Errorf("discover-os: ListServers failed: %s", err)
	}

	var addrs []string
	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		srvs, err := servers.ExtractServers(page)
		if err != nil {
			return false, err
		}
		
		l.Printf("[INFO] discover-os: Filter instances with %s=%s", tagKey, tagValue)
		for _, srv := range srvs {
			for key, value := range srv.Metadata {
				if key == tagKey && value == tagValue {
					// Loop over the server address and append any fixed one to the list
					for _, v := range srv.Addresses {
						if addrsInfo, ok := v.([]interface{}); ok {
							for _, addrInfo := range addrsInfo {
								if info, ok := addrInfo.(map[string]interface{}); ok {
									if info["OS-EXT-IPS:type"] == "fixed" {
										addrs = append(addrs, info["addr"].(string))
									}
								}
							}
						}
					}
				}
			}
		}
		return true, nil
	})
	if err != nil {
		return nil, fmt.Errorf("discover-os: ExtractServerInfo failed: %s", err)
	}

	l.Printf("[DEBUG] discover-os: Found ip addresses: %v", addrs)
	return addrs, nil
}

func newClient(args map[string]string, l *log.Logger) (*gophercloud.ServiceClient, *auth.Config, error) {
	config := auth.Config{
		Cloud:            argsOrEnv(args, "cloud", "OS_CLOUD"),
		DomainID:         argsOrEnv(args, "domain_id", "OS_DOMAIN_ID"),
		DomainName:       argsOrEnv(args, "domain_name", "OS_DOMAIN_NAME"),
		IdentityEndpoint: argsOrEnv(args, "auth_url", "OS_AUTH_URL"),
		Password:         argsOrEnv(args, "password", "OS_PASSWORD"),

		Token:                       argsOrEnv(args, "token", "OS_AUTH_TOKEN"),
		TenantID:                    argsOrEnv(args, "project_id", "OS_PROJECT_ID"),
		TenantName:                  argsOrEnv(args, "project_name", "OS_PROJECT_NAME"),
		Username:                    argsOrEnv(args, "user_name", "OS_USERNAME"),
		ApplicationCredentialID:     argsOrEnv(args, "application_credential_id", "OS_APPLICATION_CREDENTIAL_ID"),
		ApplicationCredentialName:   argsOrEnv(args, "application_credential_name", "OS_APPLICATION_CREDENTIAL_NAME"),
		ApplicationCredentialSecret: argsOrEnv(args, "application_credential_secret", "OS_APPLICATION_CREDENTIAL_SECRET"),
		TerraformVersion:            "0.11",
		SDKVersion:                  "1",
		MutexKV:                     *(mutexkv.NewMutexKV()),
	}

	region := argsOrEnv(args, "region", "OS_REGION_NAME")
	if region == "" {
		region = "RegionOne"
	}

	if argsOrEnv(args, "insecure", "OS_INSECURE") == "true" {
		insecure := true
		config.Insecure = &insecure
	}

	if err := config.LoadAndValidate(); err != nil {
		return nil, nil, err
	}

	client, _ := config.ComputeV2Client(config.Region)
	return client, &config, nil
}

func argsOrEnv(args map[string]string, key, env string) string {
	if value := args[key]; value != "" {
		return value
	}
	return os.Getenv(env)
}

// ListOpts add the project to the parameters of servers.ListOpts
type ListOpts struct {
	servers.ListOpts
	ProjectID string `q:"project_id"`
}

// ToServerListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToServerListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

func getProjectID() (string, error) {
	resp, err := http.Get("http://169.254.169.254/openstack/latest/meta_data.json")
	if err != nil {
		return "", fmt.Errorf("discover-os: Error asking metadata for project_id: %s", err)
	}
	data := struct {
		ProjectID string `json:"project_id"`
	}{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("discover-os: Can't read response body: %s", err)
	}
	resp.Body.Close()
	if err = json.Unmarshal(body, &data); err != nil {
		return "", fmt.Errorf("discover-os: Can't convert project_id: %s", err)
	}
	if data.ProjectID == "" {
		return "", fmt.Errorf("discover-os: Couln't find project_id on metadata")
	}
	return data.ProjectID, nil
}
