package chat_view

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	var content strings.Builder

	var authorStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#A32CC4")) // Purple

	var messageStyle = lipgloss.NewStyle().
		Bold(false).
		PaddingLeft(2).
		Foreground(lipgloss.Color("#999999"))
	
	// var timestampStyle = lipgloss.NewStyle().
	// 	Faint(true).
	// 	Foreground(lipgloss.Color("#888888"))

	for _, msg := range m.messages {
		content.WriteString(authorStyle.Render(msg.author))
		content.WriteString(messageStyle.Render(msg.message))
		content.WriteByte('\n')
	}

	var style = lipgloss.NewStyle()

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Left, lipgloss.Bottom,
		style.Render(content.String()),
	)
}
