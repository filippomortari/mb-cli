package cli

import (
	"fmt"

	"github.com/andreagrandi/mb-cli/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of mb-cli",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("mb-cli version %s\n", version.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
