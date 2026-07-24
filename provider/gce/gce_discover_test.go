// Copyright IBM Corp. 2017, 2026
// SPDX-License-Identifier: MPL-2.0

package gce_test

import (
	"log"
	"os"
	"strings"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/gce"
)

var _ discover.Provider = (*gce.Provider)(nil)
var _ discover.ProviderWithUserAgent = (*gce.Provider)(nil)

func testConfig(t *testing.T) discover.Config {
	t.Helper()

	projectName := os.Getenv("GOOGLE_PROJECT")
	fileContents := os.Getenv("GOOGLE_CREDENTIALS")
	if projectName == "" || fileContents == "" {
		t.Skip("Google credentials missing")
	}

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
	t.Cleanup(func() {
		_ = os.Remove(tmpCreds.Name())
	})

	return discover.Config{
		"provider":         "gce",
		"project_name":     projectName,
		"zone_pattern":     os.Getenv("GOOGLE_ZONE"),
		"credentials_file": tmpCreds.Name(),
	}
}

func TestAddrs(t *testing.T) {
	args := testConfig(t)
	args["tag_value"] = "consul-server"

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
	args := testConfig(t)
	args["label_key"] = "environment"
	args["label_value"] = "test"

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

func TestAddrsWithTagsAndLabels(t *testing.T) {
	args := testConfig(t)
	args["tag_value"] = "consul-server"
	args["label_key"] = "environment"
	args["label_value"] = "test"

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

func TestAddrsWithIncompleteLabelFilter(t *testing.T) {
	tests := []struct {
		name string
		args discover.Config
	}{
		{
			name: "missing label value",
			args: discover.Config{
				"provider":  "gce",
				"label_key": "environment",
			},
		},
		{
			name: "missing label key",
			args: discover.Config{
				"provider":    "gce",
				"label_value": "test",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := &gce.Provider{}
			_, err := p.Addrs(test.args, nil)
			if err == nil {
				t.Fatal("expected error")
			}
			if !strings.Contains(err.Error(), "label_key and label_value must both be set or both be empty") {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
