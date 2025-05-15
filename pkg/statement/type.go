package statement

import (
	"encoding/json"
	"strings"
)

// Type represents the type of a SQL statement
type Type string

// Constants for different statement types
const (
	TypeUnknown                Type = "UNKNOWN"
	TypeSelect                 Type = "SELECT"
	TypeInsert                 Type = "INSERT"
	TypeUpdate                 Type = "UPDATE"
	TypeDelete                 Type = "DELETE"
	TypeMerge                  Type = "MERGE"
	TypeCreateTable            Type = "CREATE_TABLE"
	TypeCreateView             Type = "CREATE_VIEW"
	TypeCreateIndex            Type = "CREATE_INDEX"
	TypeCreateSequence         Type = "CREATE_SEQUENCE"
	TypeCreateProcedure        Type = "CREATE_PROCEDURE"
	TypeCreateFunction         Type = "CREATE_FUNCTION"
	TypeCreatePackage          Type = "CREATE_PACKAGE"
	TypeCreatePackageBody      Type = "CREATE_PACKAGE_BODY"
	TypeCreateTrigger          Type = "CREATE_TRIGGER"
	TypeCreateType             Type = "CREATE_TYPE"
	TypeCreateTypeBody         Type = "CREATE_TYPE_BODY"
	TypeCreateMaterializedView Type = "CREATE_MATERIALIZED_VIEW"
	TypeCreateSynonym          Type = "CREATE_SYNONYM"
	TypeCreateDatabaseLink     Type = "CREATE_DATABASE_LINK"
	TypeCreate                 Type = "CREATE"
	TypeAlterTable             Type = "ALTER_TABLE"
	TypeAlterIndex             Type = "ALTER_INDEX"
	TypeAlterProcedure         Type = "ALTER_PROCEDURE"
	TypeAlterFunction          Type = "ALTER_FUNCTION"
	TypeAlterPackage           Type = "ALTER_PACKAGE"
	TypeAlterTrigger           Type = "ALTER_TRIGGER"
	TypeAlterSequence          Type = "ALTER_SEQUENCE"
	TypeAlter                  Type = "ALTER"
	TypeDropTable              Type = "DROP_TABLE"
	TypeDropIndex              Type = "DROP_INDEX"
	TypeDropProcedure          Type = "DROP_PROCEDURE"
	TypeDropFunction           Type = "DROP_FUNCTION"
	TypeDropPackage            Type = "DROP_PACKAGE"
	TypeDropTrigger            Type = "DROP_TRIGGER"
	TypeDropSequence           Type = "DROP_SEQUENCE"
	TypeDropView               Type = "DROP_VIEW"
	TypeDrop                   Type = "DROP"
	TypeTruncate               Type = "TRUNCATE"
	TypeGrant                  Type = "GRANT"
	TypeRevoke                 Type = "REVOKE"
	TypeCommit                 Type = "COMMIT"
	TypeRollback               Type = "ROLLBACK"
	TypeSavepoint              Type = "SAVEPOINT"
	TypeTransaction            Type = "TRANSACTION"
	TypePlsqlBlock             Type = "PLSQL_BLOCK"
	TypeSlash                  Type = "SLASH"
	TypeExplainPlan            Type = "EXPLAIN_PLAN"
	TypeComment                Type = "COMMENT"
	TypeSetTransaction         Type = "SET_TRANSACTION"
	TypeLockTable              Type = "LOCK_TABLE"
	TypeExecute                Type = "EXECUTE"
	TypeShow                   Type = "SHOW"
	TypeDescribe               Type = "DESCRIBE"
)

// String returns the string representation of a Type
func (st Type) String() string {
	return string(st)
}

// MarshalJSON marshals a Type to JSON
func (st Type) MarshalJSON() ([]byte, error) {
	return json.Marshal(st.String())
}

// UnmarshalJSON unmarshals JSON to a Type
func (st *Type) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*st = Parse(s)
	return nil
}

// Parse parses a string to a Type
func Parse(s string) Type {
	switch strings.ToUpper(s) {
	case "SELECT":
		return TypeSelect
	case "INSERT":
		return TypeInsert
	case "UPDATE":
		return TypeUpdate
	case "DELETE":
		return TypeDelete
	case "MERGE":
		return TypeMerge
	case "CREATE_TABLE":
		return TypeCreateTable
	case "CREATE_VIEW":
		return TypeCreateView
	case "CREATE_INDEX":
		return TypeCreateIndex
	case "CREATE_SEQUENCE":
		return TypeCreateSequence
	case "CREATE_PROCEDURE":
		return TypeCreateProcedure
	case "CREATE_FUNCTION":
		return TypeCreateFunction
	case "CREATE_PACKAGE":
		return TypeCreatePackage
	case "CREATE_PACKAGE_BODY":
		return TypeCreatePackageBody
	case "CREATE_TRIGGER":
		return TypeCreateTrigger
	case "CREATE_TYPE":
		return TypeCreateType
	case "CREATE_TYPE_BODY":
		return TypeCreateTypeBody
	case "CREATE_MATERIALIZED_VIEW":
		return TypeCreateMaterializedView
	case "CREATE_SYNONYM":
		return TypeCreateSynonym
	case "CREATE_DATABASE_LINK":
		return TypeCreateDatabaseLink
	case "CREATE":
		return TypeCreate
	case "ALTER_TABLE":
		return TypeAlterTable
	case "ALTER_INDEX":
		return TypeAlterIndex
	case "ALTER_PROCEDURE":
		return TypeAlterProcedure
	case "ALTER_FUNCTION":
		return TypeAlterFunction
	case "ALTER_PACKAGE":
		return TypeAlterPackage
	case "ALTER_TRIGGER":
		return TypeAlterTrigger
	case "ALTER_SEQUENCE":
		return TypeAlterSequence
	case "ALTER":
		return TypeAlter
	case "DROP_TABLE":
		return TypeDropTable
	case "DROP_INDEX":
		return TypeDropIndex
	case "DROP_PROCEDURE":
		return TypeDropProcedure
	case "DROP_FUNCTION":
		return TypeDropFunction
	case "DROP_PACKAGE":
		return TypeDropPackage
	case "DROP_TRIGGER":
		return TypeDropTrigger
	case "DROP_SEQUENCE":
		return TypeDropSequence
	case "DROP_VIEW":
		return TypeDropView
	case "DROP":
		return TypeDrop
	case "TRUNCATE":
		return TypeTruncate
	case "GRANT":
		return TypeGrant
	case "REVOKE":
		return TypeRevoke
	case "COMMIT":
		return TypeCommit
	case "ROLLBACK":
		return TypeRollback
	case "SAVEPOINT":
		return TypeSavepoint
	case "TRANSACTION":
		return TypeTransaction
	case "PLSQL_BLOCK":
		return TypePlsqlBlock
	case "SLASH":
		return TypeSlash
	case "EXPLAIN_PLAN":
		return TypeExplainPlan
	case "COMMENT":
		return TypeComment
	case "SET_TRANSACTION":
		return TypeSetTransaction
	case "LOCK_TABLE":
		return TypeLockTable
	case "EXECUTE":
		return TypeExecute
	case "SHOW":
		return TypeShow
	case "DESCRIBE":
		return TypeDescribe
	default:
		return TypeUnknown
	}
}

// Helper methods for statement type categorization

// IsQuery returns true if the statement is a query
func (st Type) IsQuery() bool {
	return st == TypeSelect
}

// IsDML returns true if the statement is a DML statement
func (st Type) IsDML() bool {
	return st == TypeSelect || st == TypeInsert || st == TypeUpdate || st == TypeDelete || st == TypeMerge
}

// IsDDL returns true if the statement is a DDL statement
func (st Type) IsDDL() bool {
	return strings.HasPrefix(string(st), "CREATE") ||
		strings.HasPrefix(string(st), "ALTER") ||
		strings.HasPrefix(string(st), "DROP") ||
		st == TypeTruncate || st == TypeGrant || st == TypeRevoke
}

// IsTransactional returns true if the statement is a transaction control statement
func (st Type) IsTransactional() bool {
	return st == TypeCommit || st == TypeRollback || st == TypeSavepoint || st == TypeTransaction || st == TypeSetTransaction
}

// IsPLSQL returns true if the statement is a PL/SQL block
func (st Type) IsPLSQL() bool {
	return st == TypePlsqlBlock ||
		st == TypeCreateProcedure || st == TypeCreateFunction ||
		st == TypeCreatePackage || st == TypeCreatePackageBody ||
		st == TypeCreateTrigger || st == TypeCreateType || st == TypeCreateTypeBody
}
