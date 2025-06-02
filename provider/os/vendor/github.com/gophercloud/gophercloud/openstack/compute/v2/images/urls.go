// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package images

import "github.com/gophercloud/gophercloud"

func listDetailURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("images", "detail")
}

func getURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("images", id)
}

func deleteURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("images", id)
}
