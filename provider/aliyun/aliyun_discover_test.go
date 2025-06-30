// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package aliyun_test

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/aliyun"
)

var _ discover.Provider = (*aliyun.Provider)(nil)
var _ discover.ProviderWithUserAgent = (*aliyun.Provider)(nil)

func TestAddrs(t *testing.T) {
	args := discover.Config{
		"provider":          "aliyun",
		"region":            os.Getenv("ALICLOUD_REGION"),
		"tag_key":           "consul",
		"tag_value":         "server.test",
		"access_key_id":     os.Getenv("ALICLOUD_ACCESS_KEY"),
		"access_key_secret": os.Getenv("ALICLOUD_SECRET_KEY"),
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
	if len(addrs) != 2 {
		t.Fatalf("bad: %v", addrs)
	}
}
