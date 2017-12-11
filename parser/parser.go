package parser

import (
	"fmt"
	"strconv"

	"github.com/ryym/monkey/ast"
	lx "github.com/ryym/monkey/lexer"
	tk "github.com/ryym/monkey/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

type prefixParseFn func() ast.Expression
type infixParseFn func(ast.Expression) ast.Expression

type Parser struct {
	l         *lx.Lexer
	curToken  tk.Token
	peekToken tk.Token
	errors    []string

	prefixParseFns map[tk.TokenType]prefixParseFn
	infixParseFns  map[tk.TokenType]infixParseFn
}

func New(l *lx.Lexer) *Parser {
	p := &Parser{l: l}

	p.prefixParseFns = make(map[tk.TokenType]prefixParseFn)
	p.registerPrefix(tk.IDENT, p.parseIdentifier)
	p.registerPrefix(tk.INT, p.parseIntegerLiteral)
	p.registerPrefix(tk.BANG, p.parsePrefixExpression)
	p.registerPrefix(tk.MINUS, p.parsePrefixExpression)

	// Set curToken and peekToken
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) registerPrefix(tt tk.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tt] = fn
}
func (p *Parser) registerInfix(tt tk.TokenType, fn infixParseFn) {
	p.infixParseFns[tt] = fn
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t tk.TokenType) {
	msg := fmt.Sprintf(
		"expected next token to be %s, got %s instead",
		t,
		p.peekToken.Type,
	)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	prg := &ast.Program{}
	prg.Statements = []ast.Statement{}

	for p.curToken.Type != tk.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			prg.Statements = append(prg.Statements, stmt)
		}
		p.nextToken()
	}

	return prg
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	exp.Right = p.parseExpression(PREFIX)
	return exp
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case tk.LET:
		return p.parseLetStatement()
	case tk.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(tk.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if !p.expectPeek(tk.ASSIGN) {
		return nil
	}

	// TODO: Parse expressions.

	for !p.curTokenIs(tk.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	// TODO: Parse expressions.

	for !p.curTokenIs(tk.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	// In expression, semicolon is optional.
	if p.peekTokenIs(tk.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpression(_precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	return prefix()
}

func (p *Parser) noPrefixParseFnError(t tk.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) curTokenIs(t tk.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t tk.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t tk.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}
