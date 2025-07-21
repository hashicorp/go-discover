// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// +build appengine gopherjs

package logrus

import (
	"io"
)

func checkIfTerminal(w io.Writer) bool {
	return true
}
