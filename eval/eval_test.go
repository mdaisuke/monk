package eval

import (
	"testing"

	"github.com/mdaisuke/monk/lexer"
	"github.com/mdaisuke/monk/obj"
	"github.com/mdaisuke/monk/parser"
)

func TestEvalIntegerExp(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObj(t, evaluated, tt.expected)
	}

}

func testEval(input string) obj.Obj {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := obj.NewEnv()

	return Eval(program, env)
}

func testIntegerObj(t *testing.T, o obj.Obj, expected int64) bool {
	result, ok := o.(*obj.Integer)
	if !ok {
		t.Errorf("obj. is not Integer. got=%T(%+v)", o, o)
		return false
	}
	if result.Value != expected {
		t.Errorf("obj is not %d. got=%d",
			expected, result.Value)
	}

	return true
}

func TestEvalBooleanExp(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObj(t, evaluated, tt.expected)
	}
}

func testBooleanObj(t *testing.T, o obj.Obj, expected bool) bool {
	result, ok := o.(*obj.Boolean)
	if !ok {
		t.Errorf("obj is not Boolean. got=%T(%+v)", o, o)
		return false
	}
	if result.Value != expected {
		t.Errorf("obj is not %t. got=%t", expected, result.Value)
		return false
	}
	return true
}

func TestBangOp(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObj(t, evaluated, tt.expected)
	}
}

func TestIfElseExps(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObj(t, evaluated, int64(integer))
		} else {
			testNullObj(t, evaluated)
		}
	}
}

func testNullObj(t *testing.T, o obj.Obj) bool {
	if o != NULL {
		t.Errorf("obj is not NULL. got=%T(%+v)", o, o)
		return false
	}
	return true
}

func TestReturnStmts(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{
			`
			if (10 > 1) {
				if (10 > 1) {
					return 10;
				}
				return 1;
			}
			`, 10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObj(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			" 5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			" 5 + true; 5",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown op: -BOOLEAN",
		},
		{
			"true + false",
			"unknown op: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown op: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown op: BOOLEAN + BOOLEAN",
		},
		{
			`
			if (10 > 1) {
				if (10 > 1) {
					return true + false;
				}
				return 1;
			}
			`,
			"unknown op: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			`"Hello" - "World"`,
			"unknown op: STRING - STRING",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*obj.Error)
		if !ok {
			t.Errorf("no error obj. got=%T(%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("error message is not '%q'. got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}

}

func TestLetStmts(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObj(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObj(t *testing.T) {
	input := "fn(x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*obj.Function)
	if !ok {
		t.Fatalf("obj is not Function. got=%T(%+v)",
			evaluated, evaluated)
	}

	if len(fn.Params) != 1 {
		t.Fatalf("len(fn.Params) is not 1. got=%d",
			len(fn.Params))
	}

	if fn.Params[0].String() != "x" {
		t.Fatalf("params[0] is not 'x'. got=%q", fn.Params[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q",
			expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5)", 5},
		{"let identity = fn(x) { return x; }; identity(5)", 5},
		{"let double = fn(x) { x * 2; }; double(5)", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x){ x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObj(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
	let newAdder = fn(x) {
		fn(y) { x + y };
	};

	let addTwo = newAdder(2);
	addTwo(2);
	`

	testIntegerObj(t, testEval(input), 4)
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*obj.String)
	if !ok {
		t.Fatalf("obj is not String. got=%T(%+v)",
			evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("string is not 'Hello World!'. got=%q", str.Value)
	}
}

func TestStringConcat(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*obj.String)
	if !ok {
		t.Fatalf("obj is not String. got=%T(%+v)",
			evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("str.Value is not 'Hello World!'. got=%q", str.Value)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` is not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObj(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*obj.Error)
			if !ok {
				t.Errorf("obj is not Error. got=%T(%+v)",
					evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q",
					expected, errObj.Message)
			}
		}
	}
}
