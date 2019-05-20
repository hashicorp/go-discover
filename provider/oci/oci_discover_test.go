package oci_test

import (
	"log"
	"os"
	"testing"
	"fmt"
	"os/user"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/oci"
)

var tests = []struct {
	name      string
	config    discover.Config
	addrCount int
}{
	{
		"freeform",
		discover.Config{
			"provider" : "oci",
			"tag_key"  : "discover",
			"tag_value": "me",
		},
		1,
	},
	{
		"defined",
		discover.Config{
			"provider"     : "oci",
			"tag_namespace": "defined",
			"tag_key"      : "discover",
			"tag_value"    : "me",
		},
		2,
	},
	{
		"freePartial",
		discover.Config{
			"provider": "oci",
			"tag_key" : "discover",
		},
		1,
	},
	{
		"definedPartial",
		discover.Config{
			"provider"     : "oci",
			"tag_namespace": "defined",
			"tag_key"      : "discover",
		},
		2,
	},
	{
		"definedPublic",
		discover.Config{
			"provider"     : "oci",
			"tag_namespace": "defined",
			"tag_key"      : "discover",
			"tag_value"    : "me",
			"addr_type"    : "public",
		},
		1,
	},
}

func TestAddrs(t *testing.T) {
	usr, err := user.Current()
	if err != nil {
		t.Error(err)
	}
	
	if _, err := os.Stat(fmt.Sprintf("%s/.oci/config", usr.HomeDir)); os.IsNotExist(err) {
		t.Skip("OCI config file missing.")
	}
	
	p := &oci.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)
	for _, test := range tests {
		l.Printf("[INFO] Begin Test: %s", test.name)
		addrs, err := p.Addrs(test.config, l)
		if err != nil {
			t.Error(err)
		}
		if len(addrs) != test.addrCount {
			t.Errorf("bad: %v", addrs)
		}
	}
}