package cli

import (
	"fmt"
	"os"

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
	SilenceUsage: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("format") && IsTTY() {
			cmd.Flags().Set("format", "table")
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("format", "f", "json", "Output format: json, table")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Show request details on stderr")
}
