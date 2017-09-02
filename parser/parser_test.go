package parser

import (
	"fmt"
	"testing"

	"github.com/mdaisuke/monk/ast"
	"github.com/mdaisuke/monk/lexer"
)

func TestLetStmts(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Stmts) != 1 {
			t.Fatalf("len(program.Stmts) is not 1. got=%d",
				len(program.Stmts))
		}

		stmt := program.Stmts[0]
		if !testLetStmt(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStmt).Value
		if !testLiteralExp(t, val, tt.expectedValue) {
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
			t.Errorf("stmt is not ast.ReturnStmt. got=%T", stmt)
			continue
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral() is not 'return'. got=%q",
				returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExp(t *testing.T) {
	input := "foobar;"

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

	ident, ok := stmt.Exp.(*ast.Identifier)
	if !ok {
		t.Fatalf("stmt.Exp is not ast.Identifier. got=%T",
			stmt.Exp)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value is not 'foobar'. got=%s",
			ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral() is not 'foobar'. got=%s",
			ident.TokenLiteral())
	}
}

func TestIntegerLiteralExp(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Stmts) != 1 {
		t.Fatalf("len(program.Stmts) is not 1. got=%T",
			len(program.Stmts))
	}
	stmt, ok := program.Stmts[0].(*ast.ExpStmt)
	if !ok {
		t.Fatalf("program.Stmts[0] is not ast.ExpStmt. got=%T",
			program.Stmts[0])
	}

	literal, ok := stmt.Exp.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp is not ast.IntegerLiteral. got=%T",
			stmt.Exp)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value is not 5. got=%d", literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral() is not 5. got=%s", literal.TokenLiteral())
	}
}

func TestParsingPrefixExps(t *testing.T) {
	tests := []struct {
		input        string
		op           string
		integerValue interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
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

		exp, ok := stmt.Exp.(*ast.PrefixExp)
		if !ok {
			t.Fatalf("stmt.Exp is not ast.PrefixExp. got=%T", stmt.Exp)
		}
		if exp.Op != tt.op {
			t.Fatalf("exp.Op is not '%s'. got=%s",
				tt.op, exp.Op)
		}
		if !testLiteralExp(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Exp, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il is not ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value is not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral is not %d. got=%s",
			value, integ.TokenLiteral())
		return false
	}

	return true
}

func TestParsingInfixExps(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  interface{}
		op         string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
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

		exp, ok := stmt.Exp.(*ast.InfixExp)
		if !ok {
			t.Fatalf("stmt.Exp is not ast.InfixExp. got=%T", stmt.Exp)
		}

		if exp.Op != tt.op {
			t.Fatalf("exp.Op is not '%s'. got=%s", tt.op, exp.Op)
		}

		if !testLiteralExp(t, exp.Left, tt.leftValue) {
			return
		}
		if !testLiteralExp(t, exp.Right, tt.rightValue) {
			return
		}
	}
}

func TestOpPrecedenceParsing(t *testing.T) {
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
			"5 > 4 == 3 > 4",
			"((5 > 4) == (3 > 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q",
				tt.expected, actual)
		}
	}
}

func testIdentifier(t *testing.T, exp ast.Exp, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp is not ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value is not '%s'. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral() is not '%s'. got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExp(
	t *testing.T,
	exp ast.Exp,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp is not handled. got=%T", exp)
	return false
}

func testBooleanLiteral(t *testing.T, exp ast.Exp, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp is not ast.Boolean. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value is not %t. got=%t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral is not %t. got=%s",
			value, bo.TokenLiteral())
		return false
	}

	return true
}

func testInfixExp(
	t *testing.T,
	exp ast.Exp,
	left interface{},
	op string,
	right interface{},
) bool {
	opExp, ok := exp.(*ast.InfixExp)
	if !ok {
		t.Errorf("exp is not ast.InfixExp. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExp(t, opExp.Left, left) {
		return false
	}

	if opExp.Op != op {
		t.Errorf("exp.Op is not '%s'. got=%q", op, opExp.Op)
		return false
	}

	if !testLiteralExp(t, opExp.Right, right) {
		return false
	}

	return true
}

func TestBooleanExp(t *testing.T) {
	input := "true;"

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

	b, ok := stmt.Exp.(*ast.Boolean)
	if !ok {
		t.Fatalf("stmt.Exp is not ast.Boolean. got=%T",
			stmt.Exp)
	}
	if b.Value != true {
		t.Errorf("b.Value is not 'true'. got=%t",
			b.Value)
	}
	if b.TokenLiteral() != "true" {
		t.Errorf("b.TokenLiteral() is not 'true'. got=%s",
			b.TokenLiteral())
	}
}

func TestIfExp(t *testing.T) {
	input := `if (x < y) { x }`

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
		t.Fatalf("stmts[0] is not ast.ExpStmt. got=%T",
			program.Stmts[0])
	}

	exp, ok := stmt.Exp.(*ast.IfExp)
	if !ok {
		t.Fatalf("stmt.Exp is not ast.IfExp. got=%T", stmt.Exp)
	}

	if !testInfixExp(t, exp.Cond, "x", "<", "y") {
		return
	}

	if len(exp.Conseq.Stmts) != 1 {
		t.Errorf("len(exp.Conseq.Stmts) is not 1. got=%d",
			len(exp.Conseq.Stmts))
	}

	conseq, ok := exp.Conseq.Stmts[0].(*ast.ExpStmt)
	if !ok {
		t.Fatalf("stmts[0] is not ast.ExpStmt. got=%T",
			exp.Conseq.Stmts[0])
	}

	if !testIdentifier(t, conseq.Exp, "x") {
		return
	}

	if exp.Alt != nil {
		t.Errorf("exp.Alt is not nil. got=%+v", exp.Alt)
	}
}

func TestIfElseExp(t *testing.T) {
	input := `if (x < y) { x } else { y }`

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
		t.Fatalf("stmts[0] is not ast.ExpStmt. got=%T",
			program.Stmts[0])
	}

	exp, ok := stmt.Exp.(*ast.IfExp)
	if !ok {
		t.Fatalf("stmt.Exp is not ast.IfExp. got=%T", stmt.Exp)
	}

	if !testInfixExp(t, exp.Cond, "x", "<", "y") {
		return
	}

	if len(exp.Conseq.Stmts) != 1 {
		t.Errorf("len(exp.Conseq.Stmts) is not 1. got=%d",
			len(exp.Conseq.Stmts))
	}

	conseq, ok := exp.Conseq.Stmts[0].(*ast.ExpStmt)
	if !ok {
		t.Fatalf("stmts[0] is not ast.ExpStmt. got=%T",
			exp.Conseq.Stmts[0])
	}

	if !testIdentifier(t, conseq.Exp, "x") {
		return
	}

	if len(exp.Alt.Stmts) != 1 {
		t.Errorf("len(exp.Alt.Stmts) is not 1. got=%d",
			len(exp.Alt.Stmts))
	}

	alt, ok := exp.Alt.Stmts[0].(*ast.ExpStmt)
	if !ok {
		t.Fatalf("stmts[0] is not ast.ExpStmt. got=%T",
			exp.Alt.Stmts[0])
	}

	if !testIdentifier(t, alt.Exp, "y") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

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

	function, ok := stmt.Exp.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Exp is not ast.FunctionLiteral. got=%T",
			stmt.Exp)
	}

	if len(function.Params) != 2 {
		t.Fatalf("len(function.Params) is not 2. got=%d",
			len(function.Params))
	}

	testLiteralExp(t, function.Params[0], "x")
	testLiteralExp(t, function.Params[1], "y")

	if len(function.Body.Stmts) != 1 {
		t.Fatalf("len(function.Body.Stmts) is not 1. got=%d",
			len(function.Body.Stmts))
	}

	body, ok := function.Body.Stmts[0].(*ast.ExpStmt)
	if !ok {
		t.Fatalf("function body is not ast.ExpStmt. got=%T",
			function.Body.Stmts[0])
	}

	testInfixExp(t, body.Exp, "x", "+", "y")
}

func TestFunctionParamsParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{"fn() {};", []string{}},
		{"fn(x) {};", []string{"x"}},
		{"fn(x, y, z) {};", []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Stmts[0].(*ast.ExpStmt)
		function := stmt.Exp.(*ast.FunctionLiteral)

		if len(function.Params) != len(tt.expectedParams) {
			t.Errorf("len(function.Params) is not %d. got=%d",
				len(tt.expectedParams), len(function.Params))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExp(t, function.Params[i], ident)
		}
	}

}

func TestCallExpParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

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
		t.Fatalf("stmt is not ast.ExpStmt. got=%T",
			program.Stmts[0])
	}

	exp, ok := stmt.Exp.(*ast.CallExp)
	if !ok {
		t.Fatalf("stmt.Exp is not ast.CallExp. got=%T",
			stmt.Exp)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Args) != 3 {
		t.Fatalf("len(exp.Args) is not 3. got=%d",
			len(exp.Args))
	}

	testLiteralExp(t, exp.Args[0], 1)
	testInfixExp(t, exp.Args[1], 2, "*", 3)
	testInfixExp(t, exp.Args[2], 4, "+", 5)

}
