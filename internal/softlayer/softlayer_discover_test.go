package softlayer

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/go-discover/config"
)

func TestDiscover(t *testing.T) {
	t.Parallel()
	if os.Getenv("SL_USERNAME") == "" {
		t.Skip("SL_USERNAME not set, skipping")
	}

	if os.Getenv("SL_API_KEY") == "" {
		t.Skip("SL_API_KEY not set, skipping")
	}

	cfg := fmt.Sprintf("username=%s api_key=%s datacenter=%s tag_value=%s",
		os.Getenv("SL_USERNAME"), os.Getenv("SL_API_KEY"), "dal06", "consul-server")

	m, err := config.Parse(cfg)
	if err != nil {
		t.Fatal(err)
	}

	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := Discover(m, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 3 {
		t.Fatalf("bad: %v", addrs)
	}
}
