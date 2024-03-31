package commands

import "github.com/spf13/cobra"

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List boxes or spans",
}

func init() {
	listCmd.AddCommand(listBoxesCmd)
	listCmd.AddCommand(listSpansCmd)
}
