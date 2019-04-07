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

func TestAddrsTaggedDefault(t *testing.T) {
	args := discover.Config{
		"provider":  "linode",
		"api_token": os.Getenv("LINODE_TOKEN"),
		"tag_name":  "gd-tag1",
	}

	if args["api_token"] == "" {
		t.Skip("Linode credentials missing")
	}

	p := &linode.Provider{}

	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)

	if err != nil {
		t.Fatal(err)
	}

	if len(addrs) != 2 {
		t.Fatalf("bad: %v", addrs)
	}
}

func TestAddrsTaggedPublicV6(t *testing.T) {
	args := discover.Config{
		"provider":     "linode",
		"api_token":    os.Getenv("LINODE_TOKEN"),
		"address_type": "public_v6",
		"tag_name":     "gd-tag1",
	}

	if args["api_token"] == "" {
		t.Skip("Linode credentials missing")
	}

	p := &linode.Provider{}

	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)

	if err != nil {
		t.Fatal(err)
	}

	if len(addrs) != 2 {
		t.Fatalf("bad: %v", addrs)
	}
}

func TestAddrsTaggedPublicV4(t *testing.T) {
	args := discover.Config{
		"provider":     "linode",
		"api_token":    os.Getenv("LINODE_TOKEN"),
		"address_type": "public_v4",
		"tag_name":     "gd-tag1",
	}

	if args["api_token"] == "" {
		t.Skip("Linode credentials missing")
	}

	p := &linode.Provider{}

	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)

	if err != nil {
		t.Fatal(err)
	}

	if len(addrs) != 2 {
		t.Fatalf("bad: %v", addrs)
	}
}

func TestAddrsTaggedRegion(t *testing.T) {
	args := discover.Config{
		"provider":  "linode",
		"api_token": os.Getenv("LINODE_TOKEN"),
		"tag_name":  "gd-tag1",
		"region":    "us-east",
	}

	if args["api_token"] == "" {
		t.Skip("Linode credentials missing")
	}

	p := &linode.Provider{}

	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)

	if err != nil {
		t.Fatal(err)
	}

	if len(addrs) != 1 {
		t.Fatalf("bad: %v", addrs)
	}
}
