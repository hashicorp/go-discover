package aliyun_test

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/aliyun"
)

func TestAddrs(t *testing.T) {
	args := discover.Config{
		"provider":          "aliyun",
		"region":            os.Getenv("ALIYUN_REGION"),
		"tag_key":           "consul",
		"tag_value":         "server.test",
		"access_key_id":     os.Getenv("ALIYUN_ACCESS_KEY_ID"),
		"access_key_secret": os.Getenv("ALIYUN_ACCESS_KEY_SECRET"),
	}

	if args["region"] == "" || args["access_key_id"] == "" || args["access_key_secret"] == "" {
		t.Skip("Aliyun credentials or region missing")
	}

	p := &aliyun.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) == 0 {
		t.Fatalf("bad: %v", addrs)
	}
}