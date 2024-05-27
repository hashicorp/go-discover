// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package flavors

import (
	"github.com/gophercloud/gophercloud"
)

func getURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("flavors", id)
}

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("flavors", "detail")
}

func createURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("flavors")
}

func deleteURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("flavors", id)
}
