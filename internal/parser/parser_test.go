package parser

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestDeduplicateStatements(t *testing.T) {
	// Create test cases with duplicate statements
	testCases := []struct {
		name     string
		input    []statementModel
		expected []statementModel
	}{
		{
			name:     "Empty input",
			input:    []statementModel{},
			expected: []statementModel{},
		},
		{
			name: "No duplicates",
			input: []statementModel{
				{Content: "SELECT * FROM employees", StartLine: 1, StartColumn: 0, EndLine: 1, EndColumn: 23, Type: "SELECT"},
				{Content: "UPDATE employees SET salary = 1000", StartLine: 2, StartColumn: 0, EndLine: 2, EndColumn: 32, Type: "UPDATE"},
			},
			expected: []statementModel{
				{Content: "SELECT * FROM employees", StartLine: 1, StartColumn: 0, EndLine: 1, EndColumn: 23, Type: "SELECT"},
				{Content: "UPDATE employees SET salary = 1000", StartLine: 2, StartColumn: 0, EndLine: 2, EndColumn: 32, Type: "UPDATE"},
			},
		},
		{
			name: "Duplicate positions",
			input: []statementModel{
				{Content: "SELECT * FROM employees", StartLine: 1, StartColumn: 0, EndLine: 1, EndColumn: 23, Type: "SELECT"},
				{Content: "SELECT * FROM employees", StartLine: 1, StartColumn: 0, EndLine: 1, EndColumn: 23, Type: "SELECT"},
			},
			expected: []statementModel{
				{Content: "SELECT * FROM employees", StartLine: 1, StartColumn: 0, EndLine: 1, EndColumn: 23, Type: "SELECT"},
			},
		},
		{
			name: "Improved type for duplicate",
			input: []statementModel{
				{Content: "SELECT * FROM employees", StartLine: 1, StartColumn: 0, EndLine: 1, EndColumn: 23, Type: "UNKNOWN"},
				{Content: "SELECT * FROM employees", StartLine: 1, StartColumn: 0, EndLine: 1, EndColumn: 23, Type: "SELECT"},
			},
			expected: []statementModel{
				{Content: "SELECT * FROM employees", StartLine: 1, StartColumn: 0, EndLine: 1, EndColumn: 23, Type: "SELECT"},
			},
		},
	}

	// Run the tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := deduplicateStatements(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestParseString(t *testing.T) {
	// Test cases
	testCases := []struct {
		name           string
		input          string
		expectedCount  int
		expectedTypes  []string
		expectedErrors int
	}{
		{
			name:           "Empty input",
			input:          "",
			expectedCount:  0,
			expectedTypes:  []string{},
			expectedErrors: 0,
		},
		{
			name:           "Simple SELECT",
			input:          "SELECT * FROM employees;",
			expectedCount:  1,
			expectedTypes:  []string{"SELECT"},
			expectedErrors: 0,
		},
		{
			name:           "Multiple statements",
			input:          "SELECT * FROM employees; INSERT INTO departments (id, name) VALUES (1, 'HR');",
			expectedCount:  2,
			expectedTypes:  []string{"SELECT", "INSERT"},
			expectedErrors: 0,
		},
		{
			name: "PL/SQL block",
			input: `
BEGIN
    UPDATE employees SET salary = salary * 1.1;
    COMMIT;
END;
`,
			expectedCount:  1, // With PL/SQL block depth tracking, this is now a single statement
			expectedTypes:  []string{"PLSQL_BLOCK"},
			expectedErrors: 0,
		},
		{
			name: "Invalid SQL",
			input: `
SELECT * FROM;
`,
			expectedCount:  1, // Parser finds partial statement but reports errors
			expectedTypes:  []string{"SELECT"},
			expectedErrors: 1, // Should have a syntax error
		},
		{
			name: "Oracle 12c WITH clause",
			input: `
WITH emp_avg_sal AS (
    SELECT department_id, AVG(salary) avg_sal
    FROM employees
    GROUP BY department_id
)
SELECT e.department_id, e.first_name, e.salary
FROM employees e JOIN emp_avg_sal eas ON e.department_id = eas.department_id
WHERE e.salary > eas.avg_sal;
`,
			expectedCount:  1,
			expectedTypes:  []string{"UNKNOWN"}, // Parser returns UNKNOWN for WITH clause
			expectedErrors: 0,
		},
		{
			name: "Oracle 12c Identity Column",
			input: `
CREATE TABLE departments (
    department_id NUMBER GENERATED ALWAYS AS IDENTITY,
    department_name VARCHAR2(50) NOT NULL
);
`,
			expectedCount:  1,
			expectedTypes:  []string{"CREATE_TABLE"},
			expectedErrors: 0,
		},
	}

	// Run the tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			statements, errors, err := ParseString(tc.input)

			// Check for unexpected errors
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Check statement count
			if len(statements) != tc.expectedCount {
				t.Errorf("Expected %d statements, got %d", tc.expectedCount, len(statements))
			}

			// Check statement types if we have the expected count
			if len(statements) == len(tc.expectedTypes) {
				for i, expectedType := range tc.expectedTypes {
					if i < len(statements) && statements[i].Type != expectedType {
						t.Errorf("Statement %d: expected type %s, got %s", i, expectedType, statements[i].Type)
					}
				}
			}

			// Check error count
			if len(errors) != tc.expectedErrors {
				t.Errorf("Expected %d errors, got %d", tc.expectedErrors, len(errors))
				for i, err := range errors {
					t.Logf("Error %d: %s at line %d, column %d", i, err.Message, err.Line, err.Column)
				}
			}
		})
	}
}

func TestGetDeterminedStatementType(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"SELECT * FROM employees", "SELECT"},
		{"INSERT INTO employees VALUES (1, 'John')", "INSERT"},
		{"UPDATE employees SET salary = 1000", "UPDATE"},
		{"DELETE FROM employees WHERE id = 1", "DELETE"},
		{"MERGE INTO employees USING temp_employees ON (employees.id = temp_employees.id)", "MERGE"},
		{"CREATE TABLE employees (id NUMBER, name VARCHAR2(100))", "CREATE_TABLE"},
		{"CREATE OR REPLACE PROCEDURE update_salary IS BEGIN NULL; END;", "CREATE_PROCEDURE"},
		{"CREATE OR REPLACE FUNCTION get_salary RETURN NUMBER IS BEGIN RETURN 0; END;", "CREATE_FUNCTION"},
		{"CREATE OR REPLACE PACKAGE emp_pkg IS PROCEDURE get_emp(id NUMBER); END;", "CREATE_PACKAGE"},
		{"CREATE OR REPLACE PACKAGE BODY emp_pkg IS PROCEDURE get_emp(id NUMBER) IS BEGIN NULL; END; END;", "CREATE_PACKAGE_BODY"},
		{"CREATE OR REPLACE TRIGGER emp_trg AFTER INSERT ON employees BEGIN NULL; END;", "CREATE_TRIGGER"},
		{"CREATE INDEX emp_idx ON employees(id)", "CREATE_INDEX"},
		{"ALTER TABLE employees ADD COLUMN email VARCHAR2(100)", "ALTER_TABLE"},
		{"DROP TABLE employees", "DROP_TABLE"},
		{"TRUNCATE TABLE employees", "TRUNCATE"},
		{"COMMIT", "COMMIT"},
		{"ROLLBACK", "ROLLBACK"},
		{"SAVEPOINT sp1", "SAVEPOINT"},
		{"BEGIN NULL; END;", "PLSQL_BLOCK"},
		{"/", "SLASH"},
		{"EXPLAIN PLAN FOR SELECT * FROM employees", "EXPLAIN_PLAN"},
		{"SOME_UNKNOWN_STATEMENT", "UNKNOWN"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := getDeterminedStatementType(tc.input)
			if result != tc.expected {
				t.Errorf("Expected type %s for input %s, got %s", tc.expected, tc.input, result)
			}
		})
	}
}

// TestIsVersion12Enabled tests that version 12 features are enabled by default
func TestIsVersion12Enabled(t *testing.T) {
	if !IsVersion12Enabled() {
		t.Error("Expected Version 12 to be enabled, but it was not")
	}
}

// TestCustomErrorListener_ContextLines tests the enhanced error reporting with context lines
func TestCustomErrorListener_ContextLines(t *testing.T) {
	// Test cases with different numbers of context lines
	testCases := []struct {
		name         string
		contextLines int
		input        string
		errorLine    int
		errorColumn  int
	}{
		{
			name:         "No context lines",
			contextLines: 0,
			input: `SELECT * FROM employees;
INSERT INTO employees (id, name) VALUES (1, 'John');
-- This is a syntax error line
SELECT * FROM WHERE;
UPDATE employees SET salary = 1000;
DELETE FROM employees WHERE id = 1;`,
			errorLine:   4,
			errorColumn: 14,
		},
		{
			name:         "3 context lines",
			contextLines: 3,
			input: `SELECT * FROM employees;
INSERT INTO employees (id, name) VALUES (1, 'John');
-- This is a syntax error line
SELECT * FROM WHERE;
UPDATE employees SET salary = 1000;
DELETE FROM employees WHERE id = 1;
COMMIT;
ROLLBACK;`,
			errorLine:   4,
			errorColumn: 14,
		},
		{
			name:         "Error at beginning of file",
			contextLines: 2,
			input: `SELECT * FROM WHERE;
INSERT INTO employees (id, name) VALUES (1, 'John');
UPDATE employees SET salary = 1000;`,
			errorLine:   1,
			errorColumn: 14,
		},
		{
			name:         "Error at end of file",
			contextLines: 2,
			input: `SELECT * FROM employees;
INSERT INTO employees (id, name) VALUES (1, 'John');
UPDATE employees SET salary = WHERE;`,
			errorLine:   3,
			errorColumn: 27,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a custom error listener with the specified context lines
			listener := NewCustomErrorListener(5, tc.input, tc.contextLines)

			// Simulate an error
			listener.SyntaxError(nil, nil, tc.errorLine, tc.errorColumn, "Test syntax error", nil)

			// Check that we captured the error
			if len(listener.Errors) != 1 {
				t.Fatalf("Expected 1 error, got %d", len(listener.Errors))
			}

			error := listener.Errors[0]

			// Check basic error properties
			if error.Line != tc.errorLine || error.Column != tc.errorColumn {
				t.Errorf("Expected error at line %d, column %d; got line %d, column %d",
					tc.errorLine, tc.errorColumn, error.Line, error.Column)
			}

			// Verify error marker (^) is included in the context
			if !strings.Contains(error.Context, "^") {
				t.Errorf("Error marker (^) missing from context:\n%s", error.Context)
			}

			// For context lines = 0, should still have the error line
			if tc.contextLines == 0 {
				lineNumberStr := fmt.Sprintf("%d |", tc.errorLine)
				if !strings.Contains(error.Context, lineNumberStr) {
					t.Errorf("Error line number %s missing from context with contextLines=0:\n%s",
						lineNumberStr, error.Context)
				}
			} else {
				// Calculate expected number of visible context lines
				inputLines := strings.Count(tc.input, "\n") + 1

				minLine := tc.errorLine - tc.contextLines
				if minLine < 1 {
					minLine = 1
				}

				maxLine := tc.errorLine + tc.contextLines
				if maxLine > inputLines {
					maxLine = inputLines
				}

				expectedLines := maxLine - minLine + 1

				// Add 1 for the line with the error marker (^)
				totalExpectedLines := expectedLines + 1

				// Count actual lines in context
				actualLines := strings.Count(error.Context, "\n")
				if actualLines == totalExpectedLines-1 {
					// This is fine - the last line might not have a trailing newline
					actualLines++
				}

				if actualLines != totalExpectedLines {
					t.Errorf("Expected %d total lines in context, got %d: %s",
						totalExpectedLines, actualLines, error.Context)
				}
			}

			// Verify line numbers are included for all cases
			lineNumberStr := fmt.Sprintf("%d |", tc.errorLine)
			if !strings.Contains(error.Context, lineNumberStr) {
				t.Errorf("Error line number %s missing from context:\n%s", lineNumberStr, error.Context)
			}
		})
	}
}

// TestParseStringWithOptions_ContextLines tests that the ParseStringWithOptions function properly passes context lines to the error listener
func TestParseStringWithOptions_ContextLines(t *testing.T) {
	// Test input with a syntax error
	input := `SELECT * FROM employees;
INSERT INTO employees (id, name) VALUES (1, 'John');
-- This line has a syntax error
SELECT * FROM WHERE;  -- Missing table name
UPDATE employees SET salary = 1000;`

	// Try with different context line settings
	contextLineSettings := []int{0, 1, 3, 5}

	for _, contextLines := range contextLineSettings {
		t.Run(fmt.Sprintf("ContextLines=%d", contextLines), func(t *testing.T) {
			_, errors, err := ParseStringWithOptions(input, 1, contextLines)

			// Should be no error in ParseStringWithOptions itself
			if err != nil {
				t.Fatalf("ParseStringWithOptions failed: %v", err)
			}

			// Should have syntax errors
			if len(errors) == 0 {
				t.Fatalf("Expected syntax errors, got none")
			}

			// Get the first error
			syntaxErr := errors[0]

			// Verify error location
			if syntaxErr.Line != 4 || syntaxErr.Column != 14 {
				t.Errorf("Expected error at line 4, column 14; got line %d, column %d",
					syntaxErr.Line, syntaxErr.Column)
			}

			// Check context based on settings - All cases should now have the error line and marker

			// Should have the error line
			if !strings.Contains(syntaxErr.Context, "FROM WHERE") {
				t.Errorf("Error context missing the error line: %s", syntaxErr.Context)
			}

			// Should have error marker
			if !strings.Contains(syntaxErr.Context, "^") {
				t.Errorf("Error context missing error marker (^): %s", syntaxErr.Context)
			}

			// Should have line numbers
			if !strings.Contains(syntaxErr.Context, "4 |") {
				t.Errorf("Error context missing line number: %s", syntaxErr.Context)
			}

			// Additional checks for non-zero context lines
			if contextLines > 0 {
				// Count lines with content (not the marker line)
				lines := strings.Split(syntaxErr.Context, "\n")
				contentLines := 0
				for _, line := range lines {
					if strings.Contains(line, " | ") {
						contentLines++
					}
				}

				// Calculate expected visible lines - min of total lines and (2*contextLines + 1)
				totalInputLines := strings.Count(input, "\n") + 1
				expectedVisibleLines := min(totalInputLines, 2*contextLines+1)

				if contentLines != expectedVisibleLines {
					t.Errorf("Expected %d content lines in context, got %d: %s",
						expectedVisibleLines, contentLines, syntaxErr.Context)
				}
			}
		})
	}
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestNestedPLSQLBlocks(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedCount  int
		expectedTypes  []string
		expectedErrors int
	}{
		{
			name: "Simple Nested Block",
			input: `
BEGIN
  DECLARE
    x NUMBER := 0;
  BEGIN
    UPDATE employees SET salary = salary * 1.1;
  END;
END;
`,
			expectedCount:  1,
			expectedTypes:  []string{"PLSQL_BLOCK"},
			expectedErrors: 0,
		},
		{
			name: "Procedure with Nested Blocks",
			input: `
CREATE OR REPLACE PROCEDURE update_salary AS
BEGIN
  DECLARE
    v_min_salary NUMBER := 1000;
  BEGIN
    UPDATE employees SET salary = GREATEST(salary, v_min_salary);
    COMMIT;
  END;
END update_salary;
`,
			expectedCount:  1,
			expectedTypes:  []string{"CREATE_PROCEDURE"},
			expectedErrors: 0,
		},
		{
			name: "Complex Nested Structure",
			input: `
CREATE OR REPLACE PROCEDURE process_data(p_id NUMBER) AS
BEGIN
  DECLARE
    v_count NUMBER := 0;
  BEGIN
    SELECT COUNT(*) INTO v_count FROM employees WHERE department_id = p_id;
    
    IF v_count > 0 THEN
      DECLARE
        v_avg_salary NUMBER;
      BEGIN
        SELECT AVG(salary) INTO v_avg_salary FROM employees WHERE department_id = p_id;
        
        UPDATE employees 
        SET salary = CASE
                      WHEN salary < v_avg_salary THEN salary * 1.1
                      ELSE salary
                     END
        WHERE department_id = p_id;
        
        COMMIT;
      END;
    END IF;
  END;
END process_data;
`,
			expectedCount:  1,
			expectedTypes:  []string{"CREATE_PROCEDURE"},
			expectedErrors: 0,
		},
	}

	// Run the tests
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			statements, errors, err := ParseString(tc.input)

			// Check for unexpected errors
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Check statement count
			if len(statements) != tc.expectedCount {
				t.Errorf("Expected %d statements, got %d", tc.expectedCount, len(statements))
			}

			// Check statement types if we have the expected count
			if len(statements) == len(tc.expectedTypes) {
				for i, expectedType := range tc.expectedTypes {
					if i < len(statements) && statements[i].Type != expectedType {
						t.Errorf("Statement %d: expected type %s, got %s", i, expectedType, statements[i].Type)
					}
				}
			}

			// Check error count
			if len(errors) != tc.expectedErrors {
				t.Errorf("Expected %d errors, got %d", tc.expectedErrors, len(errors))
				for _, err := range errors {
					t.Logf("Error at line %d, col %d: %s", err.Line, err.Column, err.Message)
				}
			}
		})
	}
}
