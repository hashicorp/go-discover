// Package nomad provides allocation discovery for Nomad.
package nomad

import (
	"fmt"
	"log"

	"encoding/base64"

	"github.com/hashicorp/nomad/api"
)

type Provider struct{
	// client *api.Client
}

func (p *Provider) Help() string {
	return `Nomad (Nomad):

    provider:         "nomad"
		address:    			Nomad address
		secret_id:    	  Nomad secret_id
    service_name:     Nomad service to discover allocations for
    namespace:        Namespace to search for allocations (optional)
		region:     			Nomad region to discover allocations for (optional)
`
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "nomad" {
		return nil, fmt.Errorf("discover-nomad: invalid provider " + args["provider"])
	}

	serviceName := args["service_name"]
	if serviceName == "" {
		return nil, fmt.Errorf("discover-nomad: must provider a service_name")
	}

	client, err := createNomadClient(args)
	if err != nil {
		return nil, fmt.Errorf("discover-nomad: error creating client: %v", err)
	}

	// TODO: Do I need NS and Region here if provided above?
	queryOpts := &api.QueryOptions{}
	services, _, err := client.Services().Get(serviceName, queryOpts)
	if err != nil {
		return nil, fmt.Errorf("discover-nomad: error retrieving services: %v", err)
	}

	addrs := []string{}
	for _, service := range services {
		addr := service.Address
		port := service.Port
		addrs = append(addrs, fmt.Sprintf("%s:%d", addr, port))
	}

	return addrs, nil
}

func createNomadClient(args map[string]string) (*api.Client, error) {
	clientConfig := api.DefaultConfig()

	// === BASIC CONFIG ===

	if args["address"] == "" {
		return nil, fmt.Errorf("discover-nomad: must provider an address for Nomad")
	} else {
		clientConfig.Address = args["address"]
	}
	if args["namespace"] != "" {
		clientConfig.Namespace = args["namespace"]
	}
	if args["region"] != "" {
		clientConfig.Region = args["region"]
	}
	if args["secret_id"] != "" {
		clientConfig.SecretID = args["secret_id"]
	}

	// === TLS CONFIG ===
	// TODO: Is all of this needed if using the socket and the WI token?

	if args["tls_insecure"] == "true" || args["tls_insecure"] == "1" {
		clientConfig.TLSConfig.Insecure = true
	}

	if args["ca_cert"] != "" {
		clientConfig.TLSConfig.CACert = args["ca_cert"]
	}

	if args["ca_path"] != "" {
		clientConfig.TLSConfig.CAPath = args["ca_path"]
	}

	if args["client_cert"] != "" {
		clientConfig.TLSConfig.ClientCert = args["client_cert"]
	}

	if args["client_key"] != "" {
		clientConfig.TLSConfig.ClientKey = args["client_key"]
	}

	if args["tls_server_name"] != "" {
		clientConfig.TLSConfig.TLSServerName = args["tls_server_name"]
	}

	b64CACertPEM := args["ca_cert_pem"]
	if b64CACertPEM != "" {
		decoded, err := base64.StdEncoding.DecodeString(b64CACertPEM)
		if err != nil {
			return nil, fmt.Errorf("discover-nomad: error decoding ca_cert_pem: %v", err)
		}
		clientConfig.TLSConfig.CACertPEM = decoded
	}

	b64ClientCertPEM := args["client_cert_pem"]
	if b64ClientCertPEM != "" {
		decoded, err := base64.StdEncoding.DecodeString(b64ClientCertPEM)
		if err != nil {
			return nil, fmt.Errorf("discover-nomad: error decoding client_cert_pem: %v", err)
		}
		clientConfig.TLSConfig.ClientCertPEM = decoded
	}

	b64ClientKeyPEM := args["client_key_pem"]
	if b64ClientKeyPEM != "" {
		decoded, err := base64.StdEncoding.DecodeString(b64ClientKeyPEM)
		if err != nil {
			return nil, fmt.Errorf("discover-nomad: error decoding client_key_pem: %v", err)
		}
		clientConfig.TLSConfig.ClientKeyPEM = decoded
	}

	// === HTTP AUTH CONFIG ===
	// TODO: Is all of this needed if using the socket and the WI token?

	// handle clientConfig.HttpAuth.Username input from args
	if args["http_auth_username"] != "" {
		clientConfig.HttpAuth.Username = args["http_auth_username"]
	}

	if args["http_auth_pw"] != "" {
		clientConfig.HttpAuth.Password = args["http_auth_pw"]
	}

	// === CLIENT CREATION ===

	client, err := api.NewClient(clientConfig)
	if err != nil {
		return nil, err
	}

 return client, nil
}
