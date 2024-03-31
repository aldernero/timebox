package commands

import (
	"github.com/spf13/cobra"
	"log"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a box or span",
}

func init() {
	addCmd.AddCommand(addBoxCmd)
	addCmd.AddCommand(addSpanCmd)

	// Add box command flags
	addBoxCmd.Flags().StringVarP(&cliFlags.boxName, "name", "n", "", "Name of the box")
	addBoxCmd.Flags().DurationVarP(&cliFlags.minDuration, "min", "", 0, "Minimum duration")
	addBoxCmd.Flags().DurationVarP(&cliFlags.maxDuration, "max", "", 0, "Maximum duration")
	requiredFlags := []string{"name", "min", "max"}
	for _, flag := range requiredFlags {
		err := addBoxCmd.MarkFlagRequired(flag)
		if err != nil {
			log.Fatal(err)
		}

	}

	// Add span command flags
	addSpanCmd.Flags().StringVarP(&cliFlags.boxName, "box", "b", "", "Name of the box")
	addSpanCmd.Flags().StringVarP(&cliFlags.startTime, "start", "s", "", "Start time")
	addSpanCmd.Flags().StringVarP(&cliFlags.endTime, "end", "e", "", "End time")
	requiredFlags = []string{"box", "start", "end"}
	for _, flag := range requiredFlags {
		err := addSpanCmd.MarkFlagRequired(flag)
		if err != nil {
			log.Fatal(err)
		}
	}
}
