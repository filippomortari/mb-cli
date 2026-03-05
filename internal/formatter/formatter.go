package formatter

import (
	"fmt"
	"io"
	"os"

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
