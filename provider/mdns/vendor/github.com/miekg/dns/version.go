// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dns

import "fmt"

// Version is current version of this library.
var Version = V{1, 0, 14}

// V holds the version of this library.
type V struct {
	Major, Minor, Patch int
}

func (v V) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}
