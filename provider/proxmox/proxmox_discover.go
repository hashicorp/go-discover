package proxmox

import (
	"fmt"
	"io/ioutil"
	"log"
)

// Provider for Proxmox
type Provider struct{}

// Help message generator
func (p *Provider) Help() string {
	return `Proxmox:

		provider:            "proxmox"
		api_host:            The address of the Proxmox node
		api_token_id:        The ID of the API token
		api_token_secret:    The secret of the API token
		api_skip_tls_verify: "skip" or "verify". Defaults to "verify"
		addr_type:           "v4" or "v6". Defaults to "v4".
		pool_name            Pool to get VMs from
`
}

// Addrs function to retrieve IP addresses from Proxmox
func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	if args["provider"] != "proxmox" {
		return nil, fmt.Errorf("discover-proxmox: invalid provider " + args["provider"])
	}

	if l == nil {
		l = log.New(ioutil.Discard, "", 0)
	}

	if args["api_skip_tls_verify"] != "skip" && args["api_skip_tls_verify"] != "verify" {
		l.Printf("[INFO] discover-proxmox: api_skip_tls_verify %s is not supported. Valid values are 'skip' or 'verify'. Falling back to 'verify'", args["api_skip_tls_verify"])
		args["api_skip_tls_verify"] = "verify"
	}

	if args["addr_type"] != "v4" && args["addr_type"] != "v6" {
		l.Printf("[INFO] discover-proxmox: addr_type %s is not supported. Valid values are 'v4' or 'v6'. Falling back to 'v4'", args["addr_type"])
		args["addr_type"] = "v4"
	}

	// Get all the members of the pool
	l.Printf("[DEBUG] discover-proxmox: retrieveing members of pool: %s", args["pool_name"])
	members, err := GetPoolMembers(args)
	l.Printf("[DEBUG] discover-proxmox: got %d members", len(members))
	if err != nil {
		return nil, fmt.Errorf("discover-proxmox: could not list pool members: %s", err)
	}

	// Get the network interfaces from just the members that at QEMU vm's
	l.Print("[DEBUG] discover-proxmox: retrieveing network interfaces from members")
	var interfaces []NetworkInterface
	for _, member := range members {
		if member.Type == "qemu" {
			memberInterfaces, err := GetNetworkInterfaces(args, member.Node, fmt.Sprint(member.VMID))
			if err != nil {
				return nil, fmt.Errorf(
					"discover-proxmox: could not get interfaces from pool member '%s' (ID: %d): %s",
					member.Name,
					member.VMID,
					err,
				)
			}

			// Add the first non loopback interface to the output list
			for _, memberInterface := range memberInterfaces {
				// Ignore loopback interfaces
				if memberInterface.HardwareAddress == "00:00:00:00:00:00" {
					continue
				}

				interfaces = append(interfaces, memberInterface)
				break
			}
		}
	}
	l.Printf("[DEBUG] discover-proxmox: got %d network interfaces", len(interfaces))

	// Collect the correct (ipv4 or ipv6) IP addresses from the interface
	l.Printf("[DEBUG] discover-proxmox: filtering ip addresses by type (%s)", args["addr_type"])
	var addresses []string
	for _, netInterface := range interfaces {
		for _, ipAddress := range netInterface.IPAddresses {
			if args["addr_type"] == "v4" && ipAddress.IPAddressType == "ipv4" {
				addresses = append(addresses, ipAddress.IPAddress)
				break
			}

			if args["addr_type"] == "v6" && ipAddress.IPAddressType == "ipv6" {
				addresses = append(addresses, ipAddress.IPAddress)
				break
			}
		}
	}
	l.Printf("[DEBUG] discover-proxmox: got %d ip addresses", len(addresses))

	return addresses, nil
}
