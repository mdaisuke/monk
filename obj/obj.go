package obj

import (
	"fmt"
)

type ObjType string

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL_OBJ    = "NULL"
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
