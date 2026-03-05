package tests

import (
	"encoding/json"
	"testing"

	"github.com/andreagrandi/mb-cli/internal/cli"
)

func TestGetCommandListReturnsAllCommands(t *testing.T) {
	commands := cli.GetCommandList()

	if len(commands) == 0 {
		t.Fatal("command list is empty")
	}

	expectedCommands := []string{
		"database list", "database get", "database metadata",
		"database fields", "database schemas", "database schema",
		"table list", "table get", "table metadata", "table fks", "table data",
		"field get", "field summary", "field values",
		"query sql",
		"card list", "card get", "card run",
		"search",
	}

	nameSet := make(map[string]bool)
	for _, cmd := range commands {
		nameSet[cmd.Name] = true
	}

	for _, expected := range expectedCommands {
		if !nameSet[expected] {
			t.Errorf("command list missing: %s", expected)
		}
	}
}

func TestGetCommandListIsValidJSON(t *testing.T) {
	commands := cli.GetCommandList()

	data, err := json.Marshal(commands)
	if err != nil {
		t.Fatalf("failed to marshal command list: %v", err)
	}

	var parsed []map[string]string
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("command list JSON is not valid: %v", err)
	}

	for _, cmd := range parsed {
		if cmd["name"] == "" {
			t.Error("command has empty name")
		}
		if cmd["description"] == "" {
			t.Error("command has empty description")
		}
	}
}

func TestGetSchemaForQuerySQL(t *testing.T) {
	schema, ok := cli.GetSchema("query sql")
	if !ok {
		t.Fatal("schema not found for 'query sql'")
	}

	if schema.Command != "query sql" {
		t.Errorf("expected command 'query sql', got %s", schema.Command)
	}

	flagNames := make(map[string]bool)
	for _, f := range schema.Flags {
		flagNames[f.Name] = true
	}

	requiredFlags := []string{"db", "sql"}
	for _, name := range requiredFlags {
		if !flagNames[name] {
			t.Errorf("query sql schema missing required flag: %s", name)
		}
	}

	for _, f := range schema.Flags {
		if f.Name == "db" && !f.Required {
			t.Error("db flag should be required")
		}
		if f.Name == "sql" && !f.Required {
			t.Error("sql flag should be required")
		}
		if f.Name == "export" && f.Required {
			t.Error("export flag should not be required")
		}
	}
}

func TestGetSchemaForSearch(t *testing.T) {
	schema, ok := cli.GetSchema("search")
	if !ok {
		t.Fatal("schema not found for 'search'")
	}

	if len(schema.Args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(schema.Args))
	}
	if schema.Args[0].Name != "query" {
		t.Errorf("expected arg name 'query', got %s", schema.Args[0].Name)
	}
	if !schema.Args[0].Required {
		t.Error("search query arg should be required")
	}
}

func TestGetSchemaForDatabaseSchema(t *testing.T) {
	schema, ok := cli.GetSchema("database schema")
	if !ok {
		t.Fatal("schema not found for 'database schema'")
	}

	if len(schema.Args) != 2 {
		t.Fatalf("expected 2 args, got %d", len(schema.Args))
	}
	if schema.Args[0].Name != "id" {
		t.Errorf("expected first arg 'id', got %s", schema.Args[0].Name)
	}
	if schema.Args[1].Name != "schema" {
		t.Errorf("expected second arg 'schema', got %s", schema.Args[1].Name)
	}
}

func TestGetSchemaUnknownCommand(t *testing.T) {
	_, ok := cli.GetSchema("nonexistent")
	if ok {
		t.Error("expected schema not found for unknown command")
	}
}

func TestSchemaIsValidJSON(t *testing.T) {
	schema, ok := cli.GetSchema("query sql")
	if !ok {
		t.Fatal("schema not found")
	}

	data, err := json.Marshal(schema)
	if err != nil {
		t.Fatalf("failed to marshal schema: %v", err)
	}

	var parsed map[string]any
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("schema JSON is not valid: %v", err)
	}

	if parsed["command"] != "query sql" {
		t.Errorf("expected command 'query sql' in JSON, got %v", parsed["command"])
	}
}
