package formatter

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// Formatter defines the interface for output formatting.
type Formatter interface {
	Format(data any, writer io.Writer) error
}

// NewFormatter creates a formatter based on the format string.
func NewFormatter(format string) (Formatter, error) {
	switch format {
	case "json":
		return &JSONFormatter{}, nil
	case "table":
		return &TableFormatter{}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// Output reads the --format flag from the command and writes formatted data to stdout.
func Output(cmd *cobra.Command, data any) error {
	format, _ := cmd.Flags().GetString("format")

	f, err := NewFormatter(format)
	if err != nil {
		return err
	}

	return f.Format(data, os.Stdout)
}

// FormatQueryResults formats tabular query results (columns + rows) to the given writer
// using the specified format.
func FormatQueryResults(format string, columns []string, rows [][]any, writer io.Writer) error {
	switch format {
	case "json":
		return formatQueryResultsJSON(columns, rows, writer)
	case "table":
		return formatQueryResultsTable(columns, rows, writer)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// FilterColumns filters columns and rows to include only the specified fields.
// Returns filtered columns and rows. If fields is empty, returns the original data unchanged.
func FilterColumns(columns []string, rows [][]any, fields string) ([]string, [][]any) {
	if fields == "" {
		return columns, rows
	}

	wanted := make(map[string]bool)
	for _, f := range splitFields(fields) {
		wanted[f] = true
	}

	var indices []int
	var filteredCols []string
	for i, col := range columns {
		if wanted[col] {
			indices = append(indices, i)
			filteredCols = append(filteredCols, col)
		}
	}

	if len(indices) == 0 {
		return columns, rows
	}

	filteredRows := make([][]any, len(rows))
	for r, row := range rows {
		filteredRow := make([]any, len(indices))
		for j, idx := range indices {
			if idx < len(row) {
				filteredRow[j] = row[idx]
			}
		}
		filteredRows[r] = filteredRow
	}

	return filteredCols, filteredRows
}

func splitFields(fields string) []string {
	var result []string
	for _, f := range strings.Split(fields, ",") {
		f = strings.TrimSpace(f)
		if f != "" {
			result = append(result, f)
		}
	}
	return result
}
