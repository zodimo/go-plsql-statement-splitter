// Package parser provides a parser for Oracle 11g/12c PL/SQL using ANTLR4 grammar.
//
// This package allows you to parse and validate PL/SQL code. It is based on the ANTLR4 grammar for PL/SQL and can be used to analyze, validate, or transform PL/SQL scripts programmatically.
//
// # Usage Example
//
// The following example demonstrates how to use the parser to parse a PL/SQL file and check for syntax errors:
//
//	import (
//	    "github.com/antlr4-go/antlr/v4"
//	    plsqlparser "github.com/zodimo/plsql-parser"
//	)
//
//	input, err := antlr.NewFileStream("example.sql")
//	if err != nil {
//	    // handle error
//	}
//	lexer := plsqlparser.NewPlSqlLexer(input)
//	stream := antlr.NewCommonTokenStream(lexer, 0)
//	parser := plsqlparser.NewPlSqlParser(stream)
//	parser.SetVersion12(true) // or false for 11g
//	parser.BuildParseTrees = true
//	tree := parser.Sql_script()
//	// tree is the root of the parse tree; you can now walk or inspect it
//
// To check for syntax errors, you can implement and attach a custom error listener to the lexer and parser.
// See parser_test.go for a complete example including error handling.
//
// The grammar supports a wide range of PL/SQL constructs, including anonymous blocks, procedures, functions, packages, triggers, and DDL/DML statements.
//
// For more details, see the README.md and the ANTLR4 grammar files.
package gen

import (
	"github.com/antlr4-go/antlr/v4"
)

// PlSqlParserBase implementation.
type PlSqlParserBase struct {
	*antlr.BaseParser
	_isVersion12 bool
	_isVersion10 bool
}

func (p *PlSqlParserBase) IsTableAlias() bool {
	return p.GetCurrentToken().GetTokenType() != PlSqlLexerJOIN
}

func (p *PlSqlParserBase) isVersion12() bool {
	return p._isVersion12
}

func (p *PlSqlParserBase) SetVersion12(value bool) {
	p._isVersion12 = value
}

func (p *PlSqlParserBase) isVersion10() bool {
	return p._isVersion10
}

func (p *PlSqlParserBase) SetVersion10(value bool) {
	p._isVersion10 = value
}
