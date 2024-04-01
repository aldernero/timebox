package commands

import "github.com/spf13/cobra"

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List boxes or spans",
}

func init() {
	listCmd.AddCommand(listBoxesCmd)
	listCmd.AddCommand(listSpansCmd)

	listSpansCmd.Flags().StringVarP(&cliFlags.boxName, "box", "b", "", "Name of the box")
	listSpansCmd.Flags().StringVarP(&cliFlags.startTime, "from", "f", "", "Earliest start time")
	listSpansCmd.Flags().StringVarP(&cliFlags.endTime, "to", "t", "", "Latest end time (default: now)")
}
