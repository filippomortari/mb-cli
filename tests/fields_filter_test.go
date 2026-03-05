package tests

import (
	"testing"

	"github.com/andreagrandi/mb-cli/internal/formatter"
)

func TestFilterColumnsEmpty(t *testing.T) {
	columns := []string{"id", "name", "email"}
	rows := [][]any{{1, "Alice", "a@b.com"}}

	filteredCols, filteredRows := formatter.FilterColumns(columns, rows, "")

	if len(filteredCols) != 3 {
		t.Errorf("expected 3 columns, got %d", len(filteredCols))
	}
	if len(filteredRows) != 1 || len(filteredRows[0]) != 3 {
		t.Errorf("expected unchanged rows")
	}
}

func TestFilterColumnsSingleField(t *testing.T) {
	columns := []string{"id", "name", "email"}
	rows := [][]any{
		{1, "Alice", "a@b.com"},
		{2, "Bob", "b@b.com"},
	}

	filteredCols, filteredRows := formatter.FilterColumns(columns, rows, "name")

	if len(filteredCols) != 1 || filteredCols[0] != "name" {
		t.Errorf("expected [name], got %v", filteredCols)
	}
	if len(filteredRows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(filteredRows))
	}
	if filteredRows[0][0] != "Alice" {
		t.Errorf("expected Alice, got %v", filteredRows[0][0])
	}
	if filteredRows[1][0] != "Bob" {
		t.Errorf("expected Bob, got %v", filteredRows[1][0])
	}
}

func TestFilterColumnsMultipleFields(t *testing.T) {
	columns := []string{"id", "name", "email", "age"}
	rows := [][]any{
		{1, "Alice", "a@b.com", 30},
		{2, "Bob", "b@b.com", 25},
	}

	filteredCols, filteredRows := formatter.FilterColumns(columns, rows, "id,email")

	if len(filteredCols) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(filteredCols))
	}
	if filteredCols[0] != "id" || filteredCols[1] != "email" {
		t.Errorf("expected [id, email], got %v", filteredCols)
	}
	if filteredRows[0][0] != 1 || filteredRows[0][1] != "a@b.com" {
		t.Errorf("expected [1, a@b.com], got %v", filteredRows[0])
	}
}

func TestFilterColumnsWithSpaces(t *testing.T) {
	columns := []string{"id", "name", "email"}
	rows := [][]any{{1, "Alice", "a@b.com"}}

	filteredCols, filteredRows := formatter.FilterColumns(columns, rows, "id, name")

	if len(filteredCols) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(filteredCols))
	}
	if filteredCols[0] != "id" || filteredCols[1] != "name" {
		t.Errorf("expected [id, name], got %v", filteredCols)
	}
	if filteredRows[0][0] != 1 || filteredRows[0][1] != "Alice" {
		t.Errorf("expected [1, Alice], got %v", filteredRows[0])
	}
}

func TestFilterColumnsNoMatch(t *testing.T) {
	columns := []string{"id", "name", "email"}
	rows := [][]any{{1, "Alice", "a@b.com"}}

	filteredCols, filteredRows := formatter.FilterColumns(columns, rows, "nonexistent")

	// When no fields match, return original data
	if len(filteredCols) != 3 {
		t.Errorf("expected original 3 columns when no match, got %d", len(filteredCols))
	}
	if len(filteredRows[0]) != 3 {
		t.Errorf("expected original row length when no match")
	}
}

func TestFilterColumnsPreservesOrder(t *testing.T) {
	columns := []string{"id", "name", "email", "age"}
	rows := [][]any{{1, "Alice", "a@b.com", 30}}

	// Request in different order than source - should follow source order
	filteredCols, _ := formatter.FilterColumns(columns, rows, "age,id")

	if len(filteredCols) != 2 {
		t.Fatalf("expected 2 columns, got %d", len(filteredCols))
	}
	if filteredCols[0] != "id" || filteredCols[1] != "age" {
		t.Errorf("expected columns in source order [id, age], got %v", filteredCols)
	}
}

func TestFilterColumnsEmptyRows(t *testing.T) {
	columns := []string{"id", "name"}
	var rows [][]any

	filteredCols, filteredRows := formatter.FilterColumns(columns, rows, "id")

	if len(filteredCols) != 1 || filteredCols[0] != "id" {
		t.Errorf("expected [id], got %v", filteredCols)
	}
	if len(filteredRows) != 0 {
		t.Errorf("expected 0 rows, got %d", len(filteredRows))
	}
}
