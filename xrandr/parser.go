package xrandr

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	dimensionPattern  = regexp.MustCompile(`^([0-9]+)mm$`)
	resolutionPattern = regexp.MustCompile(`^([0-9]+)x([0-9]+)i?$`)
)

type Parser struct {
	lexer  *Lexer
	token  Token
	skipWS bool
}

func NewParser() *Parser {
	return &Parser{
		skipWS: true,
	}
}

func (p *Parser) ParseProps(input []byte) (PropsOutput, error) {
	var props PropsOutput

	// This isn't thread-safe.
	p.lexer = NewLexer(input)
	p.scan()

	// "Parse" the screen, we actually just skip it entirely really.
	err := p.parseScreen()
	if err != nil {
		return props, err
	}

	for {
		output, err := p.parseOutput()
		if err != nil {
			return props, err
		}

		props.Outputs = append(props.Outputs, output)

		if p.token.Type == TokenTypeEOF {
			break
		}
	}

	return props, nil
}

func (p *Parser) parseOutputName(output *Output) error {
	p.skipWS = true

	tok, err := p.consumeType(TokenTypeName)
	if err != nil {
		return err
	}

	output.Name = tok.Literal

	return nil
}

func (p *Parser) parseOutputStatus(output *Output) error {
	p.skipWS = true

	tok, err := p.consume(TokenTypeName, "connected", "disconnected")
	if err == nil {
		output.IsConnected = tok.Literal == "connected"
	}

	if p.skip(TokenTypeName, "primary") {
		output.IsPrimary = true
	}

	return nil
}

func (p *Parser) parseResolutionAndPosition(output *Output) error {
	p.skipWS = true

	// If the output is enabled, we should see the current resolution, and the position.
	if p.token.Type == TokenTypeName {
		isRes, res := p.parseResolution(p.token.Literal)
		if !isRes {
			return nil
		}

		output.IsEnabled = true
		output.Resolution = res

		err := p.scan()
		if err != nil {
			return err
		}

		if err := p.skipWithLiteral(TokenTypePunctuator, "+"); err != nil {
			return err
		}

		tok, err := p.consume(TokenTypeIntValue)
		if err != nil {
			return err
		}

		offsetX, err := strconv.ParseInt(tok.Literal, 10, 64)
		if err != nil {
			return err
		}

		if err := p.skipWithLiteral(TokenTypePunctuator, "+"); err != nil {
			return err
		}

		tok, err = p.consume(TokenTypeIntValue)
		if err != nil {
			return err
		}

		offsetY, err := strconv.ParseInt(tok.Literal, 10, 64)
		if err != nil {
			return err
		}

		output.Position.OffsetX = int(offsetX)
		output.Position.OffsetY = int(offsetY)
	}

	return nil
}

func (p *Parser) parseOutputRotationAndReflection(output *Output) error {
	p.skipWS = true

	if p.token.Type == TokenTypeName {
		var found bool
		switch p.token.Literal {
		case "normal":
			output.Rotation = RotationNormal
			found = true
		case "left":
			output.Rotation = RotationLeft
			found = true
		case "inverted":
			output.Rotation = RotationInverted
			found = true
		case "right":
			output.Rotation = RotationRight
			found = true
		default:
			found = false
		}

		if found {
			err := p.scan()
			if err != nil {
				return err
			}
		}
	}

	if p.token.Type == TokenTypeName {
		var foundX bool
		var foundY bool

		switch p.token.Literal {
		case "x":
			foundX = true
		case "y":
			foundY = true
		}

		if foundX || foundY {
			err := p.scan()
			if err != nil {
				return err
			}

			if p.token.Type == TokenTypeName && p.token.Literal == "axis" {
				if foundX {
					output.Reflection = ReflectionXAxis
				}

				if foundY {
					output.Reflection = ReflectionYAxis
				}

				err := p.scan()
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (p *Parser) parseOutputRotationAndReflectionKey() error {
	p.skipWS = true

	if p.token.Type != TokenTypePunctuator && p.token.Literal == "(" {
		return nil
	}

	return p.all(
		p.skipWithLiteral(TokenTypePunctuator, "("),
		p.skipWithLiteral(TokenTypeName, "normal"),
		p.skipWithLiteral(TokenTypeName, "left"),
		p.skipWithLiteral(TokenTypeName, "inverted"),
		p.skipWithLiteral(TokenTypeName, "right"),
		p.skipWithLiteral(TokenTypeName, "x"),
		p.skipWithLiteral(TokenTypeName, "axis"),
		p.skipWithLiteral(TokenTypeName, "y"),
		p.skipWithLiteral(TokenTypeName, "axis"),
		p.skipWithLiteral(TokenTypePunctuator, ")"),
	)

	return nil
}

func (p *Parser) parseOutputDimensions(output *Output) error {
	p.skipWS = true

	// We probably hit the end of the line here.
	if p.token.Type != TokenTypeName {
		if p.token.Type == TokenTypeLineTerminator {
			// We _might_ hit properties next, so we have to do this in advance.
			p.skipWS = false

			if err := p.scan(); err != nil {
				return err
			}
		}

		return nil
	}

	xdim, err := p.parseOutputDimension()
	if err != nil {
		return err
	}

	err = p.skipWithLiteral(TokenTypeName, "x")
	if err != nil {
		return err
	}

	ydim, err := p.parseOutputDimension()
	if err != nil {
		return err
	}

	output.Dimensions.Width = xdim
	output.Dimensions.Height = ydim

	if p.token.Type == TokenTypeLineTerminator {
		// We might start looking for properties next, so should stop skipping whitespace.
		p.skipWS = false

		if err := p.scan(); err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) parseOutputDimension() (uint, error) {
	p.skipWS = true

	tok, err := p.consume(TokenTypeName)
	if err != nil {
		return 0, err
	}

	matches := dimensionPattern.FindStringSubmatch(tok.Literal)

	if len(matches) != 2 {
		return 0, err
	}

	dim, err := strconv.ParseUint(matches[1], 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(dim), nil
}

func (p *Parser) parseProperties(output *Output) error {
	// Stop skipping whitespace.
	p.skipWS = false

	if err := p.skipWithLiteral(TokenTypeWhiteSpace, "\t"); err != nil {
		return err
	}

	for {
		stop, err := p.parseProperty(output)
		if err != nil {
			return err
		}

		if stop {
			break
		}

		if p.token.Type != TokenTypeName {
			break
		}
	}

	return nil
}

func (p *Parser) parseProperty(output *Output) (bool, error) {
	var name string
	var value string
	var stop bool

	tok, err := p.consume(TokenTypeName)
	if err != nil {
		return stop, err
	}

	// Gather up the entire name. Including any spaces, etc. Until we hit a ':'.
	name = tok.Literal

	for {
		if p.token.Type == TokenTypePunctuator && p.token.Literal == ":" {
			break
		}

		name += p.token.Literal

		if err := p.scan(); err != nil {
			return stop, err
		}
	}

	err = p.all(
		p.skipWithLiteral(TokenTypePunctuator, ":"),
		p.skipWithLiteral(TokenTypeWhiteSpace, " "),
	)

	if err != nil {
		return stop, err
	}

	if p.token.Type == TokenTypeLineTerminator {
		for {
			if err := p.scan(); err != nil {
				return stop, err
			}

			// We're no longer processing properties if we've hit something that's not a tab at the
			// start of a new line.
			if p.token.Type != TokenTypeWhiteSpace || p.token.Literal != "\t" {
				stop = true
				break
			}

			if err := p.scan(); err != nil {
				return stop, err
			}

			// If we don't get a second tab, we've hit a new property. So, we need to bail from this
			// loop iteration.
			if p.token.Type != TokenTypeWhiteSpace || p.token.Literal != "\t" {
				break
			}

			for {
				if err := p.scan(); err != nil {
					return stop, err
				}

				if p.token.Type == TokenTypeLineTerminator {
					break
				}

				value += p.token.Literal
			}
		}
	} else if p.token.Type == TokenTypeName || p.token.Type == TokenTypeIntValue || p.token.Type == TokenTypeFloatValue {
		value += p.token.Literal

		for {
			if err := p.scan(); err != nil {
				return stop, err
			}

			// Consume the value that's on the same line.
			if p.token.Type == TokenTypeLineTerminator {
				break
			}

			value += p.token.Literal
		}

		// Then, consume everything else after it until we hit another thing that looks like a
		// new property.
		for {
			if err := p.scan(); err != nil {
				return stop, err
			}

			// We're no longer processing properties if we've hit something that's not a tab at the
			// start of a new line.
			if p.token.Type != TokenTypeWhiteSpace || p.token.Literal != "\t" {
				stop = true
				break
			}

			if err := p.scan(); err != nil {
				return stop, err
			}

			// If we don't get a second tab, we've hit a new property. So, we need to bail.
			if p.token.Type != TokenTypeWhiteSpace || p.token.Literal != "\t" {
				break
			}

			// Skip past the "value"
			for {
				if err := p.scan(); err != nil {
					return stop, err
				}

				if p.token.Type == TokenTypeLineTerminator {
					break
				}
			}
		}
	}

	output.Properties[strings.TrimSpace(name)] = strings.TrimSpace(value)

	return stop, nil
}

func (p *Parser) parseModes(output *Output) error {
	p.skipWS = true

	// Sometimes we don't have modes to parse.
	if p.token.Type != TokenTypeWhiteSpace && p.token.Literal != " " {
		return nil
	}

	if err := p.scan(); err != nil {
		return err
	}

	for {
		if p.token.Type != TokenTypeName {
			break
		}

		var mode OutputMode

		isRes, res := p.parseResolution(p.token.Literal)
		if !isRes {
			// We've probably just hit the end of resolutions at this point, and are looking at the
			// next output.
			return nil
		}

		mode.Resolution = res

		if err := p.scan(); err != nil {
			return err
		}

		for {
			if !p.nextType(TokenTypeFloatValue) {
				break
			}

			var rate Rate
			var err error

			rate.Rate, err = strconv.ParseFloat(p.token.Literal, 64)
			if err != nil {
				return err
			}

			if err := p.scan(); err != nil {
				return err
			}

			if p.skip(TokenTypePunctuator, "*") {
				rate.IsCurrent = true
			}

			if p.skip(TokenTypePunctuator, "+") {
				rate.IsPreferred = true
			}

			mode.Rates = append(mode.Rates, rate)
		}

		p.skipType(TokenTypeLineTerminator)

		output.Modes = append(output.Modes, mode)
	}

	return nil
}

func (p *Parser) parseOutput() (Output, error) {
	output := Output{}
	output.Properties = make(map[string]string)
	output.Rotation = RotationNormal
	output.Reflection = ReflectionNone

	err := p.parseOutputName(&output)
	if err != nil {
		return output, fmt.Errorf("error parsing output name: %v", err)
	}

	err = p.parseOutputStatus(&output)
	if err != nil {
		return output, fmt.Errorf("error parsing output status: %v", err)
	}

	err = p.parseResolutionAndPosition(&output)
	if err != nil {
		return output, fmt.Errorf("error parsing output resolution and position: %v", err)
	}

	err = p.parseOutputRotationAndReflection(&output)
	if err != nil {
		return output, fmt.Errorf("error parsing output rotation and reflection: %v", err)
	}

	err = p.parseOutputRotationAndReflectionKey()
	if err != nil {
		return output, fmt.Errorf("error parsing output rotation and reflection key: %v", err)
	}

	err = p.parseOutputDimensions(&output)
	if err != nil {
		return output, fmt.Errorf("error parsing output dimensions: %v", err)
	}

	err = p.parseProperties(&output)
	if err != nil {
		return output, fmt.Errorf("error parsing output properties: %v", err)
	}

	err = p.parseModes(&output)
	if err != nil {
		return output, fmt.Errorf("error parsing output modes: %v", err)
	}

	return output, nil
}

func (p *Parser) parseResolution(literal string) (bool, Resolution) {
	var res Resolution

	matches := resolutionPattern.FindStringSubmatch(literal)

	if len(matches) != 3 {
		return false, res
	}

	xres, err := strconv.ParseUint(matches[1], 10, 64)
	if err != nil {
		return false, res
	}

	yres, err := strconv.ParseUint(matches[2], 10, 64)
	if err != nil {
		return false, res
	}

	res.Width = uint(xres)
	res.Height = uint(yres)

	return true, res
}

func (p *Parser) parseScreen() error {
	// Scan, and skip all expectations here.
	return p.expectAll(
		p.expectFn(TokenTypeName, "Screen"),
		p.expectTypeFn(TokenTypeIntValue),
		p.expectTypeFn(TokenTypePunctuator),
		p.expectFn(TokenTypeName, "minimum"),
		p.expectTypeFn(TokenTypeIntValue),
		p.expectFn(TokenTypeName, "x"),
		p.expectTypeFn(TokenTypeIntValue),
		p.expectTypeFn(TokenTypePunctuator),
		p.expectFn(TokenTypeName, "current"),
		p.expectTypeFn(TokenTypeIntValue),
		p.expectFn(TokenTypeName, "x"),
		p.expectTypeFn(TokenTypeIntValue),
		p.expectTypeFn(TokenTypePunctuator),
		p.expectFn(TokenTypeName, "maximum"),
		p.expectTypeFn(TokenTypeIntValue),
		p.expectFn(TokenTypeName, "x"),
		p.expectTypeFn(TokenTypeIntValue),
		p.expectTypeFn(TokenTypeLineTerminator),
	)
}

// Parser utilities:

func (p *Parser) expectAll(fns ...func() error) error {
	for _, fn := range fns {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) consume(t TokenType, ls ...string) (Token, error) {
	token := p.token
	if token.Type != t {
		return token, p.unexpected(token, t, ls...)
	}

	for _, l := range ls {
		if token.Literal != l {
			continue
		}

		p.scan()
		return token, nil
	}

	return token, p.unexpected(token, t, ls...)
}

func (p *Parser) consumeType(t TokenType) (Token, error) {
	token := p.token
	if token.Type == t {
		p.scan()
		return token, nil
	}

	return token, p.unexpected(token, t)
}

func (p *Parser) expect(t TokenType, l string) error {
	if !p.next(t, l) {
		return p.unexpected(p.token, t, l)
	}

	return nil
}

func (p *Parser) expectType(t TokenType) error {
	if !p.nextType(t) {
		return p.unexpected(p.token, t, "")
	}

	return nil
}

func (p *Parser) expectFn(t TokenType, l string) func() error {
	return func() error {
		return p.expect(t, l)
	}
}

func (p *Parser) expectTypeFn(t TokenType) func() error {
	return func() error {
		return p.expectType(t)
	}
}

func (p *Parser) next(t TokenType, l string) bool {
	return p.token.Type == t && p.token.Literal == l
}

func (p *Parser) nextType(t TokenType) bool {
	return p.token.Type == t
}

func (p *Parser) skip(t TokenType, ls ...string) bool {
	if p.token.Type != t {
		return false
	}

	for _, l := range ls {
		if p.token.Literal != l {
			continue
		}

		p.scan()
		return true
	}

	return false
}

func (p *Parser) skipType(t TokenType) bool {
	if p.token.Type == t {
		p.scan()
		return true
	}

	return false
}

func (p *Parser) scan() {
	p.token = p.lexer.Scan()
	if p.skipWS && p.token.Type == TokenTypeWhiteSpace {
		p.scan()
	}
}

func (p *Parser) unexpected(token Token, t TokenType, ls ...string) error {
	if len(ls) == 0 {
		ls = []string{"N/A"}
	}

	return fmt.Errorf(
		"parser error: unexpected token found: %s (%q). Wanted: %s (%q). Line: %d. Column: %d",
		token.Type.String(),
		token.Literal,
		t.String(),
		strings.Join(ls, "|"),
		token.Line,
		token.Position,
	)
}
