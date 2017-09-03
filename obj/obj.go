package obj

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/mdaisuke/monk/ast"
)

type BuiltinFunction func(args ...Obj) Obj

type ObjType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
)

type Obj interface {
	Type() ObjType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjType   { return INTEGER_OBJ }
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjType   { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

type Null struct{}

func (n *Null) Type() ObjType   { return NULL_OBJ }
func (n *Null) Inspect() string { return "null" }

type ReturnValue struct {
	Value Obj
}

func (rv *ReturnValue) Type() ObjType   { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

type Error struct {
	Message string
}

func (e *Error) Type() ObjType   { return ERROR_OBJ }
func (e *Error) Inspect() string { return "ERROR: " + e.Message }

type Function struct {
	Params []*ast.Identifier
	Body   *ast.BlockStmt
	Env    *Env
}

func (f *Function) Type() ObjType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Params {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type String struct {
	Value string
}

func (s *String) Type() ObjType   { return STRING_OBJ }
func (s *String) Inspect() string { return s.Value }

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjType   { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string { return "builtin function" }
