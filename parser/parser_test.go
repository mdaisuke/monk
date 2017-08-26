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
	checkParserErrors(t, p)

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
		t.Errorf("s is not ast.LetStmt. got=%T", s)
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

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestReturnStmts(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 993322;
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Stmts) != 3 {
		t.Fatalf("len(program.Stmts) is not 3. got=%d",
			len(program.Stmts))
	}

	for _, stmt := range program.Stmts {
		returnStmt, ok := stmt.(*ast.ReturnStmt)
		if !ok {
			t.Errorf("stmt is not ast.ReturnStmt. got=%T",
				stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral is not 'return'. got=%q",
				returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExp(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	parsed := New(l)
	p := parsed.ParseProgram()
	checkParserErrors(t, parsed)

	if len(p.Stmts) != 1 {
		t.Fatalf("len(p.Stmts) is not 1. got=%d",
			len(p.Stmts))
	}

	stmt, ok := p.Stmts[0].(*ast.ExpStmt)
	if !ok {
		t.Fatalf("p.Stmts[0] is not ast.ExpStmt. got=%T",
			p.Stmts[0])
	}

	ident, ok := stmt.Exp.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp is not ast.Identifier. got=%T", stmt.Exp)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value is not foobar. got=%s", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral is not foobar. got=%s", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExp(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Stmts) != 1 {
		t.Fatalf("len(program.Stmts) is not 1. got=%d",
			len(program.Stmts))
	}
	stmt, ok := program.Stmts[0].(*ast.ExpStmt)
	if !ok {
		t.Fatalf("program.Stmts[0] is not ast.ExpStmt. got=%T",
			program.Stmts[0])
	}

	literal, ok := stmt.Exp.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("")
	}
}
