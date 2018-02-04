package mdns_test

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/mdns"
)

func TestAddrs(t *testing.T) {
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
				"domain":   "test",
				"timeout":  "1s",
				"v6":       "false",
				"v4":       "false",
			},
			false,
			0,
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

	p := &mdns.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)

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
