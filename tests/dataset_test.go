package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/andreagrandi/mb-cli/internal/cli"
	"github.com/andreagrandi/mb-cli/internal/client"
	"github.com/andreagrandi/mb-cli/internal/config"
)

func setupDatasetTestClient(handler http.HandlerFunc) (*client.Client, *httptest.Server) {
	server := httptest.NewServer(handler)
	cfg := &config.Config{
		Host:   server.URL,
		APIKey: "test-api-key",
	}
	return client.NewClient(cfg), server
}

func TestRunNativeQuery(t *testing.T) {
	c, server := setupDatasetTestClient(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/dataset/" {
			t.Errorf("expected path '/api/dataset/', got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}

		body, _ := io.ReadAll(r.Body)
		var req client.DatasetQuery
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("failed to parse request body: %v", err)
		}

		if req.Database != 1 {
			t.Errorf("expected database 1, got %d", req.Database)
		}
		if req.Type != "native" {
			t.Errorf("expected type 'native', got %s", req.Type)
		}
		if req.Native.Query != "SELECT 1" {
			t.Errorf("expected query 'SELECT 1', got %s", req.Native.Query)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"cols": []map[string]any{
					{"name": "id", "display_name": "ID", "base_type": "type/Integer"},
					{"name": "name", "display_name": "Name", "base_type": "type/Text"},
				},
				"rows": [][]any{
					{1, "Alice"},
					{2, "Bob"},
				},
			},
		})
	})
	defer server.Close()

	result, err := c.RunNativeQuery(1, "SELECT 1")
	if err != nil {
		t.Fatalf("RunNativeQuery failed: %v", err)
	}

	if len(result.Data.Columns) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(result.Data.Columns))
	}
	if result.Data.Columns[0].Name != "id" {
		t.Errorf("expected column name 'id', got %s", result.Data.Columns[0].Name)
	}
	if result.Data.Columns[1].Name != "name" {
		t.Errorf("expected column name 'name', got %s", result.Data.Columns[1].Name)
	}
	if len(result.Data.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(result.Data.Rows))
	}
}

func TestExportNativeQuery(t *testing.T) {
	c, server := setupDatasetTestClient(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/dataset/csv" {
			t.Errorf("expected path '/api/dataset/csv', got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}

		body, _ := io.ReadAll(r.Body)
		var wrapper struct {
			Query client.DatasetQuery `json:"query"`
		}
		if err := json.Unmarshal(body, &wrapper); err != nil {
			t.Fatalf("failed to parse request body: %v", err)
		}

		if wrapper.Query.Database != 1 {
			t.Errorf("expected database 1, got %d", wrapper.Query.Database)
		}
		if wrapper.Query.Native.Query != "SELECT id, name FROM users" {
			t.Errorf("expected query 'SELECT id, name FROM users', got %s", wrapper.Query.Native.Query)
		}

		w.Header().Set("Content-Type", "text/csv")
		w.Write([]byte("id,name\n1,Alice\n2,Bob\n"))
	})
	defer server.Close()

	data, err := c.ExportNativeQuery(1, "SELECT id, name FROM users", "csv")
	if err != nil {
		t.Fatalf("ExportNativeQuery failed: %v", err)
	}

	expected := "id,name\n1,Alice\n2,Bob\n"
	if string(data) != expected {
		t.Errorf("expected %q, got %q", expected, string(data))
	}
}

func TestRunNativeQueryError(t *testing.T) {
	c, server := setupDatasetTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid SQL query"}`))
	})
	defer server.Close()

	_, err := c.RunNativeQuery(1, "INVALID SQL")
	if err == nil {
		t.Fatal("expected error for bad query")
	}
}

func TestMatchDatabaseByNameExactSubstring(t *testing.T) {
	databases := []client.Database{
		{ID: 1, Name: "Production"},
		{ID: 2, Name: "Analytics"},
		{ID: 3, Name: "Staging"},
	}

	id, err := cli.MatchDatabaseByName(databases, "prod")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id != 1 {
		t.Errorf("expected 1, got %d", id)
	}
}

func TestMatchDatabaseByNameCaseInsensitive(t *testing.T) {
	databases := []client.Database{
		{ID: 1, Name: "Production"},
		{ID: 2, Name: "Analytics"},
	}

	id, err := cli.MatchDatabaseByName(databases, "ANALYTICS")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id != 2 {
		t.Errorf("expected 2, got %d", id)
	}
}

func TestMatchDatabaseByNameNoMatch(t *testing.T) {
	databases := []client.Database{
		{ID: 1, Name: "Production"},
	}

	_, err := cli.MatchDatabaseByName(databases, "nonexistent")
	if err == nil {
		t.Fatal("expected error for no match")
	}
	expected := "no database matching 'nonexistent' found"
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}

func TestMatchDatabaseByNameAmbiguous(t *testing.T) {
	databases := []client.Database{
		{ID: 1, Name: "Production DB"},
		{ID: 2, Name: "Production Replica"},
	}

	_, err := cli.MatchDatabaseByName(databases, "production")
	if err == nil {
		t.Fatal("expected error for ambiguous match")
	}
	if !strings.Contains(err.Error(), "ambiguous database name") {
		t.Errorf("expected ambiguous error, got %q", err.Error())
	}
	if !strings.Contains(err.Error(), "Production DB (id=1)") {
		t.Errorf("expected match list in error, got %q", err.Error())
	}
}
