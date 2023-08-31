package tui

import "github.com/charmbracelet/lipgloss"

const (
	ShortCutKeyColor  = "#47A4AC"
	ShortCutDescColor = "#BAEBDA"
)

var ShortcutKeyStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color(ShortCutKeyColor))
var ShortcutDescStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color(ShortCutDescColor))

type Shortcut struct {
	Key       string
	Desc      string
	KeyStyle  lipgloss.Style
	DescStyle lipgloss.Style
}

func NewShortcut(key, desc string) Shortcut {
	return Shortcut{
		Key:       key,
		Desc:      desc,
		KeyStyle:  ShortcutKeyStyle,
		DescStyle: ShortcutDescStyle,
	}
}

func (s Shortcut) KeyLength() int {
	return len(s.Key) + 2 // 2 brackets
}

func (s Shortcut) DescLength() int {
	return len(s.Desc)
}

func (s Shortcut) Length() int {
	return s.KeyLength() + s.DescLength() + 1 // 1 space
}

func (s Shortcut) String() string {
	return s.KeyStyle.Render("<"+s.Key+">") + " " + s.DescStyle.Render(s.Desc)
}

func ShortcutRow(shortcuts []Shortcut) string {
	var result string
	var maxKeyLength int
	for _, shortcut := range shortcuts {
		if shortcut.KeyLength() > maxKeyLength {
			maxKeyLength = shortcut.KeyLength()
		}
	}
	n := len(shortcuts)
	for i, shortcut := range shortcuts {
		keyLen := shortcut.KeyLength()
		keyStr := shortcut.KeyStyle.Render("<" + shortcut.Key + ">")
		for j := 0; j < maxKeyLength-keyLen; j++ {
			keyStr += " "
		}
		result += keyStr + " " + shortcut.DescStyle.Render(shortcut.Desc)
		if i < n-1 {
			result += "\n"
		}
	}
	return lipgloss.NewStyle().PaddingLeft(1).Render(result)
}
