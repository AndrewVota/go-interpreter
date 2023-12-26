package parser

import (
	"fmt"

	"github.com/andrewvota/interpreter/ast"
	"github.com/andrewvota/interpreter/lexer"
	"github.com/andrewvota/interpreter/token"
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string { // Return the errors
	return p.errors // Return the errors
}

func (p *Parser) peekError(t token.TokenType) { // Add an error to the errors slice
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type) // Create the error message
	p.errors = append(p.errors, msg)                                                        // Append the error message to the errors slice
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}              // Create the root node of our AST
	program.Statements = []ast.Statement{} // Initialize the Statements field as an empty slice

	for !p.curTokenIs(token.EOF) { // Loop until we reach the end of the input
		stmt := p.parseStatement() // Parse the current statement
		if stmt != nil {           // If the statement is not nil
			program.Statements = append(program.Statements, stmt) // Append the statement to the Statements field
		}
		p.nextToken() // Advance the tokens
	}
	return program // Return the root node of our AST
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type { // Check the type of the current token
	case token.LET: // If the current token is LET
		return p.parseLetStatement() // Parse the let statement
	default:
		return nil // Otherwise, return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken} // Create the root node of the let statement

	if !p.expectPeek(token.IDENT) { // If the next token is not an IDENT
		return nil // Return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal} // Create the identifier node

	if !p.expectPeek(token.ASSIGN) { // If the next token is not an ASSIGN
		return nil // Return nil
	}

	for !p.curTokenIs(token.SEMICOLON) { // Loop until we reach a SEMICOLON
		p.nextToken() // Advance the tokens
	}

	return stmt // Return the root node of the let statement
}

func (p *Parser) curTokenIs(t token.TokenType) bool { // Check if the current token is of a given type
	return p.curToken.Type == t // Return the result
}

func (p *Parser) peekTokenIs(t token.TokenType) bool { // Check if the next token is of a given type
	return p.peekToken.Type == t // Return the result
}

func (p *Parser) expectPeek(t token.TokenType) bool { // Check if the next token is of a given type
	if p.peekTokenIs(t) { // If the next token is of the given type
		p.nextToken() // Advance the tokens
		return true   // Return true
	} else {
		p.peekError(t)
		return false // Otherwise, return false
	}
}
