package virtualbox

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Adapter int
const (
	Nat Adapter = iota
	Host
	Bridge
	Internal
	Generic
)

func (a Adapter) String() string {
	return []string{
		"natnet",
		"hostonlyadapter",
		"bridgeadapter",
		"intnet",
		"generic",
	}[a]
}

type Nic struct {
	adapterType Adapter
	adapterName string
	mac string
	pos int
	netmask string
	broadcast string

}

type Network struct {
	ip string
	nic Nic
}

type VmIdentifier struct {
	name string
	uuid string
}

type Vm struct {
	identifier VmIdentifier
	networks []Network
}

type Provider struct {
	
}

func (p *Provider) Help() string {
	return `Virtualbox:
    provider:         "virtualbox"
    adapterType:      Filters ip to this adapter type the vm is attached to
    pos:        	  Filters ip to order of position any adapter the vm is attached to
		`
}

func (p *Provider) Addrs(args map[string]string, l *log.Logger) ([]string, error) {
	provider := args["provider"]
	if provider != "virtualbox"{
		return nil, fmt.Errorf("discover-virtualbox: invalid provider " + args["provider"])
	}

	if l == nil {
		l = log.New(ioutil.Discard, "", 0)
	}

	if _, err := exec.LookPath("vboxmanage"); err != nil {
		return nil, fmt.Errorf("discover-virtualbox: %s", err)
	}

	groupName := args["group"]

	adapterType := args["adapterType"]
	pos := args["pos"]

	vms, err := getVmList()
	if err != nil {
		return nil, fmt.Errorf("discover-virtualbox: %s", err)
	}

	var groupVms []Vm
	for _, vm := range vms {
		vmInfoOutputs, err := getCommandOutput("vboxmanage", "showvminfo", "--machinereadable", vm.uuid)
		if err != nil {
			return nil, fmt.Errorf("discover-virtualbox: %s", err)
		}

		vmInfoOutputsRows := strings.Split(vmInfoOutputs, "\n")
		vmInfo := make(map[string]string, len(vmInfoOutputsRows))
		for _, row := range vmInfoOutputsRows {
			split := strings.Split(row, "=")
			field := trimSuffixes(trimPrefixes(split[0], []string{"\""}), []string{"\""})
			value := trimSuffixes(trimPrefixes(split[1], []string{"\""}), []string{"\""})

			vmInfo[field] = value
		}

		if vmInfo["groups"] == groupName {
			var nics []Nic
			for key, value := range vmInfo {
				subMatches := regexp.MustCompile("^.+" + "([0-9]+)").FindStringSubmatch(key)
				if len(subMatches) == 2 {
					field := subMatches[0][0:len(subMatches[0]) - len(subMatches[1])]
					adapter := whichAdapter(field)
					if adapter != -1 || field == "natnet" {
						num := subMatches[1]
						pos, err := strconv.Atoi(num)
						if err != nil {
							return nil, fmt.Errorf("discover-virtualbox: %s", err)
						}

						nics = append(nics, Nic{
							adapterType: adapter,
							adapterName: value,
							mac:         vmInfo["macaddress" + num],
							pos:         pos,
						})
					}
				}
			}

			networks, err := findNetworks(vm, nics)
			if err != nil {
				return nil, fmt.Errorf("discover-virtualbox: %s", err)
			}

			groupVms = append(groupVms, Vm{
				identifier: vm,
				networks:   networks,
			})
		}
	}

	var ips []string
	for _, vm := range groupVms {
		for _, network := range vm.networks {
			if adapterType != "" && network.nic.adapterType.String() != adapterType {
				continue
			}

			if pos != "" {
				posNum, err := strconv.Atoi(pos)
				if err != nil {
					continue
				}
				if network.nic.pos != posNum {
					continue
				}
			}

			ips = append(ips, network.ip)
		}
	}

	return ips, nil
}

func whichAdapter(s string) Adapter {
	switch s {
	case Nat.String():
		return Nat
	case Host.String():
		return Nat
	case Bridge.String():
		return Nat
	case Internal.String():
		return Nat
	default:
		return -1
	}
}

func findNetworks(vm VmIdentifier, nics []Nic) ([]Network, error) {
	var networks []Network
	for _, nic := range nics {
		vmNetMacOutput, err := getCommandOutput(
			"vboxmanage",
			"guestproperty",
			"get",
			vm.uuid,
			"/VirtualBox/GuestInfo/Net/" + strconv.Itoa(nic.pos - 1) + "/MAC")
		if err != nil {
			return nil, fmt.Errorf("%s", err)
		}

		if vmNetMacOutput == "No value set!" {
			continue
		}

		if strings.Split(vmNetMacOutput, ": ")[1] != nic.mac {
			panic("MAC address of nic in position " + strconv.Itoa(nic.pos) + " - 1, did not match MAC in guestproperty query")
		}

		vmNetIPOutput, err := getCommandOutput(
			"vboxmanage",
			"guestproperty",
			"get",
			vm.uuid,
			"/VirtualBox/GuestInfo/Net/" + strconv.Itoa(nic.pos - 1) + "/V4/IP")
		if err != nil {
			return nil, fmt.Errorf("%s", err)
		}

		networks = append(networks, Network{
			ip:  strings.Split(vmNetIPOutput, ": ")[1],
			nic: nic,
		})
	}

	return networks, nil
}

func getVmList() ([]VmIdentifier, error) {
	vmListOutputs, err := getCommandOutput("vboxmanage", "list", "vms")
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	var vmList []VmIdentifier
	for _, vmListOutput := range strings.Split(vmListOutputs, "\n") {
		split := strings.Split(vmListOutput, " ")
		vmList = append(vmList, VmIdentifier{
			name: trimPrefixes(trimSuffixes(split[0], []string{"}", "\"", "\n"}), []string{"{", "\"", "\n"}),
			uuid: trimPrefixes(trimSuffixes(split[1], []string{"}", "\"", "\n"}), []string{"{", "\"", "\n"}),
		})
	}

	return vmList, nil
}

func getCommandOutput(cmdName string, args ...string) (string, error) {
	cmd := exec.Command(cmdName, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	if len(stderr.Bytes()) > 0 {
		return "", fmt.Errorf("%s", stderr.Bytes())
	}

	return trimPrefixes(trimSuffixes(stdout.String(), []string{"\n"}), []string{"\n"}), nil
}

func trimSuffixes(s string, rs []string) string {
	for _, r := range rs {
		s = strings.TrimSuffix(s, r)
	}

	return s
}

func trimPrefixes(s string, rs []string) string {
	for _, r := range rs {
		s = strings.TrimPrefix(s, r)
	}

	return s
}