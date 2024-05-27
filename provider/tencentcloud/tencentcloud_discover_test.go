// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tencentcloud_test

import (
	"log"
	"os"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/tencentcloud"
	"github.com/stretchr/testify/require"
)

var _ discover.Provider = (*tencentcloud.Provider)(nil)
var _ discover.ProviderWithUserAgent = (*tencentcloud.Provider)(nil)

func TestAddrs(t *testing.T) {
	args := discover.Config{
		"provider":          "tencentcloud",
		"access_key_id":     os.Getenv("TENCENTCLOUD_SECRET_ID"),
		"access_key_secret": os.Getenv("TENCENTCLOUD_SECRET_KEY"),
		"region":            os.Getenv("TENCENTCLOUD_REGION"),
		"tag_key":           "consul",
		"tag_value":         "test",
		"address_type":      "private_v4",
	}

	if args["access_key_id"] == "" || args["access_key_secret"] == "" || args["region"] == "" {
		t.Skip("TencentCloud credentials or region info missing")
	}

	p := &tencentcloud.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)

	addrs, err := p.Addrs(args, l)
	require.NoError(t, err)
	require.Equal(t, len(addrs), 2)
}
