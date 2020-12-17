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
		"provider":  "scaleway",
		"addr_type": "private_v4",
	}

	p := &scaleway.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) == 0 {
		t.Fatalf("bad: %v", addrs)
	}
}
