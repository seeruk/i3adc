package xrandr

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Parser struct {
	lexer  *Lexer
	token  Token
	debug  bool
	skipWS bool
}

func NewParser(debug bool) *Parser {
	return &Parser{
		debug:  debug,
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
	tok, err := p.consume(TokenTypeName)
	if err != nil {
		return err
	}

	output.Name = tok.Literal

	return nil
}

func (p *Parser) parseOutputStatus(output *Output) error {
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
	// If the output is enabled, we should see the current resolution, and the position.
	if p.next(TokenTypeName) {
		isRes, res := p.parseResolution(p.token.Literal)
		if !isRes {
			return nil
		}

		output.IsEnabled = true
		output.Resolution = res

		p.scan()

		if err := p.expect(TokenTypePunctuator, "+"); err != nil {
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

		if err := p.expect(TokenTypePunctuator, "+"); err != nil {
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
	// We can ignore this error, we might not have any rotation status.
	if tok, err := p.consume(TokenTypeName, "normal", "left", "inverted", "right"); err == nil {
		switch tok.Literal {
		case "normal":
			output.Rotation = RotationNormal
		case "left":
			output.Rotation = RotationLeft
		case "inverted":
			output.Rotation = RotationInverted
		case "right":
			output.Rotation = RotationRight
		}

		p.scan()
	}

	if tok, err := p.consume(TokenTypeName, "x", "y"); err == nil {
		// If we get 'x' or 'y' we always expect the word 'axis' to follow.
		if _, err := p.consume(TokenTypeName, "axis"); err != nil {
			return err
		}

		switch tok.Literal {
		case "x":
			output.Reflection = ReflectionXAxis
		case "y":
			output.Reflection = ReflectionYAxis
		}
	}

	return nil
}

func (p *Parser) parseOutputRotationAndReflectionKey() error {
	if p.token.Type != TokenTypePunctuator && p.token.Literal == "(" {
		return nil
	}

	return p.expectAll(
		p.expectFn(TokenTypePunctuator, "("),
		p.expectFn(TokenTypeName, "normal"),
		p.expectFn(TokenTypeName, "left"),
		p.expectFn(TokenTypeName, "inverted"),
		p.expectFn(TokenTypeName, "right"),
		p.expectFn(TokenTypeName, "x"),
		p.expectFn(TokenTypeName, "axis"),
		p.expectFn(TokenTypeName, "y"),
		p.expectFn(TokenTypeName, "axis"),
		p.expectFn(TokenTypePunctuator, ")"),
	)
}

func (p *Parser) parseOutputDimensions(output *Output) error {
	// We probably hit the end of the line here.
	if !p.next(TokenTypeName) {
		// We _might_ hit properties next, so we have to do this in advance. If we do, we'll skip
		// past the line terminator too so we're in the right place for property parsing.
		p.skipWS = false
		p.skip(TokenTypeLineTerminator)

		return nil
	}

	xdim, err := p.parseOutputDimension()
	if err != nil {
		return err
	}

	if err = p.expect(TokenTypeName, "x"); err != nil {
		return err
	}

	ydim, err := p.parseOutputDimension()
	if err != nil {
		return err
	}

	output.Dimensions.Width = xdim
	output.Dimensions.Height = ydim

	// We might hit properties next here too, so also turn whitespace skipping off, and move past
	// the new line, if there is one.
	p.skipWS = false
	p.skip(TokenTypeLineTerminator)

	return nil
}

func (p *Parser) parseOutputDimension() (uint, error) {
	tok, err := p.consume(TokenTypeName)
	if err != nil {
		return 0, err
	}

	dimLen := len(tok.Literal)
	if dimLen < 3 {
		return 0, p.unexpected(tok, TokenTypeName, "xxxmm")
	}

	// We can move a byte at a time, because we should only have single-byte runes to deal with.
	for i := 0; i < dimLen; i++ {
		r := rune(tok.Literal[i])

		if i >= dimLen-2 && r != 'm' {
			return 0, p.unexpected(tok, TokenTypeName, "xxxmm")
		} else if i < dimLen-2 && (r < '0' || r > '9') {
			return 0, p.unexpected(tok, TokenTypeName, "xxxmm")
		}
	}

	dim, err := strconv.ParseUint(tok.Literal[:dimLen-2], 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(dim), nil
}

func (p *Parser) parseProperties(output *Output) error {
	if err := p.expect(TokenTypeWhiteSpace, "\t"); err != nil {
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

	for !p.next(TokenTypePunctuator, ":") {
		name += p.token.Literal
		p.scan()
	}

	err = p.expectAll(
		p.expectFn(TokenTypePunctuator, ":"),
		p.expectFn(TokenTypeWhiteSpace, " "),
	)

	if err != nil {
		return stop, err
	}

	if p.skip(TokenTypeLineTerminator) {
		// If we hit a line terminator, then the value should follow on the next line(s).
		for {
			// We're no longer processing properties if we've hit something that's not a tab at the
			// start of a new line.
			if !p.skip(TokenTypeWhiteSpace, "\t") {
				stop = true
				break
			}

			// If we don't get a second tab, we've hit a new property. So, we need to bail from this
			// loop iteration.
			if !p.skip(TokenTypeWhiteSpace, "\t") {
				break
			}

			for {
				value += p.token.Literal

				p.scan()
				if p.skip(TokenTypeLineTerminator) {
					break
				}
			}
		}
	} else if p.next(TokenTypeName) || p.next(TokenTypeIntValue) || p.next(TokenTypeFloatValue) {
		// If instead we hit more tokens after the property name on the same line, then we'll take
		// the value from there, and when we hit a new line, we'll skip anything that isn't a new
		// property, assuming that's it's more like documentation.
		value += p.token.Literal

		// Read the value until we hit the line terminator.
		for !p.skip(TokenTypeLineTerminator) {
			p.scan()
			value += p.token.Literal
		}

		// Then, consume everything else after it until we hit another thing that looks like a
		// new property, or the end of properties altogether.
		for {
			// We're no longer processing properties if we've hit something that's not a tab at the
			// start of a new line.
			if !p.skip(TokenTypeWhiteSpace, "\t") {
				stop = true
				break
			}

			// If we don't get a second tab, we've hit a new property. So, we need to bail.
			if p.token.Type != TokenTypeWhiteSpace || p.token.Literal != "\t" {
				break
			}

			// Skip past the "value"
			for !p.skip(TokenTypeLineTerminator) {
				p.scan()
			}
		}
	}

	output.Properties[strings.TrimSpace(name)] = strings.TrimSpace(value)

	return stop, nil
}

func (p *Parser) parseModes(output *Output) error {
	p.skipWS = true

	// Sometimes we don't have modes to parse.
	if !p.skip(TokenTypeWhiteSpace, " ") {
		return nil
	}

	for p.next(TokenTypeName) {
		var mode OutputMode

		isRes, res := p.parseResolution(p.token.Literal)
		if !isRes {
			// We've probably just hit the end of resolutions at this point, and are looking at the
			// next output.
			return nil
		}

		mode.Resolution = res
		p.scan()

		for p.next(TokenTypeFloatValue) {
			var rate Rate
			var err error

			rate.Rate, err = strconv.ParseFloat(p.token.Literal, 64)
			if err != nil {
				return err
			}

			p.scan()

			if p.skip(TokenTypePunctuator, "*") {
				rate.IsCurrent = true
			}

			if p.skip(TokenTypePunctuator, "+") {
				rate.IsPreferred = true
			}

			mode.Rates = append(mode.Rates, rate)
		}

		p.skip(TokenTypeLineTerminator)

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

	litLen := len(literal)
	if litLen < 3 {
		return false, res
	}

	// Get the first and last runes, so we can check they're numbers.
	firstRune := rune(literal[0])
	lastRuneIdx := litLen - 1
	lastRune := rune(literal[lastRuneIdx])

	// If the last rune is 'i', make sure we have a literal long enough for that to be valid.
	if lastRune == 'i' && litLen < 4 {
		return false, res
	}

	// If the last rune is 'i', the last number should be the character before the 'i'.
	if lastRune == 'i' {
		lastRuneIdx--
		lastRune = rune(literal[lastRuneIdx])
	}

	if firstRune < '0' || firstRune > '9' {
		return false, res
	}

	if lastRune < '0' || lastRune > '9' {
		return false, res
	}

	// At this point, we know the first and last rune is a number, and we have at least 3 characters
	// including those things. So, we now need loop through and find where the 'x' is. If there is
	// no 'x', then it's not a resolution. We already know it's not at the ends at this point.

	xIdx := -1
	for i := 1; i < lastRuneIdx; i++ {
		if rune(literal[i]) == 'x' {
			xIdx = i
		}
	}

	// We didn't find an x, so it's not a resolution.
	if xIdx == -1 {
		return false, res
	}

	// Extract the x and y values from the literal string.
	xVal := literal[0:xIdx]
	yVal := literal[xIdx+1 : lastRuneIdx+1]

	xres, err := strconv.ParseUint(xVal, 10, 64)
	if err != nil {
		return false, res
	}

	yres, err := strconv.ParseUint(yVal, 10, 64)
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
		p.expectFn(TokenTypeIntValue),
		p.expectFn(TokenTypePunctuator),
		p.expectFn(TokenTypeName, "minimum"),
		p.expectFn(TokenTypeIntValue),
		p.expectFn(TokenTypeName, "x"),
		p.expectFn(TokenTypeIntValue),
		p.expectFn(TokenTypePunctuator),
		p.expectFn(TokenTypeName, "current"),
		p.expectFn(TokenTypeIntValue),
		p.expectFn(TokenTypeName, "x"),
		p.expectFn(TokenTypeIntValue),
		p.expectFn(TokenTypePunctuator),
		p.expectFn(TokenTypeName, "maximum"),
		p.expectFn(TokenTypeIntValue),
		p.expectFn(TokenTypeName, "x"),
		p.expectFn(TokenTypeIntValue),
		p.expectFn(TokenTypeLineTerminator),
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
	tok := p.token
	if tok.Type != t {
		return tok, p.unexpected(tok, t, ls...)
	}

	if len(ls) == 0 {
		p.scan()
		return tok, nil
	}

	for _, l := range ls {
		if tok.Literal != l {
			continue
		}

		p.scan()
		return tok, nil
	}

	return tok, p.unexpected(tok, t, ls...)
}

func (p *Parser) expect(t TokenType, ls ...string) error {
	if !p.next(t, ls...) {
		return p.unexpected(p.token, t, ls...)
	}

	p.scan()

	return nil
}

func (p *Parser) expectFn(t TokenType, ls ...string) func() error {
	return func() error {
		return p.expect(t, ls...)
	}
}

func (p *Parser) next(t TokenType, ls ...string) bool {
	if p.token.Type != t {
		return false
	}

	if len(ls) == 0 {
		return true
	}

	for _, l := range ls {
		if p.token.Literal == l {
			return true
		}
	}

	return false
}

func (p *Parser) skip(t TokenType, ls ...string) bool {
	_, err := p.consume(t, ls...)
	if err != nil {
		return false
	}

	return true
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

	// This is as nasty as I'm willing to make this right now. But this is the slowest function in
	// the parser by far, because of the allocations it has to do, simply because it's generating
	// this message.
	// TODO(seeruk): Revisit this, it can almost definitely be improved.
	// TODO(seeruk): Don't call unexpected when it's not absolutely necessary. We can not pass
	// around errors if we don't need to (i.e. if we want to consume without caring about the error,
	// like if we just care about whether or not we did consume something).
	buf := bytes.Buffer{}
	buf.WriteString("parser error: unexpected token found: ")
	buf.WriteString(token.Type.String())
	buf.WriteString(" (")
	buf.WriteString(token.Literal)
	buf.WriteString("). Wanted: ")
	buf.WriteString(t.String())
	buf.WriteString(" (")
	buf.WriteString(strings.Join(ls, "|"))
	buf.WriteString("). Line: ")
	buf.WriteString(strconv.Itoa(token.Line))
	buf.WriteString(". Column: ")
	buf.WriteString(strconv.Itoa(token.Position))

	return errors.New(btos(buf.Bytes()))
}
