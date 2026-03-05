package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andreagrandi/mb-cli/internal/client"
	"github.com/andreagrandi/mb-cli/internal/config"
)

func setupCardTestClient(handler http.HandlerFunc) (*client.Client, *httptest.Server) {
	server := httptest.NewServer(handler)
	cfg := &config.Config{
		Host:   server.URL,
		APIKey: "test-api-key",
	}
	return client.NewClient(cfg), server
}

func TestListCards(t *testing.T) {
	c, server := setupCardTestClient(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/card/" {
			t.Errorf("expected path '/api/card/', got %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]any{
			{"id": 1, "name": "User Count", "description": "Total users", "database_id": 1, "display": "scalar", "query_type": "native", "archived": false},
			{"id": 2, "name": "Revenue by Month", "description": "Monthly revenue", "database_id": 1, "display": "bar", "query_type": "query", "archived": false},
		})
	})
	defer server.Close()

	cards, err := c.ListCards()
	if err != nil {
		t.Fatalf("ListCards failed: %v", err)
	}

	if len(cards) != 2 {
		t.Fatalf("expected 2 cards, got %d", len(cards))
	}
	if cards[0].ID != 1 {
		t.Errorf("expected ID 1, got %d", cards[0].ID)
	}
	if cards[0].Name != "User Count" {
		t.Errorf("expected name 'User Count', got %s", cards[0].Name)
	}
	if cards[0].Display != "scalar" {
		t.Errorf("expected display 'scalar', got %s", cards[0].Display)
	}
	if cards[1].Name != "Revenue by Month" {
		t.Errorf("expected name 'Revenue by Month', got %s", cards[1].Name)
	}
}

func TestGetCard(t *testing.T) {
	c, server := setupCardTestClient(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/card/1" {
			t.Errorf("expected path '/api/card/1', got %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id": 1, "name": "User Count", "description": "Total users",
			"database_id": 1, "display": "scalar", "query_type": "native", "archived": false,
		})
	})
	defer server.Close()

	card, err := c.GetCard(1)
	if err != nil {
		t.Fatalf("GetCard failed: %v", err)
	}

	if card.ID != 1 {
		t.Errorf("expected ID 1, got %d", card.ID)
	}
	if card.Name != "User Count" {
		t.Errorf("expected name 'User Count', got %s", card.Name)
	}
	if card.DatabaseID != 1 {
		t.Errorf("expected database_id 1, got %d", card.DatabaseID)
	}
}

func TestRunCard(t *testing.T) {
	c, server := setupCardTestClient(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/card/1/query" {
			t.Errorf("expected path '/api/card/1/query', got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"cols": []map[string]any{
					{"name": "count", "display_name": "Count", "base_type": "type/Integer"},
				},
				"rows": [][]any{
					{42},
				},
			},
		})
	})
	defer server.Close()

	result, err := c.RunCard(1)
	if err != nil {
		t.Fatalf("RunCard failed: %v", err)
	}

	if len(result.Data.Columns) != 1 {
		t.Fatalf("expected 1 column, got %d", len(result.Data.Columns))
	}
	if result.Data.Columns[0].Name != "count" {
		t.Errorf("expected column name 'count', got %s", result.Data.Columns[0].Name)
	}
	if len(result.Data.Rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(result.Data.Rows))
	}
}

func TestGetCardNotFound(t *testing.T) {
	c, server := setupCardTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Not found"}`))
	})
	defer server.Close()

	_, err := c.GetCard(999)
	if err == nil {
		t.Fatal("expected error for 404 response")
	}
}

func TestRunCardNotFound(t *testing.T) {
	c, server := setupCardTestClient(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Not found"}`))
	})
	defer server.Close()

	_, err := c.RunCard(999)
	if err == nil {
		t.Fatal("expected error for 404 response")
	}
}
