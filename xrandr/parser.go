package xrandr

import (
	"fmt"
)

type Parser struct {
	lexer *Lexer
	token Token
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseProps(input []byte) (PropsOutput, error) {
	var output PropsOutput

	// This isn't thread-safe.
	p.lexer = NewLexer(input)

	err := p.parseScreen()
	if err != nil {
		return output, err
	}

	return output, nil
}

func (p *Parser) parseScreen() error {
	// Scan, and skip all expectations here.
	return p.all(
		p.expectWithLiteral(TokenTypeName, "Screen"),
		p.expect(TokenTypeWhiteSpace),
		p.expect(TokenTypeIntValue),
		p.expect(TokenTypePunctuator),
		p.expect(TokenTypeWhiteSpace),
		p.expectWithLiteral(TokenTypeName, "minimum"),
	)
}

// Parser utilities:

func (p *Parser) all(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) expect(t TokenType) error {
	err := p.scan()
	if err != nil {
		return err
	}

	if p.token.Type != t {
		return fmt.Errorf(
			"syntax error: unexpected token found: %s (%q). Wanted: %s",
			p.token.Type.String(),
			p.token.Literal,
			t.String(),
		)
	}

	return nil
}

func (p *Parser) expectWithLiteral(t TokenType, l string) error {
	err := p.expect(t)
	if err != nil {
		return err
	}

	if p.token.Literal != l {
		return fmt.Errorf(
			"syntax error: unexpected literal found for token type %s: %s",
			t.String(),
			l,
		)
	}

	return nil
}

func (p *Parser) scan() (err error) {
	p.token, err = p.lexer.Scan()
	return err
}
