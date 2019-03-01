package linode_test

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/linode"
)

var _ discover.Provider = (*linode.Provider)(nil)
var _ discover.ProviderWithUserAgent = (*linode.Provider)(nil)

func TestAddrs(t *testing.T) {
	args := discover.Config{
		"provider":  "linode",
		"tag_name":  "go-discover-test-tag",
		"region":    "us-east",
		"api_token": os.Getenv("LINODE_TOKEN"),
	}
	if args["api_token"] == "" {
		t.Skip("Linode credentials missing")
	}

	l := log.New(os.Stderr, "", log.LstdFlags)
	p := &linode.Provider{}
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 2 {
		t.Fatalf("bad: %v", addrs)
	}
}
