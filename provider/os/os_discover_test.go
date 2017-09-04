package os_test

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	openstack "github.com/hashicorp/go-discover/provider/os"
)

func TestAddrs(t *testing.T) {
	args := discover.Config{
		"provider":  "os",
		"tag_key":   "consul",
		"tag_value": "server",
	}

	l := log.New(os.Stderr, "", log.LstdFlags)
	p := &openstack.Provider{}
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 3 {
		t.Fatalf("bad: %v", addrs)
	}
}
