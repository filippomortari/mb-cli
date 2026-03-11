package cli

import (
	"os"
	"strconv"

	"github.com/andreagrandi/mb-cli/internal/client"
	"github.com/andreagrandi/mb-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Dashboard commands",
}

var dashboardListCmd = &cobra.Command{
	Use:   "list",
	Short: "List dashboards",
	Args:  cobra.NoArgs,
	RunE:  runDashboardList,
}

var dashboardGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get dashboard details",
	Args:  cobra.ExactArgs(1),
	RunE:  runDashboardGet,
}

var dashboardCardsCmd = &cobra.Command{
	Use:   "cards <id>",
	Short: "List cards in a dashboard",
	Args:  cobra.ExactArgs(1),
	RunE:  runDashboardCards,
}

type dashboardListRow struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Archived    bool   `json:"archived"`
}

type dashboardCardRow struct {
	DashcardID int    `json:"dashcard_id"`
	Tab        string `json:"tab"`
	CardID     string `json:"card_id,omitempty"`
	Name       string `json:"name,omitempty"`
	QueryType  string `json:"query_type,omitempty"`
	Display    string `json:"display,omitempty"`
}

func init() {
	rootCmd.AddCommand(dashboardCmd)

	dashboardCmd.AddCommand(dashboardListCmd)
	dashboardCmd.AddCommand(dashboardGetCmd)
	dashboardCmd.AddCommand(dashboardCardsCmd)
}

func runDashboardList(cmd *cobra.Command, args []string) error {
	c, err := newClient(cmd)
	if err != nil {
		return err
	}

	dashboards, err := c.ListDashboards()
	if err != nil {
		return err
	}

	rows := make([]dashboardListRow, 0, len(dashboards))
	for _, dashboard := range dashboards {
		rows = append(rows, dashboardListRow{
			ID:          dashboard.ID,
			Name:        dashboard.Name,
			Description: dashboard.Description,
			Archived:    dashboard.Archived,
		})
	}

	return formatter.Output(cmd, rows)
}

func runDashboardGet(cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	c, err := newClient(cmd)
	if err != nil {
		return err
	}

	dashboard, err := c.GetDashboard(id)
	if err != nil {
		return err
	}

	format, _ := cmd.Flags().GetString("format")
	if format == "json" {
		return formatter.Output(cmd, dashboard)
	}

	return formatter.FormatDashboardTable(dashboard, os.Stdout)
}

func runDashboardCards(cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	c, err := newClient(cmd)
	if err != nil {
		return err
	}

	dashboard, err := c.GetDashboard(id)
	if err != nil {
		return err
	}

	return formatter.Output(cmd, buildDashboardCardRows(dashboard))
}

func buildDashboardCardRows(dashboard *client.Dashboard) []dashboardCardRow {
	if dashboard == nil {
		return nil
	}

	tabNames := make(map[int]string, len(dashboard.Tabs))
	for _, tab := range dashboard.Tabs {
		tabNames[tab.ID] = tab.Name
	}

	rows := make([]dashboardCardRow, 0, len(dashboard.DashCards))
	for _, dashCard := range dashboard.DashCards {
		row := dashboardCardRow{
			DashcardID: dashCard.ID,
			Tab:        "Ungrouped",
		}
		if dashCard.TabID != nil {
			if name, ok := tabNames[*dashCard.TabID]; ok {
				row.Tab = name
			}
		}
		if dashCard.CardID != nil {
			row.CardID = strconv.Itoa(*dashCard.CardID)
		}
		if dashCard.Card != nil {
			row.Name = dashCard.Card.Name
			row.QueryType = dashCard.Card.QueryType
			row.Display = dashCard.Card.Display
		}
		rows = append(rows, row)
	}

	return rows
}
