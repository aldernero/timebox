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
