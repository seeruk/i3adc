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
	TokenTypeIllegal:        "illegal",
	TokenTypeEOF:            "eof",
	TokenTypePunctuator:     "punctuator",
	TokenTypeName:           "name",
	TokenTypeIntValue:       "int",
	TokenTypeFloatValue:     "float",
	TokenTypeWhiteSpace:     "whitespace",
	TokenTypeLineTerminator: "line terminator",
}

// TokenType represents the type of a lexical token.
type TokenType int

func (t TokenType) String() string {
	return fmt.Sprintf("%s (%d)", typeNames[t], t)
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
