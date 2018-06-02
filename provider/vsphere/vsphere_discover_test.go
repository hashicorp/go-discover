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

func TestAddrsEnv(t *testing.T) {
	testPreCheck(t)

	args := discover.Config{
		"provider": "vsphere",
	}

	os.Setenv("VSPHERE_TAG_NAME", "go-discover-test-tag")
	os.Setenv("VSPHERE_CATEGORY_NAME", "go-discover-test-category")
	os.Setenv("VSPHERE_TIMEOUT", "20m")
	defer func() {
		_ = os.Unsetenv("VSPHERE_TAG_NAME")
		_ = os.Unsetenv("VSPHERE_CATEGORY_NAME")
		_ = os.Unsetenv("VSPHERE_TIMEOUT")
	}()

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
