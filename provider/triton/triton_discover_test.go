// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package triton_test

import (
	"log"
	"os"
	"testing"

	"github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/triton"
)

func TestAddrs(t *testing.T) {
	args := discover.Config{
		"provider":  "triton",
		"account":   os.Getenv("TRITON_ACCOUNT"),
		"url":       os.Getenv("TRITON_URL"),
		"key_id":    os.Getenv("TRITON_KEY_ID"),
		"tag_key":   "consul-role",
		"tag_value": "server",
	}

	if args["account"] == "" || args["url"] == "" || args["key_id"] == "" {
		t.Skip("Triton credentials or url missing")
	}

	l := log.New(os.Stderr, "", log.LstdFlags)
	p := &triton.Provider{}
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 2 {
		t.Fatalf("bad: %v", addrs)
	}
}
