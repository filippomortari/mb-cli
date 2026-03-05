package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andreagrandi/mb-cli/internal/client"
	"github.com/andreagrandi/mb-cli/internal/config"
)

func TestIsPIISemanticType(t *testing.T) {
	tests := []struct {
		semanticType string
		expected     bool
	}{
		{"type/Email", true},
		{"type/Name", true},
		{"type/Phone", true},
		{"type/Address", true},
		{"type/City", true},
		{"type/State", true},
		{"type/ZipCode", true},
		{"type/Country", true},
		{"type/Latitude", true},
		{"type/Longitude", true},
		{"type/Birthdate", true},
		{"type/AvatarURL", true},
		{"type/URL", true},
		{"type/ImageURL", true},
		{"type/Company", true},
		{"type/FK", false},
		{"type/PK", false},
		{"type/Category", false},
		{"type/Number", false},
		{"type/Description", false},
		{"", false},
		{"type/Unknown", false},
	}

	for _, tt := range tests {
		t.Run(tt.semanticType, func(t *testing.T) {
			result := client.IsPIISemanticType(tt.semanticType)
			if result != tt.expected {
				t.Errorf("IsPIISemanticType(%q) = %v, want %v", tt.semanticType, result, tt.expected)
			}
		})
	}
}

func TestRedactQueryResult(t *testing.T) {
	tests := []struct {
		name     string
		input    client.QueryResult
		expected [][]any
	}{
		{
			name: "mixed PII and non-PII columns",
			input: client.QueryResult{
				Data: client.QueryResultData{
					Columns: []client.ResultColumn{
						{Name: "id", BaseType: "type/Integer"},
						{Name: "email", BaseType: "type/Text", SemanticType: "type/Email"},
						{Name: "name", BaseType: "type/Text", SemanticType: "type/Name"},
					},
					Rows: [][]any{
						{1, "alice@example.com", "Alice"},
						{2, "bob@example.com", "Bob"},
					},
				},
			},
			expected: [][]any{
				{1, client.RedactedValue, client.RedactedValue},
				{2, client.RedactedValue, client.RedactedValue},
			},
		},
		{
			name: "no PII columns",
			input: client.QueryResult{
				Data: client.QueryResultData{
					Columns: []client.ResultColumn{
						{Name: "id", BaseType: "type/Integer"},
						{Name: "count", BaseType: "type/Integer", SemanticType: "type/Number"},
					},
					Rows: [][]any{
						{1, 42},
						{2, 99},
					},
				},
			},
			expected: [][]any{
				{1, 42},
				{2, 99},
			},
		},
		{
			name: "all PII columns",
			input: client.QueryResult{
				Data: client.QueryResultData{
					Columns: []client.ResultColumn{
						{Name: "email", BaseType: "type/Text", SemanticType: "type/Email"},
						{Name: "phone", BaseType: "type/Text", SemanticType: "type/Phone"},
					},
					Rows: [][]any{
						{"alice@example.com", "+1234567890"},
					},
				},
			},
			expected: [][]any{
				{client.RedactedValue, client.RedactedValue},
			},
		},
		{
			name: "empty rows",
			input: client.QueryResult{
				Data: client.QueryResultData{
					Columns: []client.ResultColumn{
						{Name: "email", BaseType: "type/Text", SemanticType: "type/Email"},
					},
					Rows: [][]any{},
				},
			},
			expected: [][]any{},
		},
		{
			name: "nil values in PII columns",
			input: client.QueryResult{
				Data: client.QueryResultData{
					Columns: []client.ResultColumn{
						{Name: "id", BaseType: "type/Integer"},
						{Name: "email", BaseType: "type/Text", SemanticType: "type/Email"},
					},
					Rows: [][]any{
						{1, nil},
					},
				},
			},
			expected: [][]any{
				{1, client.RedactedValue},
			},
		},
		{
			name: "columns with empty semantic type",
			input: client.QueryResult{
				Data: client.QueryResultData{
					Columns: []client.ResultColumn{
						{Name: "data", BaseType: "type/Text"},
					},
					Rows: [][]any{
						{"some data"},
					},
				},
			},
			expected: [][]any{
				{"some data"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client.RedactQueryResult(&tt.input)
			if len(tt.input.Data.Rows) != len(tt.expected) {
				t.Fatalf("row count = %d, want %d", len(tt.input.Data.Rows), len(tt.expected))
			}
			for i, row := range tt.input.Data.Rows {
				for j, val := range row {
					if fmt.Sprintf("%v", val) != fmt.Sprintf("%v", tt.expected[i][j]) {
						t.Errorf("row[%d][%d] = %v, want %v", i, j, val, tt.expected[i][j])
					}
				}
			}
		})
	}
}

func TestRedactFieldValues(t *testing.T) {
	tests := []struct {
		name     string
		input    client.FieldValues
		expected [][]any
	}{
		{
			name: "basic redaction",
			input: client.FieldValues{
				FieldID: 1,
				Values:  [][]any{{"alice@example.com"}, {"bob@example.com"}},
			},
			expected: [][]any{{client.RedactedValue}, {client.RedactedValue}},
		},
		{
			name: "empty values",
			input: client.FieldValues{
				FieldID: 1,
				Values:  [][]any{},
			},
			expected: [][]any{},
		},
		{
			name: "multi-element value arrays",
			input: client.FieldValues{
				FieldID: 1,
				Values:  [][]any{{"alice@example.com", "Alice"}, {"bob@example.com", "Bob"}},
			},
			expected: [][]any{{client.RedactedValue, client.RedactedValue}, {client.RedactedValue, client.RedactedValue}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client.RedactFieldValues(&tt.input)
			if len(tt.input.Values) != len(tt.expected) {
				t.Fatalf("values count = %d, want %d", len(tt.input.Values), len(tt.expected))
			}
			for i, vals := range tt.input.Values {
				for j, val := range vals {
					if fmt.Sprintf("%v", val) != fmt.Sprintf("%v", tt.expected[i][j]) {
						t.Errorf("values[%d][%d] = %v, want %v", i, j, val, tt.expected[i][j])
					}
				}
			}
		})
	}
}

func TestRedactPIIIntegration(t *testing.T) {
	queryResponse := client.QueryResult{
		Data: client.QueryResultData{
			Columns: []client.ResultColumn{
				{Name: "id", DisplayName: "ID", BaseType: "type/Integer"},
				{Name: "email", DisplayName: "Email", BaseType: "type/Text", SemanticType: "type/Email"},
				{Name: "name", DisplayName: "Name", BaseType: "type/Text", SemanticType: "type/Name"},
			},
			Rows: [][]any{
				{float64(1), "alice@example.com", "Alice"},
				{float64(2), "bob@example.com", "Bob"},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(queryResponse)
	}))
	defer server.Close()

	t.Run("redaction enabled", func(t *testing.T) {
		c := client.NewClient(&config.Config{Host: server.URL, APIKey: "test"})
		c.RedactPII = true

		result, err := c.RunNativeQuery(1, "SELECT id, email, name FROM users")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		for _, row := range result.Data.Rows {
			if row[1] != client.RedactedValue {
				t.Errorf("email not redacted: got %v", row[1])
			}
			if row[2] != client.RedactedValue {
				t.Errorf("name not redacted: got %v", row[2])
			}
		}

		if result.Data.Rows[0][0] != float64(1) {
			t.Errorf("id should not be redacted: got %v", result.Data.Rows[0][0])
		}
	})

	t.Run("redaction disabled", func(t *testing.T) {
		c := client.NewClient(&config.Config{Host: server.URL, APIKey: "test"})
		c.RedactPII = false

		result, err := c.RunNativeQuery(1, "SELECT id, email, name FROM users")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.Data.Rows[0][1] != "alice@example.com" {
			t.Errorf("email should not be redacted: got %v", result.Data.Rows[0][1])
		}
	})

	t.Run("export blocked when redaction enabled", func(t *testing.T) {
		c := client.NewClient(&config.Config{Host: server.URL, APIKey: "test"})
		c.RedactPII = true

		_, err := c.ExportNativeQuery(1, "SELECT 1", "csv")
		if err == nil {
			t.Fatal("expected error for export with redaction enabled")
		}

		expectedMsg := "export is not supported when PII redaction is enabled"
		if err.Error() != expectedMsg+" (use JSON or table format instead)" {
			t.Errorf("unexpected error message: %v", err)
		}
	})

	t.Run("export allowed when redaction disabled", func(t *testing.T) {
		exportServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/csv")
			w.Write([]byte("id,email\n1,alice@example.com\n"))
		}))
		defer exportServer.Close()

		c := client.NewClient(&config.Config{Host: exportServer.URL, APIKey: "test"})
		c.RedactPII = false

		data, err := c.ExportNativeQuery(1, "SELECT 1", "csv")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(data) == 0 {
			t.Error("expected non-empty export data")
		}
	})
}
