package upcloud

import (
	"log"
	"os"
	"testing"

	"github.com/hashicorp/go-discover"
)

func TestAddrs(t *testing.T) {
	args := discover.Config{
		"provider": "upcloud",
		"tag":      "vault-dev",
		"username": os.Getenv("UPCLOUD_API_USERNAME"),
		"password": os.Getenv("UPCLOUD_API_PASSWORD"),
	}

	if args["username"] == "" || args["password"] == "" {
		t.Skip("UpCloud credentials missing")
	}

	p := &Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("UpCloud found the following IP Addresses: %v", addrs)
}
