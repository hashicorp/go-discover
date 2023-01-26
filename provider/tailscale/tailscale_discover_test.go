package tailscale_test

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/tailscale"
)

var _ discover.Provider = (*tailscale.Provider)(nil)

func TestAcc(t *testing.T) {
	args := discover.Config{
		"provider":   "tailscale",
		"tag_regexp": "^tag:server$",
	}
	l := log.New(os.Stderr, "", log.LstdFlags)
	p := &tailscale.Provider{}
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Addrs: %v", addrs)
	if len(addrs) != 5 {
		t.Fatalf("bad: %v", addrs)
	}
}
