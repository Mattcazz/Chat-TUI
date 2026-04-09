package chat_view

import (
	"strings"

	"github.com/Mattcazz/Chat-TUI/client/internal/config"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	var content strings.Builder

	var authorStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(config.Configuration.Colors.Username)) // Purple

	var messageStyle = lipgloss.NewStyle().
		Bold(false).
		PaddingLeft(2).
		Foreground(lipgloss.Color(config.Configuration.Colors.Text))
	
	// var timestampStyle = lipgloss.NewStyle().
	// 	Faint(true).
	// 	Foreground(lipgloss.Color("#888888"))

	// order of messages is "backwards" (see model.go)
	for _, msg := range m.messages {
		content.WriteString(authorStyle.Render(msg.Author))
		content.WriteString(messageStyle.Render(msg.Message))
		content.WriteByte('\n')
	}

	var style = lipgloss.NewStyle()

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Left, lipgloss.Bottom,
		style.Render(content.String()),
	)
}
