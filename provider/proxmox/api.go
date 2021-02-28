package proxmox

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-discover"
)

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

type member struct {
	ID     string `json:"id"`
	Node   string `json:"node"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Type   string `json:"type"`
	VMID   int    `json:"vmid"`
}

type poolData struct {
	Members []member `json:"members"`
}

type nodesAPIResponse struct {
	Data poolData `json:"data"`
}

func getPoolMembers(args discover.Config) ([]member, error) {
	res, err := makeRequest(args, "/pools/"+args["pool_name"])
	if err != nil {
		return nil, err
	}

	var nodes = new(nodesAPIResponse)
	jsonErr := json.NewDecoder(res.Body).Decode(&nodes)
	if jsonErr != nil {
		return nil, err
	}

	return nodes.Data.Members, nil
}

type ipAddresses struct {
	IPAddress     string `json:"ip-address"`
	IPAddressType string `json:"ip-address-type"`
	Prefix        int    `json:"prefix"`
}

type statistics struct {
	RxBytes   int `json:"rx-bytes"`
	RxDropped int `json:"rx-dropped"`
	RxErrs    int `json:"rx-errs"`
	RxPackets int `json:"rx-packets"`
	TxBytes   int `json:"tx-bytes"`
	TxDropped int `json:"tx-dropped"`
	TxErrs    int `json:"tx-errs"`
	TxPackets int `json:"tx-packets"`
}

type networkInterface struct {
	HardwareAddress string        `json:"hardware-address"`
	IPAddresses     []ipAddresses `json:"ip-addresses"`
	Name            string        `json:"name"`
	Statistics      statistics    `json:"statistics"`
}

type data struct {
	Result []networkInterface `json:"result"`
}

type getNetworkInterfacesResponse struct {
	Data data `json:"data"`
}

func getNetworkInterfaces(args discover.Config, node string, vmID string) ([]networkInterface, error) {
	res, err := makeRequest(args, "/nodes/"+node+"/qemu/"+vmID+"/agent/network-get-interfaces")
	if err != nil {
		return nil, err
	}

	var interfaces = new(getNetworkInterfacesResponse)
	jsonErr := json.NewDecoder(res.Body).Decode(&interfaces)
	if jsonErr != nil {
		return nil, err
	}

	return interfaces.Data.Result, nil
}
