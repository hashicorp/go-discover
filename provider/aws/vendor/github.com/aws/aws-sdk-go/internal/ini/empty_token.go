// Copyright IBM Corp. 2017, 2025
// SPDX-License-Identifier: MPL-2.0

package ini

// emptyToken is used to satisfy the Token interface
var emptyToken = newToken(TokenNone, []rune{}, NoneType)
