// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mdns_test

import (
	"log"
	"net"
	"os"
	"testing"

	"github.com/hashicorp/mdns"

	discover "github.com/hashicorp/go-discover"
	provider "github.com/hashicorp/go-discover/provider/mdns"
)

func newTestServer() (*mdns.Server, error) {
	zone, err := mdns.NewMDNSService(
		"localhost",
		"_test-service._noop",
		"local.",
		"",
		1234,
		[]net.IP{net.IPv4(127, 0, 0, 1)},
		[]string{"testing123"},
	)
	if err != nil {
		return nil, err
	}
	return mdns.NewServer(&mdns.Config{Zone: zone})
}

func TestDiscover(t *testing.T) {
	type testCases []struct {
		desc  string
		args  discover.Config
		fail  bool
		addrs int
	}

	cases := testCases{
		{
			"valid config - no addresses",
			discover.Config{
				"provider": "mdns",
				"service":  "_fake-service._noop",
				"domain":   "local",
				"timeout":  "1s",
				"v6":       "false",
				"v4":       "false",
			},
			false,
			0,
		},
		{
			"valid config - one address",
			discover.Config{
				"provider": "mdns",
				"service":  "_test-service._noop",
				"domain":   "local",
				"timeout":  "10s",
				"v6":       "true",
				"v4":       "true",
			},
			false,
			1,
		},
		{
			"invalid config - missing service option",
			discover.Config{
				"provider": "mdns",
				"service":  "",
				"domain":   "test",
				"timeout":  "1s",
				"v6":       "false",
				"v4":       "false",
			},
			true,
			0,
		},
		{
			"invalid config - bad timeout option",
			discover.Config{
				"provider": "mdns",
				"service":  "_fake-service._noop",
				"domain":   "test",
				"timeout":  "1z",
				"v6":       "false",
				"v4":       "false",
			},
			true,
			0,
		},
		{
			"invalid config - bad v6 option",
			discover.Config{
				"provider": "mdns",
				"service":  "_fake-service._noop",
				"domain":   "test",
				"timeout":  "1s",
				"v6":       "xxxx",
				"v4":       "false",
			},
			true,
			0,
		},
		{
			"invalid config - bad v4 option",
			discover.Config{
				"provider": "mdns",
				"service":  "_fake-service._noop",
				"domain":   "test",
				"timeout":  "1s",
				"v6":       "false",
				"v4":       "xxxx",
			},
			true,
			0,
		},
	}

	p := &provider.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)

	svr, err := newTestServer()
	if err != nil {
		t.Fatalf("Failed to start test server: %s", err.Error())
		return
	}
	defer svr.Shutdown()

	for idx, tc := range cases {
		addrs, err := p.Addrs(tc.args, l)

		if !tc.fail && err != nil {
			t.Fatalf("FAIL [%d/%d] %s -> %s",
				idx, len(cases), tc.desc, err)
		}

		if len(addrs) != tc.addrs {
			t.Fatalf("FAIL [%d/%d] %s -> wrong addr count: expected %d / got %d",
				idx, len(cases), tc.desc, tc.addrs, len(addrs))
		}

		t.Logf("PASS [%d/%d] %s", idx, len(cases), tc.desc)
	}
}
