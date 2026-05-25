// Copyright IBM Corp. 2017, 2026

package ini

var commaRunes = []rune(",")

func isComma(b rune) bool {
	return b == ','
}

func newCommaToken() Token {
	return newToken(TokenComma, commaRunes, NoneType)
}
