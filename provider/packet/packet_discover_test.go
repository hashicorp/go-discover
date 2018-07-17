package packet_test

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/packet"
)

var _ discover.Provider = (*packet.Provider)(nil)
var _ discover.ProviderWithUserAgent = (*packet.Provider)(nil)

func TestAddrsDefault(t *testing.T) {
	args := discover.Config{
		"provider":          "packet",
		"packet_auth_token": os.Getenv("PACKET_TOKEN"),
		"packet_project":    os.Getenv("PACKET_PROJECT"),
	}

	if args["packet_auth_token"] == "" {
		t.Skip("Packet credentials missing")
	}

	if args["packet_project"] == "" {
		t.Skip("Packet project UUID missing")
	}

	p := packet.Provider{}

	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)

	if err != nil {
		t.Fatal(err)
	}

	if len(addrs) != 2 {
		t.Fatalf("bad: %v", addrs)
	}
}

func TestAddrsPublicV6(t *testing.T) {
	args := discover.Config{
		"provider":          "packet",
		"packet_auth_token": os.Getenv("PACKET_TOKEN"),
		"packet_project":    os.Getenv("PACKET_PROJECT"),
		"address_type":      "public_v6",
	}

	if args["packet_auth_token"] == "" {
		t.Skip("Packet credentials missing")
	}

	if args["packet_project"] == "" {
		t.Skip("Packet project UUID missing")
	}

	p := packet.Provider{}

	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)

	if err != nil {
		t.Fatal(err)
	}

	if len(addrs) != 2 {
		t.Fatalf("bad: %v", addrs)
	}
}

func TestAddrsPublicV4(t *testing.T) {
	args := discover.Config{
		"provider":          "packet",
		"packet_auth_token": os.Getenv("PACKET_TOKEN"),
		"packet_project":    os.Getenv("PACKET_PROJECT"),
		"address_type":      "public_v4",
	}

	if args["packet_auth_token"] == "" {
		t.Skip("Packet credentials missing")
	}

	if args["packet_project"] == "" {
		t.Skip("Packet project UUID missing")
	}

	p := packet.Provider{}

	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)

	if err != nil {
		t.Fatal(err)
	}

	if len(addrs) != 2 {
		t.Fatalf("bad: %v", addrs)
	}
}
