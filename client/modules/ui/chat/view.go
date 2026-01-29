package chat

import (
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	content := m.chat_input.View()

	style := lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("#bbbbbb"))

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Bottom,
		style.Render(content),
	)
}
