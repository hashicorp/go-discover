// Copyright IBM Corp. 2017, 2025
// SPDX-License-Identifier: MPL-2.0

// +build appengine gopherjs

package logrus

import (
	"io"
)

func checkIfTerminal(w io.Writer) bool {
	return true
}
