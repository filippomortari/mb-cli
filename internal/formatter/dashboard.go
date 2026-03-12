package formatter

import (
	"encoding/json"
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
		if err := renderDashboardCardsTable(cards, writer); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(writer); err != nil {
			return err
		}
	}

	if len(grouped) == 0 {
		return renderDashboardCardsTable(nil, writer)
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
		if err := renderDashboardCardsTable(grouped[tabName], writer); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(writer); err != nil {
			return err
		}
	}

	return nil
}

// FormatDashboardParameterValuesTable renders parameter metadata, mappings, and values.
func FormatDashboardParameterValuesTable(dashboard *client.Dashboard, parameter *client.DashParameter, values *client.ParameterValues, writer io.Writer) error {
	tw := tabwriter.NewWriter(writer, 0, 4, 2, ' ', 0)
	fmt.Fprintf(tw, "dashboard_id\t%d\n", dashboard.ID)
	if parameter != nil {
		fmt.Fprintf(tw, "parameter_id\t%s\n", parameter.ID)
		fmt.Fprintf(tw, "parameter_name\t%s\n", parameter.Name)
		fmt.Fprintf(tw, "parameter_slug\t%s\n", parameter.Slug)
		fmt.Fprintf(tw, "parameter_type\t%s\n", parameter.Type)
	}
	fmt.Fprintf(tw, "has_more_values\t%t\n", values.HasMoreValues)
	if err := tw.Flush(); err != nil {
		return err
	}

	if _, err := fmt.Fprintln(writer); err != nil {
		return err
	}

	if err := formatDashboardParameterMappings(dashboard, parameter, writer); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(writer); err != nil {
		return err
	}

	if _, err := fmt.Fprintln(writer, "Values"); err != nil {
		return err
	}
	tw = tabwriter.NewWriter(writer, 0, 4, 2, ' ', 0)
	if _, err := fmt.Fprintln(tw, "value\tlabel"); err != nil {
		return err
	}
	if len(values.Values) == 0 {
		if _, err := fmt.Fprintln(tw, "-\t-"); err != nil {
			return err
		}
		return tw.Flush()
	}
	for _, value := range values.Values {
		if _, err := fmt.Fprintf(tw, "%s\t%s\n", stringify(value.Value), value.Label); err != nil {
			return err
		}
	}
	return tw.Flush()
}

func formatDashboardParameterMappings(dashboard *client.Dashboard, parameter *client.DashParameter, writer io.Writer) error {
	if _, err := fmt.Fprintln(writer, "Mapped Cards"); err != nil {
		return err
	}

	tw := tabwriter.NewWriter(writer, 0, 4, 2, ' ', 0)
	if _, err := fmt.Fprintln(tw, "dashcard_id\tcard_id\tcard_name\ttarget"); err != nil {
		return err
	}
	if parameter == nil {
		if _, err := fmt.Fprintln(tw, "-\t-\t-\t-"); err != nil {
			return err
		}
		return tw.Flush()
	}

	found := false
	for _, dashCard := range dashboard.DashCards {
		for _, mapping := range dashCard.ParameterMappings {
			if mapping.ParameterID != parameter.ID {
				continue
			}

			found = true
			cardID := ""
			if dashCard.CardID != nil {
				cardID = fmt.Sprintf("%d", *dashCard.CardID)
			}
			cardName := ""
			if dashCard.Card != nil {
				cardName = dashCard.Card.Name
			}
			target := stringifyMappingTarget(mapping.Target)
			if _, err := fmt.Fprintf(tw, "%d\t%s\t%s\t%s\n", dashCard.ID, cardID, cardName, target); err != nil {
				return err
			}
		}
	}

	if !found {
		if _, err := fmt.Fprintln(tw, "-\t-\t-\t-"); err != nil {
			return err
		}
	}

	return tw.Flush()
}

func renderDashboardCardsTable(cards []client.DashCard, writer io.Writer) error {
	tw := tabwriter.NewWriter(writer, 0, 4, 2, ' ', 0)
	if _, err := fmt.Fprintln(tw, "dashcard_id\tcard_id\tname\tquery_type\tdisplay"); err != nil {
		return err
	}
	if len(cards) == 0 {
		if _, err := fmt.Fprintln(tw, "-\t-\t-\t-\t-"); err != nil {
			return err
		}
		return tw.Flush()
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

	return tw.Flush()
}

func stringifyMappingTarget(target []any) string {
	if len(target) == 0 {
		return ""
	}

	data, err := json.Marshal(target)
	if err != nil {
		return fmt.Sprintf("%v", target)
	}

	return string(data)
}
