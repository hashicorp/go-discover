package tailscale

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"

	"tailscale.com/client/tailscale"
)

type Provider struct {
	client tailscale.LocalClient
}

func (p *Provider) Help() string {
	return `Tailscale:

    provider:        "tailscale"
    socket:          The optional alternative socket path to tailscaled.
    use_socket_only: If true, tries to only connect to tailscaled via
                     the Unix socket and not via fallback mechanisms.
    tag_regexp:      The regular expression to use to match host tags.
    include_ipv4:    If true, include ipv4 addresses.
    include_ipv6:    If true, include ipv6 addresses.

    For Tailscale discovery, a local Tailscale daemon must be available.
`
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "tailscale" {
		return nil, fmt.Errorf("discover-tailscale: invalid provider " + args["provider"])
	}

	if l == nil {
		l = log.New(ioutil.Discard, "", 0)
	}

	if socket := args["socket"]; socket != "" {
		p.client.Socket = socket
	}

	useSocketOnly, err := boolArg("use_socket_only", args, false)
	if err != nil {
		return nil, err
	}
	if useSocketOnly {
		p.client.UseSocketOnly = useSocketOnly
	}

	includeIp4, err := boolArg("include_ipv4", args, true)
	if err != nil {
		return nil, err
	}
	includeIp6, err := boolArg("include_ipv6", args, false)
	if err != nil {
		return nil, err
	}

	reTag := regexp.MustCompile(args["tag_regexp"])

	st, err := p.client.Status(context.Background())
	if err != nil {
		return nil, err
	}

	var addrs []string
	for _, peer := range st.Peer {
		tags := peer.Tags
		if tags != nil {
			for t := 0; t < tags.Len(); t++ {
				tag := tags.At(t)
				if reTag.MatchString(tag) {
					for i := range peer.TailscaleIPs {
						addr := peer.TailscaleIPs[i]
						if (includeIp4 && addr.Is4()) || (includeIp6 && addr.Is6()) {
							addrs = append(addrs, addr.String())
						}
					}
				}
			}
		}
	}
	return addrs, nil
}

func boolArg(key string, args map[string]string, defaultValue bool) (bool, error) {
	ret := defaultValue
	if v := args[key]; v != "" {
		var err error
		ret, err = strconv.ParseBool(v)
		if err != nil {
			return false, fmt.Errorf("discover-tailscale: %s must be boolean value: %s", key, err)
		}
	}
	return ret, nil
}
