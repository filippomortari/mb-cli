package client

// Database represents a Metabase database.
type Database struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Engine  string `json:"engine"`
	Details any    `json:"details,omitempty"`
	Tables  []Table `json:"tables,omitempty"`
}

// DatabaseMetadata represents full database metadata including tables and fields.
type DatabaseMetadata struct {
	ID     int             `json:"id"`
	Name   string          `json:"name"`
	Engine string          `json:"engine"`
	Tables []TableMetadata `json:"tables,omitempty"`
}

// Table represents a Metabase table.
type Table struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Schema      string `json:"schema"`
	DBId        int    `json:"db_id"`
	EntityType  string `json:"entity_type,omitempty"`
}

// TableMetadata represents table metadata including field details.
type TableMetadata struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	DisplayName string  `json:"display_name"`
	Schema      string  `json:"schema"`
	DBId        int     `json:"db_id"`
	Fields      []Field `json:"fields,omitempty"`
}

// Field represents a Metabase field (column).
type Field struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	DisplayName   string `json:"display_name"`
	BaseType      string `json:"base_type"`
	DatabaseType  string `json:"database_type"`
	SemanticType  string `json:"semantic_type,omitempty"`
	TableID       int    `json:"table_id"`
	TableName     string `json:"table_name,omitempty"`
}

// ForeignKey represents a foreign key relationship.
type ForeignKey struct {
	Relationship string         `json:"relationship"`
	Origin       FKFieldRef     `json:"origin"`
	Destination  FKFieldRef     `json:"destination"`
}

// FKFieldRef represents a field reference in a foreign key.
type FKFieldRef struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Table FKTableRef `json:"table"`
}

// FKTableRef represents a table reference in a foreign key.
type FKTableRef struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// QueryResult represents the result of a dataset query.
type QueryResult struct {
	Data QueryResultData `json:"data"`
}

// QueryResultData holds the columns and rows of a query result.
type QueryResultData struct {
	Columns []ResultColumn `json:"cols"`
	Rows    [][]any        `json:"rows"`
}

// ResultColumn describes a column in a query result.
type ResultColumn struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	BaseType    string `json:"base_type"`
}

// DatasetQuery represents a query request to the Metabase dataset API.
type DatasetQuery struct {
	Database int         `json:"database"`
	Type     string      `json:"type"`
	Native   NativeQuery `json:"native"`
}

// NativeQuery represents the native SQL query part of a dataset query.
type NativeQuery struct {
	Query string `json:"query"`
}

// FieldSummary represents summary statistics for a field.
type FieldSummary struct {
	Type  string `json:"type"`
	Value any    `json:"value"`
}

// FieldValues represents distinct values for a field.
type FieldValues struct {
	FieldID int      `json:"field_id"`
	Values  [][]any  `json:"values"`
}
