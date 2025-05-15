package splitter

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	internalParser "github.com/zodimo/go-plsql-statement-splitter/internal/parser"
	"github.com/zodimo/go-plsql-statement-splitter/pkg/statement"
)

// Splitter is responsible for splitting PL/SQL scripts into individual statements
type Splitter struct {
	// Configuration options
	includePosition       bool
	verboseErrors         bool
	maxErrors             int
	includeContext        bool
	includeErrorStatement bool
	contextLines          int // Number of context lines to show before and after the error
}

// NewSplitter creates a new Splitter instance with the provided options
func NewSplitter(options ...Option) *Splitter {
	s := &Splitter{
		includePosition:       true,
		verboseErrors:         false,
		maxErrors:             1,
		includeContext:        false,
		includeErrorStatement: false,
		contextLines:          3, // Default to 3 lines of context before and after
	}

	// Apply options
	for _, option := range options {
		option(s)
	}

	return s
}

// Option represents a configuration option for the Splitter
type Option func(*Splitter)

// WithPositionInfo configures whether position information is included
func WithPositionInfo(include bool) Option {
	return func(s *Splitter) {
		s.includePosition = include
	}
}

// WithVerboseErrors configures whether detailed error information is included
func WithVerboseErrors(verbose bool) Option {
	return func(s *Splitter) {
		s.verboseErrors = verbose
	}
}

// WithMaxErrors configures the maximum number of errors to return
func WithMaxErrors(max int) Option {
	return func(s *Splitter) {
		s.maxErrors = max
	}
}

// WithErrorContext configures whether error context (line and error pointer) is included
func WithErrorContext(include bool) Option {
	return func(s *Splitter) {
		s.includeContext = include
	}
}

// WithErrorStatement configures whether the full statement with the error is included
func WithErrorStatement(include bool) Option {
	return func(s *Splitter) {
		s.includeErrorStatement = include
	}
}

// WithErrorContextLines configures the number of context lines to include before and after errors
func WithErrorContextLines(lines int) Option {
	return func(s *Splitter) {
		if lines >= 0 {
			s.contextLines = lines
		}
	}
}

// SplitFile splits a PL/SQL script file into individual statements
func SplitFile(filePath string) ([]Statement, error) {
	splitter := NewSplitter()
	return splitter.SplitFile(filePath)
}

// SplitReader splits a PL/SQL script from an io.Reader into individual statements
func SplitReader(reader io.Reader) ([]Statement, error) {
	splitter := NewSplitter()
	return splitter.SplitReader(reader)
}

// SplitString splits a PL/SQL script string into individual statements
func SplitString(content string) ([]Statement, error) {
	splitter := NewSplitter()
	return splitter.SplitString(content)
}

// SplitFile splits a PL/SQL script file into individual statements
func (s *Splitter) SplitFile(filePath string) ([]Statement, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return s.SplitString(string(data))
}

// SplitReader splits a PL/SQL script from an io.Reader into individual statements
func (s *Splitter) SplitReader(reader io.Reader) ([]Statement, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}

	return s.SplitString(string(data))
}

// SplitString splits a PL/SQL script string into individual statements
func (s *Splitter) SplitString(content string) ([]Statement, error) {
	if strings.TrimSpace(content) == "" {
		return []Statement{}, nil
	}

	// Use the ANTLR4 parser to parse the SQL
	parsedStatements, syntaxErrors, err := internalParser.ParseStringWithOptions(content, s.maxErrors, s.contextLines)
	if err != nil {
		return nil, fmt.Errorf("parser error: %w", err)
	}

	// If there are syntax errors, return an error
	if len(syntaxErrors) > 0 {
		if s.verboseErrors {
			// Return all errors up to the maximum
			var errorMessages []string
			maxErrors := s.maxErrors
			if maxErrors <= 0 || maxErrors > len(syntaxErrors) {
				maxErrors = len(syntaxErrors)
			}

			for i := 0; i < maxErrors; i++ {
				err := syntaxErrors[i]
				errMsg := fmt.Sprintf("Line %d, Column %d: %s", err.Line, err.Column, err.Message)

				// Include context if configured
				if s.includeContext && err.Context != "" {
					errMsg += "\n" + err.Context
				}

				errorMessages = append(errorMessages, errMsg)
			}

			if len(syntaxErrors) > maxErrors {
				errorMessages = append(errorMessages, fmt.Sprintf("... and %d more errors", len(syntaxErrors)-maxErrors))
			}

			// Create the syntax error with detailed information
			syntaxErr := &SyntaxError{
				Message: strings.Join(errorMessages, "\n"),
				Line:    syntaxErrors[0].Line,
				Column:  syntaxErrors[0].Column,
				Context: syntaxErrors[0].Context,
			}

			// Include statement content if configured and available
			if s.includeErrorStatement && syntaxErrors[0].Context != "" {
				syntaxErr.Statement = strings.Split(syntaxErrors[0].Context, "\n")[0]
			}

			return nil, syntaxErr
		} else {
			// Just return the first error
			firstError := syntaxErrors[0]

			// Create the syntax error with basic information
			syntaxErr := &SyntaxError{
				Message: firstError.Message,
				Line:    firstError.Line,
				Column:  firstError.Column,
				Context: firstError.Context,
			}

			// Include statement content if configured and available
			if s.includeErrorStatement && firstError.Context != "" {
				syntaxErr.Statement = strings.Split(firstError.Context, "\n")[0]
			}

			return nil, syntaxErr
		}
	}

	// Convert internal statement representation to public model
	statements := make([]Statement, 0, len(parsedStatements))
	for _, stmt := range parsedStatements {
		statement := Statement{
			Content: stmt.Content,
			Type:    statement.Parse(stmt.Type),
		}

		// Include position information if configured
		if s.includePosition {
			statement.StartLine = stmt.StartLine
			statement.EndLine = stmt.EndLine
			statement.StartColumn = stmt.StartColumn
			statement.EndColumn = stmt.EndColumn
		}

		statements = append(statements, statement)
	}

	return statements, nil
}

// Error messages
var (
	ErrEmptyInput = errors.New("empty input")
	ErrSyntax     = errors.New("syntax error")
	ErrReadFile   = errors.New("error reading file")
	ErrParsing    = errors.New("error parsing SQL")
)

// Error implements the error interface for SyntaxError
func (e *SyntaxError) Error() string {
	if e.Context != "" {
		return fmt.Sprintf("syntax error at line %d, column %d: %s\n%s", e.Line, e.Column, e.Message, e.Context)
	}

	// Handle case where error message is too verbose (containing the full token list)
	if strings.Contains(e.Message, "expecting {") && len(e.Message) > 200 {
		// Extract just the first part of the message before the token list
		parts := strings.SplitN(e.Message, "expecting {", 2)
		if len(parts) == 2 {
			return fmt.Sprintf("syntax error at line %d, column %d: %s expecting {...}",
				e.Line, e.Column, strings.TrimSpace(parts[0]))
		}
	}

	return fmt.Sprintf("syntax error at line %d, column %d: %s", e.Line, e.Column, e.Message)
}

// GetSyntaxErrors returns all syntax errors that were encountered during parsing
func (s *Splitter) GetSyntaxErrors(content string) ([]SyntaxError, error) {
	if strings.TrimSpace(content) == "" {
		return []SyntaxError{}, nil
	}

	// Use the ANTLR4 parser to parse the SQL
	_, syntaxErrors, err := internalParser.ParseStringWithOptions(content, s.maxErrors, s.contextLines)
	if err != nil {
		return nil, fmt.Errorf("parser error: %w", err)
	}

	// Convert internal syntax errors to public model
	errors := make([]SyntaxError, 0, len(syntaxErrors))
	for _, err := range syntaxErrors {
		syntaxErr := SyntaxError{
			Message: err.Message,
			Line:    err.Line,
			Column:  err.Column,
			Context: err.Context,
		}

		// Include statement content if configured and available
		if s.includeErrorStatement && err.Context != "" {
			syntaxErr.Statement = strings.Split(err.Context, "\n")[0]
		}

		errors = append(errors, syntaxErr)
	}

	return errors, nil
}

// GetAllSyntaxErrors returns all syntax errors without limiting to maxErrors
func (s *Splitter) GetAllSyntaxErrors(content string) ([]SyntaxError, error) {
	// Create a temporary splitter with unlimited errors
	tempSplitter := NewSplitter(
		WithMaxErrors(9999), // Use a large number to effectively make it unlimited
		WithErrorStatement(s.includeErrorStatement),
		WithErrorContext(s.includeContext),
		WithErrorContextLines(s.contextLines),
	)

	return tempSplitter.GetSyntaxErrors(content)
}

// FileExists returns whether the given file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// SplitFileWithPosition is a convenience function that includes position information
func SplitFileWithPosition(filePath string) ([]Statement, error) {
	splitter := NewSplitter(WithPositionInfo(true))
	return splitter.SplitFile(filePath)
}

// SplitStringWithPosition is a convenience function that includes position information
func SplitStringWithPosition(content string) ([]Statement, error) {
	splitter := NewSplitter(WithPositionInfo(true))
	return splitter.SplitString(content)
}

// SplitReaderWithPosition is a convenience function that includes position information
func SplitReaderWithPosition(reader io.Reader) ([]Statement, error) {
	splitter := NewSplitter(WithPositionInfo(true))
	return splitter.SplitReader(reader)
}
