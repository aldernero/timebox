package commands

import (
	"fmt"
	"github.com/aldernero/timebox/pkg/util"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
	"log"
)

var listSpansCmd = &cobra.Command{
	Use:   "spans",
	Short: "List spans",
	Run: func(cmd *cobra.Command, args []string) {
		var rows [][]string
		for k, v := range tb.Spans {
			for _, s := range v.Spans {
				start := s.Start.Format("2006-01-02 15:04:05")
				end := s.End.Format("2006-01-02 15:04:05")
				rows = append(rows, []string{k, start, end})
			}
		}
		t := table.New().
			Border(lipgloss.NormalBorder()).
			Headers("Box", "Start", "End").
			StyleFunc(func(row, col int) lipgloss.Style {
				return lipgloss.NewStyle().Margin(0, 1)
			}).
			Rows(rows...)
		fmt.Println(t.Render())
	},
}

var addSpanCmd = &cobra.Command{
	Use:   "span",
	Short: "Add a new span",
	Run: func(cmd *cobra.Command, args []string) {
		start, err := util.ParseTime(cliFlags.startTime)
		if err != nil {
			log.Fatal(err)
		}
		end, err := util.ParseTime(cliFlags.endTime)
		if err != nil {
			log.Fatal(err)
		}
		if cliFlags.boxName == "" {
			log.Fatal("box name is required")
		}
		if start.IsZero() || end.IsZero() {
			log.Fatal("start and end times are required")
		}
		if start.After(end) {
			log.Fatal("start time must be before end time")
		}
		span := util.Span{
			Start: start,
			End:   end,
			Box:   cliFlags.boxName,
		}
		err = tb.AddSpan(span, cliFlags.boxName)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var deleteSpanCmd = &cobra.Command{
	Use:   "span",
	Short: "Delete a span",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var updateSpanCmd = &cobra.Command{
	Use:   "span",
	Short: "Update a span",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
