package dns_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/miekg/dns"

	discover "github.com/hashicorp/go-discover"
	provider "github.com/hashicorp/go-discover/provider/dns"
)

const (
	testService = "test-service"
)

var (
	testAddresses = []string{"127.0.0.1"}
)

func newTestServer() (*dns.Server, error) {
	var rr []dns.RR

	for _, testAddress := range testAddresses {
		a, _ := dns.NewRR(fmt.Sprintf("tasks.%s. IN A %s", testService, testAddress))
		rr = append(rr, a)
	}

	server := &dns.Server{Addr: ":5300", Net: "udp"}
	go server.ListenAndServe()
	dns.HandleFunc(fmt.Sprintf("tasks.%s.", testService), func(w dns.ResponseWriter, r *dns.Msg) {
		m := new(dns.Msg)
		m.SetReply(r)
		m.Insert(rr)
		w.WriteMsg(m)
	})

	return server, nil
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
				"provider": "dns",
				"query":    "tasks.fake-service",
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
				"provider": "dns",
				"query":    "tasks.test-service",
				"timeout":  "10s",
				"v6":       "true",
				"v4":       "true",
			},
			false,
			1,
		},
		{
			"invalid config - missing query option",
			discover.Config{
				"provider": "dns",
				"query":    "",
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
				"provider": "dns",
				"query":    "tasks.fake-service",
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
				"provider": "dns",
				"service":  "tasks.fake-service",
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
				"provider": "dns",
				"service":  "tasks.fake-service",
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
