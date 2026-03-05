package tests

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/andreagrandi/mb-cli/internal/formatter"
)

func TestNewFormatterJSON(t *testing.T) {
	f, err := formatter.NewFormatter("json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil formatter")
	}
}

func TestNewFormatterTable(t *testing.T) {
	f, err := formatter.NewFormatter("table")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil formatter")
	}
}

func TestNewFormatterUnsupported(t *testing.T) {
	_, err := formatter.NewFormatter("xml")
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
	if !strings.Contains(err.Error(), "unsupported format") {
		t.Fatalf("expected 'unsupported format' error, got: %v", err)
	}
}

// --- JSON Formatter Tests ---

func TestJSONFormatMap(t *testing.T) {
	f, _ := formatter.NewFormatter("json")
	var buf bytes.Buffer

	data := map[string]string{"name": "test_db", "engine": "postgres"}
	err := f.Format(data, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result map[string]string
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if result["name"] != "test_db" {
		t.Fatalf("expected name=test_db, got %s", result["name"])
	}
}

func TestJSONFormatSlice(t *testing.T) {
	f, _ := formatter.NewFormatter("json")
	var buf bytes.Buffer

	data := []map[string]int{{"id": 1}, {"id": 2}}
	err := f.Format(data, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result []map[string]int
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 items, got %d", len(result))
	}
}

func TestJSONFormatStruct(t *testing.T) {
	type DB struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	f, _ := formatter.NewFormatter("json")
	var buf bytes.Buffer

	err := f.Format(DB{ID: 1, Name: "prod"}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if result["name"] != "prod" {
		t.Fatalf("expected name=prod, got %v", result["name"])
	}
}

func TestJSONFormatPrettyPrinted(t *testing.T) {
	f, _ := formatter.NewFormatter("json")
	var buf bytes.Buffer

	err := f.Format(map[string]string{"key": "value"}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "\n") {
		t.Fatal("expected pretty-printed JSON with newlines")
	}
	if !strings.Contains(output, "  ") {
		t.Fatal("expected indented JSON output")
	}
}

// --- Table Formatter Tests ---

func TestTableFormatStruct(t *testing.T) {
	type DB struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	f, _ := formatter.NewFormatter("table")
	var buf bytes.Buffer

	err := f.Format(DB{ID: 1, Name: "prod"}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "id") {
		t.Fatalf("expected 'id' in output, got: %s", output)
	}
	if !strings.Contains(output, "prod") {
		t.Fatalf("expected 'prod' in output, got: %s", output)
	}
}

func TestTableFormatStructSlice(t *testing.T) {
	type DB struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	f, _ := formatter.NewFormatter("table")
	var buf bytes.Buffer

	data := []DB{{ID: 1, Name: "prod"}, {ID: 2, Name: "staging"}}
	err := f.Format(data, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	// Should have header row with field names
	if !strings.Contains(output, "id") {
		t.Fatalf("expected 'id' header, got: %s", output)
	}
	if !strings.Contains(output, "name") {
		t.Fatalf("expected 'name' header, got: %s", output)
	}
	// Should have data rows
	if !strings.Contains(output, "prod") {
		t.Fatalf("expected 'prod' in output, got: %s", output)
	}
	if !strings.Contains(output, "staging") {
		t.Fatalf("expected 'staging' in output, got: %s", output)
	}
}

func TestTableFormatNil(t *testing.T) {
	f, _ := formatter.NewFormatter("table")
	var buf bytes.Buffer

	err := f.Format(nil, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No data") {
		t.Fatalf("expected 'No data' for nil input, got: %s", buf.String())
	}
}

func TestTableFormatEmptySlice(t *testing.T) {
	type DB struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	f, _ := formatter.NewFormatter("table")
	var buf bytes.Buffer

	err := f.Format([]DB{}, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No data") {
		t.Fatalf("expected 'No data' for empty slice, got: %s", buf.String())
	}
}

// --- Query Results Tests ---

func TestFormatQueryResultsJSON(t *testing.T) {
	var buf bytes.Buffer
	columns := []string{"id", "name", "email"}
	rows := [][]any{
		{1, "Alice", "alice@example.com"},
		{2, "Bob", "bob@example.com"},
	}

	err := formatter.FormatQueryResults("json", columns, rows, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result []map[string]any
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(result))
	}
	if result[0]["name"] != "Alice" {
		t.Fatalf("expected name=Alice, got %v", result[0]["name"])
	}
	if result[1]["email"] != "bob@example.com" {
		t.Fatalf("expected email=bob@example.com, got %v", result[1]["email"])
	}
}

func TestFormatQueryResultsTable(t *testing.T) {
	var buf bytes.Buffer
	columns := []string{"id", "name"}
	rows := [][]any{
		{1, "Alice"},
		{2, "Bob"},
	}

	err := formatter.FormatQueryResults("table", columns, rows, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "id") {
		t.Fatalf("expected 'id' column header, got: %s", output)
	}
	if !strings.Contains(output, "name") {
		t.Fatalf("expected 'name' column header, got: %s", output)
	}
	if !strings.Contains(output, "Alice") {
		t.Fatalf("expected 'Alice' in output, got: %s", output)
	}
	if !strings.Contains(output, "Bob") {
		t.Fatalf("expected 'Bob' in output, got: %s", output)
	}
}

func TestFormatQueryResultsEmptyRows(t *testing.T) {
	var buf bytes.Buffer
	columns := []string{"id", "name"}
	rows := [][]any{}

	err := formatter.FormatQueryResults("json", columns, rows, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result []map[string]any
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected empty array, got %d items", len(result))
	}
}

func TestFormatQueryResultsUnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	err := formatter.FormatQueryResults("xml", []string{"id"}, [][]any{{1}}, &buf)
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestFormatQueryResultsNilValues(t *testing.T) {
	var buf bytes.Buffer
	columns := []string{"id", "name"}
	rows := [][]any{
		{1, nil},
	}

	err := formatter.FormatQueryResults("json", columns, rows, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result []map[string]any
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if result[0]["name"] != nil {
		t.Fatalf("expected nil value for name, got %v", result[0]["name"])
	}
}

func TestFormatQueryResultsTableNilValues(t *testing.T) {
	var buf bytes.Buffer
	columns := []string{"id", "name"}
	rows := [][]any{
		{1, nil},
	}

	err := formatter.FormatQueryResults("table", columns, rows, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "id") {
		t.Fatalf("expected 'id' header, got: %s", output)
	}
}
