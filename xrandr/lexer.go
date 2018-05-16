package xrandr

import (
	"fmt"
	"unicode/utf8"
	"unsafe"
)

const (
	// er represents an "empty" rune, but is also an invalid one.
	er = rune(-1)
	// eof represents the end of input.
	eof = rune(0)
)

// Lexer holds the state for lexical analysis of xrandr output.
type Lexer struct {
	input    []byte // Raw input is just a byte slice. It is expected to be UTF-8 encoded characters.
	inputLen int    // Length of the input, in bytes.

	// Positional information.
	pos  int // The start position of the last rune read, in bytes.
	lpos int // The start position of the last rune read, in runes, on the current line.
	line int // The current line number.

	// Previous read information.
	lrw int // The width of the last rune read.
}

// NewLexer returns a new lexer, for lexically analysing xrandr output.
func NewLexer(input []byte) *Lexer {
	return &Lexer{
		input:    input,
		inputLen: len(input),
		line:     1,
	}
}

// Scan attempts to read the next significant token from the input. Tokens that are not understood
// will yield an "illegal" token.
func (l *Lexer) Scan() (token Token, err error) {
	r, _ := l.read()

	switch {
	// Names:
	case (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' || r == '-':
		token, err = l.scanNameOrNumber(r)
	// Whitespace:
	case r == ' ' || r == rune(0x0009):
		token, err = l.scanWhitespace(r)
	// Punctuators:
	case r == ':' || r == '(' || r == ')' || r == ',' || r == '*' || r == '+':
		token, err = l.scanPunctuator(r)
	// Line Terminators:
	case r == '\n' || r == '\r':
		token, err = l.scanLineTerminator(r)
	case r == eof:
		token = Token{
			Type:     TokenTypeEOF,
			Position: l.lpos,
			Line:     l.line,
		}
	default:
		token = Token{
			Type:     TokenTypeIllegal,
			Literal:  btos(l.input[l.pos-1 : l.pos]),
			Position: l.lpos,
			Line:     l.line,
		}
	}

	return
}

func (l *Lexer) scanLineTerminator(r rune) (Token, error) {
	byteStart := l.pos - 1
	runeStart := l.lpos

	// If we got a carriage return, we might expect a newline next.
	if r == '\r' {
		r, _ := l.read()
		if r != '\n' {
			// If we don't get what we expected, unread it.
			l.unread()
		}
	}

	// Increment line number.
	l.line++
	l.lpos = 0

	return Token{
		Type:     TokenTypeLineTerminator,
		Literal:  btos(l.input[byteStart:l.pos]),
		Position: runeStart,
		Line:     l.line,
	}, nil
}

func (l *Lexer) scanNameOrNumber(r rune) (Token, error) {
	byteStart := l.pos - 1
	runeStart := l.lpos

	isNumber := r >= '0' || r <= '9' // Until we encounter a letter.
	isFloat := false                 // Until we encounter a '.' and are still a number.

	var done bool
	for !done {
		r, _ := l.read()

		switch {
		case r == eof:
			done = true
		case isNumber && !isFloat && r == '.':
			isFloat = true
			continue
		case r >= '0' && r <= '9':
			// Numbers...
			continue
		case (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r == '_' || r == '-':
			isNumber = false
			// Names...
			continue
		default:
			l.unread()
			done = true
		}
	}

	tokenType := TokenTypeName
	if isNumber {
		tokenType = TokenTypeIntValue
		if isFloat {
			tokenType = TokenTypeFloatValue
		}
	}

	return Token{
		Type:     tokenType,
		Literal:  btos(l.input[byteStart:l.pos]),
		Position: runeStart,
		Line:     l.line,
	}, nil
}

func (l *Lexer) scanNumber(r rune) (Token, error) {
	byteStart := l.pos - 1
	runeStart := l.lpos

	var float bool
	var err error

	r, err = l.readDigits(r)
	if err != nil {
		return Token{}, err
	}

	if r == '.' {
		float = true

		r, _ = l.read()
		r, err = l.readDigits(r)
		if err != nil {
			return Token{}, err
		}
	}

	if r != eof {
		l.unread()
	}

	tokenType := TokenTypeIntValue
	if float {
		tokenType = TokenTypeFloatValue
	}

	return Token{
		Type:     tokenType,
		Literal:  btos(l.input[byteStart:l.pos]),
		Position: runeStart,
		Line:     l.line,
	}, nil
}

func (l *Lexer) scanPunctuator(r rune) (Token, error) {
	byteStart := l.pos - 1
	runeStart := l.lpos

	return Token{
		Type:     TokenTypePunctuator,
		Literal:  btos(l.input[byteStart:l.pos]),
		Position: runeStart,
		Line:     l.line,
	}, nil
}

func (l *Lexer) scanWhitespace(r rune) (Token, error) {
	byteStart := l.pos - 1
	runeStart := l.lpos

	return Token{
		Type:     TokenTypeWhiteSpace,
		Literal:  btos(l.input[byteStart:l.pos]),
		Position: runeStart,
		Line:     l.line,
	}, nil
}

func (l *Lexer) readDigits(r rune) (rune, error) {
	if !(r >= '0' && r <= '9') {
		return eof, fmt.Errorf("invalid number, expected digit but got: %q", r)
	}

	var done bool
	for !done {
		r, _ = l.read()

		switch {
		case r >= '0' && r <= '9':
			continue
		default:
			// No need to unread here. We actually want to read the character after the numbers.
			done = true
		}
	}

	return r, nil
}

// read moves forward in the input, and returns the next rune available. This function also updates
// the position(s) that the lexer keeps track of in the input so the next read continues from where
// the last left off. Returns the EOF rune if we hit the end of the input.
func (l *Lexer) read() (rune, int) {
	if l.pos >= l.inputLen {
		return eof, 0
	}

	var r rune
	var w int
	if sbr := l.input[l.pos]; sbr < utf8.RuneSelf {
		r = rune(sbr)
		w = 1
	} else {
		r, w = utf8.DecodeRune(l.input[l.pos:])
	}

	l.pos += w
	l.lpos++

	l.lrw = w

	return r, w
}

// unread goes back one rune's worth of bytes in the input, changing the
// positions we keep track of.
// Does not currently go back a line.
func (l *Lexer) unread() {
	l.pos -= l.lrw

	if l.pos > 0 {
		// update rune width for further rewind
		_, l.lrw = utf8.DecodeLastRune(l.input[:l.pos])
	} else {
		// If we're already at the start, set this to so we don't end up with a negative position.
		l.lrw = 0
	}

	if l.lpos > 0 {
		l.lpos--
	}
}

// btos takes the given bytes, and turns them into a string, but without allocations.
func btos(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}
