package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/andreagrandi/mb-cli/internal/client"
	"github.com/andreagrandi/mb-cli/internal/formatter"
	"github.com/andreagrandi/mb-cli/internal/validation"
	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query commands",
}

var querySQLCmd = &cobra.Command{
	Use:   "sql",
	Short: "Run a native SQL query",
	Args:  cobra.NoArgs,
	RunE:  runQuerySQL,
}

func init() {
	rootCmd.AddCommand(queryCmd)
	queryCmd.AddCommand(querySQLCmd)

	querySQLCmd.Flags().String("db", "", "Database ID or name substring (required)")
	querySQLCmd.Flags().String("sql", "", "SQL query to execute (required)")
	querySQLCmd.Flags().String("export", "", "Export format: csv, json, xlsx")
	querySQLCmd.Flags().Int("limit", 0, "Append LIMIT to SQL query")
	querySQLCmd.MarkFlagRequired("db")
	querySQLCmd.MarkFlagRequired("sql")
}

func runQuerySQL(cmd *cobra.Command, args []string) error {
	c, err := newClient(cmd)
	if err != nil {
		return err
	}

	dbFlag, _ := cmd.Flags().GetString("db")
	sql, _ := cmd.Flags().GetString("sql")
	export, _ := cmd.Flags().GetString("export")
	limit, _ := cmd.Flags().GetInt("limit")

	if err := validation.ValidateSQL(sql); err != nil {
		return err
	}

	dbID, err := resolveDatabaseID(c, dbFlag)
	if err != nil {
		return err
	}

	if limit > 0 {
		sql = fmt.Sprintf("%s LIMIT %d", sql, limit)
	}

	if export != "" {
		data, err := c.ExportNativeQuery(dbID, sql, export)
		if err != nil {
			return err
		}
		_, err = os.Stdout.Write(data)
		return err
	}

	result, err := c.RunNativeQuery(dbID, sql)
	if err != nil {
		return err
	}

	format, _ := cmd.Flags().GetString("format")

	columns := make([]string, len(result.Data.Columns))
	for i, col := range result.Data.Columns {
		columns[i] = col.Name
	}

	return formatter.FormatQueryResults(format, columns, result.Data.Rows, os.Stdout)
}

func resolveDatabaseID(c *client.Client, dbFlag string) (int, error) {
	if id, err := strconv.Atoi(dbFlag); err == nil {
		return id, nil
	}

	databases, err := c.ListDatabases(false)
	if err != nil {
		return 0, fmt.Errorf("failed to list databases for name resolution: %w", err)
	}

	return matchDatabaseByName(databases, dbFlag)
}

// MatchDatabaseByName finds a database by case-insensitive substring match.
// Exported for testing.
func MatchDatabaseByName(databases []client.Database, name string) (int, error) {
	return matchDatabaseByName(databases, name)
}

func matchDatabaseByName(databases []client.Database, name string) (int, error) {
	var matches []client.Database
	search := strings.ToLower(name)
	for _, db := range databases {
		if strings.Contains(strings.ToLower(db.Name), search) {
			matches = append(matches, db)
		}
	}

	switch len(matches) {
	case 0:
		return 0, fmt.Errorf("no database matching '%s' found", name)
	case 1:
		return matches[0].ID, nil
	default:
		names := make([]string, len(matches))
		for i, db := range matches {
			names[i] = fmt.Sprintf("%s (id=%d)", db.Name, db.ID)
		}
		return 0, fmt.Errorf("ambiguous database name '%s', matches: %s. Use database ID instead", name, strings.Join(names, ", "))
	}
}
