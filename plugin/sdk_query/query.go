package sdk_query

import "context"

// =======================================================================
// Define structs identical to "github.com/kolide/osquery-go/plugin/table"
// allows easier integration with osquery extensions such as CloudQuery.
// =======================================================================

type Table struct {
	Name     string
	Columns  []ColumnDefinition
	Generate GenerateFunc
}

type ColumnDefinition struct {
	Name string
	Type ColumnType
}

//GenerateFunc callback function that fetches data
type GenerateFunc func(ctx context.Context, queryContext QueryContext) ([]map[string]string, error)

//QueryContext passed to plugins as the relevant WHERE parts
type QueryContext struct {
	// Constraints is a map from column Name to the details of the
	// constraints on that column.
	Constraints map[string]ConstraintList `json:"constraints,omitempty"`

	// limit, offset, order by
	Limit   int      `json:"limit,omitempty"`
	Offset  int      `json:"offset,omitempty"`
	OrderBy []string `json:"order_by,omitempty"`
	Desc    bool     `json:"desc,omitempty"`
	Columns []string `json:"columns,omitempty"`

	// Limit the number of results to protect our RAM usage
	MaxRows int `json:"max_rows,omitempty"`
}

// ConstraintList contains the details of the constraints for the given column.
type ConstraintList struct {
	Affinity    ColumnType   `json:"affinity,omitempty"`
	Constraints []Constraint `json:"constraints,omitempty"`
}

// Constraint contains both an operator and an expression that are applied as
// constraints in the query.
type Constraint struct {
	Operator   Op       `json:"operator,omitempty"`
	Expression []string `json:"expression,omitempty"`
}

// Op is type of operations.
type Op uint8

// Op mean identity of operations.
const (
	OpEQ         Op = 2
	OpGT            = 4
	OpLE            = 8
	OpLT            = 16
	OpGE            = 32
	OpMATCH         = 64
	OpLIKE          = 65 /* 3.10.0 and later only */
	OpGLOB          = 66 /* 3.10.0 and later only */
	OpREGEXP        = 67 /* 3.10.0 and later only */
	OpScanUnique    = 1  /* Scan visits at most 1 row */
)

// ColumnType is a strongly typed representation of the data type string for a
// column definition. The named constants should be used.
type ColumnType string

const (
	ColumnTypeText    ColumnType = "TEXT"
	ColumnTypeInteger            = "INTEGER"
	ColumnTypeNumeric            = "NUMERIC"
	ColumnTypeReal               = "REAL"
)
