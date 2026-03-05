package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/andreagrandi/mb-cli/internal/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "mb-cli",
	Short:   "A read-only CLI for the Metabase API",
	Version: version.Version,
	Long: `mb-cli is a read-only command-line interface for querying Metabase databases.
It allows you to list databases, inspect schemas, run SQL queries, and explore
saved questions directly from your terminal.

Before using mb-cli, set your environment variables:
  export MB_HOST=https://your-metabase-instance.com
  export MB_API_KEY=your-api-key`,
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("format") && IsTTY() {
			cmd.Flags().Set("format", "table")
		}

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		errorFormat, _ := rootCmd.PersistentFlags().GetString("error-format")
		if errorFormat == "json" {
			writeJSONError(err)
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("format", "f", "json", "Output format: json, table")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Show request details on stderr")
	rootCmd.PersistentFlags().String("error-format", "text", "Error output format: text, json")
	rootCmd.PersistentFlags().Bool("redact-pii", true, "Redact PII values in query results (disable with --redact-pii=false)")
}

type jsonError struct {
	Error jsonErrorDetail `json:"error"`
}

type jsonErrorDetail struct {
	Type       string `json:"type"`
	Message    string `json:"message"`
	Suggestion string `json:"suggestion,omitempty"`
	ExitCode   int    `json:"exit_code"`
}

// ClassifyError determines the error type and suggestion for structured error output. Exported for testing.
func ClassifyError(err error) (errorType, suggestion string) {
	return classifyError(err)
}

func classifyError(err error) (errorType, suggestion string) {
	msg := err.Error()

	switch {
	case strings.Contains(msg, "MB_HOST") || strings.Contains(msg, "MB_API_KEY"):
		return "CONFIG_ERROR", "Set MB_HOST and MB_API_KEY environment variables"
	case strings.Contains(msg, "API request failed with status 401"),
		strings.Contains(msg, "API request failed with status 403"):
		return "AUTH_ERROR", "Check that MB_API_KEY is valid"
	case strings.Contains(msg, "API request failed with status"):
		return "API_ERROR", ""
	case strings.Contains(msg, "no database matching"),
		strings.Contains(msg, "ambiguous database name"):
		return "RESOLUTION_ERROR", "Use a database ID instead of a name"
	case strings.Contains(msg, "no table matching"),
		strings.Contains(msg, "ambiguous table name"):
		return "RESOLUTION_ERROR", "Use a table ID instead of a name"
	case strings.Contains(msg, "no field matching"):
		return "RESOLUTION_ERROR", "Check field names with 'mb-cli table metadata <id>'"
	default:
		return "GENERAL_ERROR", ""
	}
}

func writeJSONError(err error) {
	errorType, suggestion := classifyError(err)
	je := jsonError{
		Error: jsonErrorDetail{
			Type:       errorType,
			Message:    err.Error(),
			Suggestion: suggestion,
			ExitCode:   1,
		},
	}
	data, _ := json.Marshal(je)
	fmt.Fprintln(os.Stderr, string(data))
}
