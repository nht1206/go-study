package fqp

import "regexp"

func existInSlice(ch rune, s []rune) bool {
	for _, e := range s {
		if ch == e {
			return true
		}
	}

	return false
}

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

func isJoin(literal string) bool {
	switch literal {
	case "&&", "||":
		return true
	default:
		return false
	}
}
