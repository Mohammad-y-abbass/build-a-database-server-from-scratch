package parser

import (
	"fmt"
	"strconv"
	"strings"
)

type Parser struct {
	input  string
	output []any
	cursor int
}

func New(input string) *Parser {
	return &Parser{
		input:  input,
		cursor: 0,
	}
}

func (p *Parser) Next() {
	p.cursor++
}

func (p *Parser) eatWhitespace() {
	for p.cursor < len(p.input) && p.input[p.cursor] == ' ' {
		p.Next()
	}
}

func (p *Parser) expect(expected string) error {
	p.eatWhitespace()

	// 1. Boundary check
	if p.cursor+len(expected) > len(p.input) {
		return fmt.Errorf("unexpected end of input: expected %q", expected)
	}

	// 2. Content check
	if p.input[p.cursor:p.cursor+len(expected)] != expected {
		return fmt.Errorf("expected %q at position %d", expected, p.cursor)
	}

	p.cursor += len(expected)

	p.eatWhitespace()

	return nil
}
func (p *Parser) parseInt() (int, error) {
	p.eatWhitespace()
	sign := 1
	if p.cursor < len(p.input) && p.input[p.cursor] == '-' {
		sign = -1
		p.Next()
	}

	start := p.cursor
	for p.cursor < len(p.input) && p.input[p.cursor] >= '0' && p.input[p.cursor] <= '9' {
		p.Next()
	}

	if p.cursor == start {
		return 0, fmt.Errorf("expected integer")
	}

	val, _ := strconv.Atoi(p.input[start:p.cursor])
	if val == 0 {
		sign = 1
	} // -0 becomes 0 per your test

	result := val * sign
	p.output = append(p.output, result)
	return result, nil
}

func (p *Parser) parseBool() (bool, error) {
	p.eatWhitespace()
	remaining := p.input[p.cursor:]

	var val bool
	if strings.HasPrefix(remaining, "TRUE") {
		val = true
		p.cursor += 4
	} else if strings.HasPrefix(remaining, "FALSE") {
		val = false
		p.cursor += 5
	} else {
		return false, fmt.Errorf("expected boolean at %d", p.cursor)
	}

	p.output = append(p.output, val)
	return val, nil
}

func (p *Parser) parseExpression() (any, error) {
	if p.cursor >= len(p.input) {
		return nil, fmt.Errorf("unexpected end of input")
	}

	char := p.input[p.cursor]

	if char >= '0' && char <= '9' {
		return p.parseInt()
	} else {
		return p.parseBool()
	}
}

func (p *Parser) parseExpressionList() error {
	p.eatWhitespace()
	// Check if we hit the semicolon immediately
	if p.cursor < len(p.input) && p.input[p.cursor] == ';' {
		return nil
	}

	for {
		_, err := p.parseExpression()
		if err != nil {
			return err
		}

		p.eatWhitespace()
		if p.cursor >= len(p.input) || p.input[p.cursor] != ',' {
			break
		}
		p.Next() // consume ','
	}
	return nil
}

func (p *Parser) parseSelectStatement() {
	p.expect("SELECT")

	p.parseExpressionList()

	p.expect(";")
}
func (p *Parser) Parse() error {
	p.parseSelectStatement()

	p.eatWhitespace()
	if p.cursor < len(p.input) {
		return fmt.Errorf("parsing_error: unexpected characters after statement")
	}
	return nil
}

func (p *Parser) GetOutput() []any {
	return p.output
}
