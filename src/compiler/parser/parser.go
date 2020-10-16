package parser

import . "github.com/tdkr/go-luavm/src/compiler/ast"
import . "github.com/tdkr/go-luavm/src/compiler/lexer"

/* recursive descent parser */

func Parse(chunk, chunkName string) *Block {
	lexer := NewLexer(chunk, chunkName)
	block := parseBlock(lexer)
	lexer.NextTokenOfKind(TOKEN_EOF)
	return block
}
