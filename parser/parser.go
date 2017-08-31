package parser

import (
	"fmt"
	"strconv"

	"github.com/mdaisuke/monk/ast"
	"github.com/mdaisuke/monk/lexer"
	"github.com/mdaisuke/monk/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

type (
	nud func() ast.Exp
	led func(ast.Exp) ast.Exp
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	nuds map[token.TokenType]nud
	leds map[token.TokenType]led
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.nuds = make(map[token.TokenType]nud)
	p.registerNud(token.IDENT, p.parseIdentifier)
	p.registerNud(token.INT, p.parseIntegerLiteral)
	p.registerNud(token.BANG, p.parsePrefixExp)
	p.registerNud(token.MINUS, p.parsePrefixExp)
	p.registerNud(token.TRUE, p.parseBoolean)
	p.registerNud(token.FALSE, p.parseBoolean)
	p.registerNud(token.LPAREN, p.parseGroupedExp)
	p.leds = make(map[token.TokenType]led)
	p.registerLed(token.PLUS, p.parseInfixExp)
	p.registerLed(token.MINUS, p.parseInfixExp)
	p.registerLed(token.SLASH, p.parseInfixExp)
	p.registerLed(token.ASTERISK, p.parseInfixExp)
	p.registerLed(token.EQ, p.parseInfixExp)
	p.registerLed(token.NOT_EQ, p.parseInfixExp)
	p.registerLed(token.LT, p.parseInfixExp)
	p.registerLed(token.GT, p.parseInfixExp)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got=%s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Stmts = []ast.Stmt{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStmt()
		if stmt != nil {
			program.Stmts = append(program.Stmts, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStmt() ast.Stmt {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStmt()
	case token.RETURN:
		return p.parseReturnStmt()
	default:
		return p.parseExpStmt()
	}
}

func (p *Parser) parseLetStmt() *ast.LetStmt {
	stmt := &ast.LetStmt{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) parseReturnStmt() *ast.ReturnStmt {
	stmt := &ast.ReturnStmt{Token: p.curToken}

	p.nextToken()

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) registerNud(tokenType token.TokenType, fn nud) {
	p.nuds[tokenType] = fn
}

func (p *Parser) registerLed(tokenType token.TokenType, fn led) {
	p.leds[tokenType] = fn
}

func (p *Parser) parseExpStmt() *ast.ExpStmt {
	stmt := &ast.ExpStmt{Token: p.curToken}

	stmt.Exp = p.parseExp(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExp(precedence int) ast.Exp {
	prefix := p.nuds[p.curToken.Type]
	if prefix == nil {
		p.noNud(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.leds[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Exp {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Exp {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) noNud(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parsePrefixExp() ast.Exp {
	exp := &ast.PrefixExp{
		Token: p.curToken,
		Op:    p.curToken.Literal,
	}

	p.nextToken()

	exp.Right = p.parseExp(PREFIX)

	return exp
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parseInfixExp(left ast.Exp) ast.Exp {
	exp := &ast.InfixExp{
		Token: p.curToken,
		Op:    p.curToken.Literal,
		Left:  left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	exp.Right = p.parseExp(precedence)

	return exp
}

func (p *Parser) parseBoolean() ast.Exp {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExp() ast.Exp {
	p.nextToken()

	exp := p.parseExp(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}
