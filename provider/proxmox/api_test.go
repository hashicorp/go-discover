package proxmox

import (
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
)

func TestGetNodes(t *testing.T) {
	args := discover.Config{
		"provider":            "proxmox",
		"api_host":            os.Getenv("PROXMOX_API_HOST"),
		"api_token_id":        os.Getenv("PROXMOX_API_ID"),
		"api_token_secret":    os.Getenv("PROXMOX_API_SECRET"),
		"api_skip_tls_verify": "skip",
	}

	if args["api_host"] == "" || args["api_token_id"] == "" || args["api_token_secret"] == "" {
		t.Skip("Proxmox credentials missing")
	}

	nodes, err := getNodes(args)
	if err != nil {
		t.Fatalf("bad: %v", err)
	}

	if len(nodes) <= 0 {
		t.Fatal("Zero nodes retrieved")
	}
}
