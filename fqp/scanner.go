package fqp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

const (
	eof = rune(0)
	dot = rune('.')
)

var (
	defaultWhitespaces     = []rune(` \t\n`)
	defaultIdentifierStart = []rune(`@_#`)
	defaultTextStart       = []rune(`'"`)
	defaultOperator        = []rune(`=?!><~`)
	defaultJoin            = []rune(`&|`)
	defaultGroupStart      = []rune(`(`)
	defaultGroupEnd        = []rune(`)`)
)

type Scanner struct {
	r               *bufio.Reader
	nextToken       *Token
	whiteSpaces     []rune
	identifierStart []rune
	textStart       []rune
	operator        []rune
	join            []rune
	groupStart      []rune
	groupEnd        []rune
	Err             error
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		r:               bufio.NewReader(r),
		whiteSpaces:     defaultWhitespaces,
		identifierStart: defaultIdentifierStart,
		textStart:       defaultTextStart,
		operator:        defaultOperator,
		join:            defaultJoin,
		groupStart:      defaultGroupStart,
		groupEnd:        defaultGroupEnd,
	}
}

func (s *Scanner) HasNext() bool {
	if s.nextToken != nil {
		return true
	}

	s.nextToken = s.scan()

	return s.nextToken.Type != TokenEOF && s.Err == nil
}

func (s *Scanner) Scan() *Token {
	if s.nextToken == nil {
		return s.scan()
	}

	defer func() {
		s.nextToken = nil
	}()

	return s.nextToken
}

func (s *Scanner) scan() *Token {
	ch := s.read()

	if ch == eof {
		return &Token{Type: TokenEOF, Literal: string(ch)}
	}

	if s.isWhitespace(ch) {
		s.unread()
		return s.scanWhitespaces()
	}

	if s.isTextStart(ch) {
		s.unread()
		return s.scanText()
	}

	if isNumberStart(ch) {
		s.unread()
		return s.scanNumber()
	}

	if s.isIdentifierStart(ch) {
		s.unread()
		return s.scanIdentifier()
	}

	if s.isOperator(ch) {
		s.unread()
		return s.scanOperator()
	}

	if s.isJoin(ch) {
		s.unread()
		return s.scanJoin()
	}

	if s.isGroupStart(ch) {
		s.unread()
		return s.scanGroup()
	}

	return &Token{Type: TokenUnsupported, Literal: string(ch)}
}

func (s *Scanner) scanWhitespaces() *Token {
	var buf bytes.Buffer
	for {
		ch := s.read()

		if ch == eof {
			break
		}

		if !s.isWhitespace(ch) {
			s.unread()
			break
		}

		buf.WriteRune(ch)
	}

	return &Token{Type: TokenWhiteSpace, Literal: buf.String()}
}

func (s *Scanner) scanIdentifier() *Token {
	var buf bytes.Buffer
	identifierStart := s.read()
	buf.WriteRune(identifierStart)

	for {
		ch := s.read()

		if ch == eof {
			break
		}

		if !isLetterRune(ch) && !isDigitRune(ch) && ch != dot {
			s.unread()
			break
		}

		buf.WriteRune(ch)
	}

	literal := buf.String()

	if !isIdentifier(literal) {
		s.Err = fmt.Errorf("invalid identifier %q", literal)
	}

	return &Token{Type: TokenIdentifier, Literal: literal}
}

func (s *Scanner) scanText() *Token {
	var buf bytes.Buffer

	quoteCh := s.read()
	buf.WriteRune(quoteCh)
	var prevCh rune
	hasCloseQuote := false

	for {
		ch := s.read()

		if ch == eof {
			break
		}

		buf.WriteRune(ch)

		if ch == quoteCh && prevCh != '\\' {
			hasCloseQuote = true
			break
		}
	}

	literal := buf.String()

	if !hasCloseQuote {
		s.Err = fmt.Errorf("invalid quoted text %q", literal)
	} else {
		literal = literal[1 : len(literal)-1]
		quoteStr := string(quoteCh)
		literal = strings.Replace(literal, `\`+quoteStr, quoteStr, -1)
	}

	return &Token{Type: TokenText, Literal: literal}
}

func (s *Scanner) scanNumber() *Token {
	var buf bytes.Buffer

	ch := s.read()
	if isNumberStart(ch) {
		buf.WriteRune(ch)
	}

	for {
		ch := s.read()

		if ch == eof {
			break
		}

		if !isDigitRune(ch) {
			s.unread()
			break
		}

		buf.WriteRune(ch)
	}

	return &Token{Type: TokenNumber, Literal: buf.String()}
}

func (s *Scanner) scanOperator() *Token {
	var buf bytes.Buffer

	for {
		ch := s.read()

		if ch == eof {
			break
		}

		if !s.isOperator(ch) {
			s.unread()
			break
		}

		buf.WriteRune(ch)
	}

	return &Token{Type: TokenOperator, Literal: buf.String()}
}

func (s *Scanner) scanJoin() *Token {
	var buf bytes.Buffer

	for {
		ch := s.read()

		if ch == eof {
			break
		}

		if !s.isJoin(ch) {
			s.unread()
			break
		}

		buf.WriteRune(ch)
	}

	return &Token{Type: TokenOperator, Literal: buf.String()}
}

func (s *Scanner) read() rune {

	r, _, err := s.r.ReadRune()
	if err != nil {
		if err == io.EOF {
			return eof
		}
		s.Err = err
	}

	return r
}

func (s *Scanner) scanGroup() *Token {
	var buf bytes.Buffer

	startGroup := s.read()
	buf.WriteRune(startGroup)
	hasEndGroup := false

	for {
		ch := s.read()

		if ch == eof {
			break
		}

		buf.WriteRune(ch)
		if s.isGroupEnd(ch) {
			hasEndGroup = true
			break
		}
	}

	literal := buf.String()

	if !hasEndGroup {
		s.Err = fmt.Errorf("invalid group %q", literal)
	}

	return &Token{Type: TokenGroup, Literal: literal}
}

func (s *Scanner) unread() error {
	return s.r.UnreadRune()
}

func (s *Scanner) isWhitespace(ch rune) bool {
	return existInSlice(ch, s.whiteSpaces)
}

func (s *Scanner) isIdentifierStart(ch rune) bool {
	return isLetterRune(ch) || existInSlice(ch, s.identifierStart)
}

func (s *Scanner) isTextStart(ch rune) bool {
	return existInSlice(ch, s.textStart)
}

func (s *Scanner) isOperator(ch rune) bool {
	return existInSlice(ch, s.operator)
}

func (s *Scanner) isJoin(ch rune) bool {
	return existInSlice(ch, s.join)
}

func (s *Scanner) isGroupStart(ch rune) bool {
	return existInSlice(ch, s.groupStart)
}

func (s *Scanner) isGroupEnd(ch rune) bool {
	return existInSlice(ch, s.groupEnd)
}
