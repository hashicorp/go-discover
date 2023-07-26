package nomad_test

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/nomad"
)

var _ discover.Provider = (*nomad.Provider)(nil)

// Acceptance test against a local dev cluster
func TestAcc(t *testing.T) {
	args := discover.Config{
		"address": "http://127.0.0.1:4646",
		"provider":       "nomad",
		// "namespace":     "default",
		"service_name":   "consul",
	}

	l := log.New(os.Stderr, "", log.LstdFlags)
	p := &nomad.Provider{}

	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Addrs: %v", addrs)
	if len(addrs) != 3 {
		t.Fatalf("Bad Response (wanted 3): %v", addrs)
	}
}

// Ideally test that it formats the return addrs correctly
// Test that it passes in region and namespace correctly
// Test passing in token properly