package commands

import (
	"fmt"
	"github.com/aldernero/timebox/pkg/util"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
	"log"
	"strconv"
	"time"
)

var listSpansCmd = &cobra.Command{
	Use:   "spans",
	Short: "List spans",
	Run: func(cmd *cobra.Command, args []string) {
		var filterSpan util.Span
		if cliFlags.startTime == "" {
			filterSpan.Start = time.Time{}
		} else {
			from, err := util.ParseDurationOrTime(cliFlags.startTime)
			if err != nil {
				log.Fatal(err)
			}
			filterSpan.Start = from
		}
		if cliFlags.endTime == "" {
			filterSpan.End = time.Now()
		} else {
			to, err := util.ParseDurationOrTime(cliFlags.endTime)
			if err != nil {
				log.Fatal(err)
			}
			filterSpan.End = to
		}
		var spanset util.SpanSet
		if cliFlags.boxName != "" {
			// check if box exists
			if _, ok := tb.Boxes[cliFlags.boxName]; !ok {
				log.Fatalf("box \"%s\" does not exist", cliFlags.boxName)
			}
			fullset := tb.GetSpansForTimespan(filterSpan)
			spanset = util.NewSpanSet()
			for _, s := range fullset.Spans {
				if s.Box == cliFlags.boxName {
					spanset.Add(s)
				}
			}
		} else {
			spanset = tb.GetSpansForTimespan(filterSpan)
		}
		var rows [][]string
		for _, s := range spanset.Spans {
			id := fmt.Sprintf("%d", s.ID)
			start := s.Start.Format("2006-01-02 15:04:05")
			end := s.End.Format("2006-01-02 15:04:05")
			dur := s.End.Sub(s.Start).String()
			rows = append(rows, []string{id, s.Box, start, end, dur})
		}
		t := table.New().
			Border(lipgloss.NormalBorder()).
			Headers("ID", "Box", "Start", "End", "Duration").
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
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatal(err)
		}
		var confirmed bool
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Are you sure you want to delete the span?").
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
				fmt.Println("Cancelling span delete")
				return
			}
		}
		err = tb.DeleteSpanByID(int64(id))
		if err != nil {
			log.Fatal(err)
		}
	},
}

var updateSpanCmd = &cobra.Command{
	Use:   "span",
	Short: "Update a span",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
