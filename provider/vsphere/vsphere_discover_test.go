// Copyright IBM Corp. 2017, 2026
// SPDX-License-Identifier: MPL-2.0

package vsphere_test

import (
	"context"
	"log"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/simulator"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vapi/tags"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/types"

	// Importing vapi/simulator registers the REST/tags endpoints on the
	// in-process simulator via its init() function.
	_ "github.com/vmware/govmomi/vapi/simulator"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/vsphere"
)

var _ discover.Provider = (*vsphere.Provider)(nil)

func testPreCheck(t *testing.T) {
	if v := os.Getenv("VSPHERE_USER"); v == "" {
		t.Skip("VSPHERE_USER must be set for acceptance tests")
	}

	if v := os.Getenv("VSPHERE_PASSWORD"); v == "" {
		t.Skip("VSPHERE_PASSWORD must be set for acceptance tests")
	}

	if v := os.Getenv("VSPHERE_SERVER"); v == "" {
		t.Skip("VSPHERE_SERVER must be set for acceptance tests")
	}
}

// vsphereExpectedIPs returns the list of IP addresses that acceptance tests
// should assert are present in the discovery output. The list is read from the
// VSPHERE_EXPECTED_IPS environment variable as a comma-separated string
// (e.g. "10.0.0.10,10.0.0.11"). If the variable is unset or empty, nil is
// returned and the caller should fall back to asserting at least one IP.
func vsphereExpectedIPs() []string {
	v := os.Getenv("VSPHERE_EXPECTED_IPS")
	if v == "" {
		return nil
	}
	return strings.Split(v, ",")
}

// assertDiscoveredIPs checks that all expected IPs are present in addrs.
// If no expected IPs are configured, it asserts that at least one address
// was returned.
func assertDiscoveredIPs(t *testing.T, addrs []string) {
	t.Helper()
	expected := vsphereExpectedIPs()
	if len(expected) == 0 {
		if len(addrs) == 0 {
			t.Fatal("expected at least one address in discovery output, got none")
		}
		return
	}
	set := make(map[string]bool, len(addrs))
	for _, a := range addrs {
		set[a] = true
	}
	for _, ip := range expected {
		if !set[ip] {
			t.Fatalf("IP address %s is missing from discovery output", ip)
		}
	}
}

func TestAddrs(t *testing.T) {
	testPreCheck(t)

	args := discover.Config{
		"provider":      "vsphere",
		"tag_name":      "go-discover-test-tag",
		"category_name": "go-discover-test-category",
		"host":          os.Getenv("VSPHERE_SERVER"),
		"user":          os.Getenv("VSPHERE_USER"),
		"password":      os.Getenv("VSPHERE_PASSWORD"),
		"insecure_ssl":  os.Getenv("VSPHERE_ALLOW_UNVERIFIED_SSL"),
		"timeout":       "20m",
	}

	l := log.New(os.Stderr, "", log.LstdFlags)
	p := &vsphere.Provider{}
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}

	assertDiscoveredIPs(t, addrs)
}

// TestAddrsEnv tests to make sure that we can lean on the environment for
// credentials automatically. User credential environment variables are not set
// using Setenv, leaving them to be fetched from the environment 100%.
func TestAddrsEnv(t *testing.T) {
	testPreCheck(t)

	args := discover.Config{
		"provider":      "vsphere",
		"tag_name":      "go-discover-test-tag",
		"category_name": "go-discover-test-category",
		"timeout":       "20m",
	}

	l := log.New(os.Stderr, "", log.LstdFlags)
	p := &vsphere.Provider{}
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}

	assertDiscoveredIPs(t, addrs)
}

// TestAddrsSimulator exercises the full Addrs() discovery path against an
// in-process govmomi vSphere simulator. It runs unconditionally — no live
// vSphere environment or environment variables required.
func TestAddrsSimulator(t *testing.T) {
	const (
		categoryName = "go-discover-test-category"
		tagName      = "go-discover-test-tag"
		testIP       = "10.0.0.100"
	)

	model := simulator.VPX()
	model.Machine = 1 // ensure at least one VM with a NIC is created

	// simulator.Test uses a func(ctx, c) callback so t.Fatal/t.Skip work
	// directly — no sentinel error or outer error check needed.
	simulator.Test(func(ctx context.Context, c *vim25.Client) {
		// Set up REST tags client.
		rc := rest.NewClient(c)
		if err := rc.Login(ctx, simulator.DefaultLogin); err != nil {
			t.Fatal(err)
		}
		defer func() { _ = rc.Logout(ctx) }()

		mgr := tags.NewManager(rc)

		// Create a tag category and tag.
		catID, err := mgr.CreateCategory(ctx, &tags.Category{
			Name:        categoryName,
			Cardinality: "SINGLE",
		})
		if err != nil {
			t.Fatal(err)
		}
		tagID, err := mgr.CreateTag(ctx, &tags.Tag{Name: tagName, CategoryID: catID})
		if err != nil {
			t.Fatal(err)
		}

		// Find the first VM in the simulator inventory.
		finder := find.NewFinder(c, false)
		dc, err := finder.DefaultDatacenter(ctx)
		if err != nil {
			t.Fatal(err)
		}
		finder.SetDatacenter(dc)
		vms, err := finder.VirtualMachineList(ctx, "*")
		if err != nil {
			t.Fatal(err)
		}
		if len(vms) == 0 {
			t.Skip("simulator produced no virtual machines")
		}
		vm := vms[0]

		// Inject a guest IP directly into the simulator's in-memory registry.
		// This avoids the power-off → CustomizeVM → power-on task chain that
		// would otherwise be required to populate guest.net[].IpConfig.
		simVM := simulator.Map(ctx).Get(vm.Reference()).(*simulator.VirtualMachine)
		simVM.Guest.Net = []types.GuestNicInfo{{
			IpConfig: &types.NetIpConfigInfo{
				IpAddress: []types.NetIpConfigInfoIpAddress{{IpAddress: testIP}},
			},
		}}

		// Attach the tag to the VM and run discovery.
		if err := mgr.AttachTag(ctx, tagID, vm.Reference()); err != nil {
			t.Fatal(err)
		}

		pass, _ := simulator.DefaultLogin.Password()
		addrs, err := (&vsphere.Provider{}).Addrs(discover.Config{
			"provider":      "vsphere",
			"tag_name":      tagName,
			"category_name": categoryName,
			"host":          c.URL().Host,
			"user":          simulator.DefaultLogin.Username(),
			"password":      pass,
			"insecure_ssl":  "true",
			"timeout":       "2m",
		}, log.New(os.Stderr, "", log.LstdFlags))
		if err != nil {
			t.Fatal(err)
		}

		if !slices.Contains(addrs, testIP) {
			t.Errorf("expected IP %s in addrs %v", testIP, addrs)
		}
	}, model)
}
