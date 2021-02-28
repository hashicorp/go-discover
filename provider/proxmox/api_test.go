package proxmox

import (
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
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

	members, err := getPoolMembers(args)
	if err != nil {
		t.Fatalf("bad: %v", err)
	}

	// Assume the pool has at least one member
	if len(members) <= 0 {
		t.Fatal("Zero members retrieved")
	}
}
