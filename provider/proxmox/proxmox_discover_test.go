package proxmox

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"testing"

	discover "github.com/hashicorp/go-discover"
)

var _ discover.Provider = (*Provider)(nil)

var (
	ipv4Regex, _ = regexp.Compile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`)
	ipv6Regex, _ = regexp.Compile(`^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`)
)

func TestAddrs(t *testing.T) {
	args := discover.Config{
		"provider":            "proxmox",
		"api_host":            os.Getenv("PROXMOX_API_HOST"),
		"api_token_id":        os.Getenv("PROXMOX_API_ID"),
		"api_token_secret":    os.Getenv("PROXMOX_API_SECRET"),
		"api_skip_tls_verify": "skip",
		"pool_name":           os.Getenv("PROXMOX_POOL_NAME"),
	}

	if args["api_host"] == "" || args["api_token_id"] == "" || args["api_token_secret"] == "" || args["pool_name"] == "" {
		t.Skip("Proxmox credentials missing")
	}

	logger := log.New(os.Stderr, "", log.LstdFlags)
	p := &Provider{}
	addrs, err := p.Addrs(args, logger)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Addresses recieved:", addrs)

	// Assume that at least one address was returned
	if len(addrs) <= 0 {
		t.Fatal("Zero addresses returned")
	}

	// Ensure the addresses returned were ipv4 addresses
	for _, address := range addrs {
		if !ipv4Regex.MatchString(address) {
			t.Fatalf("IP address is not a valid ipv4 address: %s", address)
		}
	}
}

func TestAddrsIPv6(t *testing.T) {
	args := discover.Config{
		"provider":            "proxmox",
		"api_host":            os.Getenv("PROXMOX_API_HOST"),
		"api_token_id":        os.Getenv("PROXMOX_API_ID"),
		"api_token_secret":    os.Getenv("PROXMOX_API_SECRET"),
		"api_skip_tls_verify": "skip",
		"pool_name":           os.Getenv("PROXMOX_POOL_NAME"),
		"addr_type":           "v6",
	}

	if args["api_host"] == "" || args["api_token_id"] == "" || args["api_token_secret"] == "" || args["pool_name"] == "" {
		t.Skip("Proxmox credentials missing")
	}

	logger := log.New(os.Stderr, "", log.LstdFlags)
	p := &Provider{}
	addrs, err := p.Addrs(args, logger)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Addresses recieved:", addrs)

	// Assume that at least one address was returned
	if len(addrs) <= 0 {
		t.Fatal("Zero addresses returned")
	}

	// Ensure the addresses returned were ipv4 addresses
	for _, address := range addrs {
		if !ipv6Regex.MatchString(address) {
			t.Fatalf("IP address is not a valid ipv6 address: %s", address)
		}
	}
}
