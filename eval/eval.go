package eval

import (
	"github.com/mdaisuke/monk/ast"
	"github.com/mdaisuke/monk/obj"
)

var (
	NULL  = &obj.Null{}
	TRUE  = &obj.Boolean{Value: true}
	FALSE = &obj.Boolean{Value: false}
)

func Eval(node ast.Node) obj.Obj {
	switch node := node.(type) {

	case *ast.Program:
		return evalStmts(node.Stmts)

	case *ast.ExpStmt:
		return Eval(node.Exp)

	case *ast.PrefixExp:
		right := Eval(node.Right)
		return evalPrefixExp(node.Op, right)
	case *ast.InfixExp:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExp(node.Op, left, right)

	case *ast.IntegerLiteral:
		return &obj.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObj(node.Value)
	}

	return nil
}

func evalStmts(stmts []ast.Stmt) obj.Obj {
	var result obj.Obj

	for _, stmt := range stmts {
		result = Eval(stmt)
	}

	return result
}

func nativeBoolToBooleanObj(input bool) obj.Obj {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExp(op string, right obj.Obj) obj.Obj {
	switch op {
	case "!":
		return evalBangOpExp(right)
	case "-":
		return evalMinusPrefixOpExp(right)
	default:
		return NULL
	}
}

func evalBangOpExp(right obj.Obj) obj.Obj {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOpExp(right obj.Obj) obj.Obj {
	if right.Type() != obj.INTEGER_OBJ {
		return NULL
	}

	value := right.(*obj.Integer).Value
	return &obj.Integer{Value: -value}
}

func evalInfixExp(
	op string,
	left, right obj.Obj,
) obj.Obj {
	switch {
	case left.Type() == obj.INTEGER_OBJ && right.Type() == obj.INTEGER_OBJ:
		return evalIntegerInfixExp(op, left, right)
	case op == "==":
		return nativeBoolToBooleanObj(left == right)
	case op == "!=":
		return nativeBoolToBooleanObj(left != right)
	default:
		return NULL
	}
}

func evalIntegerInfixExp(
	op string,
	left, right obj.Obj,
) obj.Obj {
	leftVal := left.(*obj.Integer).Value
	rightVal := right.(*obj.Integer).Value

	switch op {
	case "+":
		return &obj.Integer{Value: leftVal + rightVal}
	case "-":
		return &obj.Integer{Value: leftVal - rightVal}
	case "*":
		return &obj.Integer{Value: leftVal * rightVal}
	case "/":
		return &obj.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObj(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObj(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObj(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObj(leftVal != rightVal)
	default:
		return NULL
	}
}
