package dns_test

import (
	"fmt"
	"log"
	"net"
	"os"
	"testing"
	"time"

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

func pickUnusedUDPPort() (int, error) {
	addr, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenUDP("udp", addr)
	if err != nil {
		return 0, err
	}
	port := l.LocalAddr().(*net.UDPAddr).Port
	if err := l.Close(); err != nil {
		return 0, err
	}
	return port, nil
}

func newTestServer(port int) (*dns.Server, error) {
	var rr []dns.RR

	for _, testAddress := range testAddresses {
		a, _ := dns.NewRR(fmt.Sprintf("tasks.%s. IN A %s", testService, testAddress))
		rr = append(rr, a)
	}

	server := &dns.Server{Addr: fmt.Sprintf(":%d", port), Net: "udp"}
	go server.ListenAndServe()
	dns.HandleFunc(fmt.Sprintf("tasks.%s.", testService), func(w dns.ResponseWriter, r *dns.Msg) {
		m := new(dns.Msg)
		m.SetReply(r)
		m.Answer = rr
		w.WriteMsg(m)
	})

	// XXX: This test harness seems flakey without a small timeout
	time.Sleep(time.Second * 1)

	return server, nil
}

func TestDiscover(t *testing.T) {
	port, err := pickUnusedUDPPort()
	if err != nil {
		t.Fatalf("unable to find free udp port for test dns server harness: %s", err)
		return
	}
	testPort := fmt.Sprintf("%d", port)

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
				"server":   "127.0.0.1",
				"port":     testPort,
				"query":    "tasks.fake-service.",
				"timeout":  "1s",
			},
			false,
			0,
		},
		{
			"valid config - one address",
			discover.Config{
				"provider": "dns",
				"server":   "127.0.0.1",
				"port":     testPort,
				"query":    "tasks.test-service.",
				"timeout":  "10s",
			},
			false,
			1,
		},
		{
			"invalid config - missing query option",
			discover.Config{
				"provider": "dns",
				"server":   "127.0.0.1",
				"port":     testPort,
				"query":    "",
				"timeout":  "1s",
			},
			true,
			0,
		},
		{
			"invalid config - bad timeout option",
			discover.Config{
				"provider": "dns",
				"server":   "127.0.0.1",
				"port":     testPort,
				"query":    "tasks.fake-service.",
				"timeout":  "1z",
			},
			true,
			0,
		},
		{
			"invalid config - bad v6 option",
			discover.Config{
				"provider": "dns",
				"server":   "127.0.0.1",
				"port":     testPort,
				"query":    "tasks.fake-service.",
				"timeout":  "1s",
			},
			true,
			0,
		},
		{
			"invalid config - bad v4 option",
			discover.Config{
				"provider": "dns",
				"server":   "127.0.0.1",
				"port":     testPort,
				"query":    "tasks.fake-service.",
				"timeout":  "1s",
			},
			true,
			0,
		},
	}

	p := &provider.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)

	svr, err := newTestServer(port)
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
