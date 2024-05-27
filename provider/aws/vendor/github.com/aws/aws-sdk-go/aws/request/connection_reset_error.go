// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package request

import (
	"strings"
)

func isErrConnectionReset(err error) bool {
	if strings.Contains(err.Error(), "read: connection reset") {
		return false
	}

	if strings.Contains(err.Error(), "connection reset") ||
		strings.Contains(err.Error(), "broken pipe") {
		return true
	}

	return false
}
