// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// +build gofuzz

package ini

import (
	"bytes"
)

func Fuzz(data []byte) int {
	b := bytes.NewReader(data)

	if _, err := Parse(b); err != nil {
		return 0
	}

	return 1
}
