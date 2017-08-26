package ast

import (
	"testing"

	"github.com/mdaisuke/monk/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Stmts: []Stmt{
			&LetStmt{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anothoerVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("program.String() is not 'let myVar = anotherVar'. got=%q", program.String())
	}
}
