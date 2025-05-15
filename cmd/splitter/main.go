package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/zodimo/go-plsql-statement-splitter/pkg/splitter"
)

func main() {
	// Define flags
	var (
		outputFormat        string
		noPositionInfo      bool
		verboseErrors       bool
		maxErrors           int
		outputFile          string
		printStatements     bool
		printStatementTypes bool
		includeErrorContext bool
		includeErrorStmt    bool
		allErrors           bool
		jsonPretty          bool
		jsonIndent          string
		contextLines        int
	)

	flag.StringVar(&outputFormat, "format", "text", "Output format: text or json")
	flag.BoolVar(&noPositionInfo, "no-position", false, "Don't include position information")
	flag.BoolVar(&verboseErrors, "verbose-errors", false, "Show detailed error information")
	flag.IntVar(&maxErrors, "max-errors", 5, "Maximum number of errors to report")
	flag.StringVar(&outputFile, "output", "", "Output file (works with any format)")
	flag.BoolVar(&printStatements, "print-statements", true, "Print the statements")
	flag.BoolVar(&printStatementTypes, "print-types", true, "Print statement types")
	flag.BoolVar(&includeErrorContext, "error-context", false, "Include context lines for errors")
	flag.BoolVar(&includeErrorStmt, "error-statement", false, "Include full statement with errors")
	flag.BoolVar(&allErrors, "all-errors", false, "Show all errors, ignoring max-errors setting")
	flag.BoolVar(&jsonPretty, "pretty", true, "Pretty print JSON output")
	flag.StringVar(&jsonIndent, "indent", "  ", "Indentation for JSON output")
	flag.IntVar(&contextLines, "context-lines", 3, "Number of context lines to show before and after errors")
	flag.Parse()

	// Check if a file path was provided
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Usage:")
		fmt.Println("  splitter [options] <file>")
		fmt.Println("  If no file is provided, a demo will be run")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		fmt.Println("\nExamples:")
		fmt.Println("  splitter script.sql")
		fmt.Println("  splitter -format=json script.sql")
		fmt.Println("  splitter -format=json -output=result.json script.sql")
		fmt.Println("  splitter -verbose-errors script.sql")
		fmt.Println("  splitter -all-errors -error-context -context-lines=5 invalid.sql")

		fmt.Println("\nRunning demo...")
		demoSplitString()
		return
	}

	// Create splitter with options
	splitterOpts := []splitter.Option{}
	if noPositionInfo {
		splitterOpts = append(splitterOpts, splitter.WithPositionInfo(false))
	}
	if verboseErrors {
		splitterOpts = append(splitterOpts, splitter.WithVerboseErrors(true))
	}
	if maxErrors > 0 && !allErrors {
		splitterOpts = append(splitterOpts, splitter.WithMaxErrors(maxErrors))
	} else if allErrors {
		splitterOpts = append(splitterOpts, splitter.WithMaxErrors(0)) // 0 means unlimited
	}
	if includeErrorContext {
		splitterOpts = append(splitterOpts, splitter.WithErrorContext(true))
		splitterOpts = append(splitterOpts, splitter.WithErrorContextLines(contextLines))
	}
	if includeErrorStmt {
		splitterOpts = append(splitterOpts, splitter.WithErrorStatement(true))
	}

	s := splitter.NewSplitter(splitterOpts...)

	// Split statements from a file
	filePath := args[0]
	fmt.Printf("Splitting SQL statements from file: %s\n", filePath)

	// Check if file exists
	if !splitter.FileExists(filePath) {
		log.Fatalf("File not found: %s", filePath)
	}

	statements, err := s.SplitFile(filePath)
	if err != nil {
		syntaxErr, isSyntaxErr := err.(*splitter.SyntaxError)
		if isSyntaxErr {
			if verboseErrors || includeErrorContext {
				// Print the full error message with context if available
				log.Fatalf("Syntax error:\n%s", syntaxErr.Error())
			} else {
				log.Fatalf("Syntax error at line %d, column %d: %s",
					syntaxErr.Line, syntaxErr.Column, syntaxErr.Message)
			}
		} else {
			log.Fatalf("Error splitting file: %v", err)
		}
	}

	// Output according to format
	if outputFormat == "json" {
		outputJSON(statements, outputFile, jsonPretty, jsonIndent)
	} else {
		outputText(statements, printStatements, printStatementTypes, outputFile)
	}
}

func outputJSON(statements []splitter.Statement, outputFile string, pretty bool, indent string) {
	var data []byte
	var err error

	if pretty {
		data, err = json.MarshalIndent(statements, "", indent)
	} else {
		data, err = json.Marshal(statements)
	}

	if err != nil {
		log.Fatalf("Error marshalling to JSON: %v", err)
	}

	if outputFile != "" {
		if err := os.WriteFile(outputFile, data, 0644); err != nil {
			log.Fatalf("Error writing JSON to file: %v", err)
		}
		fmt.Printf("JSON output written to %s\n", outputFile)
	} else {
		fmt.Println(string(data))
	}
}

func outputText(statements []splitter.Statement, printStatements, printTypes bool, outputFile string) {
	// Prepare the output
	var output strings.Builder

	fmt.Fprintf(&output, "Found %d statements:\n\n", len(statements))

	for i, stmt := range statements {
		fmt.Fprintf(&output, "Statement %d", i+1)
		if printTypes {
			fmt.Fprintf(&output, " (%s)", stmt.Type)
		}
		fmt.Fprintf(&output, ":\n")

		fmt.Fprintf(&output, "  Position: %d:%d to %d:%d\n",
			stmt.StartLine, stmt.StartColumn, stmt.EndLine, stmt.EndColumn)

		if printStatements {
			// Only print up to 100 characters of content, with ellipsis if truncated
			content := stmt.Content
			if len(content) > 100 {
				content = content[:97] + "..."
			}
			fmt.Fprintf(&output, "  Content: %s\n", content)
		}
		fmt.Fprintf(&output, "\n")
	}

	// Output to file or stdout
	if outputFile != "" {
		if err := os.WriteFile(outputFile, []byte(output.String()), 0644); err != nil {
			log.Fatalf("Error writing text to file: %v", err)
		}
		fmt.Printf("Text output written to %s\n", outputFile)
	} else {
		fmt.Print(output.String())
	}
}

func demoSplitString() {
	// Example PL/SQL script
	script := `
-- Simple SQL statements
SELECT * FROM employees;
INSERT INTO employees (id, name) VALUES (1, 'John');

-- PL/SQL block
BEGIN
    FOR emp IN (SELECT * FROM employees) LOOP
        DBMS_OUTPUT.PUT_LINE('Employee: ' || emp.name);
    END LOOP;
END;
/

-- Another SQL statement
UPDATE employees SET salary = salary * 1.1;
`

	fmt.Println("Splitting SQL statements from example string:")
	statements, err := splitter.SplitString(script)
	if err != nil {
		log.Fatalf("Error splitting string: %v", err)
	}

	// Print the statements
	fmt.Printf("Found %d statements:\n\n", len(statements))
	for i, stmt := range statements {
		fmt.Printf("Statement %d (%s):\n", i+1, stmt.Type)
		fmt.Printf("  Position: %d:%d to %d:%d\n", stmt.StartLine, stmt.StartColumn, stmt.EndLine, stmt.EndColumn)

		// Only print up to 100 characters of content, with ellipsis if truncated
		content := stmt.Content
		if len(content) > 100 {
			content = content[:97] + "..."
		}
		fmt.Printf("  Content: %s\n\n", content)
	}
}

// demoSyntaxErrors shows how to handle syntax errors
func demoSyntaxErrors() {
	// Example invalid PL/SQL script
	script := `
-- Invalid SQL with syntax errors
SELECT * FROM employees WHERE;

-- Unclosed string
INSERT INTO employees (id, name) VALUES (1, 'John);
`

	// Create a splitter with custom error handling options
	s := splitter.NewSplitter(
		splitter.WithVerboseErrors(true),
		splitter.WithErrorContext(true),
		splitter.WithErrorContextLines(5),
		splitter.WithMaxErrors(5),
	)

	// Try to split the invalid script
	_, err := s.SplitString(script)
	if err != nil {
		// Check if it's a syntax error
		syntaxErr, isSyntaxErr := err.(*splitter.SyntaxError)
		if isSyntaxErr {
			fmt.Printf("Syntax error detected:\n%s\n", syntaxErr.Error())
		} else {
			fmt.Printf("Error: %v\n", err)
		}
	}

	// Get all syntax errors
	errors, err := s.GetAllSyntaxErrors(script)
	if err != nil {
		fmt.Printf("Error getting syntax errors: %v\n", err)
	} else {
		fmt.Printf("Found %d syntax errors:\n", len(errors))
		for i, syntaxErr := range errors {
			fmt.Printf("%d. Line %d, Column %d: %s\n", i+1, syntaxErr.Line, syntaxErr.Column, syntaxErr.Message)
		}
	}
}
