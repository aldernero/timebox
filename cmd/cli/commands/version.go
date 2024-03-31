package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Timebox",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Timebox v0.1")
	},
}
