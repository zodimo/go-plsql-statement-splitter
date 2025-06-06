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
	Errors       []SyntaxError
	MaxErrors    int    // Maximum number of errors to capture
	SourceText   string // The original source text for context
	ContextLines int    // Number of context lines to include before and after the error
}

// NewCustomErrorListener creates a new error listener with the given max errors and source text
func NewCustomErrorListener(maxErrors int, sourceText string, contextLines int) *CustomErrorListener {
	if maxErrors <= 0 {
		maxErrors = 10 // Default to 10 if not specified or invalid
	}
	if contextLines < 0 {
		contextLines = 0 // Ensure non-negative
	}
	return &CustomErrorListener{
		Errors:       make([]SyntaxError, 0),
		MaxErrors:    maxErrors,
		SourceText:   sourceText,
		ContextLines: contextLines,
	}
}

// SyntaxError captures syntax errors during parsing with enhanced context
func (l *CustomErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	// Check if we've reached the error limit
	if len(l.Errors) >= l.MaxErrors {
		return
	}

	// Extract token text from the offending symbol if possible
	var tokenText string
	if symbol, ok := offendingSymbol.(antlr.Token); ok {
		tokenText = symbol.GetText()
	}

	// Enhance error message for common nested block errors
	enhancedMsg := msg
	if strings.Contains(msg, "no viable alternative") && strings.Contains(tokenText, "BEGIN") {
		enhancedMsg += " - This might be an issue with nested PL/SQL blocks. Check the block structure and ensure BEGIN/END pairs match."
	}

	// Extract context lines if source text is available
	context := ""
	if l.SourceText != "" {
		// Always generate context, even if contextLines is 0
		context = l.extractErrorContext(line, column)
	}

	// Create and store the error
	l.Errors = append(l.Errors, SyntaxError{
		Line:      line,
		Column:    column,
		Message:   enhancedMsg,
		TokenText: tokenText,
		Context:   context,
	})
}

// extractErrorContext extracts the source code context around an error location
func (l *CustomErrorListener) extractErrorContext(line, column int) string {
	lines := strings.Split(l.SourceText, "\n")
	if line <= 0 || line > len(lines) {
		return ""
	}

	// Determine the range of lines to include in context
	startLine := line - l.ContextLines
	if startLine < 1 {
		startLine = 1
	}

	endLine := line + l.ContextLines
	if endLine > len(lines) {
		endLine = len(lines)
	}

	// For contextLines=0, ensure we at least include the error line
	if l.ContextLines == 0 {
		startLine = line
		endLine = line
	}

	// Generate context with line numbers
	var contextBuilder strings.Builder
	padding := len(fmt.Sprintf("%d", endLine)) // Calculate padding for line numbers

	for i := startLine; i <= endLine; i++ {
		lineContent := lines[i-1]
		lineNumber := fmt.Sprintf("%*d |", padding, i)

		if i == line {
			// This is the error line, add a marker
			contextBuilder.WriteString(lineNumber + " " + lineContent + "\n")
			markerPosition := padding + 3 // Length of "nnn | " prefix
			if column > 0 && column <= len(lineContent) {
				contextBuilder.WriteString(strings.Repeat(" ", markerPosition+column-1) + "^\n")
			} else {
				// If column is out of range, mark the start of the line
				contextBuilder.WriteString(strings.Repeat(" ", markerPosition) + "^\n")
			}
		} else {
			// Regular context line
			contextBuilder.WriteString(lineNumber + " " + lineContent + "\n")
		}
	}

	return contextBuilder.String()
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

// ParseString parses a SQL string with default error handling
func ParseString(input string) ([]Statement, []SyntaxError, error) {
	return ParseStringWithOptions(input, 1, 3)
}

// IsVersion12Enabled returns whether Oracle 12c features are enabled in the parser
func IsVersion12Enabled() bool {
	// Create a dummy parser to check its configuration
	inputStream := antlr.NewInputStream("")
	lexer := gen.NewPlSqlLexer(inputStream)
	tokenStream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := gen.NewPlSqlParser(tokenStream)

	// Set version 12 features (this is normally done in ParseStringWithOptions)
	parser.SetVersion12(true)

	// The method is unexported in the generated code, so we'll just return true
	// since we just set it above
	return true
}

// ParseStringWithOptions parses a SQL string with configurable error handling options
func ParseStringWithOptions(input string, maxErrors int, contextLines int) ([]Statement, []SyntaxError, error) {
	// Setup the ANTLR lexer and parser
	inputStream := antlr.NewInputStream(input)
	lexer := gen.NewPlSqlLexer(inputStream)

	// Create token stream with error recovery
	tokenStream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Create the parser with error recovery
	parser := gen.NewPlSqlParser(tokenStream)

	// Enable version 12 features by default
	parser.SetVersion12(true)

	// Set error recovery strategy
	if maxErrors <= 1 {
		// If we only want one error, use bail strategy for better performance
		parser.SetErrorHandler(antlr.NewBailErrorStrategy())
	} else {
		// For multiple errors, use default strategy with recovery
		parser.SetErrorHandler(antlr.NewDefaultErrorStrategy())
	}

	// Add custom error listener
	errorListener := NewCustomErrorListener(maxErrors, input, contextLines)
	parser.RemoveErrorListeners()
	parser.AddErrorListener(errorListener)

	// Create the statement listener
	listener := NewStatementListener(parser, tokenStream)

	// Start parsing
	antlr.ParseTreeWalkerDefault.Walk(listener, parser.Sql_script())

	// Process the statements
	statements := listener.Statements

	// Remove duplicates
	statements = deduplicateStatements(statements)

	// Convert internal statementModel to public Statement
	result := make([]Statement, len(statements))
	for i, stmt := range statements {
		result[i] = Statement{
			Content:     stmt.Content,
			StartLine:   stmt.StartLine,
			EndLine:     stmt.EndLine,
			StartColumn: stmt.StartColumn,
			EndColumn:   stmt.EndColumn,
			Type:        stmt.Type,
		}
	}

	// Return the statements and any syntax errors
	return result, errorListener.Errors, nil
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
	Statements      []statementModel
	parser          *gen.PlSqlParser
	tokenStream     *antlr.CommonTokenStream
	currentType     string // Track the current statement type
	plsqlBlockDepth int    // Track the nesting level of PL/SQL blocks
}

// NewStatementListener creates a new statement listener
func NewStatementListener(parser *gen.PlSqlParser, tokenStream *antlr.CommonTokenStream) *StatementListener {
	return &StatementListener{
		BasePlSqlParserListener: &gen.BasePlSqlParserListener{},
		Statements:              make([]statementModel, 0),
		parser:                  parser,
		tokenStream:             tokenStream,
		currentType:             "",
		plsqlBlockDepth:         0,
	}
}

// EnterUnit_statement is called when entering a unit_statement rule
func (l *StatementListener) EnterUnit_statement(ctx *gen.Unit_statementContext) {
	// Skip if we are inside a PL/SQL block
	if l.plsqlBlockDepth > 0 {
		return
	}

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
	// Skip if we are inside a PL/SQL block
	if l.plsqlBlockDepth > 0 {
		return
	}

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
	l.plsqlBlockDepth++
}

// ExitCreate_procedure_body is called when exiting a create_procedure_body rule
func (l *StatementListener) ExitCreate_procedure_body(ctx *gen.Create_procedure_bodyContext) {
	l.plsqlBlockDepth--
}

// EnterCreate_function_body is called when entering a create_function_body rule
func (l *StatementListener) EnterCreate_function_body(ctx *gen.Create_function_bodyContext) {
	l.currentType = "CREATE_FUNCTION"
	l.plsqlBlockDepth++
}

// ExitCreate_function_body is called when exiting a create_function_body rule
func (l *StatementListener) ExitCreate_function_body(ctx *gen.Create_function_bodyContext) {
	l.plsqlBlockDepth--
}

// EnterCreate_package is called when entering a create_package rule
func (l *StatementListener) EnterCreate_package(ctx *gen.Create_packageContext) {
	l.currentType = "CREATE_PACKAGE"
	l.plsqlBlockDepth++
}

// ExitCreate_package is called when exiting a create_package rule
func (l *StatementListener) ExitCreate_package(ctx *gen.Create_packageContext) {
	l.plsqlBlockDepth--
}

// EnterCreate_package_body is called when entering a create_package_body rule
func (l *StatementListener) EnterCreate_package_body(ctx *gen.Create_package_bodyContext) {
	l.currentType = "CREATE_PACKAGE_BODY"
	l.plsqlBlockDepth++
}

// ExitCreate_package_body is called when exiting a create_package_body rule
func (l *StatementListener) ExitCreate_package_body(ctx *gen.Create_package_bodyContext) {
	l.plsqlBlockDepth--
}

// EnterCreate_trigger is called when entering a create_trigger rule
func (l *StatementListener) EnterCreate_trigger(ctx *gen.Create_triggerContext) {
	l.currentType = "CREATE_TRIGGER"
	l.plsqlBlockDepth++
}

// ExitCreate_trigger is called when exiting a create_trigger rule
func (l *StatementListener) ExitCreate_trigger(ctx *gen.Create_triggerContext) {
	l.plsqlBlockDepth--
}

// EnterCreate_type is called when entering a create_type rule
func (l *StatementListener) EnterCreate_type(ctx *gen.Create_typeContext) {
	l.currentType = "CREATE_TYPE"
	l.plsqlBlockDepth++
}

// ExitCreate_type is called when exiting a create_type rule
func (l *StatementListener) ExitCreate_type(ctx *gen.Create_typeContext) {
	l.plsqlBlockDepth--
}

// EnterCreate_type_body is called when the CreateTypeBody rule is encountered
// Note: Update this if the exact gen.CreateTypeBodyContext name is different
func (l *StatementListener) EnterCreate_type_body(ctx interface{}) {
	l.currentType = "CREATE_TYPE_BODY"
	l.plsqlBlockDepth++
}

// ExitCreate_type_body is called when exiting a create_type_body rule
func (l *StatementListener) ExitCreate_type_body(ctx interface{}) {
	l.plsqlBlockDepth--
}

// EnterAnonymous_block is called when entering an anonymous_block rule
func (l *StatementListener) EnterAnonymous_block(ctx *gen.Anonymous_blockContext) {
	// Increment the block depth counter
	l.plsqlBlockDepth++
	// fmt.Printf("Entering anonymous block, depth now: %d\n", l.plsqlBlockDepth)

	// Skip adding the statement if it's a nested block
	if l.plsqlBlockDepth > 1 {
		return
	}

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

// ExitAnonymous_block is called when exiting an anonymous_block rule
func (l *StatementListener) ExitAnonymous_block(ctx *gen.Anonymous_blockContext) {
	// Decrement the block depth counter
	l.plsqlBlockDepth--
	// Ensure we don't go negative
	if l.plsqlBlockDepth < 0 {
		l.plsqlBlockDepth = 0
	}
	// fmt.Printf("Exiting anonymous block, depth now: %d\n", l.plsqlBlockDepth)
}

// EnterTransaction_control_statements is called when entering a transaction_control_statements rule
func (l *StatementListener) EnterTransaction_control_statements(ctx *gen.Transaction_control_statementsContext) {
	// Skip if we are inside a PL/SQL block
	if l.plsqlBlockDepth > 0 {
		return
	}

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

// EnterBlock is called when entering a block statement rule
func (l *StatementListener) EnterBlock(ctx *gen.BlockContext) {
	l.plsqlBlockDepth++
	// Add debug logging if needed
	// fmt.Printf("Entering block, depth now: %d\n", l.plsqlBlockDepth)
}

// ExitBlock is called when exiting a block statement rule
func (l *StatementListener) ExitBlock(ctx *gen.BlockContext) {
	l.plsqlBlockDepth--
	// Ensure we don't go negative
	if l.plsqlBlockDepth < 0 {
		l.plsqlBlockDepth = 0
	}
	// Add debug logging if needed
	// fmt.Printf("Exiting block, depth now: %d\n", l.plsqlBlockDepth)
}

// EnterSeq_of_declare_specs is called when entering a sequence of declare specifications
func (l *StatementListener) EnterSeq_of_declare_specs(ctx *gen.Seq_of_declare_specsContext) {
	// Track that we're in a declare section
	// fmt.Printf("Entering declare specs, block depth: %d\n", l.plsqlBlockDepth)
}

// ExitSeq_of_declare_specs is called when exiting a sequence of declare specifications
func (l *StatementListener) ExitSeq_of_declare_specs(ctx *gen.Seq_of_declare_specsContext) {
	// fmt.Printf("Exiting declare specs, block depth: %d\n", l.plsqlBlockDepth)
}

// EnterDeclare_block is called when entering a declare_block rule
func (l *StatementListener) EnterDeclare_block(ctx *gen.Declare_blockContext) {
	// Increment the block depth counter
	l.plsqlBlockDepth++
	// fmt.Printf("Entering declare block, depth now: %d\n", l.plsqlBlockDepth)
}

// ExitDeclare_block is called when exiting a declare_block rule
func (l *StatementListener) ExitDeclare_block(ctx *gen.Declare_blockContext) {
	// Decrement the block depth counter
	l.plsqlBlockDepth--
	// Ensure we don't go negative
	if l.plsqlBlockDepth < 0 {
		l.plsqlBlockDepth = 0
	}
	// fmt.Printf("Exiting declare block, depth now: %d\n", l.plsqlBlockDepth)
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
