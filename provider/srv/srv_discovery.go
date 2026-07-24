// Copyright IBM Corp. 2017, 2026
// SPDX-License-Identifier: MPL-2.0

// Package srv provides node discovery for SRV dns entries.
package srv

import (
	"fmt"
	"io"
	"log"
	"net"
)

type Provider struct {
	userAgent string
}

func (p *Provider) SetUserAgent(s string) {
	p.userAgent = s
}

func (p *Provider) Help() string {
	return `SRV:

    provider:   "srv"
    service:    The SRV service to filter on
    proto:      The protocol to filter on
    domain:     The SRV domain to filter on
`
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "srv" {
		return nil, fmt.Errorf("discover-srv: invalid provider %s", args["provider"])
	}

	if l == nil {
		l = log.New(io.Discard, "", 0)
	}
	proto := args["proto"]
	if proto == "" {
		proto = "tcp"
	}
	domain := args["domain"]
	service := args["service"]
	if domain == "" || service == "" {
		return nil, fmt.Errorf("discover-srv: service or domain is required")
	}
	l.Printf("[INFO] srv: Using service=%s proto=%s domain=%s", service, proto, domain)

	_, records, err := net.LookupSRV(service, proto, domain)
	if err != nil {
		return nil, err
	}

	var addrs []string
	for _, r := range records {
		l.Printf("[INFO] discover-srv: %s:%d", r.Target, r.Port)
		addrs = append(addrs, fmt.Sprintf("%s:%d", r.Target, r.Port))
	}
	return addrs, nil
}
