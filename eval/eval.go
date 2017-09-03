package eval

import (
	"fmt"

	"github.com/mdaisuke/monk/ast"
	"github.com/mdaisuke/monk/obj"
)

var (
	NULL  = &obj.Null{}
	TRUE  = &obj.Boolean{Value: true}
	FALSE = &obj.Boolean{Value: false}
)

func Eval(node ast.Node, env *obj.Env) obj.Obj {
	switch node := node.(type) {

	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpStmt:
		return Eval(node.Exp, env)
	case *ast.BlockStmt:
		return evalBlockStmt(node, env)
	case *ast.ReturnStmt:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &obj.ReturnValue{Value: val}
	case *ast.LetStmt:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	case *ast.PrefixExp:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExp(node.Op, right)
	case *ast.InfixExp:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExp(node.Op, left, right)
	case *ast.IfExp:
		return evalIfExp(node, env)

	case *ast.IntegerLiteral:
		return &obj.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObj(node.Value)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FunctionLiteral:
		params := node.Params
		body := node.Body
		return &obj.Function{Params: params, Env: env, Body: body}
	case *ast.CallExp:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExps(node.Args, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	case *ast.StringLiteral:
		return &obj.String{Value: node.Value}
	}

	return nil
}

func evalStmts(stmts []ast.Stmt, env *obj.Env) obj.Obj {
	var result obj.Obj

	for _, stmt := range stmts {
		result = Eval(stmt, env)

		if returnValue, ok := result.(*obj.ReturnValue); ok {
			return returnValue.Value
		}
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
		return newError("unknown op: %s%s", op, right.Type())
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
		return newError("unknown op: -%s",
			right.Type())
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
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",
			left.Type(), op, right.Type())
	case left.Type() == obj.STRING_OBJ && right.Type() == obj.STRING_OBJ:
		return evalStringInfixExp(op, left, right)
	default:
		return newError("unknown op: %s %s %s",
			left.Type(), op, right.Type())
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
		return newError("unknown op: %s %s %s",
			left.Type(), op, right.Type())
	}
}

func evalIfExp(ie *ast.IfExp, env *obj.Env) obj.Obj {
	cond := Eval(ie.Cond, env)
	if isError(cond) {
		return cond
	}

	if isTruthy(cond) {
		return Eval(ie.Conseq, env)
	} else if ie.Alt != nil {
		return Eval(ie.Alt, env)
	} else {
		return NULL
	}
}

func isTruthy(o obj.Obj) bool {
	switch o {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalProgram(program *ast.Program, env *obj.Env) obj.Obj {
	var result obj.Obj

	for _, stmt := range program.Stmts {
		result = Eval(stmt, env)

		switch result := result.(type) {
		case *obj.ReturnValue:
			return result.Value
		case *obj.Error:
			return result
		}
	}

	return result
}

func evalBlockStmt(block *ast.BlockStmt, env *obj.Env) obj.Obj {
	var result obj.Obj

	for _, stmt := range block.Stmts {
		result = Eval(stmt, env)

		if result != nil {
			rt := result.Type()
			if rt == obj.RETURN_VALUE_OBJ || rt == obj.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func newError(format string, a ...interface{}) *obj.Error {
	return &obj.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(o obj.Obj) bool {
	if o != nil {
		return o.Type() == obj.ERROR_OBJ
	}
	return false
}

func evalIdentifier(
	node *ast.Identifier,
	env *obj.Env,
) obj.Obj {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
}

func evalExps(
	exps []ast.Exp,
	env *obj.Env,
) []obj.Obj {
	var result []obj.Obj

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []obj.Obj{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn obj.Obj, args []obj.Obj) obj.Obj {
	switch fn := fn.(type) {
	case *obj.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *obj.Builtin:
		return fn.Fn(args...)

	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(
	fn *obj.Function,
	args []obj.Obj,
) *obj.Env {
	env := obj.NewEnclosedEnv(fn.Env)

	for i, param := range fn.Params {
		env.Set(param.Value, args[i])
	}

	return env
}

func unwrapReturnValue(o obj.Obj) obj.Obj {
	if returnValue, ok := o.(*obj.ReturnValue); ok {
		return returnValue.Value
	}

	return o
}

func evalStringInfixExp(
	op string,
	left, right obj.Obj,
) obj.Obj {
	if op != "+" {
		return newError("unknown op: %s %s %s",
			left.Type(), op, right.Type())
	}

	leftVal := left.(*obj.String).Value
	rightVal := right.(*obj.String).Value
	return &obj.String{Value: leftVal + rightVal}
}
