package parser

import (
	"testing"

	"github.com/ryym/monkey/ast"
	"github.com/ryym/monkey/lexer"
)

func checkParserErrors(t *testing.T, p *Parser) {
	errs := p.Errors()
	if len(errs) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errs))
	for _, msg := range errs {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestLetStatements(t *testing.T) {
	input := `let x=5; let y=10; let foobar=838383;`

	p := New(lexer.New(input))

	prg := p.ParseProgram()
	if prg == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	checkParserErrors(t, p)

	if len(prg.Statements) != 3 {
		t.Fatalf(
			"Statements length is not 3. got=%d",
			len(prg.Statements),
		)
	}

	tests := []struct {
		wantIdent string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}
	for i, tt := range tests {
		stmt := prg.Statements[i]
		if !testLetStatement(t, stmt, tt.wantIdent) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	let, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if let.Name.Value != name {
		t.Errorf("let.Name.Value not '%s'. got=%s", name, let.Name.Value)
		return false
	}

	if let.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, let.Name)
		return false
	}

	return true
}
