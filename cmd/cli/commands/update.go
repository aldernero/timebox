package commands

import (
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a box or span",
}

func init() {
	updateCmd.AddCommand(updateBoxCmd)
	updateCmd.AddCommand(updateSpanCmd)

	// Update box command flags
	updateBoxCmd.Flags().DurationVarP(&cliFlags.minDuration, "min", "", 0, "Minimum duration")
	updateBoxCmd.Flags().DurationVarP(&cliFlags.maxDuration, "max", "", 0, "Maximum duration")
	updateBoxCmd.MarkFlagsOneRequired("min", "max")

	// Update span command flags
	updateSpanCmd.Flags().StringVarP(&cliFlags.boxName, "box", "b", "", "Name of the box")
	updateSpanCmd.Flags().StringVarP(&cliFlags.startTime, "start", "s", "", "Start time")
	updateSpanCmd.Flags().StringVarP(&cliFlags.endTime, "end", "e", "", "End time")
	updateSpanCmd.MarkFlagsOneRequired("box", "start", "end")
}
