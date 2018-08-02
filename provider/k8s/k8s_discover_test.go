package k8s_test

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/k8s"
)

var _ discover.Provider = (*k8s.Provider)(nil)

func TestAddrs(t *testing.T) {
	args := discover.Config{
		"provider":       "k8s",
		"label_selector": "app=consul-server",
	}

	l := log.New(os.Stderr, "", log.LstdFlags)
	p := &k8s.Provider{}
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Result: %v", addrs)

	// This is a weak assertion, but given the dynamic scheduling of
	// Consul in a K8S cluster, its hard to expect specific IP addresses.
	if len(addrs) != 3 {
		t.Fatalf("expected 3 results, got %v", addrs)
	}
}
