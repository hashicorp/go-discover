package exoscale_test

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/exoscale"
)

var _ discover.Provider = (*exoscale.Provider)(nil)

func TestAddrs(t *testing.T) {
	args := discover.Config{
		"provider":   "exoscale",
		"api_key":    os.Getenv("EXOSCALE_API_KEY"),
		"api_secret": os.Getenv("EXOSCALE_API_SECRET"),
	}

	p := &exoscale.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 2 {
		t.Fatalf("bad: %v", addrs)
	}
}

func TestAddrsZone(t *testing.T) {
	args := discover.Config{
		"provider":   "exoscale",
		"api_key":    os.Getenv("EXOSCALE_API_KEY"),
		"api_secret": os.Getenv("EXOSCALE_API_SECRET"),
		"zone":       "ch-gva-2",
	}

	p := &exoscale.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 1 {
		t.Fatalf("bad: %v", addrs)
	}
}

func TestAddrsTags(t *testing.T) {
	args := discover.Config{
		"provider":   "exoscale",
		"api_key":    os.Getenv("EXOSCALE_API_KEY"),
		"api_secret": os.Getenv("EXOSCALE_API_SECRET"),
		"tag_key":    "test-01",
		"tag_value":  "consul",
	}

	p := &exoscale.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 1 {
		t.Fatalf("bad: %v", addrs)
	}
}

func TestAddrsIPv6(t *testing.T) {
	args := discover.Config{
		"provider":   "exoscale",
		"api_key":    os.Getenv("EXOSCALE_API_KEY"),
		"api_secret": os.Getenv("EXOSCALE_API_SECRET"),
		"addr_type":  "public_v6",
	}

	p := &exoscale.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 1 {
		t.Fatalf("bad: %v", addrs)
	}
}
