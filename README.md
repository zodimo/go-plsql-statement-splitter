# go-plsql-statement-splitter

A Go library for accurately splitting Oracle PL/SQL scripts into individual statements with precise boundary detection using ANTLR4 for parsing.

# DISCLAIMER
This project was created with the help of AI assistance:
- Development performed using [Cursor IDE](https://cursor.com), an AI-powered code editor
- Project structure and development methodology based on the [VAN Memory Bank](https://github.com/vanzan01/cursor-memory-bank) framework

## Overview

The goal of this library is to help developers extract individual SQL statements from PL/SQL scripts with 100% accurate boundary detection. It uses ANTLR4 for parsing and provides precise source location tracking for each statement.

## Features

- Split PL/SQL scripts into individual statements with accurate boundary detection
- Track line and column numbers for each statement
- Properly handle both single-line and multi-line comments
- Process input from both files and strings
- Provide detailed syntax error reporting
- JSON marshalling support for all output structures

## Requirements

- Go 1.21 or later
- ANTLR4 runtime for Go

## Installation

```bash
go get github.com/zodimo/go-plsql-statement-splitter
```

## Usage

Basic usage example:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/zodimo/go-plsql-statement-splitter/pkg/splitter"
)

func main() {
    // Split statements from a file
    statements, err := splitter.SplitFile("path/to/script.sql")
    if err != nil {
        log.Fatalf("Error splitting file: %v", err)
    }
    
    for i, stmt := range statements {
        fmt.Printf("Statement %d: %s\n", i+1, stmt.Content)
        fmt.Printf("  Position: %d:%d to %d:%d\n", 
            stmt.StartLine, stmt.StartColumn, 
            stmt.EndLine, stmt.EndColumn)
    }
    
    // Split statements from a string
    sqlContent := `
    SELECT * FROM employees;
    
    CREATE OR REPLACE PROCEDURE hello_world IS
    BEGIN
        DBMS_OUTPUT.PUT_LINE('Hello, World!');
    END;
    /
    `
    
    statements, err = splitter.SplitString(sqlContent)
    if err != nil {
        log.Fatalf("Error splitting string: %v", err)
    }
    
    for i, stmt := range statements {
        fmt.Printf("Statement %d: %s\n", i+1, stmt.Content)
    }
}
```

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Implementation Details

This library uses ANTLR4 for PL/SQL parsing and statement boundary detection. It incorporates the ANTLR4 grammar from the [zodimo/plsql-parser](https://github.com/zodimo/plsql-parser) repository, which provides comprehensive support for Oracle PL/SQL syntax.

### Key Features

- Accurate statement boundary detection using a formal grammar approach
- Support for complex PL/SQL constructs
- Detailed positional information for each statement
- Statement type classification (SELECT, INSERT, etc.)

### Building from Source

To build the project from source, you need to have ANTLR4 installed:

```bash
# Generate the parser code
cd internal/parser
./generate.sh
```

## Advanced Usage

### Customizing Splitter Behavior

The library provides several configuration options to customize the behavior of the splitter:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/zodimo/go-plsql-statement-splitter/pkg/splitter"
)

func main() {
    // Create a splitter with custom options
    s := splitter.NewSplitter(
        splitter.WithPositionInfo(true),          // Include position information (default: true)
        splitter.WithVerboseErrors(true),         // Include detailed error messages (default: false)
        splitter.WithMaxErrors(5),                // Maximum number of errors to report (default: 1)
        splitter.WithErrorContext(true),          // Include error context (default: false)
        splitter.WithErrorStatement(true),        // Include statement causing error (default: false)
    )
    
    // Split statements from a file with custom options
    statements, err := s.SplitFile("path/to/script.sql")
    if err != nil {
        syntaxErr, isSyntaxErr := err.(*splitter.SyntaxError)
        if isSyntaxErr {
            fmt.Printf("Syntax error at line %d, column %d: %s\n", 
                syntaxErr.Line, syntaxErr.Column, syntaxErr.Message)
        } else {
            log.Fatalf("Error splitting file: %v", err)
        }
    }
    
    for i, stmt := range statements {
        fmt.Printf("Statement %d (%s): %s\n", i+1, stmt.Type, stmt.Content)
    }
}
```

### Working with io.Reader

The library supports reading from an io.Reader:

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/zodimo/go-plsql-statement-splitter/pkg/splitter"
)

func main() {
    // Open a file
    file, err := os.Open("path/to/script.sql")
    if err != nil {
        log.Fatalf("Error opening file: %v", err)
    }
    defer file.Close()
    
    // Split statements from an io.Reader
    statements, err := splitter.SplitReader(file)
    if err != nil {
        log.Fatalf("Error splitting file: %v", err)
    }
    
    for i, stmt := range statements {
        fmt.Printf("Statement %d: %s\n", i+1, stmt.Content)
    }
}
```

### Getting All Syntax Errors

To get all syntax errors in a script:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/zodimo/go-plsql-statement-splitter/pkg/splitter"
)

func main() {
    s := splitter.NewSplitter(
        splitter.WithErrorContext(true),
        splitter.WithErrorStatement(true),
    )
    
    // Get all syntax errors in a script
    content := `
    SELECT * FROM;
    INSERT INTO employees (id name) VALUES (1, 'John');
    `
    
    errors, err := s.GetAllSyntaxErrors(content)
    if err != nil {
        log.Fatalf("Error getting syntax errors: %v", err)
    }
    
    for i, syntaxErr := range errors {
        fmt.Printf("Error %d: Line %d, Column %d: %s\n", 
            i+1, syntaxErr.Line, syntaxErr.Column, syntaxErr.Message)
        if syntaxErr.Context != "" {
            fmt.Printf("  Context: %s\n", syntaxErr.Context)
        }
    }
}
```

## Command Line Interface

The library includes a command-line interface for splitting SQL scripts:

```bash
# Basic usage
go run cmd/splitter/main.go script.sql

# Specify output format
go run cmd/splitter/main.go -format=json script.sql

# Save output to a file
go run cmd/splitter/main.go -format=json -output=output.json script.sql

# Include verbose error messages
go run cmd/splitter/main.go -verbose-errors script.sql

# Show all syntax errors with context
go run cmd/splitter/main.go -all-errors -error-context script.sql
```

Available CLI options:

```
  -all-errors
        Show all errors, ignoring max-errors setting
  -error-context
        Include context lines for errors
  -error-statement
        Include full statement with errors
  -format string
        Output format: text or json (default "text")
  -indent string
        Indentation for JSON output (default "  ")
  -max-errors int
        Maximum number of errors to report (default 5)
  -no-position
        Don't include position information
  -output string
        Output file (works with any format)
  -pretty
        Pretty print JSON output (default true)
  -print-statements
        Print the statements (default true)
  -print-types
        Print statement types (default true)
  -verbose-errors
        Show detailed error information
```

## Implementation Details

### Parser Architecture

The library uses a two-phase approach:

1. ANTLR4 parsing of the entire PL/SQL script
2. Visitor/listener pattern to extract individual statements with position information

### Supported Statement Types

The library can identify the following statement types:

- DML: SELECT, INSERT, UPDATE, DELETE, MERGE
- DDL: CREATE_TABLE, CREATE_VIEW, CREATE_INDEX, etc.
- PL/SQL: PLSQL_BLOCK, CREATE_PROCEDURE, CREATE_FUNCTION, etc.
- Transaction control: COMMIT, ROLLBACK, SAVEPOINT
- Other: EXPLAIN_PLAN, LOCK_TABLE, etc.

### Dependencies

- github.com/antlr4-go/antlr/v4: ANTLR4 runtime for Go
- zodimo/plsql-parser: PL/SQL grammar files (embedded in the project)

### Testing

The project has a comprehensive testing workflow to ensure code quality:

```bash
# Run the test script with all checks
./scripts/test.sh -a

# Run basic tests
go test ./...

# Run tests with race detector
go test -race ./...

# Run tests with coverage
go test -cover ./...
```

For more details about testing, see [docs/testing.md](docs/testing.md).

#### Continuous Integration

This project uses GitHub Actions for continuous integration:

- **Test workflow**: Runs tests and linting on push and pull requests
- **Coverage workflow**: Generates and uploads code coverage reports 
- **Security scanning**: Checks for security issues in the code and dependencies

#### Contributing Tests

When adding new features or fixing bugs, please include appropriate tests:

- Add unit tests for new functionality
- Ensure existing tests continue to pass
- Consider adding benchmark tests for performance-critical code

## Future Enhancements

Planned enhancements for future releases:

- Streaming support for processing very large files
- Enhanced error recovery for incomplete statements
- Support for additional Oracle-specific syntax
- Performance optimizations for large scripts

### Advanced Error Handling

The library provides detailed error reporting capabilities to help diagnose and fix syntax errors:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/zodimo/go-plsql-statement-splitter/pkg/splitter"
)

func main() {
    // Create a splitter with enhanced error reporting
    s := splitter.NewSplitter(
        splitter.WithVerboseErrors(true),         // Include detailed error messages
        splitter.WithMaxErrors(5),                // Maximum number of errors to report
        splitter.WithErrorContext(true),          // Include error context
        splitter.WithErrorContextLines(5),        // Show 5 lines before and after each error
        splitter.WithErrorStatement(true),        // Include statement causing error
    )
    
    // Try to split statements and handle errors
    _, err := s.SplitFile("path/to/script.sql")
    if err != nil {
        syntaxErr, isSyntaxErr := err.(*splitter.SyntaxError)
        if isSyntaxErr {
            // This will print the error with context showing 5 lines before and after
            fmt.Printf("Syntax error detected:\n%s\n", syntaxErr.Error())
        } else {
            log.Fatalf("Error: %v", err)
        }
    }
}
```

An example of the enhanced error output:

```
Syntax error at line 42, column 10: mismatched input 'END' expecting {';', ','}
40 |     FOR emp IN (SELECT * FROM employees) LOOP
41 |         DBMS_OUTPUT.PUT_LINE('Employee: ' || emp.name)
42 |     END LOOP
           ^
43 | END;
44 | /
```

The error output includes:
- Error message with line and column number
- Context showing several lines before and after the error
- A marker (^) pointing precisely to the error position

## Acknowledgements

### Development Tools

- [Cursor](https://cursor.com) - An advanced IDE powered by AI that was used for the development of this project.

### Development Methodology

This project utilizes the structured memory bank system developed by [vanzan01 - cursor-memory-bank](https://github.com/vanzan01/cursor-memory-bank), which provides an organized framework for AI-assisted software development.
