// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ini

var commaRunes = []rune(",")

func isComma(b rune) bool {
	return b == ','
}

func newCommaToken() Token {
	return newToken(TokenComma, commaRunes, NoneType)
}
