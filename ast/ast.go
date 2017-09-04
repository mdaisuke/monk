package ast

import (
	"bytes"
	"strings"

	"github.com/mdaisuke/monk/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Stmt interface {
	Node
	stmtNode()
}

type Exp interface {
	Node
	expNode()
}

type Program struct {
	Stmts []Stmt
}

func (p *Program) TokenLiteral() string {
	if len(p.Stmts) > 0 {
		return p.Stmts[0].TokenLiteral()
	} else {
		return ""
	}
}
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Stmts {
		out.WriteString(s.String())
	}

	return out.String()
}

type LetStmt struct {
	Token token.Token
	Name  *Identifier
	Value Exp
}

func (ls *LetStmt) stmtNode()            {}
func (ls *LetStmt) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStmt) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expNode()             {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type ReturnStmt struct {
	Token       token.Token
	ReturnValue Exp
}

func (rs *ReturnStmt) stmtNode()            {}
func (rs *ReturnStmt) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStmt) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

type ExpStmt struct {
	Token token.Token
	Exp   Exp
}

func (es *ExpStmt) stmtNode()            {}
func (es *ExpStmt) TokenLiteral() string { return es.Token.Literal }
func (es *ExpStmt) String() string {
	if es.Exp != nil {
		return es.Exp.String()
	}
	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expNode()             {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type PrefixExp struct {
	Token token.Token
	Op    string
	Right Exp
}

func (pe *PrefixExp) expNode()             {}
func (pe *PrefixExp) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExp) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Op)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExp struct {
	Token token.Token
	Left  Exp
	Op    string
	Right Exp
}

func (ie *InfixExp) expNode()             {}
func (ie *InfixExp) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExp) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Op + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expNode()             {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

type IfExp struct {
	Token  token.Token
	Cond   Exp
	Conseq *BlockStmt
	Alt    *BlockStmt
}

func (ie *IfExp) expNode()             {}
func (ie *IfExp) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExp) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Cond.String())
	out.WriteString(" ")
	out.WriteString(ie.Conseq.String())

	if ie.Alt != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alt.String())
	}

	return out.String()
}

type BlockStmt struct {
	Token token.Token
	Stmts []Stmt
}

func (bs *BlockStmt) stmtNode()            {}
func (bs *BlockStmt) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStmt) String() string {
	var out bytes.Buffer

	for _, s := range bs.Stmts {
		out.WriteString(s.String())
	}

	return out.String()
}

type FunctionLiteral struct {
	Token  token.Token
	Params []*Identifier
	Body   *BlockStmt
}

func (fl *FunctionLiteral) expNode()             {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, param := range fl.Params {
		params = append(params, param.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(fl.Body.String())

	return out.String()
}

type CallExp struct {
	Token    token.Token
	Function Exp
	Args     []Exp
}

func (ce *CallExp) expNode()             {}
func (ce *CallExp) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExp) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Args {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expNode()             {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

type ArrayLiteral struct {
	Token token.Token
	Elems []Exp
}

func (al *ArrayLiteral) expNode()             {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elems := []string{}
	for _, el := range al.Elems {
		elems = append(elems, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elems, ", "))
	out.WriteString("]")

	return out.String()
}

type IndexExp struct {
	Token token.Token
	Left  Exp
	Index Exp
}

func (ie *IndexExp) expNode()             {}
func (ie *IndexExp) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExp) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}
