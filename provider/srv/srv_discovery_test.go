// Copyright IBM Corp. 2017, 2026
// SPDX-License-Identifier: MPL-2.0

package srv_test

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/srv"
)

var _ discover.Provider = (*srv.Provider)(nil)
var _ discover.ProviderWithUserAgent = (*srv.Provider)(nil)

func TestSRVAddrs(t *testing.T) {
	args := discover.Config{
		"provider": "srv",
		"service":  "ldap",
		"domain":   "google.com",
	}

	l := log.New(os.Stderr, "", log.LstdFlags)
	p := &srv.Provider{}
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 1 {
		t.Fatalf("bad: %v", addrs)
	}
}
