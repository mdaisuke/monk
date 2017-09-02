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

	return Eval(program)
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
