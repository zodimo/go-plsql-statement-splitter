package parser

import (
	"fmt"
	"sort"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/zodimo/go-plsql-statement-splitter/internal/parser/gen"
)

// SqlStatementContext is a placeholder until the ANTLR4 parser is generated
type SqlStatementContext struct {
	antlr.ParserRuleContext
}

// GetStart returns the start token
func (c *SqlStatementContext) GetStart() antlr.Token {
	return nil
}

// GetStop returns the stop token
func (c *SqlStatementContext) GetStop() antlr.Token {
	return nil
}

// GetText returns the text of the SQL statement
func (c *SqlStatementContext) GetText() string {
	return ""
}

// DmlStatement returns the DML statement context
func (c *SqlStatementContext) DmlStatement() DmlStatementContext {
	return DmlStatementContext{}
}

// DdlStatement returns the DDL statement context
func (c *SqlStatementContext) DdlStatement() DdlStatementContext {
	return DdlStatementContext{}
}

// TransactionStatement returns the transaction statement context
func (c *SqlStatementContext) TransactionStatement() TransactionStatementContext {
	return TransactionStatementContext{}
}

// DmlStatementContext is a placeholder
type DmlStatementContext struct{}

// SelectStatement returns nil
func (c DmlStatementContext) SelectStatement() interface{} {
	return nil
}

// InsertStatement returns nil
func (c DmlStatementContext) InsertStatement() interface{} {
	return nil
}

// UpdateStatement returns nil
func (c DmlStatementContext) UpdateStatement() interface{} {
	return nil
}

// DeleteStatement returns nil
func (c DmlStatementContext) DeleteStatement() interface{} {
	return nil
}

// MergeStatement returns nil
func (c DmlStatementContext) MergeStatement() interface{} {
	return nil
}

// DdlStatementContext is a placeholder
type DdlStatementContext struct{}

// CreateStatement returns nil
func (c DdlStatementContext) CreateStatement() interface{} {
	return nil
}

// AlterStatement returns nil
func (c DdlStatementContext) AlterStatement() interface{} {
	return nil
}

// DropStatement returns nil
func (c DdlStatementContext) DropStatement() interface{} {
	return nil
}

// TruncateStatement returns nil
func (c DdlStatementContext) TruncateStatement() interface{} {
	return nil
}

// GrantStatement returns nil
func (c DdlStatementContext) GrantStatement() interface{} {
	return nil
}

// RevokeStatement returns nil
func (c DdlStatementContext) RevokeStatement() interface{} {
	return nil
}

// TransactionStatementContext is a placeholder
type TransactionStatementContext struct{}

// CommitStatement returns nil
func (c TransactionStatementContext) CommitStatement() interface{} {
	return nil
}

// RollbackStatement returns nil
func (c TransactionStatementContext) RollbackStatement() interface{} {
	return nil
}

// SavePointStatement returns nil
func (c TransactionStatementContext) SavePointStatement() interface{} {
	return nil
}

// PlsqlBlockContext is a placeholder
type PlsqlBlockContext struct {
	antlr.ParserRuleContext
}

// GetStart returns the start token
func (c *PlsqlBlockContext) GetStart() antlr.Token {
	return nil
}

// GetStop returns the stop token
func (c *PlsqlBlockContext) GetStop() antlr.Token {
	return nil
}

// GetText returns the text of the PL/SQL block
func (c *PlsqlBlockContext) GetText() string {
	return ""
}

// AnonymousPlsqlBlockContext is a placeholder
type AnonymousPlsqlBlockContext struct {
	antlr.ParserRuleContext
}

// GetStart returns the start token
func (c *AnonymousPlsqlBlockContext) GetStart() antlr.Token {
	return nil
}

// GetStop returns the stop token
func (c *AnonymousPlsqlBlockContext) GetStop() antlr.Token {
	return nil
}

// GetText returns the text of the anonymous PL/SQL block
func (c *AnonymousPlsqlBlockContext) GetText() string {
	return ""
}

// TerminalNode is a placeholder
type TerminalNode interface {
	GetSymbol() antlr.Token
}

// Statement represents a PL/SQL statement with position information
type Statement struct {
	Content     string
	StartLine   int
	EndLine     int
	StartColumn int
	EndColumn   int
	Type        string
}

// SyntaxError represents a syntax error that occurred during parsing
type SyntaxError struct {
	Line      int
	Column    int
	Message   string
	TokenText string // Text of the offending token, if available
	Context   string // Surrounding context for better error reporting
}

// CustomErrorListener captures syntax errors during parsing
type CustomErrorListener struct {
	antlr.DefaultErrorListener
	Errors     []SyntaxError
	MaxErrors  int    // Maximum number of errors to capture
	SourceText string // The original source text for context
}

// NewCustomErrorListener creates a new error listener with the given max errors and source text
func NewCustomErrorListener(maxErrors int, sourceText string) *CustomErrorListener {
	if maxErrors <= 0 {
		maxErrors = 10 // Default to 10 if not specified or invalid
	}
	return &CustomErrorListener{
		Errors:     make([]SyntaxError, 0),
		MaxErrors:  maxErrors,
		SourceText: sourceText,
	}
}

// SyntaxError is called by the parser when a syntax error is encountered
func (l *CustomErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	// Only capture up to MaxErrors
	if len(l.Errors) >= l.MaxErrors {
		return
	}

	// Extract token text if available
	tokenText := ""
	if offendingSymbol != nil {
		if token, ok := offendingSymbol.(antlr.Token); ok {
			tokenText = token.GetText()
		}
	}

	// Extract context from source text if available
	context := ""
	if l.SourceText != "" {
		lines := strings.Split(l.SourceText, "\n")
		if line > 0 && line <= len(lines) {
			// Get the line with the error
			errorLine := lines[line-1]

			// Add a marker for the error position
			if column > 0 && column <= len(errorLine) {
				context = errorLine + "\n" + strings.Repeat(" ", column-1) + "^"
			} else {
				context = errorLine
			}
		}
	}

	l.Errors = append(l.Errors, SyntaxError{
		Line:      line,
		Column:    column,
		Message:   msg,
		TokenText: tokenText,
		Context:   context,
	})
}

// deduplicateStatements removes duplicate statements with the same position
func deduplicateStatements(statements []statementModel) []statementModel {
	if len(statements) == 0 {
		return statements
	}

	// Use a map to track unique statements by their position
	uniqueMap := make(map[string]statementModel)

	for _, stmt := range statements {
		// Create a key using the position information
		key := fmt.Sprintf("%d:%d:%d:%d", stmt.StartLine, stmt.StartColumn, stmt.EndLine, stmt.EndColumn)

		// Only add if it's not already in the map or if the current has a more specific type
		existing, exists := uniqueMap[key]
		if !exists || (exists && existing.Type == "UNKNOWN" && stmt.Type != "UNKNOWN") {
			uniqueMap[key] = stmt
		}
	}

	// Convert map back to slice
	result := make([]statementModel, 0, len(uniqueMap))
	for _, stmt := range uniqueMap {
		result = append(result, stmt)
	}

	// Sort by start position for consistency
	sort.Slice(result, func(i, j int) bool {
		if result[i].StartLine != result[j].StartLine {
			return result[i].StartLine < result[j].StartLine
		}
		return result[i].StartColumn < result[j].StartColumn
	})

	return result
}

func ParseString(input string) ([]Statement, []SyntaxError, error) {
	return ParseStringWithOptions(input, 1)
}

func ParseStringWithOptions(input string, maxErrors int) ([]Statement, []SyntaxError, error) {
	// Create an input stream from the input string
	inputStream := antlr.NewInputStream(input)

	// Create the lexer
	lexer := gen.NewPlSqlLexer(inputStream)

	// Create the token stream
	tokenStream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Create the parser
	parser := gen.NewPlSqlParser(tokenStream)
	parser.SetVersion12(true)

	// Create the error listener
	errorListener := NewCustomErrorListener(maxErrors, input)
	parser.RemoveErrorListeners()
	parser.AddErrorListener(errorListener)

	// Parse the input
	tree := parser.Sql_script()

	// Create the statement listener
	listener := NewStatementListener(parser, tokenStream)

	// Walk the parse tree
	antlr.ParseTreeWalkerDefault.Walk(listener, tree)

	// Deduplicate statements
	deduplicatedStatements := deduplicateStatements(listener.Statements)

	// Convert to Statement type
	statements := make([]Statement, 0, len(deduplicatedStatements))
	for _, stmt := range deduplicatedStatements {
		statements = append(statements, Statement{
			Content:     stmt.Content,
			StartLine:   stmt.StartLine,
			EndLine:     stmt.EndLine,
			StartColumn: stmt.StartColumn,
			EndColumn:   stmt.EndColumn,
			Type:        stmt.Type,
		})
	}

	// Return the statements and any errors
	return statements, errorListener.Errors, nil
}

// Internal statement model
type statementModel struct {
	Content     string
	StartLine   int
	EndLine     int
	StartColumn int
	EndColumn   int
	Type        string
}

// StatementListener listens for statements in the parse tree
type StatementListener struct {
	*gen.BasePlSqlParserListener
	Statements  []statementModel
	parser      *gen.PlSqlParser
	tokenStream *antlr.CommonTokenStream
	currentType string // Track the current statement type
}

// NewStatementListener creates a new statement listener
func NewStatementListener(parser *gen.PlSqlParser, tokenStream *antlr.CommonTokenStream) *StatementListener {
	return &StatementListener{
		BasePlSqlParserListener: &gen.BasePlSqlParserListener{},
		Statements:              make([]statementModel, 0),
		parser:                  parser,
		tokenStream:             tokenStream,
		currentType:             "",
	}
}

// EnterUnit_statement is called when entering a unit_statement rule
func (l *StatementListener) EnterUnit_statement(ctx *gen.Unit_statementContext) {
	// Get the start and stop tokens
	start := ctx.GetStart()
	stop := ctx.GetStop()

	if start == nil || stop == nil {
		return
	}

	// Get the statement text
	content := l.tokenStream.GetTextFromTokens(start, stop)

	// Get position information
	startLine := start.GetLine()
	startColumn := start.GetColumn()
	endLine := stop.GetLine()
	endColumn := stop.GetColumn() + len(stop.GetText())

	// Determine the statement type
	stmtType := getDeterminedStatementType(content)

	// Add the statement to the list
	l.Statements = append(l.Statements, statementModel{
		Content:     content,
		StartLine:   startLine,
		EndLine:     endLine,
		StartColumn: startColumn,
		EndColumn:   endColumn,
		Type:        stmtType,
	})
}

// EnterSql_statement is called when entering a sql_statement rule
func (l *StatementListener) EnterSql_statement(ctx *gen.Sql_statementContext) {
	// Get the start and stop tokens
	start := ctx.GetStart()
	stop := ctx.GetStop()

	if start == nil || stop == nil {
		return
	}

	// Get the statement text
	content := l.tokenStream.GetTextFromTokens(start, stop)

	// Get position information
	startLine := start.GetLine()
	startColumn := start.GetColumn()
	endLine := stop.GetLine()
	endColumn := stop.GetColumn() + len(stop.GetText())

	// Determine the statement type
	stmtType := getDeterminedStatementType(content)

	// Add the statement to the list
	l.Statements = append(l.Statements, statementModel{
		Content:     content,
		StartLine:   startLine,
		EndLine:     endLine,
		StartColumn: startColumn,
		EndColumn:   endColumn,
		Type:        stmtType,
	})
}

// EnterCreate_procedure_body is called when entering a create_procedure_body rule
func (l *StatementListener) EnterCreate_procedure_body(ctx *gen.Create_procedure_bodyContext) {
	l.currentType = "CREATE_PROCEDURE"
}

// EnterCreate_function_body is called when entering a create_function_body rule
func (l *StatementListener) EnterCreate_function_body(ctx *gen.Create_function_bodyContext) {
	l.currentType = "CREATE_FUNCTION"
}

// EnterCreate_package is called when entering a create_package rule
func (l *StatementListener) EnterCreate_package(ctx *gen.Create_packageContext) {
	l.currentType = "CREATE_PACKAGE"
}

// EnterCreate_package_body is called when entering a create_package_body rule
func (l *StatementListener) EnterCreate_package_body(ctx *gen.Create_package_bodyContext) {
	l.currentType = "CREATE_PACKAGE_BODY"
}

// EnterCreate_trigger is called when entering a create_trigger rule
func (l *StatementListener) EnterCreate_trigger(ctx *gen.Create_triggerContext) {
	l.currentType = "CREATE_TRIGGER"
}

// EnterCreate_type is called when entering a create_type rule
func (l *StatementListener) EnterCreate_type(ctx *gen.Create_typeContext) {
	l.currentType = "CREATE_TYPE"
}

// EnterCreate_type_body is called when the CreateTypeBody rule is encountered
// Note: Update this if the exact gen.CreateTypeBodyContext name is different
func (l *StatementListener) EnterCreate_type_body(ctx interface{}) {
	l.currentType = "CREATE_TYPE_BODY"
}

// EnterAnonymous_block is called when entering an anonymous_block rule
func (l *StatementListener) EnterAnonymous_block(ctx *gen.Anonymous_blockContext) {
	// Get the start and stop tokens
	start := ctx.GetStart()
	stop := ctx.GetStop()

	if start == nil || stop == nil {
		return
	}

	// Get the statement text
	content := l.tokenStream.GetTextFromTokens(start, stop)

	// Get position information
	startLine := start.GetLine()
	startColumn := start.GetColumn()
	endLine := stop.GetLine()
	endColumn := stop.GetColumn() + len(stop.GetText())

	// Add the statement to the list
	l.Statements = append(l.Statements, statementModel{
		Content:     content,
		StartLine:   startLine,
		EndLine:     endLine,
		StartColumn: startColumn,
		EndColumn:   endColumn,
		Type:        "PLSQL_BLOCK",
	})
}

// EnterTransaction_control_statements is called when entering a transaction_control_statements rule
func (l *StatementListener) EnterTransaction_control_statements(ctx *gen.Transaction_control_statementsContext) {
	// Get the start and stop tokens
	start := ctx.GetStart()
	stop := ctx.GetStop()

	if start == nil || stop == nil {
		return
	}

	// Get the statement text
	content := l.tokenStream.GetTextFromTokens(start, stop)

	// Get position information
	startLine := start.GetLine()
	startColumn := start.GetColumn()
	endLine := stop.GetLine()
	endColumn := stop.GetColumn() + len(stop.GetText())

	// Determine transaction statement type
	stmtType := "TRANSACTION"
	upperContent := strings.ToUpper(strings.TrimSpace(content))
	if strings.HasPrefix(upperContent, "COMMIT") {
		stmtType = "COMMIT"
	} else if strings.HasPrefix(upperContent, "ROLLBACK") {
		stmtType = "ROLLBACK"
	} else if strings.HasPrefix(upperContent, "SAVEPOINT") {
		stmtType = "SAVEPOINT"
	}

	// Add the statement to the list
	l.Statements = append(l.Statements, statementModel{
		Content:     content,
		StartLine:   startLine,
		EndLine:     endLine,
		StartColumn: startColumn,
		EndColumn:   endColumn,
		Type:        stmtType,
	})
}

// getDeterminedStatementType identifies the type of SQL statement
func getDeterminedStatementType(text string) string {
	text = strings.ToUpper(strings.TrimSpace(text))

	if strings.HasPrefix(text, "SELECT") {
		return "SELECT"
	} else if strings.HasPrefix(text, "INSERT") {
		return "INSERT"
	} else if strings.HasPrefix(text, "UPDATE") {
		return "UPDATE"
	} else if strings.HasPrefix(text, "DELETE") {
		return "DELETE"
	} else if strings.HasPrefix(text, "MERGE") {
		return "MERGE"
	} else if strings.HasPrefix(text, "CREATE") {
		if strings.Contains(text, "PACKAGE BODY") {
			return "CREATE_PACKAGE_BODY"
		} else if strings.Contains(text, "PACKAGE") {
			return "CREATE_PACKAGE"
		} else if strings.Contains(text, "PROCEDURE") {
			return "CREATE_PROCEDURE"
		} else if strings.Contains(text, "FUNCTION") {
			return "CREATE_FUNCTION"
		} else if strings.Contains(text, "TRIGGER") {
			return "CREATE_TRIGGER"
		} else if strings.Contains(text, "TABLE") {
			return "CREATE_TABLE"
		} else if strings.Contains(text, "VIEW") {
			return "CREATE_VIEW"
		} else if strings.Contains(text, "INDEX") {
			return "CREATE_INDEX"
		} else if strings.Contains(text, "SEQUENCE") {
			return "CREATE_SEQUENCE"
		} else if strings.Contains(text, "TYPE BODY") {
			return "CREATE_TYPE_BODY"
		} else if strings.Contains(text, "TYPE") {
			return "CREATE_TYPE"
		} else if strings.Contains(text, "MATERIALIZED VIEW") {
			return "CREATE_MATERIALIZED_VIEW"
		} else if strings.Contains(text, "SYNONYM") {
			return "CREATE_SYNONYM"
		} else if strings.Contains(text, "DATABASE LINK") {
			return "CREATE_DATABASE_LINK"
		} else {
			return "CREATE"
		}
	} else if strings.HasPrefix(text, "ALTER") {
		if strings.Contains(text, "TABLE") {
			return "ALTER_TABLE"
		} else if strings.Contains(text, "INDEX") {
			return "ALTER_INDEX"
		} else if strings.Contains(text, "PROCEDURE") {
			return "ALTER_PROCEDURE"
		} else if strings.Contains(text, "FUNCTION") {
			return "ALTER_FUNCTION"
		} else if strings.Contains(text, "PACKAGE") {
			return "ALTER_PACKAGE"
		} else if strings.Contains(text, "TRIGGER") {
			return "ALTER_TRIGGER"
		} else if strings.Contains(text, "SEQUENCE") {
			return "ALTER_SEQUENCE"
		} else {
			return "ALTER"
		}
	} else if strings.HasPrefix(text, "DROP") {
		if strings.Contains(text, "TABLE") {
			return "DROP_TABLE"
		} else if strings.Contains(text, "INDEX") {
			return "DROP_INDEX"
		} else if strings.Contains(text, "PROCEDURE") {
			return "DROP_PROCEDURE"
		} else if strings.Contains(text, "FUNCTION") {
			return "DROP_FUNCTION"
		} else if strings.Contains(text, "PACKAGE") {
			return "DROP_PACKAGE"
		} else if strings.Contains(text, "TRIGGER") {
			return "DROP_TRIGGER"
		} else if strings.Contains(text, "SEQUENCE") {
			return "DROP_SEQUENCE"
		} else if strings.Contains(text, "VIEW") {
			return "DROP_VIEW"
		} else {
			return "DROP"
		}
	} else if strings.HasPrefix(text, "TRUNCATE") {
		return "TRUNCATE"
	} else if strings.HasPrefix(text, "GRANT") {
		return "GRANT"
	} else if strings.HasPrefix(text, "REVOKE") {
		return "REVOKE"
	} else if strings.HasPrefix(text, "COMMIT") {
		return "COMMIT"
	} else if strings.HasPrefix(text, "ROLLBACK") {
		return "ROLLBACK"
	} else if strings.HasPrefix(text, "SAVEPOINT") {
		return "SAVEPOINT"
	} else if strings.HasPrefix(text, "BEGIN") || strings.HasPrefix(text, "DECLARE") {
		return "PLSQL_BLOCK"
	} else if text == "/" {
		return "SLASH"
	} else if strings.HasPrefix(text, "EXPLAIN PLAN") {
		return "EXPLAIN_PLAN"
	} else if strings.HasPrefix(text, "COMMENT ON") {
		return "COMMENT"
	} else if strings.HasPrefix(text, "SET TRANSACTION") {
		return "SET_TRANSACTION"
	} else if strings.HasPrefix(text, "LOCK TABLE") {
		return "LOCK_TABLE"
	} else if strings.HasPrefix(text, "EXECUTE") || strings.HasPrefix(text, "EXEC") {
		return "EXECUTE"
	} else if strings.HasPrefix(text, "SHOW") {
		return "SHOW"
	} else if strings.HasPrefix(text, "DESC") || strings.HasPrefix(text, "DESCRIBE") {
		return "DESCRIBE"
	}

	return "UNKNOWN"
}
