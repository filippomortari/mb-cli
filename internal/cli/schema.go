package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var schemaPretty bool

var schemaCmd = &cobra.Command{
	Use:   "schema [command]",
	Short: "Print JSON schema for command inputs and outputs",
	Long:  "Prints JSON schema describing a command's parameters, types, defaults, and valid values. With no arguments, lists all commands.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runSchema,
}

func init() {
	rootCmd.AddCommand(schemaCmd)
	schemaCmd.Flags().BoolVar(&schemaPretty, "pretty", false, "Pretty-print JSON output")
}

type commandSummary struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type paramSchema struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Required    bool     `json:"required"`
	Default     any      `json:"default,omitempty"`
	Enum        []string `json:"enum,omitempty"`
	Description string   `json:"description"`
}

type commandSchema struct {
	Command    string        `json:"command"`
	Args       []paramSchema `json:"args,omitempty"`
	Flags      []paramSchema `json:"flags"`
	OutputKeys []string      `json:"output_keys,omitempty"`
}

func runSchema(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return printSchemaJSON(getCommandList())
	}

	schema, ok := schemas[args[0]]
	if !ok {
		return fmt.Errorf("unknown command: %s", args[0])
	}

	return printSchemaJSON(schema)
}

func printSchemaJSON(data any) error {
	var output []byte
	var err error

	if schemaPretty {
		output, err = json.MarshalIndent(data, "", "  ")
	} else {
		output, err = json.Marshal(data)
	}

	if err != nil {
		return fmt.Errorf("failed to marshal schema JSON: %w", err)
	}

	fmt.Fprintln(os.Stdout, string(output))
	return nil
}

// GetCommandList returns the list of command summaries. Exported for testing.
func GetCommandList() []commandSummary {
	return getCommandList()
}

// GetSchema returns the schema for a command. Exported for testing.
func GetSchema(name string) (commandSchema, bool) {
	s, ok := schemas[name]
	return s, ok
}

func getCommandList() []commandSummary {
	return []commandSummary{
		{Name: "database list", Description: "List all databases"},
		{Name: "database get", Description: "Get database details"},
		{Name: "database metadata", Description: "Full metadata (tables + fields)"},
		{Name: "database fields", Description: "List all fields in database"},
		{Name: "database schemas", Description: "List schema names"},
		{Name: "database schema", Description: "Tables in a specific schema"},
		{Name: "table list", Description: "List all tables"},
		{Name: "table get", Description: "Get table details"},
		{Name: "table metadata", Description: "Table metadata with fields"},
		{Name: "table fks", Description: "Foreign key relationships"},
		{Name: "table data", Description: "Get raw table data"},
		{Name: "field get", Description: "Get field details"},
		{Name: "field summary", Description: "Summary statistics for a field"},
		{Name: "field values", Description: "Distinct values for a field"},
		{Name: "query sql", Description: "Run a native SQL query"},
		{Name: "card list", Description: "List saved questions"},
		{Name: "card get", Description: "Get card details"},
		{Name: "card run", Description: "Execute a saved question"},
		{Name: "search", Description: "Search across Metabase items"},
	}
}

var schemas = map[string]commandSchema{
	"database list": {
		Command: "database list",
		Flags:   []paramSchema{},
	},
	"database get": {
		Command: "database get",
		Args: []paramSchema{
			{Name: "id", Type: "integer", Required: true, Description: "Database ID"},
		},
		Flags:      []paramSchema{},
		OutputKeys: []string{"id", "name", "engine", "description", "details"},
	},
	"database metadata": {
		Command: "database metadata",
		Args: []paramSchema{
			{Name: "id", Type: "integer", Required: true, Description: "Database ID"},
		},
		Flags:      []paramSchema{},
		OutputKeys: []string{"id", "name", "tables"},
	},
	"database fields": {
		Command: "database fields",
		Args: []paramSchema{
			{Name: "id", Type: "integer", Required: true, Description: "Database ID"},
		},
		Flags:      []paramSchema{},
		OutputKeys: []string{"id", "name", "base_type", "table_name"},
	},
	"database schemas": {
		Command: "database schemas",
		Args: []paramSchema{
			{Name: "id", Type: "integer", Required: true, Description: "Database ID"},
		},
		Flags: []paramSchema{},
	},
	"database schema": {
		Command: "database schema",
		Args: []paramSchema{
			{Name: "id", Type: "integer", Required: true, Description: "Database ID"},
			{Name: "schema", Type: "string", Required: true, Description: "Schema name"},
		},
		Flags: []paramSchema{},
	},
	"table list": {
		Command: "table list",
		Flags:   []paramSchema{},
	},
	"table get": {
		Command: "table get",
		Args: []paramSchema{
			{Name: "id", Type: "integer", Required: true, Description: "Table ID"},
		},
		Flags:      []paramSchema{},
		OutputKeys: []string{"id", "name", "db_id", "schema", "description"},
	},
	"table metadata": {
		Command: "table metadata",
		Args: []paramSchema{
			{Name: "id", Type: "integer", Required: true, Description: "Table ID"},
		},
		Flags:      []paramSchema{},
		OutputKeys: []string{"id", "name", "fields"},
	},
	"table fks": {
		Command: "table fks",
		Args: []paramSchema{
			{Name: "id", Type: "integer", Required: true, Description: "Table ID"},
		},
		Flags: []paramSchema{},
	},
	"table data": {
		Command: "table data",
		Args: []paramSchema{
			{Name: "id", Type: "integer", Required: true, Description: "Table ID"},
		},
		Flags: []paramSchema{},
	},
	"field get": {
		Command: "field get",
		Args: []paramSchema{
			{Name: "id", Type: "integer", Required: true, Description: "Field ID"},
		},
		Flags:      []paramSchema{},
		OutputKeys: []string{"id", "name", "base_type", "semantic_type", "table_id"},
	},
	"field summary": {
		Command: "field summary",
		Args: []paramSchema{
			{Name: "id", Type: "integer", Required: true, Description: "Field ID"},
		},
		Flags: []paramSchema{},
	},
	"field values": {
		Command: "field values",
		Args: []paramSchema{
			{Name: "id", Type: "integer", Required: true, Description: "Field ID"},
		},
		Flags: []paramSchema{},
	},
	"query sql": {
		Command: "query sql",
		Flags: []paramSchema{
			{Name: "db", Type: "string", Required: true, Description: "Database ID or name substring"},
			{Name: "sql", Type: "string", Required: true, Description: "SQL query to execute"},
			{Name: "export", Type: "string", Required: false, Enum: []string{"csv", "json", "xlsx"}, Description: "Export format"},
			{Name: "limit", Type: "integer", Required: false, Default: 0, Description: "Append LIMIT to SQL query"},
			{Name: "fields", Type: "string", Required: false, Description: "Comma-separated columns to include in output"},
		},
	},
	"card list": {
		Command: "card list",
		Flags:   []paramSchema{},
	},
	"card get": {
		Command: "card get",
		Args: []paramSchema{
			{Name: "id", Type: "integer", Required: true, Description: "Card ID"},
		},
		Flags:      []paramSchema{},
		OutputKeys: []string{"id", "name", "description", "database_id", "display", "query_type"},
	},
	"card run": {
		Command: "card run",
		Args: []paramSchema{
			{Name: "id", Type: "integer", Required: true, Description: "Card ID"},
		},
		Flags: []paramSchema{
			{Name: "fields", Type: "string", Required: false, Description: "Comma-separated columns to include in output"},
		},
	},
	"search": {
		Command: "search",
		Args: []paramSchema{
			{Name: "query", Type: "string", Required: true, Description: "Search query string"},
		},
		Flags: []paramSchema{
			{Name: "models", Type: "string", Required: false, Enum: []string{"table", "card", "database", "dashboard", "collection", "metric"}, Description: "Filter by type (comma-separated)"},
		},
		OutputKeys: []string{"id", "name", "model", "database_id", "description"},
	},
}
