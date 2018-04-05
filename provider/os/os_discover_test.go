package os_test

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	openstack "github.com/hashicorp/go-discover/provider/os"
)

var _ discover.Provider = (*openstack.Provider)(nil)
var _ discover.ProviderWithUserAgent = (*openstack.Provider)(nil)

func TestAddrs(t *testing.T) {
	// todo: maybe check for http://169.254.169.254/openstack/latest/meta_data.json first
	t.Skip("Skipping Openstack test in non-openstack env. Please enable manually")

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
