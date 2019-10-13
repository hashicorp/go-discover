package vultr_test

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/vultr"
)

var _ discover.Provider = (*vultr.Provider)(nil)
var _ discover.ProviderWithUserAgent = (*vultr.Provider)(nil)

func TestAddrs(t *testing.T) {
	var addrs []string
	var err error

	args := discover.Config{
		"provider":  "vultr",
		"tag_name":  "go-discover-test-tag",
		"region":    "1",
		"api_token": os.Getenv("VULTR_API_KEY"),
	}
	if args["api_token"] == "" {
		t.Skip("Vultr API Key Missing")
	}

	l := log.New(os.Stderr, "", log.LstdFlags)
	p := &vultr.Provider{}
	addrs, err = p.Addrs(args, l)
	if err != nil{
		t.Fatal(err)
	}
	if len(addrs) != 2{
		t.Fatalf("bad: %v", addrs)
	}
}
