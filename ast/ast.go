package ast

import (
	"bytes"

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
