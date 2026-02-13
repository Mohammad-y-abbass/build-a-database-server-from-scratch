package parser

import (
	"fmt"
	"strings"

	"github.com/Mohammad-y-abbass/moDB/internal/ast"
	"github.com/Mohammad-y-abbass/moDB/internal/lexer"
)

type Parser struct {
	l            *lexer.Lexer
	currentToken lexer.Token
	peekToken    lexer.Token
	errors       []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}
	// Read two tokens to fill currentToken and peekToken
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, msg)
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{Statements: []ast.Statement{}}

	for p.currentToken.Type != lexer.EOF_TOKEN {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case lexer.SELECT_TOKEN:
		return p.parseSelectStatement()
	case lexer.ILLEGAL:
		p.addError(fmt.Sprintf("Illegal character '%s' at line %d, column %d",
			p.currentToken.Value, p.currentToken.Line, p.currentToken.Col))
		return nil
	case lexer.EOF_TOKEN:
		return nil
	default:
		p.addError(fmt.Sprintf("Unexpected token '%s' at line %d, column %d. Expected a statement (e.g., SELECT)",
			p.currentToken.Value, p.currentToken.Line, p.currentToken.Col))
		return nil
	}
}

func (p *Parser) parseSelectStatement() *ast.SelectStatement {
	stmt := &ast.SelectStatement{Token: p.currentToken}

	p.nextToken()

	// Check for columns or asterisk
	switch p.currentToken.Type {
	case lexer.ASTERISK:
		stmt.Columns = []string{"*"}
		p.nextToken()
	case lexer.IDENTIFIER:
		stmt.Columns = p.parseColumns()
	default:
		p.addError(fmt.Sprintf("Expected column name or '*' after SELECT at line %d, column %d, but got '%s'",
			p.currentToken.Line, p.currentToken.Col, p.currentToken.Value))
		return nil
	}

	// Expect FROM keyword
	if p.currentToken.Type != lexer.FROM_TOKEN {
		p.addError(fmt.Sprintf("Expected FROM keyword at line %d, column %d, but got '%s'",
			p.currentToken.Line, p.currentToken.Col, p.currentToken.Value))
		return nil
	}

	p.nextToken()

	// Expect table name
	if p.currentToken.Type != lexer.IDENTIFIER {
		p.addError(fmt.Sprintf("Expected table name after FROM at line %d, column %d, but got '%s'",
			p.currentToken.Line, p.currentToken.Col, p.currentToken.Value))
		return nil
	}

	stmt.Table = p.currentToken.Value

	return stmt
}

func (p *Parser) parseColumns() []string {
	var columns []string

	if p.currentToken.Type == lexer.IDENTIFIER {
		columns = append(columns, p.currentToken.Value)
	}

	for p.peekToken.Type == lexer.COMMA {
		p.nextToken() // Move to comma
		p.nextToken() // Move to next identifier

		if p.currentToken.Type != lexer.IDENTIFIER {
			p.addError(fmt.Sprintf("Expected column name after comma at line %d, column %d, but got '%s'",
				p.currentToken.Line, p.currentToken.Col, p.currentToken.Value))
			break
		}

		columns = append(columns, p.currentToken.Value)
	}

	p.nextToken()
	return columns
}

// GetErrorMessage returns the first parsing error if any
func (p *Parser) GetErrorMessage() string {
	if len(p.errors) == 0 {
		return ""
	}

	return fmt.Sprintf("Parsing error: %s", p.errors[0])
}

// FormatAST returns a formatted tree representation of the AST
func (p *Parser) FormatAST(program *ast.Program) string {
	if program == nil || len(program.Statements) == 0 {
		return "Program {\n  Statements: []\n}"
	}

	var builder strings.Builder
	builder.WriteString("Program {\n")
	builder.WriteString("  Statements: [\n")

	for i, stmt := range program.Statements {
		builder.WriteString(p.formatStatement(stmt, 4))
		if i < len(program.Statements)-1 {
			builder.WriteString(",\n")
		} else {
			builder.WriteString("\n")
		}
	}

	builder.WriteString("  ]\n")
	builder.WriteString("}")

	return builder.String()
}

func (p *Parser) formatStatement(stmt ast.Statement, indent int) string {
	indentStr := strings.Repeat(" ", indent)

	switch s := stmt.(type) {
	case *ast.SelectStatement:
		var builder strings.Builder
		builder.WriteString(indentStr + "SelectStatement {\n")
		builder.WriteString(indentStr + "  Token: " + s.Token.Value + ",\n")
		builder.WriteString(indentStr + "  Columns: [")

		if len(s.Columns) > 0 {
			builder.WriteString("\n")
			for i, col := range s.Columns {
				builder.WriteString(indentStr + "    \"" + col + "\"")
				if i < len(s.Columns)-1 {
					builder.WriteString(",\n")
				} else {
					builder.WriteString("\n")
				}
			}
			builder.WriteString(indentStr + "  ],\n")
		} else {
			builder.WriteString("],\n")
		}

		builder.WriteString(indentStr + "  Table: \"" + s.Table + "\"\n")
		builder.WriteString(indentStr + "}")

		return builder.String()
	default:
		return indentStr + "UnknownStatement {}"
	}
}
