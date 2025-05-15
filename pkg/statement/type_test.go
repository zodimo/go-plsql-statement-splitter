package statement

import (
	"encoding/json"
	"testing"
)

func TestTypeString(t *testing.T) {
	testCases := []struct {
		name     string
		stmtType Type
		expected string
	}{
		{"SELECT", TypeSelect, "SELECT"},
		{"INSERT", TypeInsert, "INSERT"},
		{"UPDATE", TypeUpdate, "UPDATE"},
		{"DELETE", TypeDelete, "DELETE"},
		{"MERGE", TypeMerge, "MERGE"},
		{"UNKNOWN", TypeUnknown, "UNKNOWN"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.stmtType.String()
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestTypeMarshalJSON(t *testing.T) {
	testCases := []struct {
		name     string
		stmtType Type
		expected string
	}{
		{"SELECT", TypeSelect, `"SELECT"`},
		{"INSERT", TypeInsert, `"INSERT"`},
		{"UNKNOWN", TypeUnknown, `"UNKNOWN"`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bytes, err := json.Marshal(tc.stmtType)
			if err != nil {
				t.Fatalf("Failed to marshal: %v", err)
			}
			result := string(bytes)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestTypeUnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name     string
		json     string
		expected Type
	}{
		{"SELECT", `"SELECT"`, TypeSelect},
		{"INSERT", `"INSERT"`, TypeInsert},
		{"Unknown", `"UNKNOWN"`, TypeUnknown},
		{"Invalid", `"INVALID"`, TypeUnknown},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var result Type
			err := json.Unmarshal([]byte(tc.json), &result)
			if err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestParse(t *testing.T) {
	testCases := []struct {
		input    string
		expected Type
	}{
		{"SELECT", TypeSelect},
		{"INSERT", TypeInsert},
		{"UPDATE", TypeUpdate},
		{"DELETE", TypeDelete},
		{"MERGE", TypeMerge},
		{"CREATE_TABLE", TypeCreateTable},
		{"CREATE_PROCEDURE", TypeCreateProcedure},
		{"CREATE_FUNCTION", TypeCreateFunction},
		{"CREATE_PACKAGE", TypeCreatePackage},
		{"CREATE_PACKAGE_BODY", TypeCreatePackageBody},
		{"CREATE_TRIGGER", TypeCreateTrigger},
		{"ALTER_TABLE", TypeAlterTable},
		{"DROP_TABLE", TypeDropTable},
		{"TRUNCATE", TypeTruncate},
		{"GRANT", TypeGrant},
		{"COMMIT", TypeCommit},
		{"ROLLBACK", TypeRollback},
		{"PLSQL_BLOCK", TypePlsqlBlock},
		{"INVALID", TypeUnknown},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := Parse(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %s for input %s, got %s", tc.expected, tc.input, result)
			}
		})
	}
}

func TestTypeHelperMethods(t *testing.T) {
	// Test IsQuery
	if !TypeSelect.IsQuery() {
		t.Errorf("Expected SELECT to be a query")
	}
	if TypeInsert.IsQuery() {
		t.Errorf("Expected INSERT not to be a query")
	}

	// Test IsDML
	dmlTypes := []Type{TypeSelect, TypeInsert, TypeUpdate, TypeDelete, TypeMerge}
	for _, st := range dmlTypes {
		if !st.IsDML() {
			t.Errorf("Expected %s to be DML", st)
		}
	}
	if TypeCreateTable.IsDML() {
		t.Errorf("Expected CREATE_TABLE not to be DML")
	}

	// Test IsDDL
	ddlTypes := []Type{
		TypeCreateTable, TypeCreateProcedure, TypeAlterTable,
		TypeDropTable, TypeTruncate, TypeGrant, TypeRevoke,
	}
	for _, st := range ddlTypes {
		if !st.IsDDL() {
			t.Errorf("Expected %s to be DDL", st)
		}
	}
	if TypeSelect.IsDDL() {
		t.Errorf("Expected SELECT not to be DDL")
	}

	// Test IsTransactional
	transTypes := []Type{TypeCommit, TypeRollback, TypeSavepoint, TypeTransaction, TypeSetTransaction}
	for _, st := range transTypes {
		if !st.IsTransactional() {
			t.Errorf("Expected %s to be transactional", st)
		}
	}
	if TypeSelect.IsTransactional() {
		t.Errorf("Expected SELECT not to be transactional")
	}

	// Test IsPLSQL
	plsqlTypes := []Type{
		TypePlsqlBlock, TypeCreateProcedure, TypeCreateFunction,
		TypeCreatePackage, TypeCreatePackageBody,
		TypeCreateTrigger, TypeCreateType, TypeCreateTypeBody,
	}
	for _, st := range plsqlTypes {
		if !st.IsPLSQL() {
			t.Errorf("Expected %s to be PL/SQL", st)
		}
	}
	if TypeSelect.IsPLSQL() {
		t.Errorf("Expected SELECT not to be PL/SQL")
	}
}
