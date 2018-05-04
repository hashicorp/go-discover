package gce_test

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/gce"
)

var _ discover.Provider = (*gce.Provider)(nil)
var _ discover.ProviderWithUserAgent = (*gce.Provider)(nil)

func TestAddrs(t *testing.T) {
	// assume the google credentials file contents are in the environment,
	// as with the terraform provider
	fileContents := os.Getenv("GOOGLE_CREDENTIALS")
	tmpCreds, err := ioutil.TempFile("", "gce-credentials")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tmpCreds.WriteString(fileContents); err != nil {
		t.Fatal(err)
	}
	if err := tmpCreds.Close(); err != nil {
		t.Fatal(err)
	}
	// remove credentials file
	defer os.Remove(tmpCreds.Name())

	args := discover.Config{
		"provider":         "gce",
		"project_name":     os.Getenv("GOOGLE_PROJECT"),
		"zone_pattern":     os.Getenv("GOOGLE_ZONE"),
		"tag_value":        "consul-server",
		"credentials_file": tmpCreds.Name(),
	}

	if args["project_name"] == "" || args["credentials_file"] == "" {
		t.Skip("Google credentials missing")
	}

	l := log.New(os.Stderr, "", log.LstdFlags)
	p := &gce.Provider{}
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 2 {
		t.Fatalf("bad: %v", addrs)
	}
}
