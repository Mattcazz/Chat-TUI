package styles

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	Title lipgloss.Style
	Text lipgloss.Style
	Border lipgloss.Style
	Error lipgloss.Style
}

var Default = Styles{
	Title: lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#CCCCCC")),
	Text: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#BBBBBB")),
	Border: lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#BBBBBB")),
	Error: lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000")),
}
