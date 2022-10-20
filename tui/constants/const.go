package constants

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	TUIWidth              = 79
	defaultInputWidth     = 28
	SummaryPageSize       = 10
	DetailPageSize        = 16
	inputTimeFormShort    = "2006-01-02"
	inputTimeFormLong     = "2006-01-02 15:04:05"
	ColorError            = "#CF002E"
	ColorDetailTitle      = "#DE3E93"
	ColorPromptBorder     = "#DE3E93"
	ColorTextLightGray    = "#FFFDF5"
	ColorLogo             = "#FFDF80"
	ColorTableBorder      = "#47A4AC"
	ColorTableText        = "#BAEBDA"
	ColorPeriodForeground = "#BAEBDA"
	ColorPeriodHighlight  = "#DE3E93"
	ColorHelpText         = "#ABABAB"
	ColorShortcuts        = "#47A4AC"
	ColorSession          = "#FFDF80"
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
var PromptStyle = lipgloss.NewStyle().
	Width(TUIWidth).
	Align(lipgloss.Center)
var ErrStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorError)).Render
var NoStyle = lipgloss.NewStyle()
var FocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorPromptBorder))
var BlurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorTableText))
var PeriodStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color(ColorPeriodForeground)).
	Padding(0, 1)
var CurrentPeriodStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color(ColorPeriodForeground)).
	Background(lipgloss.Color(ColorPeriodHighlight)).
	Padding(0, 1)
var PeriodPickerStyle = lipgloss.NewStyle().
	PaddingRight(1).
	Align(lipgloss.Right)
var HelpStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color(ColorHelpText)).
	Padding(0, 1)
var ShortcutStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color(ColorShortcuts)).
	Padding(0, 1)
var SessionShortcutStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color(ColorSession))
var SessionSelectedStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color(ColorPeriodForeground)).
	Background(lipgloss.Color(ColorSession))
var SessionTextStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color(ColorHelpText)).
	PaddingRight(2)
var HelpBlockStyle = lipgloss.NewStyle().
	PaddingLeft(8).
	PaddingRight(1)
