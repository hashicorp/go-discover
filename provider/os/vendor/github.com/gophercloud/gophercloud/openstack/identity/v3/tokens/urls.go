// Copyright IBM Corp. 2017, 2025
// SPDX-License-Identifier: MPL-2.0

package tokens

import "github.com/gophercloud/gophercloud"

func tokenURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("auth", "tokens")
}
