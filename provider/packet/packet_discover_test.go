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

func TestAddrs(t *testing.T) {
	args := discover.Config{
		"provider":          "packet",
		"packet_auth_token": os.Getenv("PACKET_TOKEN"),
		"packet_project":    os.Getenv("PACKET_PROJECT"),
	}

	if args["packet_auth_token"] == "" {
		t.Skip("Packet credentials missing")
	}
	p := packet.Provider{}

	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)

	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) == 0 {
		t.Fatalf("bad: %v", addrs)
	}
}
