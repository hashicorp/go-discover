package scaleway_test

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/scaleway"
)

var _ discover.Provider = (*scaleway.Provider)(nil)

func TestAddrs(t *testing.T) {
	args := discover.Config{
		"provider":     "scaleway",
		"organization": os.Getenv("SCALEWAY_ORGANIZATION"),
		"token":        os.Getenv("SCALEWAY_TOKEN"),
		"tag_name":     "consul-server",
		"region":       os.Getenv("SCALEWAY_REGION"),
	}
	if args["organization"] == "" {
		t.Skip("Scaleway organization missing")
	}

	if args["token"] == "" {
		t.Skip("Scaleway token missing")
	}

	p := &scaleway.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 2 {
		t.Fatalf("bad: %v", addrs)
	}
}
