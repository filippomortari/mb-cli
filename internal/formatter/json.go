package formatter

import (
	"encoding/json"
	"io"
)

// JSONFormatter outputs data as pretty-printed JSON.
type JSONFormatter struct{}

func (f *JSONFormatter) Format(data any, writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// formatQueryResultsJSON formats query results as an array of row objects keyed by column name.
func formatQueryResultsJSON(columns []string, rows [][]any, writer io.Writer) error {
	result := make([]map[string]any, 0, len(rows))

	for _, row := range rows {
		obj := make(map[string]any, len(columns))
		for i, col := range columns {
			if i < len(row) {
				obj[col] = row[i]
			}
		}
		result = append(result, obj)
	}

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}
