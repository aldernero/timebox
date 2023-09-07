package tui

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	UIWidth            = 79
	defaultInputWidth  = 28
	SummaryPageSize    = 10
	DetailPageSize     = 16
	inputTimeFormShort = "2006-01-02"
	inputTimeFormLong  = "2006-01-02 15:04:05"
	ColorError         = "#CF002E"
	ColorDetailTitle   = "#DE3E93"
	ColorPromptBorder  = "#DE3E93"
	ColorTextLightGray = "#FFFDF5"
	ColorLogo          = "#FFDF80"
	ColorTableBorder   = "#47A4AC"
	ColorTableText     = "#BAEBDA"
)

var LogoStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color(ColorLogo)).
	Padding(0, 1)
var TableStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color(ColorTableText)).
	BorderForeground(lipgloss.Color(ColorTableBorder)).
	Align(lipgloss.Right)
var InputTitleStyle = lipgloss.NewStyle().
	Width(defaultInputWidth).
	Foreground(lipgloss.Color(ColorTextLightGray)).
	Background(lipgloss.Color(ColorDetailTitle)).
	Padding(0, 1).
	Align(lipgloss.Center)
var InputStyle = lipgloss.NewStyle().
	Margin(1, 1).
	Padding(1, 2).
	Border(lipgloss.RoundedBorder(), true, true, true, true).
	BorderForeground(lipgloss.Color(ColorPromptBorder)).
	Render
var DeleteStyle = lipgloss.NewStyle().
	Margin(1, 1).
	Padding(1, 2).
	Foreground(lipgloss.Color(ColorPromptBorder)).
	Align(lipgloss.Center).
	Border(lipgloss.RoundedBorder(), true, true, true, true).
	BorderForeground(lipgloss.Color(ColorPromptBorder)).
	Render
var PromptStyle = lipgloss.NewStyle().
	Width(UIWidth).
	Align(lipgloss.Center)
var ErrStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorError)).Render
var NoStyle = lipgloss.NewStyle()
var FocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorPromptBorder))
var BlurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorTableText))
