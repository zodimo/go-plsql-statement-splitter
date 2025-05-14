package parser

import (
	"reflect"
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
