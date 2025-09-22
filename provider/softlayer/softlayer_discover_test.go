// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package softlayer_test

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/softlayer"
)

func TestAddrs(t *testing.T) {
	args := discover.Config{
		"provider":   "softlayer",
		"username":   os.Getenv("SL_USERNAME"),
		"api_key":    os.Getenv("SL_API_KEY"),
		"datacenter": "dal06",
		"tag_value":  "consul-server",
	}
	if args["username"] == "" || args["api_key"] == "" {
		t.Skip("SoftLayer credentials missing")
	}

	p := &softlayer.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 3 {
		t.Fatalf("bad: %v", addrs)
	}
}
