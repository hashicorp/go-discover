package linode_test

import (
	"context"
	"log"
	"net/http"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/linode"
	"github.com/linode/linodego"
	"golang.org/x/oauth2"
)

var _ discover.Provider = (*linode.Provider)(nil)
var _ discover.ProviderWithUserAgent = (*linode.Provider)(nil)

func TestAddrsTaggedDefault(t *testing.T) {
	opts := linodego.InstanceCreateOptions{}
	opts.Label = "go-discover-test-default"
	opts.Tags = []string{"gd-tag1"}
	opts.Region = "us-southeast"
	opts.PrivateIP = true

	_, destroy, err := buildInstance(t, opts)
	if err != nil {
		t.Fatal(err)
	}
	defer destroy()
	args := discover.Config{
		"provider":  "linode",
		"api_token": os.Getenv("LINODE_TOKEN"),
		"tag_name":  "gd-tag1",
	}

	if args["api_token"] == "" {
		t.Skip("Linode credentials missing")
	}

	p := &linode.Provider{}

	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)

	if err != nil {
		t.Fatal(err)
	}

	if len(addrs) != 1 {
		t.Fatalf("bad: %v", addrs)
	}
}

func TestAddrsTaggedPublicV6(t *testing.T) {
	opts := linodego.InstanceCreateOptions{}
	opts.Label = "go-discover-test-public-v6"
	opts.Tags = []string{"gd-tag1"}
	opts.Region = "us-southeast"

	_, destroy, err := buildInstance(t, opts)
	if err != nil {
		t.Fatal(err)
	}
	defer destroy()

	args := discover.Config{
		"provider":     "linode",
		"api_token":    os.Getenv("LINODE_TOKEN"),
		"address_type": "public_v6",
		"tag_name":     "gd-tag1",
	}

	if args["api_token"] == "" {
		t.Skip("Linode credentials missing")
	}

	p := &linode.Provider{}

	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)

	if err != nil {
		t.Fatal(err)
	}

	if len(addrs) != 1 {
		t.Fatalf("bad: %v", addrs)
	}
}

func TestAddrsTaggedPublicV4(t *testing.T) {
	opts := linodego.InstanceCreateOptions{}
	opts.Label = "go-discover-test-public-v4"
	opts.Tags = []string{"gd-tag1"}
	opts.Region = "us-southeast"

	_, destroy, err := buildInstance(t, opts)
	if err != nil {
		t.Fatal(err)
	}
	defer destroy()

	args := discover.Config{
		"provider":     "linode",
		"api_token":    os.Getenv("LINODE_TOKEN"),
		"address_type": "public_v4",
		"tag_name":     "gd-tag1",
	}

	if args["api_token"] == "" {
		t.Skip("Linode credentials missing")
	}

	p := &linode.Provider{}

	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)

	if err != nil {
		t.Fatal(err)
	}

	if len(addrs) != 1 {
		t.Fatalf("bad: %v", addrs)
	}
}

func TestAddrsTaggedRegion(t *testing.T) {
	opts := linodego.InstanceCreateOptions{}
	opts.Label = "go-discover-test-tagged-region"
	opts.Tags = []string{"gd-tag1"}
	opts.Region = "us-southeast"
	opts.PrivateIP = true

	_, destroy, err := buildInstance(t, opts)
	if err != nil {
		t.Fatal(err)
	}
	defer destroy()

	args := discover.Config{
		"region":    "us-southeast",
		"provider":  "linode",
		"tag_name":  "gd-tag1",
		"api_token": os.Getenv("LINODE_TOKEN"),
	}

	if args["api_token"] == "" {
		t.Skip("Linode credentials missing")
	}

	p := &linode.Provider{}

	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)

	if err != nil {
		t.Fatal(err)
	}

	if len(addrs) != 1 {
		t.Fatalf("bad: %v", addrs)
	}
}

func TestAddrsTaggedVLAN(t *testing.T) {
	opts := linodego.InstanceCreateOptions{}
	opts.Interfaces = []linodego.InstanceConfigInterface{
		{
			Label:       "go-discover-test-vlan",
			Purpose:     linodego.InterfacePurposeVLAN,
			IPAMAddress: "10.0.0.1/24",
		},
	}
	opts.Label = "go-discover-test-vlan"
	opts.Tags = []string{"gd-tag1"}
	opts.Region = "us-southeast"
	opts.Image = "linode/alpine3.16"
	opts.RootPass = "supercoolpasspleasedontsteal"

	_, destroy, err := buildInstance(t, opts)
	if err != nil {
		t.Fatal(err)
	}
	defer destroy()

	args := discover.Config{
		"region":     "us-southeast",
		"provider":   "linode",
		"tag_name":   "gd-tag1",
		"vlan_label": "go-discover-test-vlan",
		"api_token":  os.Getenv("LINODE_TOKEN"),
	}

	if args["api_token"] == "" {
		t.Skip("Linode credentials missing")
	}

	p := &linode.Provider{}

	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)

	if err != nil {
		t.Fatal(err)
	}

	if len(addrs) != 1 {
		t.Fatalf("bad: %v", addrs)
	}
}

func buildInstance(t *testing.T, opts linodego.InstanceCreateOptions) (*linodego.Instance, func(), error) {
	t.Helper()
	client := getLinodeClient(t)

	falseBoot := true
	opts.Type = "g6-nanode-1"
	opts.Booted = &falseBoot
	instance, err := client.CreateInstance(context.Background(), opts)
	if err != nil {
		return nil, nil, err
	}

	teardown := func() {
		if terr := client.DeleteInstance(context.Background(), instance.ID); terr != nil {
			t.Errorf("Error deleting test Instance: %s", terr)
		}
	}

	return instance, teardown, nil
}

func getLinodeClient(t *testing.T) *linodego.Client {
	t.Helper()
	apiToken := os.Getenv("LINODE_API_TOKEN")
	if apiToken == "" {
		t.Fatal("failed to get $LINODE_API_TOKEN")
	}
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: apiToken})
	oauth2Client := &http.Client{
		Transport: &oauth2.Transport{
			Source: tokenSource,
		},
	}
	client := linodego.NewClient(oauth2Client)
	return &client
}
