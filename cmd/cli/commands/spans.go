package commands

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
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
