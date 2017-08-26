package ast

import (
	"github.com/mdaisuke/monk/token"
)

type Node interface {
	TokenLiteral() string
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

type LetStmt struct {
	Token token.Token
	Name  *Identifier
	Value Exp
}

func (ls *LetStmt) stmtNode()            {}
func (ls *LetStmt) TokenLiteral() string { return ls.Token.Literal }

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expNode()             {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
