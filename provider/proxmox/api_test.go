package proxmox_test

import (
	"fmt"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/proxmox"
)

func TestGetPoolMembers(t *testing.T) {
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

	members, err := proxmox.GetPoolMembers(args)
	if err != nil {
		t.Fatalf("bad: %v", err)
	}

	// Assume the pool has at least one member
	if len(members) <= 0 {
		t.Fatal("Zero members retrieved")
	}
}

func TestGetNetworkInterfaces(t *testing.T) {
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

	members, err := proxmox.GetPoolMembers(args)
	if err != nil {
		t.Fatalf("bad: %v", err)
	}

	for _, member := range members {
		interfaces, err := proxmox.GetNetworkInterfaces(args, member.Node, fmt.Sprint(member.VMID))
		if err != nil {
			t.Fatalf("bad: %v", err)
		}

		// Assume each member has at least one network interface
		if len(interfaces) <= 0 {
			t.Fatal("Zero network interfaces retrieved")
		}
	}
}
