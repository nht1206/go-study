package fqp

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
