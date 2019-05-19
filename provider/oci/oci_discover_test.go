package oci_test

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/oci"
)

func TestAddrs(t *testing.T) {
	args := discover.Config{
		"provider"     : "oci",
		"tag_namespace": "discovery",
		"tag_key"      : "consul",
	}

	p := &oci.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 1 {
		t.Fatalf("bad: %v", addrs)
	}
}