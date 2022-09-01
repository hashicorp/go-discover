package equinixmetal_test

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/equinixmetal"
)

var _ discover.Provider = (*equinixmetal.Provider)(nil)
var _ discover.ProviderWithUserAgent = (*equinixmetal.Provider)(nil)

func TestAddrsDefault(t *testing.T) {
	args := discover.Config{
		"provider":   "metal",
		"auth_token": os.Getenv("METAL_AUTH_TOKEN"),
		"project":    os.Getenv("METAL_PROJECT"),
	}

	if args["auth_token"] == "" {
		t.Skip("Equinix Metal credentials missing")
	}

	if args["project"] == "" {
		t.Skip("Equinix Metal project UUID missing")
	}

	p := equinixmetal.Provider{}

	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)

	if err != nil {
		t.Fatal(err)
	}

	if len(addrs) != 4 {
		t.Fatalf("bad: %v", addrs)
	}
}

func TestAddrsPublicV6(t *testing.T) {
	args := discover.Config{
		"provider":     "metal",
		"auth_token":   os.Getenv("METAL_AUTH_TOKEN"),
		"project":      os.Getenv("METAL_PROJECT"),
		"address_type": "public_v6",
	}

	if args["auth_token"] == "" {
		t.Skip("Equinix Metal credentials missing")
	}

	if args["project"] == "" {
		t.Skip("Equinix Metal project UUID missing")
	}

	p := equinixmetal.Provider{}

	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)

	if err != nil {
		t.Fatal(err)
	}

	if len(addrs) != 4 {
		t.Fatalf("bad: %v", addrs)
	}
}

func TestAddrsPublicV4(t *testing.T) {
	args := discover.Config{
		"provider":     "metal",
		"auth_token":   os.Getenv("METAL_AUTH_TOKEN"),
		"project":      os.Getenv("METAL_PROJECT"),
		"address_type": "public_v4",
	}

	if args["auth_token"] == "" {
		t.Skip("Equinix Metal credentials missing")
	}

	if args["project"] == "" {
		t.Skip("Equinix Metal project UUID missing")
	}

	p := equinixmetal.Provider{}

	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)

	if err != nil {
		t.Fatal(err)
	}

	if len(addrs) != 4 {
		t.Fatalf("bad: %v", addrs)
	}
}

func TestAddrsFacilityInclude(t *testing.T) {
	args := discover.Config{
		"provider":     "metal",
		"auth_token":   os.Getenv("METAL_AUTH_TOKEN"),
		"project":      os.Getenv("METAL_PROJECT"),
		"address_type": "private_v4",
		"facility":     "ewr1",
	}

	if args["auth_token"] == "" {
		t.Skip("Equinix Metal credentials missing")
	}

	if args["project"] == "" {
		t.Skip("Equinix Metal project UUID missing")
	}

	p := equinixmetal.Provider{}

	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)

	if err != nil {
		t.Fatal(err)
	}

	if len(addrs) != 1 {
		t.Fatalf("bad: %v", addrs)
	}
}

func TestAddrsFacilityIncludeMulti(t *testing.T) {
	args := discover.Config{
		"provider":     "metal",
		"auth_token":   os.Getenv("METAL_AUTH_TOKEN"),
		"project":      os.Getenv("METAL_PROJECT"),
		"address_type": "private_v4",
		"facility":     "ewr1,ams1",
	}

	if args["auth_token"] == "" {
		t.Skip("Equinix Metal credentials missing")
	}

	if args["project"] == "" {
		t.Skip("Equinix Metal project UUID missing")
	}

	p := equinixmetal.Provider{}

	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)

	if err != nil {
		t.Fatal(err)
	}

	if len(addrs) != 2 {
		t.Fatalf("bad: %v", addrs)
	}
}

func TestAddrsTagInclude(t *testing.T) {
	args := discover.Config{
		"provider":     "metal",
		"auth_token":   os.Getenv("METAL_AUTH_TOKEN"),
		"project":      os.Getenv("METAL_PROJECT"),
		"address_type": "private_v4",
		"tag":          "tag1",
	}

	if args["auth_token"] == "" {
		t.Skip("Equinix Metal credentials missing")
	}

	if args["project"] == "" {
		t.Skip("Equinix Metal project UUID missing")
	}

	p := equinixmetal.Provider{}

	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)

	if err != nil {
		t.Fatal(err)
	}

	if len(addrs) != 2 {
		t.Fatalf("bad: %v", addrs)
	}
}

func TestAddrsTagIncludeMulti(t *testing.T) {
	args := discover.Config{
		"provider":     "metal",
		"auth_token":   os.Getenv("METAL_AUTH_TOKEN"),
		"project":      os.Getenv("METAL_PROJECT"),
		"address_type": "private_v4",
		"tag":          "tag1,tag2",
	}

	if args["auth_token"] == "" {
		t.Skip("Equinix Metal credentials missing")
	}

	if args["project"] == "" {
		t.Skip("Equinix Metal project UUID missing")
	}

	p := equinixmetal.Provider{}

	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)

	if err != nil {
		t.Fatal(err)
	}

	if len(addrs) != 3 {
		t.Fatalf("bad: %v", addrs)
	}
}
