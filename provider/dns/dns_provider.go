// Package dns provides node discovery via DNS.
package dns

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/miekg/dns"
)

// Provider implements the Provider interface.
type Provider struct{}

// Help returns help information for the DNS package.
func (p *Provider) Help() string {
	return `DNS:

    provider:          "dns"
    query:             The DNS query to lookup.  Required.
    server:            The DNS server to query.  Default "8.8.8.8" (Google DNS).
    port:              The DNS port to use.      Default: "53" (standard DNS port).
    timeout:           The mDNS lookup timeout.  Default "5s" (five seconds).
`
}

// Addrs returns discovered addresses for the mDNS package.
func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	var err error

	// default to null logger
	if l == nil {
		l = log.New(io.Discard, "", 0)
	}

	// validate and set service record
	if args["query"] == "" {
		return nil, fmt.Errorf("discover-dns: Query not provided")
	}

	var server string
	// validate and set server
	if args["server"] != "" {
		server = args["server"]
	} else {
		server = "8.8.8.8"
	}

	var port int
	// validate and set port
	if args["port"] != "" {
		if port, err = strconv.Atoi(args["port"]); err != nil {
			return nil, fmt.Errorf("discover-dns: Failed to parse port: %s", err)
		}
	} else {
		port = 53
	}

	var timeout time.Duration
	// validate and set timeout
	if args["timeout"] != "" {
		if timeout, err = time.ParseDuration(args["timeout"]); err != nil {
			return nil, fmt.Errorf("discover-dns: Failed to parse timeout: %s", err)
		}
	} else {
		timeout = 5 * time.Second
	}

	var addrs []string

	// lookup and return
	m1 := new(dns.Msg)
	m1.Id = dns.Id()
	m1.RecursionDesired = true
	m1.Question = make([]dns.Question, 1)
	m1.Question[0] = dns.Question{
		Name:   args["query"],
		Qtype:  dns.TypeA,
		Qclass: dns.ClassINET,
	}

	c := new(dns.Client)

	laddr := net.UDPAddr{
		IP:   net.ParseIP("[::1]"),
		Port: 0,
		Zone: "",
	}
	c.Dialer = &net.Dialer{
		Timeout:   timeout,
		LocalAddr: &laddr,
	}
	raddr := fmt.Sprintf("%s:%d", server, port)
	in, _, err := c.Exchange(m1, raddr)
	if err != nil {
		return nil, fmt.Errorf("discover-dns: Failed to process query: %s", err)
	}

	for _, answer := range in.Answer {
		if t, ok := answer.(*dns.A); ok {
			addrs = append(addrs, t.A.String())
		}
	}
	l.Printf("[DEBUG] discover-dns: found addrs %q", addrs)

	return addrs, err
}
