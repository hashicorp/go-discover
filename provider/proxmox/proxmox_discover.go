package proxmox

import (
	"fmt"
	"log"
	"net/http"
)

// Provider for Proxmox
type Provider struct{}

// Help message generator
func (p *Provider) Help() string {
	return `Proxmox:

		provider:            "proxmox"
		api_host:            The address of the Proxmox node
		api_token_id:        The ID of the API token
		api_token_secret:    The secret of the API token
		api_skip_tls_verify: "skip" or "verify". Defaults to "verify"
		addr_type:           "v4", "v6" or "both". Defaults to "v4".
		filter_name_prefix   Filter VMs by name prefix
`
}

// Addrs function to retrieve IP addresses from Proxmox
func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "proxmox" {
		return nil, fmt.Errorf("discover-proxmox: invalid provider " + args["provider"])
	}

	if args["api_skip_tls_verify"] != "skip" && args["api_skip_tls_verify"] != "verify" {
		args["api_skip_tls_verify"] = "verify"
	}

	_, err := http.Get("")
	if err != nil {
		return nil, fmt.Errorf("discover-proxmox: %s", err)
	}

	addrs := []string{"haha", "123"}

	return addrs, nil
}
