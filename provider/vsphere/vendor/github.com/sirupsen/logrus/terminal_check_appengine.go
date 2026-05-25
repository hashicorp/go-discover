// Copyright IBM Corp. 2017, 2026

// +build appengine gopherjs

package logrus

import (
	"io"
)

func checkIfTerminal(w io.Writer) bool {
	return true
}
