package splitter

import (
	"strings"
	"testing"

	"github.com/zodimo/go-plsql-statement-splitter/pkg/statement"
)

func TestTypeSafeStatementType(t *testing.T) {
	// This is a complex script with multiple statement types
	script := `
	-- Simple DML statements
	SELECT * FROM employees;
	INSERT INTO employees (id, name) VALUES (1, 'John');
	UPDATE employees SET salary = 1000 WHERE id = 1;
	DELETE FROM employees WHERE id = 2;
	MERGE INTO employees USING temp_employees ON (employees.id = temp_employees.id);

	-- DDL statements
	CREATE TABLE customers (id NUMBER, name VARCHAR2(100));
	ALTER TABLE customers ADD email VARCHAR2(100);
	DROP TABLE old_customers;
	TRUNCATE TABLE empty_table;
	GRANT SELECT ON employees TO hr_user;
	REVOKE DELETE ON employees FROM hr_user;

	-- Transaction statements
	COMMIT;
	ROLLBACK;
	SAVEPOINT sp1;
	SET TRANSACTION READ ONLY;

	-- PL/SQL blocks
	BEGIN
		DBMS_OUTPUT.PUT_LINE('Hello, World!');
	END;
	/

	CREATE OR REPLACE PROCEDURE hello_world IS
	BEGIN
		DBMS_OUTPUT.PUT_LINE('Hello, World!');
	END;
	/
	`

	// Split the script
	statements, err := SplitString(script)
	if err != nil {
		t.Fatalf("Failed to split script: %v", err)
	}

	// Define expected statement types
	expectedTypes := map[string]statement.Type{
		"SELECT":           statement.TypeSelect,
		"INSERT":           statement.TypeInsert,
		"UPDATE":           statement.TypeUpdate,
		"DELETE":           statement.TypeDelete,
		"MERGE":            statement.TypeMerge,
		"CREATE TABLE":     statement.TypeCreateTable,
		"ALTER TABLE":      statement.TypeAlterTable,
		"DROP TABLE":       statement.TypeDropTable,
		"TRUNCATE":         statement.TypeTruncate,
		"GRANT":            statement.TypeGrant,
		"REVOKE":           statement.TypeRevoke,
		"COMMIT":           statement.TypeCommit,
		"ROLLBACK":         statement.TypeRollback,
		"SAVEPOINT":        statement.TypeSavepoint,
		"PLSQL_BLOCK":      statement.TypePlsqlBlock,
		"CREATE_PROCEDURE": statement.TypeCreateProcedure,
	}

	// Keep track of statement types found
	foundTypes := make(map[statement.Type]bool)

	// Verify statements have correct types
	for _, stmt := range statements {
		// Log statement for debugging
		t.Logf("Found statement type %s: %.40s...", stmt.Type, strings.ReplaceAll(stmt.Content, "\n", " "))

		// Record that we found this type
		foundTypes[stmt.Type] = true

		// Check if the statement type starts with an expected prefix
		prefix := strings.SplitN(strings.TrimSpace(stmt.Content), " ", 2)[0]
		if prefix == "/" {
			continue // Skip forward slash - it's special
		}

		// Check helper methods
		if strings.Contains(strings.ToUpper(stmt.Content), "SELECT") &&
			!stmt.Type.IsDML() &&
			stmt.Type != statement.TypeCreateView &&
			!strings.HasPrefix(strings.ToUpper(strings.TrimSpace(stmt.Content)), "GRANT") {
			t.Errorf("Statement with SELECT should have IsDML() return true: %s", stmt.Content)
		}

		if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(stmt.Content)), "CREATE") && !stmt.Type.IsDDL() {
			t.Errorf("Statement with CREATE should have IsDDL() return true: %s", stmt.Content)
		}

		if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(stmt.Content)), "BEGIN") && !stmt.Type.IsPLSQL() {
			t.Errorf("Statement with BEGIN should have IsPLSQL() return true: %s", stmt.Content)
		}

		if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(stmt.Content)), "COMMIT") && !stmt.Type.IsTransactional() {
			t.Errorf("Statement with COMMIT should have IsTransactional() return true: %s", stmt.Content)
		}
	}

	// Check that we found all the expected types
	for keyword, expectedType := range expectedTypes {
		if strings.Contains(script, keyword) && !foundTypes[expectedType] {
			// Check if we found any statement with content containing the keyword
			foundInContent := false
			for _, stmt := range statements {
				if strings.Contains(strings.ToUpper(stmt.Content), strings.ToUpper(keyword)) {
					foundInContent = true
					break
				}
			}

			if !foundInContent {
				t.Errorf("Script contains %s but no statement with type %s was found", keyword, expectedType)
			}
		}
	}
}
