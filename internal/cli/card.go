package cli

import (
	"os"
	"strconv"

	"github.com/andreagrandi/mb-cli/internal/formatter"
	"github.com/spf13/cobra"
)

var cardCmd = &cobra.Command{
	Use:   "card",
	Short: "Saved question commands",
}

var cardListCmd = &cobra.Command{
	Use:   "list",
	Short: "List saved questions",
	Args:  cobra.NoArgs,
	RunE:  runCardList,
}

var cardGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get card details",
	Args:  cobra.ExactArgs(1),
	RunE:  runCardGet,
}

var cardRunCmd = &cobra.Command{
	Use:   "run <id>",
	Short: "Execute a saved question",
	Args:  cobra.ExactArgs(1),
	RunE:  runCardRun,
}

func init() {
	rootCmd.AddCommand(cardCmd)

	cardCmd.AddCommand(cardListCmd)
	cardCmd.AddCommand(cardGetCmd)
	cardCmd.AddCommand(cardRunCmd)
}

func runCardList(cmd *cobra.Command, args []string) error {
	c, err := newClient(cmd)
	if err != nil {
		return err
	}

	cards, err := c.ListCards()
	if err != nil {
		return err
	}

	return formatter.Output(cmd, cards)
}

func runCardGet(cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	c, err := newClient(cmd)
	if err != nil {
		return err
	}

	card, err := c.GetCard(id)
	if err != nil {
		return err
	}

	return formatter.Output(cmd, card)
}

func runCardRun(cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	c, err := newClient(cmd)
	if err != nil {
		return err
	}

	result, err := c.RunCard(id)
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
