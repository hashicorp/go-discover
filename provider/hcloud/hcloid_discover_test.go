package hcloud_test

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/hcloud"
)

var _ discover.Provider = (*hcloud.Provider)(nil)

func TestAddrsPublicV4(t *testing.T) {
	args := discover.Config{
		"provider":  "hcloud",
		"label_selector":  "go-discover-test-tag",
		"api_token": os.Getenv("HCLOUD_TOKEN"),
		"addr_type": "public_v4",
	}

	if args["api_token"] == "" {
		t.Skip("Hetzner Cloud credentials missing")
	}

	l := log.New(os.Stderr, "", log.LstdFlags)
	p := &hcloud.Provider{}
	addrs, err := p.Addrs(args, l)

	if err != nil {
		t.Fatal(err)
	}

	if len(addrs) != 2 {
		t.Fatalf("bad: %v", addrs)
	}
}

func TestAddrsPublicV6(t *testing.T) {
	args := discover.Config{
		"provider":  "hcloud",
		"label_selector":  "go-discover-test-tag",
		"api_token": os.Getenv("HCLOUD_TOKEN"),
		"addr_type": "public_v6",
	}

	if args["api_token"] == "" {
		t.Skip("Hetzner Cloud credentials missing")
	}

	l := log.New(os.Stderr, "", log.LstdFlags)
	p := &hcloud.Provider{}
	addrs, err := p.Addrs(args, l)

	if err != nil {
		t.Fatal(err)
	}

	if len(addrs) != 2 {
		t.Fatalf("bad: %v", addrs)
	}
}

func TestAddrsPrivateV4(t *testing.T) {
	args := discover.Config{
		"provider":  "hcloud",
		"label_selector":  "go-discover-test-tag",
		"api_token": os.Getenv("HCLOUD_TOKEN"),
		"addr_type": "private_v4",
	}

	if args["api_token"] == "" {
		t.Skip("Hetzner Cloud credentials missing")
	}

	l := log.New(os.Stderr, "", log.LstdFlags)
	p := &hcloud.Provider{}
	addrs, err := p.Addrs(args, l)

	if err != nil {
		t.Fatal(err)
	}

	if len(addrs) != 2 {
		t.Fatalf("bad: %v", addrs)
	}
}
