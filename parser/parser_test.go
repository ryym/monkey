package parser

import (
	"fmt"
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

func checkStatementLen(t *testing.T, prg *ast.Program, expectedLen int) {
	if len(prg.Statements) != expectedLen {
		t.Fatalf(
			"program does not contain %d statements. got=%d",
			expectedLen,
			len(prg.Statements),
		)
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf(
			"integ.TokenLiteral not %d. got=%s",
			value,
			integ.TokenLiteral(),
		)
		return false
	}

	return true
}

func TestLetStatements(t *testing.T) {
	input := `let x=5; let y=10; let foobar=838383;`

	p := New(lexer.New(input))

	prg := p.ParseProgram()
	if prg == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	checkParserErrors(t, p)
	checkStatementLen(t, prg, 3)

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

func TestReturnStatement(t *testing.T) {
	input := `return 5; return 10; return 993322;`

	p := New(lexer.New(input))
	prg := p.ParseProgram()
	checkParserErrors(t, p)
	checkStatementLen(t, prg, 3)

	for _, stmt := range prg.Statements {
		_, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stm not *ast.ReturnStatement, got=%T", stmt)
			continue
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	p := New(lexer.New(input))
	prg := p.ParseProgram()
	checkParserErrors(t, p)
	checkStatementLen(t, prg, 1)

	stmt, ok := prg.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"prg.Statements[0] is not ast.Expression statement. got=%T",
			prg.Statements[0],
		)
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("expression not *ast.Identifier. got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not foobar. got=%s", ident.Value)
	}
}

func TestIntegerExpression(t *testing.T) {
	input := "5;"

	p := New(lexer.New(input))
	prg := p.ParseProgram()
	checkParserErrors(t, p)
	checkStatementLen(t, prg, 1)

	stmt, ok := prg.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"prg.Statements[0] is not ast.Expression statement. got=%T",
			prg.Statements[0],
		)
	}

	il, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expression not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}

	if il.Value != 5 {
		t.Errorf("il.Value not 5. got=%d", il.Value)
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		p := New(lexer.New(tt.input))
		prg := p.ParseProgram()
		checkParserErrors(t, p)
		checkStatementLen(t, prg, 1)

		stmt, ok := prg.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"prg.Statements[0] is not ast.Expression statement. got=%T",
				prg.Statements[0],
			)
		}

		bo, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("expression not *ast.Boolean. got=%T", stmt.Expression)
		}

		if bo.Value != tt.expected {
			t.Errorf("bo.Value not %v. got=%v", tt.expected, bo.Value)
		}
	}
}

func TestPrefixExpressions(t *testing.T) {
	infixTests := []struct {
		input    string
		operator string
		intVal   int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}
	for _, tt := range infixTests {
		p := New(lexer.New(tt.input))
		prg := p.ParseProgram()
		checkParserErrors(t, p)
		checkStatementLen(t, prg, 1)

		stmt, ok := prg.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"prg.Statements[0] is not ast.Expression statement. got=%T",
				prg.Statements[0],
			)
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}
		if !testIntegerLiteral(t, exp.Right, tt.intVal) {
			return
		}
	}
}

func TestInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input    string
		leftVal  int64
		operator string
		rightVal int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}
	for _, tt := range infixTests {
		p := New(lexer.New(tt.input))
		prg := p.ParseProgram()
		checkParserErrors(t, p)
		checkStatementLen(t, prg, 1)

		stmt, ok := prg.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"prg.Statements[0] is not ast.Expression statement. got=%T",
				prg.Statements[0],
			)
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}

		if !testIntegerLiteral(t, exp.Left, tt.leftVal) {
			return
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}
		if !testIntegerLiteral(t, exp.Right, tt.rightVal) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}

	for _, tt := range tests {
		p := New(lexer.New(tt.input))
		prg := p.ParseProgram()
		checkParserErrors(t, p)

		actual := prg.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}
