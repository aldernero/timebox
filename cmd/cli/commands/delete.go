package commands

import "github.com/spf13/cobra"

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a box or span",
}

func init() {
	deleteCmd.AddCommand(deleteBoxCmd)
	deleteCmd.AddCommand(deleteSpanCmd)

	// Delete box command flags
	deleteBoxCmd.Flags().BoolVarP(&cliFlags.force, "force", "f", false, "Force delete")

	// Delete span command flags
	deleteSpanCmd.Flags().BoolVarP(&cliFlags.force, "force", "f", false, "Force delete")
}
