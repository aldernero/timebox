package tui

type crudState int

const (
	add crudState = iota
	nav           // read
	edit
	del
)

type viewMode int

const (
	boxSummary viewMode = iota
	boxView
	timeline
)

// Common shortcuts
var (
	addShortcut        = NewShortcut("a", "Add")
	delShortcut        = NewShortcut("d", "Delete")
	editShortcut       = NewShortcut("e", "Edit")
	deleteShortcut     = NewShortcut("d", "Delete")
	quitShortcut       = NewShortcut("q", "Quit")
	periodShortcut     = NewShortcut("Tab", "Period")
	enterShortcut      = NewShortcut("Enter", "Spans")
	backShortcut       = NewShortcut("Esc", "Back")
	boxSummaryShortcut = NewShortcut("b", "Boxes")
	timelineShortcut   = NewShortcut("t", "Timeline")
)

func printCrudState(s crudState) string {
	switch s {
	case add:
		return "Add"
	case nav:
		return "Read"
	case edit:
		return "Edit"
	case del:
		return "Delete"
	default:
		return "Unknown"
	}
}

func printViewMode(v viewMode) string {
	switch v {
	case boxSummary:
		return "Box Summary"
	case boxView:
		return "Box View"
	case timeline:
		return "Timeline"
	default:
		return "Unknown"
	}
}
