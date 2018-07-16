package packet

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
)

func TestAddrs(t *testing.T) {
	args := discover.Config{
		"provider":          "packet",
		"packet_auth_token": os.Getenv("PACKET_TOKEN"),
		"packet_project":    "93125c2a-8b78-4d4f-a3c4-7367d6b7cca8",
	}

	p := Provider{}

	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)

	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) == 0 {
		t.Fatalf("bad: %v", addrs)
	}
}
