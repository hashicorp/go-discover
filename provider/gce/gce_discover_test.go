package gce_test

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/gce"
)

func TestAddrs(t *testing.T) {
	args := discover.Config{
		"provider":         "gce",
		"project_name":     os.Getenv("GCE_PROJECT"),
		"zone_pattern":     os.Getenv("GCE_ZONE"),
		"tag_value":        "consul-server",
		"credentials_file": os.Getenv("GCE_CONFIG_CREDENTIALS"),
	}
	if args["project_name"] == "" || args["credentials_file"] == "" {
		t.Skip("GCE credentials missing")
	}

	l := log.New(os.Stderr, "", log.LstdFlags)
	p := &gce.Provider{}
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 3 {
		t.Fatalf("bad: %v", addrs)
	}
}
