package gce_test

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
)

func TestAddrs(t *testing.T) {
	cfg := discover.Config{
		"provider":         "gce",
		"project_name":     os.Getenv("GCE_PROJECT"),
		"zone_pattern":     os.Getenv("GCE_ZONE"),
		"tag_value":        "consul-server",
		"credentials_file": os.Getenv("GCE_CONFIG_CREDENTIALS"),
	}
	if cfg["project_name"] == "" || cfg["credentials_file"] == "" {
		t.Skip("GCE credentials missing")
	}

	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := discover.Addrs(cfg.String(), l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 3 {
		t.Fatalf("bad: %v", addrs)
	}
}
