package commands

import (
	"fmt"
	"github.com/aldernero/timebox/pkg/util"
	"github.com/charmbracelet/huh"
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
		for _, name := range tb.Names {
			box := tb.Boxes[name]
			minTime := util.DurationParser(box.MinTime)
			maxTime := util.DurationParser(box.MaxTime)
			rows = append(rows, []string{name, minTime, maxTime})
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
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		box := args[0]
		_, ok := tb.Boxes[box]
		if !ok {
			log.Fatalf("box \"%s\" does not exist", box)
		}
		var confirmed bool
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Are you sure you want to delete the box?").
					Affirmative("Yes").
					Negative("No").
					Value(&confirmed),
			),
		)
		if !cliFlags.force {
			err := form.Run()
			if err != nil {
				log.Fatal(err)
			}
			if !confirmed {
				fmt.Println("Cancelling box delete")
				return
			}
		}
		err := tb.DeleteBoxAndSpans(box)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var updateBoxCmd = &cobra.Command{
	Use:   "box",
	Short: "Update a box",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		boxName := args[0]
		box, ok := tb.Boxes[boxName]
		if !ok {
			log.Fatalf("box \"%s\" does not exist", boxName)
		}
		if cliFlags.maxDuration > 0 && (cliFlags.minDuration >= cliFlags.maxDuration) {
			log.Fatal("min duration must be less than max duration")
		}
		if cliFlags.minDuration == 0 && cliFlags.maxDuration == 0 {
			fmt.Println("No changes made")
			return
		}
		if cliFlags.minDuration != 0 {
			box.MinTime = cliFlags.minDuration
		}
		if cliFlags.maxDuration != 0 {
			box.MaxTime = cliFlags.maxDuration
		}
		err := tb.UpdateBox(box)
		if err != nil {
			log.Fatal(err)
		}
	},
}
