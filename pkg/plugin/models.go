package plugin

// PluginSettings holds the datasource configuration.
type PluginSettings struct {
	ConsistencyLevel string `json:"consistencyLevel"`
	Timeout          string `json:"timeout"`
}

// QueryModel represents a query from the frontend.
type QueryModel struct {
	RawSQL      string   `json:"rawSql"`
	Format      string   `json:"format"` // "table" or "time_series"
	TimeColumns []string `json:"timeColumns"`

	// Visual builder fields
	EditorMode  string            `json:"editorMode"` // "code" or "builder"
	Table       string            `json:"table"`
	Columns     []ColumnSelection `json:"columns"`
	WhereClause []WhereCondition  `json:"whereClause"`
	GroupBy     []string          `json:"groupBy"`
	OrderBy     []OrderByClause   `json:"orderBy"`
	Limit       string            `json:"limit"`
}

// ColumnSelection represents a column with an optional aggregation.
type ColumnSelection struct {
	Name        string `json:"name"`
	Aggregation string `json:"aggregation"` // "", "COUNT", "SUM", "AVG", "MIN", "MAX"
}

// WhereCondition represents a single WHERE clause condition.
type WhereCondition struct {
	Column   string `json:"column"`
	Operator string `json:"operator"` // "=", "!=", "<", ">", "<=", ">=", "LIKE", "IN", "IS NULL", "IS NOT NULL"
	Value    string `json:"value"`
}

// OrderByClause represents an ORDER BY clause.
type OrderByClause struct {
	Column    string `json:"column"`
	Direction string `json:"direction"` // "ASC" or "DESC"
}

// RqliteQueryRequest is the request body sent to rqlite's /db/query endpoint.
type RqliteQueryRequest []string

// RqliteQueryResponse is the response from rqlite's /db/query endpoint.
type RqliteQueryResponse struct {
	Results []RqliteResult `json:"results"`
}

// RqliteResult is a single result set from rqlite.
type RqliteResult struct {
	Columns []string        `json:"columns"`
	Types   []string        `json:"types"`
	Values  [][]interface{} `json:"values"`
	Error   string          `json:"error"`
}
