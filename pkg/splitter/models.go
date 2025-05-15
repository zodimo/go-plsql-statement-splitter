package splitter

// Statement represents a single PL/SQL statement with position information
type Statement struct {
	Content     string `json:"content"`
	StartLine   int    `json:"startLine"`
	EndLine     int    `json:"endLine"`
	StartColumn int    `json:"startColumn"`
	EndColumn   int    `json:"endColumn"`
	Type        string `json:"type,omitempty"` // If available from ANTLR parser
}

// SyntaxError represents a syntax error in a PL/SQL script
type SyntaxError struct {
	Line      int    `json:"line"`      // Line number where the error occurred
	Column    int    `json:"column"`    // Column number where the error occurred
	Message   string `json:"message"`   // Error message
	Statement string `json:"statement"` // The statement that caused the error
	Context   string `json:"context"`   // Context lines showing the error in context
}
