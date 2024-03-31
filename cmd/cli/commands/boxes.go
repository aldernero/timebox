package commands

import (
	"fmt"
	"github.com/aldernero/timebox/pkg/util"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
	"log"
)

var headerStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("205")).
	Background(lipgloss.Color("240")).
	Bold(true)

var evenRowStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("236")).
	Background(lipgloss.Color("238"))

var oddRowStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("237")).
	Background(lipgloss.Color("238"))

var listBoxesCmd = &cobra.Command{
	Use:   "boxes",
	Short: "List boxes",
	Run: func(cmd *cobra.Command, args []string) {
		var rows [][]string
		for k, v := range tb.Boxes {
			minTime := util.DurationParser(v.MinTime)
			maxTime := util.DurationParser(v.MaxTime)
			rows = append(rows, []string{k, minTime, maxTime})
		}
		t := table.New().
			Border(lipgloss.NormalBorder()).
			Headers("Box", "Min", "Max").
			StyleFunc(func(row, col int) lipgloss.Style {
				return lipgloss.NewStyle().Margin(0, 1)
			}).
			Rows(rows...)
		fmt.Println(t.Render())
	},
}

var addBoxCmd = &cobra.Command{
	Use:   "box",
	Short: "Add a new box",
	Run: func(cmd *cobra.Command, args []string) {
		if cliFlags.boxName == "" {
			log.Fatal("box name is required")
		}
		if ok := tb.Boxes[cliFlags.boxName]; ok.Name != "" {
			log.Fatalf("box \"%s\" already exists", cliFlags.boxName)
		}
		if cliFlags.minDuration >= cliFlags.maxDuration {
			log.Fatal("min duration must be less than max duration")
		}
		box := util.Box{
			Name:    cliFlags.boxName,
			MinTime: cliFlags.minDuration,
			MaxTime: cliFlags.maxDuration,
		}
		err := tb.AddBox(box)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var deleteBoxCmd = &cobra.Command{
	Use:   "box",
	Short: "Delete a box",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var updateBoxCmd = &cobra.Command{
	Use:   "box",
	Short: "Update a box",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
