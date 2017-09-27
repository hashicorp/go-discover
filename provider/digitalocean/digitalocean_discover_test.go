package digitalocean_test

import (
	"log"
	"os"
	"testing"

	"github.com/hashicorp/go-discover/provider/digitalocean"
	discover "github.com/hashicorp/go-discover"
)

func TestAddrs(t *testing.T) {
	args := discover.Config{
		"provider":  "digitalocean",
		"api_key":   os.Getenv("DIGITALOCEAN_API_KEY"),
		"region":    "nyc3",
		"tag_value": "consul-server",
		"addr_type": "public_v4",
	}
	if args["api_key"] == "" {
		t.Skip("DigitalOcean credentials missing")
	}

	p := &digitalocean.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 3 {
		t.Fatalf("bad: %v", addrs)
	}
}
