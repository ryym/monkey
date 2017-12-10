package ast

import (
	"testing"

	tk "github.com/ryym/monkey/token"
)

func TestString(t *testing.T) {
	prg := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: tk.Token{tk.LET, "let"},
				Name: &Identifier{
					Token: tk.Token{tk.IDENT, "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: tk.Token{tk.IDENT, "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if prg.String() != "let myVar = anotherVar;" {
		t.Errorf("prg.String() wrong. got=%q", prg.String())
	}
}
