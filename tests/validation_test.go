package tests

import (
	"strings"
	"testing"

	"github.com/andreagrandi/mb-cli/internal/validation"
)

func TestValidateNoControlChars(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid ascii", "SELECT * FROM users", false},
		{"valid with tab", "col1\tcol2", false},
		{"valid with newline", "line1\nline2", false},
		{"valid with CR", "line1\r\nline2", false},
		{"null byte", "SELECT \x00 FROM users", true},
		{"bell character", "SELECT \x07 FROM users", true},
		{"backspace", "SELECT \x08 FROM users", true},
		{"escape", "SELECT \x1b FROM users", true},
		{"DEL character", "SELECT \x7f FROM users", true},
		{"control char at start", "\x01SELECT 1", true},
		{"empty string", "", false},
		{"unicode valid", "SELECT name FROM users WHERE city = 'Zürich'", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidateNoControlChars(tt.input, "test field")
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNoControlChars() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && !strings.Contains(err.Error(), "control character") {
				t.Errorf("error should mention 'control character', got: %v", err)
			}
		})
	}
}

func TestValidateSQL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{"valid query", "SELECT * FROM users LIMIT 10", false, ""},
		{"empty query", "", false, ""},
		{"control char", "SELECT \x01 FROM users", true, "control character"},
		{"too long", strings.Repeat("x", 10001), true, "too long"},
		{"max length ok", strings.Repeat("x", 10000), false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidateSQL(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSQL() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("error should contain '%s', got: %v", tt.errMsg, err)
			}
		})
	}
}

func TestValidateSearchQuery(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid search", "users table", false},
		{"valid with special chars", "user-data_v2", false},
		{"control char", "users\x00table", true},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidateSearchQuery(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSearchQuery() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
