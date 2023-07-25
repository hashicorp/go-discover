package hcloud_test

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/hcloud"
)

var _ discover.Provider = (*hcloud.Provider)(nil)
var addrTests = map[string]struct {
	addrType string
	location string
	outLen   int
}{
	"public ipv4 all locations":    {"public_v4", "", 2},
	"public ipv6 all locations":    {"public_v6", "", 2},
	"private ipv4 all locations":   {"private_v4", "", 2},
	"private ipv4 fsn1 datacenter": {"private_v4", "fsn1", 1},
	"public ipv6 nbg1 datacenter":  {"public_v6", "nbg1", 1},
	"public ipv4 hel1 datacenter":  {"public_v4", "hel1", 0},
}

func TestAddrs(t *testing.T) {
	l := log.New(os.Stderr, "", log.LstdFlags)
	for name, at := range addrTests {
		t.Run(name, func(t *testing.T) {
			args := discover.Config{
				"provider":       "hcloud",
				"label_selector": "go-discover-test-tag",
				"address_type":   at.addrType,
			}
			p := &hcloud.Provider{}
			addrs, err := p.Addrs(args, l)

			if err != nil {
				t.Fatal(err)
			}

			if len(addrs) != at.outLen {
				t.Fatalf("expected: %d, got: %d", at.outLen, len(addrs))
			}
		})
	}
}
