package parser

import (
	"testing"

	"github.com/mdaisuke/monk/ast"
	"github.com/mdaisuke/monk/lexer"
)

func TestLetStmts(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	if program == nil {
		t.Fatalf("ParseProgram returned nil")
	}

	if len(program.Stmts) != 3 {
		t.Fatalf("len(program.Stmts) is not 3. got=%d",
			len(program.Stmts))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Stmts[i]
		if !testLetStmt(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStmt(t *testing.T, s ast.Stmt, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral is not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStmt)
	if !ok {
		t.Errorf("s is not ast.LetStmt. got=%T",
			s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value is not '%s'. got=%s",
			name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral is not '%s'. got=%s",
			name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}
