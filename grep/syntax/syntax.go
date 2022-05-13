package syntax

import (
	"github.com/vela-security/vela-public/grep/syntax/ast"
	"github.com/vela-security/vela-public/grep/syntax/lexer"
)

func Parse(s string) (*ast.Node, error) {
	return ast.Parse(lexer.NewLexer(s))
}

func Special(b byte) bool {
	return lexer.Special(b)
}
