package chat

import (
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	chat_view_content := m.chat_view.View()
	chat_input_box_content := m.chat_input.View()

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		chat_view_content,
		chat_input_box_content,
	)

	style := lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("#bbbbbb"))

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		style.Render(content),
	)
}
