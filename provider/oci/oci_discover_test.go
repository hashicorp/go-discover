package oci_test

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/oci"
)

func TestAddrs(t *testing.T) {
	freeform := discover.Config{
		"provider" : "oci",
		"tag_key"  : "discover",
		"tag_value": "me",
	}

	defined := discover.Config{
		"provider"     : "oci",
		"tag_namespace": "defined",
		"tag_key"      : "discover",
		"tag_value"    : "me",
	}

	freePartial := discover.Config{
		"provider": "oci",
		"tag_key" : "discover",
	}

	definedPartial := discover.Config{
		"provider"     : "oci",
		"tag_namespace": "defined",
		"tag_key"      : "discover",
	}

	p := &oci.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)
	
	// Testing freeform tags
	addrs, err := p.Addrs(freeform, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 1 {
		t.Fatalf("bad: %v", addrs)
	}

	// Testing defined tags
	addrs, err = p.Addrs(defined, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 2 {
		t.Fatalf("bad: %v", addrs)
	}

	// Testing freeform partial tags
	addrs, err = p.Addrs(freePartial, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 1 {
		t.Fatalf("bad: %v", addrs)
	}

	// Testing defined partial tags
	addrs, err = p.Addrs(definedPartial, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 2 {
		t.Fatalf("bad: %v", addrs)
	}
}