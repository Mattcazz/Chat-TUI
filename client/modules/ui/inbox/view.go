package inbox

import (
	"github.com/Mattcazz/Chat-TUI/client/internal/config"
	"github.com/Mattcazz/Chat-TUI/client/styles"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	var content string = lipgloss.JoinVertical(
		lipgloss.Center,
		"Inbox",
		m.conversationList.View(),
	)

	// append error
	if len(m.errorMsg) > 0 {
		content = lipgloss.JoinVertical(
			lipgloss.Center,
			content,
			m.errorMsg,
		)
	}

	style := styles.Default.Border
	style.BorderForeground(lipgloss.Color(config.Configuration.Colors.Border))

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		style.Render(content),
	)
}
