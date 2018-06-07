package xrandr

import "fmt"

const (
	TokenTypeIllegal TokenType = iota
	TokenTypeEOF

	// Lexical tokens, all are significant.
	TokenTypePunctuator     // One of: ():,
	TokenTypeName           // Regex: /[A-Za-z0-9_-]+/
	TokenTypeIntValue       // Regex: /[0-9]+/
	TokenTypeFloatValue     // Regex: /([0-9]+)(\.[0-9]+)?/
	TokenTypeWhiteSpace     // Regex: /\s/
	TokenTypeLineTerminator // Unicode: "\u000A", "\u000D"
)

var typeNames = map[TokenType]string{
	TokenTypeIllegal:        fmt.Sprintf("illegal (%d)", TokenTypeIllegal),
	TokenTypeEOF:            fmt.Sprintf("eof (%d)", TokenTypeEOF),
	TokenTypePunctuator:     fmt.Sprintf("punctuator (%d)", TokenTypePunctuator),
	TokenTypeName:           fmt.Sprintf("name (%d)", TokenTypeName),
	TokenTypeIntValue:       fmt.Sprintf("int (%d)", TokenTypeIntValue),
	TokenTypeFloatValue:     fmt.Sprintf("float (%d)", TokenTypeFloatValue),
	TokenTypeWhiteSpace:     fmt.Sprintf("whitespace (%d)", TokenTypeWhiteSpace),
	TokenTypeLineTerminator: fmt.Sprintf("line terminator (%d)", TokenTypeLineTerminator),
}

// TokenType represents the type of a lexical token.
type TokenType int

func (t TokenType) String() string {
	return typeNames[t]
}

// Token represents a lexical token, used to assist parsing.
type Token struct {
	// Type is the type of the token that has been lexed.
	Type TokenType
	// Literal is the raw literal string.
	Literal string
	// The starting position of this token in the input, in runes.
	Position int
	// The line number at the start of this token.
	Line int
}
