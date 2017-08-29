package digitalocean_test

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/digitalocean"
)

func TestAddrs(t *testing.T) {
	args := discover.Config{
		"provider":  "digitalocean",
		"tag_name":  "go-discover-test-tag",
		"region":    "nyc3",
		"api_token": os.Getenv("DIGITALOCEAN_TOKEN"),
	}
	if args["api_token"] == "" {
		t.Skip("DigitalOcean credentials missing")
	}

	l := log.New(os.Stderr, "", log.LstdFlags)
	p := &digitalocean.Provider{}
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 1 {
		t.Fatalf("bad: %v", addrs)
	}
}
