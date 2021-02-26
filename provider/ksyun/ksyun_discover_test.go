package ksyun_test

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/ksyun"
)

var _ discover.Provider = (*ksyun.Provider)(nil)
var _ discover.ProviderWithUserAgent = (*ksyun.Provider)(nil)

func TestAddrs(t *testing.T) {
	args := discover.Config{
		"provider":          "ksyun",
		"region":            os.Getenv("KSYUN_REGION"),
		"tag_key":           "ChargeType",
		"tag_value":         "Daily",
		"access_key_id":     os.Getenv("KSYUN_ACCESS_KEY"),
		"access_key_secret": os.Getenv("KSYUN_SECRET_KEY"),
	}

	if args["region"] == "" || args["access_key_id"] == "" || args["access_key_secret"] == "" {
		t.Skip("Ksyun credentials or region missing")
	}

	p := &ksyun.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 2 {
		t.Fatalf("bad: %v", addrs)
	}
}
