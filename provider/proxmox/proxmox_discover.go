package proxmox

import (
	"fmt"
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

	if args["api_skip_tls_verify"] != "skip" && args["api_skip_tls_verify"] != "verify" {
		args["api_skip_tls_verify"] = "verify"
	}

	if args["addr_type"] != "v4" && args["addr_type"] != "v6" {
		args["addr_type"] = "v4"
	}

	// Get all the members of the pool
	members, err := getPoolMembers(args)
	if err != nil {
		return nil, fmt.Errorf("discover-proxmox: could not list pool members: %s", err)
	}

	// Get the network interfaces from just the members that at QEMU vm's
	var interfaces []networkInterface
	for _, member := range members {
		if member.Type == "qemu" {
			memberInterfaces, err := getNetworkInterfaces(args, member.Node, fmt.Sprint(member.VMID))
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

	// Collect the correct (ipv4 or ipv6) IP addresses from the interface
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

	return addresses, nil
}
