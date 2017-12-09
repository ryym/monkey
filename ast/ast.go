package ast

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode() // dummy method
}

type Expression interface {
	Node
	expressionNode() // dummy method
}

// Program is a series of statements.
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}
