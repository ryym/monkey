package parser

import (
	"github.com/ryym/monkey/ast"
	lx "github.com/ryym/monkey/lexer"
	tk "github.com/ryym/monkey/token"
)

type Parser struct {
	l         *lx.Lexer
	curToken  tk.Token
	peekToken tk.Token
}

func New(l *lx.Lexer) *Parser {
	p := &Parser{l: l}

	// Set curToken and peekToken
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	return nil // TODO
}
