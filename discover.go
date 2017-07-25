// Package discover provides functions to get metadata for different
// cloud environments.
package discover

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/hashicorp/go-discover/internal/aws"
	"github.com/hashicorp/go-discover/internal/azure"
	"github.com/hashicorp/go-discover/internal/gce"
	"github.com/hashicorp/go-discover/internal/softlayer"
)

// HelpDiscoverAddrs describes the format of the configuration
// string for address discovery and the various provider specific
// options.
var HelpDiscoverAddrs = `
  The options for discovering ip addresses are provided as a
  single string value in "key=value key=value ..." format where
  the values are URL encoded.

    provider=aws region=eu-west-1 ...

  The options are provider specific and are listed below.

  Amazon AWS:

    provider:          "aws"
    region:            The AWS region. Default to region of instance.
    tag_key:           The tag key to filter on
    tag_value:         The tag value to filter on
    access_key_id:     The AWS access key to use
    secret_access_key: The AWS secret access key to use

  Microsoft Azure:

   provider:          "azure"
   tenant_id:         The id of the tenant
   client_id:         The id of the client
   subscription_id:   The id of the subscription
   secret_access_key: The authentication credential
   tag_name:          The name of the tag to filter on
   tag_value:         The value of the tag to filter on

  Google Cloud:

    provider:         "gce"
    project_name:     The name of the project. discovered if not set
    zone_pattern:     A RE2 regular expression for filtering zones, e.g. us-west1-.*, or us-(?west|east).*
    tag_value:        The tag value for filtering instances
    credentials_file: The path to the credentials file. See below for more details

    Authentication is handled in the following order:

     1. Use credentials from "credentials_file", if provided.
     2. Use JSON file from GOOGLE_APPLICATION_CREDENTIALS environment variable.
     3. Use JSON file in a location known to the gcloud command-line tool.
        On Windows, this is %APPDATA%/gcloud/application_default_credentials.json.
        On other systems, $HOME/.config/gcloud/application_default_credentials.json.
     4. On Google Compute Engine, use credentials from the metadata
        server. In this final case any provided scopes are ignored.

  Softlayer:

    provider:   "softlayer"
    datacenter: The SoftLayer datacenter to filter on
    tag_value:  The tag value to filter on
    username:   The SoftLayer username to use
    api_key:    The SoftLayer api key to use
`

// Addrs discovers ip addresses of nodes that match the given filter
// criteria. It is a convenience function for &Discoverer{l}.Addrs(cfg).
func Addrs(cfg string, l *log.Logger) ([]string, error) {
	return (&Discoverer{l}).Addrs(cfg)
}

// Discoverer provides functions for getting metadata in cloud
// environments.
type Discoverer struct {
	Log *log.Logger
}

// Addrs discovers ip addresses of nodes that match the given filter
// criteria. The configuration is provider specific and is described in
// HelpDiscoverAddrs.
func (d *Discoverer) Addrs(cfg string) ([]string, error) {
	m, err := Parse(cfg)
	if err != nil {
		return nil, fmt.Errorf("discover: %s", err)
	}
	l := d.Log
	if l == nil {
		l = log.New(ioutil.Discard, "", 0)
	}
	p := m["provider"]
	switch p {
	case "aws":
		return aws.Discover(m, l)
	case "gce":
		return gce.Discover(m, l)
	case "azure":
		return azure.Discover(m, l)
	case "softlayer":
		return softlayer.Discover(m, l)
	default:
		return nil, fmt.Errorf("discover: unknown provider %q", p)
	}
}
