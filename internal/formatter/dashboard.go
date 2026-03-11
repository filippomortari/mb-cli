package formatter

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/andreagrandi/mb-cli/internal/client"
)

const ungroupedDashboardTab = "Ungrouped"

// FormatDashboardTable renders a dashboard in a human-readable table layout.
func FormatDashboardTable(dashboard *client.Dashboard, writer io.Writer) error {
	if dashboard == nil {
		_, err := fmt.Fprintln(writer, "No data")
		return err
	}

	tw := tabwriter.NewWriter(writer, 0, 4, 2, ' ', 0)
	fmt.Fprintf(tw, "id\t%d\n", dashboard.ID)
	fmt.Fprintf(tw, "name\t%s\n", dashboard.Name)
	fmt.Fprintf(tw, "description\t%s\n", dashboard.Description)
	fmt.Fprintf(tw, "archived\t%t\n", dashboard.Archived)
	if err := tw.Flush(); err != nil {
		return err
	}

	if _, err := fmt.Fprintln(writer); err != nil {
		return err
	}

	if err := formatDashboardTabs(dashboard.Tabs, writer); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(writer); err != nil {
		return err
	}

	if err := formatDashboardParameters(dashboard.Parameters, writer); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(writer); err != nil {
		return err
	}

	return formatDashboardCardsByTab(dashboard, writer)
}

func formatDashboardTabs(tabs []client.DashTab, writer io.Writer) error {
	if _, err := fmt.Fprintln(writer, "Tabs"); err != nil {
		return err
	}

	tw := tabwriter.NewWriter(writer, 0, 4, 2, ' ', 0)
	if len(tabs) == 0 {
		if _, err := fmt.Fprintln(tw, "name\tUngrouped"); err != nil {
			return err
		}
		return tw.Flush()
	}

	if _, err := fmt.Fprintln(tw, "id\tname"); err != nil {
		return err
	}
	for _, tab := range tabs {
		if _, err := fmt.Fprintf(tw, "%d\t%s\n", tab.ID, tab.Name); err != nil {
			return err
		}
	}

	return tw.Flush()
}

func formatDashboardParameters(params []client.DashParameter, writer io.Writer) error {
	if _, err := fmt.Fprintln(writer, "Parameters"); err != nil {
		return err
	}

	tw := tabwriter.NewWriter(writer, 0, 4, 2, ' ', 0)
	if _, err := fmt.Fprintln(tw, "id\tname\tslug\ttype"); err != nil {
		return err
	}
	if len(params) == 0 {
		if _, err := fmt.Fprintln(tw, "-\t-\t-\t-"); err != nil {
			return err
		}
		return tw.Flush()
	}

	for _, param := range params {
		if _, err := fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", param.ID, param.Name, param.Slug, param.Type); err != nil {
			return err
		}
	}

	return tw.Flush()
}

func formatDashboardCardsByTab(dashboard *client.Dashboard, writer io.Writer) error {
	if _, err := fmt.Fprintln(writer, "Cards"); err != nil {
		return err
	}

	tabNames := make(map[int]string, len(dashboard.Tabs))
	orderedTabs := make([]string, 0, len(dashboard.Tabs)+1)
	for _, tab := range dashboard.Tabs {
		tabNames[tab.ID] = tab.Name
		orderedTabs = append(orderedTabs, tab.Name)
	}

	grouped := make(map[string][]client.DashCard)
	for _, dashCard := range dashboard.DashCards {
		groupName := ungroupedDashboardTab
		if dashCard.TabID != nil {
			if name, ok := tabNames[*dashCard.TabID]; ok {
				groupName = name
			}
		}
		grouped[groupName] = append(grouped[groupName], dashCard)
	}

	if _, ok := grouped[ungroupedDashboardTab]; ok || len(orderedTabs) == 0 {
		orderedTabs = append(orderedTabs, ungroupedDashboardTab)
	}

	seen := make(map[string]bool, len(orderedTabs))
	for _, tabName := range orderedTabs {
		if seen[tabName] {
			continue
		}
		seen[tabName] = true
		cards := grouped[tabName]
		if len(cards) == 0 {
			continue
		}

		if _, err := fmt.Fprintf(writer, "[%s]\n", tabName); err != nil {
			return err
		}

		tw := tabwriter.NewWriter(writer, 0, 4, 2, ' ', 0)
		if _, err := fmt.Fprintln(tw, "dashcard_id\tcard_id\tname\tquery_type\tdisplay"); err != nil {
			return err
		}
		for _, dashCard := range cards {
			cardID := ""
			if dashCard.CardID != nil {
				cardID = fmt.Sprintf("%d", *dashCard.CardID)
			}

			name := ""
			queryType := ""
			display := ""
			if dashCard.Card != nil {
				name = dashCard.Card.Name
				queryType = dashCard.Card.QueryType
				display = dashCard.Card.Display
			}

			if _, err := fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%s\n", dashCard.ID, cardID, name, queryType, display); err != nil {
				return err
			}
		}
		if err := tw.Flush(); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(writer); err != nil {
			return err
		}
	}

	if len(grouped) == 0 {
		tw := tabwriter.NewWriter(writer, 0, 4, 2, ' ', 0)
		if _, err := fmt.Fprintln(tw, "dashcard_id\tcard_id\tname\tquery_type\tdisplay"); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(tw, "-\t-\t-\t-\t-"); err != nil {
			return err
		}
		return tw.Flush()
	}

	leftovers := make([]string, 0, len(grouped))
	for tabName := range grouped {
		if !seen[tabName] {
			leftovers = append(leftovers, tabName)
		}
	}
	sort.Strings(leftovers)
	for _, tabName := range leftovers {
		if _, err := fmt.Fprintf(writer, "[%s]\n", tabName); err != nil {
			return err
		}
		tw := tabwriter.NewWriter(writer, 0, 4, 2, ' ', 0)
		if _, err := fmt.Fprintln(tw, "dashcard_id\tcard_id\tname\tquery_type\tdisplay"); err != nil {
			return err
		}
		for _, dashCard := range grouped[tabName] {
			cardID := ""
			if dashCard.CardID != nil {
				cardID = fmt.Sprintf("%d", *dashCard.CardID)
			}

			name := ""
			queryType := ""
			display := ""
			if dashCard.Card != nil {
				name = dashCard.Card.Name
				queryType = dashCard.Card.QueryType
				display = dashCard.Card.Display
			}

			if _, err := fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%s\n", dashCard.ID, cardID, name, queryType, display); err != nil {
				return err
			}
		}
		if err := tw.Flush(); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(writer); err != nil {
			return err
		}
	}

	return nil
}
