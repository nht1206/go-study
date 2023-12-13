package tokenizer

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
	defaultWhitespaces          = []rune(" \t\n")
	defaultIdentifierStartRunes = []rune(`@_#`)
	defaultTextStartRunes       = []rune(`'"`)
	defaultOperatorStartRunes   = []rune(`=?!><~`)
	defaultJoinStartRunes       = []rune(`&|`)
	defaultGroupStartRunes      = []rune(`(`)
	defaultGroupEndRunes        = []rune(`)`)
)

type TokenType string

const (
	TokenWhiteSpace  TokenType = "whitespace"
	TokenText        TokenType = "text"
	TokenNumber      TokenType = "number"
	TokenIdentifier  TokenType = "identifier"
	TokenOperator    TokenType = "operator"
	TokenJoin        TokenType = "join"
	TokenGroup       TokenType = "group"
	TokenEOF         TokenType = "eof"
	TokenUnsupported TokenType = "unsupported"
)

type Token struct {
	Type    TokenType
	Literal string
}

type Tokenizer struct {
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

func NewTokenizer(r io.Reader) *Tokenizer {
	return &Tokenizer{
		r:               bufio.NewReader(r),
		whiteSpaces:     defaultWhitespaces,
		identifierStart: defaultIdentifierStartRunes,
		textStart:       defaultTextStartRunes,
		operator:        defaultOperatorStartRunes,
		join:            defaultJoinStartRunes,
		groupStart:      defaultGroupStartRunes,
		groupEnd:        defaultGroupEndRunes,
	}
}

func (s *Tokenizer) HasNext() bool {
	if s.nextToken != nil {
		return true
	}

	s.nextToken = s.scan()

	return s.nextToken.Type != TokenEOF && s.Err == nil
}

func (s *Tokenizer) Scan() *Token {
	if s.nextToken == nil {
		return s.scan()
	}

	defer func() {
		s.nextToken = nil
	}()

	return s.nextToken
}

func (s *Tokenizer) scan() *Token {
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

func (s *Tokenizer) scanWhitespaces() *Token {
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

func (s *Tokenizer) scanIdentifier() *Token {
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

func (s *Tokenizer) scanText() *Token {
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

func (s *Tokenizer) scanNumber() *Token {
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

func (s *Tokenizer) scanOperator() *Token {
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

func (s *Tokenizer) scanJoin() *Token {
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

	return &Token{Type: TokenJoin, Literal: buf.String()}
}

func (s *Tokenizer) read() rune {

	r, _, err := s.r.ReadRune()
	if err != nil {
		if err == io.EOF {
			return eof
		}
		s.Err = err
	}

	return r
}

func (s *Tokenizer) scanGroup() *Token {
	var buf bytes.Buffer

	// Skip the group start rune
	s.read()
	hasEndGroup := false

	for {
		ch := s.read()

		if ch == eof {
			break
		}

		if s.isGroupEnd(ch) {
			hasEndGroup = true
			break
		}

		buf.WriteRune(ch)
	}

	literal := buf.String()

	if !hasEndGroup {
		s.Err = fmt.Errorf("invalid group %q", literal)
	}

	return &Token{Type: TokenGroup, Literal: literal}
}

func (s *Tokenizer) unread() error {
	return s.r.UnreadRune()
}

func (s *Tokenizer) isWhitespace(ch rune) bool {
	return existInSlice(ch, s.whiteSpaces)
}

func (s *Tokenizer) isIdentifierStart(ch rune) bool {
	return isLetterRune(ch) || existInSlice(ch, s.identifierStart)
}

func (s *Tokenizer) isTextStart(ch rune) bool {
	return existInSlice(ch, s.textStart)
}

func (s *Tokenizer) isOperator(ch rune) bool {
	return existInSlice(ch, s.operator)
}

func (s *Tokenizer) isJoin(ch rune) bool {
	return existInSlice(ch, s.join)
}

func (s *Tokenizer) isGroupStart(ch rune) bool {
	return existInSlice(ch, s.groupStart)
}

func (s *Tokenizer) isGroupEnd(ch rune) bool {
	return existInSlice(ch, s.groupEnd)
}
