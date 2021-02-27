package proxmox

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	discover "github.com/hashicorp/go-discover"
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

func makeRequest(args discover.Config, apiPath string) (*http.Response, error) {
	apiBase := "/api2/json"
	apiURL, err := url.Parse(args["api_host"] + apiBase + apiPath)
	if err != nil {
		return nil, err
	}

	// Allow skipping certificate since many Proxmox users use self-signed and untrusted certs
	var transport *http.Transport = &http.Transport{}
	if args["api_skip_tls_verify"] == "skip" {
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	httpClient := &http.Client{
		Transport: transport,
	}

	req, err := http.NewRequest("GET", apiURL.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("PVEAPIToken=%s=%s", args["api_token_id"], args["api_token_secret"]))

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

type node struct {
	ID             string `json:"id"`
	Level          string `json:"level"`
	Node           string `json:"node"`
	SSLFingerprint string `json:"ssl_fingerprint"`
	Status         string `json:"status"`
	Type           string `json:"type"`
}

type nodesAPIResponse struct {
	Data []node `json:"data"`
}

func getNodes(args discover.Config) ([]node, error) {
	res, err := makeRequest(args, "/nodes")
	if err != nil {
		return nil, err
	}

	var nodes = new(nodesAPIResponse)
	jsonErr := json.NewDecoder(res.Body).Decode(&nodes)
	if jsonErr != nil {
		return nil, err
	}

	return nodes.Data, nil
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
