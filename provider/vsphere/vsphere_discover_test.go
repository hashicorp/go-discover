package vsphere_test

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/vsphere"
)

var _ discover.Provider = (*vsphere.Provider)(nil)

func testPreCheck(t *testing.T) {
	if v := os.Getenv("VSPHERE_USER"); v == "" {
		t.Skip("VSPHERE_USER must be set for acceptance tests")
	}

	if v := os.Getenv("VSPHERE_PASSWORD"); v == "" {
		t.Skip("VSPHERE_PASSWORD must be set for acceptance tests")
	}

	if v := os.Getenv("VSPHERE_SERVER"); v == "" {
		t.Skip("VSPHERE_SERVER must be set for acceptance tests")
	}
}

func TestAddrs(t *testing.T) {
	testPreCheck(t)

	args := discover.Config{
		"provider":      "vsphere",
		"tag_name":      "go-discover-test-tag",
		"category_name": "go-discover-test-category",
		"host":          os.Getenv("VSPHERE_SERVER"),
		"user":          os.Getenv("VSPHERE_USER"),
		"password":      os.Getenv("VSPHERE_PASSWORD"),
		"insecure_ssl":  os.Getenv("VSPHERE_ALLOW_UNVERIFIED_SSL"),
		"timeout":       "20m",
	}

	l := log.New(os.Stderr, "", log.LstdFlags)
	p := &vsphere.Provider{}
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}

	actual := map[string]bool{
		"10.0.0.10": false,
		"10.0.0.11": false,
	}

	for _, addr := range addrs {
		if addr == "10.0.0.10" || addr == "10.0.0.11" {
			actual[addr] = true
		}
	}

	for k, v := range actual {
		if !v {
			t.Fatalf("IP address %s is missing from discovery output", k)
		}
	}
}

// TestAddrsEnv tests to make sure that we can lean on the environment for
// credentials automatically. User credential environment variables are not set
// using Setenv, leaving them to be fetched from the environment 100%.
func TestAddrsEnv(t *testing.T) {
	testPreCheck(t)

	args := discover.Config{
		"provider":      "vsphere",
		"tag_name":      "go-discover-test-tag",
		"category_name": "go-discover-test-category",
		"timeout":       "20m",
	}

	l := log.New(os.Stderr, "", log.LstdFlags)
	p := &vsphere.Provider{}
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}

	actual := map[string]bool{
		"10.0.0.10": false,
		"10.0.0.11": false,
	}

	for _, addr := range addrs {
		if addr == "10.0.0.10" || addr == "10.0.0.11" {
			actual[addr] = true
		}
	}

	for k, v := range actual {
		if !v {
			t.Fatalf("IP address %s is missing from discovery output", k)
		}
	}
}
