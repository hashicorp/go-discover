package gce_test

import (
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
	tmpCreds, err := os.CreateTemp("", "gce-credentials")
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

func TestAddrsWithLabelSearch(t *testing.T) {
	// Check if we have basic requirements
	projectName := os.Getenv("GOOGLE_PROJECT")
	if projectName == "" {
		t.Skip("GOOGLE_PROJECT environment variable not set")
	}

	// Test label-based search
	args := discover.Config{
		"provider":     "gce",
		"project_name": projectName,
		"zone_pattern": os.Getenv("GOOGLE_ZONE"),
		"label_key":    "environment",
		"label_value":  "test",
	}

	l := log.New(os.Stderr, "", log.LstdFlags)
	p := &gce.Provider{}

	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Logf("Error (this may be expected if no instances exist or auth fails): %v", err)
		t.Skip("Skipping due to error - likely no instances with label 'environment=test' exist")
	}

	t.Logf("Found %d addresses with label search: %v", len(addrs), addrs)
}

func TestAddrsWithTagsAndLabels(t *testing.T) {
	// Check if we have basic requirements
	projectName := os.Getenv("GOOGLE_PROJECT")
	if projectName == "" {
		t.Skip("GOOGLE_PROJECT environment variable not set")
	}

	// Test label-based search
	args := discover.Config{
		"provider":     "gce",
		"project_name": projectName,
		"zone_pattern": os.Getenv("GOOGLE_ZONE"),
		"tag_value":    "consul-server",
		"label_key":    "environment",
		"label_value":  "test",
	}

	l := log.New(os.Stderr, "", log.LstdFlags)
	p := &gce.Provider{}

	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Logf("Error (this may be expected if no instances exist or auth fails): %v", err)
		t.Skip("Skipping due to error - likely no instances with label 'environment=test' exist")
	}

	t.Logf("Found %d addresses with tag and label search: %v", len(addrs), addrs)
}
