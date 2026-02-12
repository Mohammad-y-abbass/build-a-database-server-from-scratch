package ast

import "github.com/Mohammad-y-abbass/build-a-database-server-from-scratch/internal/lexer"

type Identifier struct {
	Token lexer.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Value }
func (i *Identifier) String() string       { return i.Value }
