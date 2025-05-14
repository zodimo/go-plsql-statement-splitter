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

// SyntaxError represents a syntax error that occurred during parsing
type SyntaxError struct {
	Message   string `json:"message"`
	Line      int    `json:"line"`
	Column    int    `json:"column"`
	Statement string `json:"statement,omitempty"`
}
