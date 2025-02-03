// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// +build !appengine,!gopherjs

package logrus

import (
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func checkIfTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return terminal.IsTerminal(int(v.Fd()))
	default:
		return false
	}
}
