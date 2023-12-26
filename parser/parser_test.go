package parser

import (
	"testing"

	"github.com/andrewvota/interpreter/ast"
	"github.com/andrewvota/interpreter/lexer"
)

func TestLetStatements(t *testing.T) {
	// input := `
	// let x = 5;
	// let y = 10;
	// let foobar = 838383;
	// `

	input := `
	let x 5;
	let = 10;
	let 838383;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p) // Check for parser errors
	if program == nil {
		t.Fatal("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}
	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) { // Check for parser errors
	errors := p.Errors()  // Get the errors
	if len(errors) == 0 { // If there are no errors
		return // Return
	}
	t.Errorf("parser has %d errors", len(errors)) // Print the number of errors
	for _, msg := range errors {                  // Loop through the errors
		t.Errorf("parser error: %q", msg) // Print each error
	}
	t.FailNow() // Fail the test
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}
