package tokenizer

import "regexp"

func isLetterRune(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigitRune(ch rune) bool {
	return (ch >= '0' && ch <= '9')
}

func isNumberStart(ch rune) bool {
	return ch == '-' || isDigitRune(ch)
}

var identifierRegex = regexp.MustCompile(`^[\@\#\_]?[\w\.\:]*\w+$`)

func isIdentifier(literal string) bool {
	return identifierRegex.MatchString(literal)
}
