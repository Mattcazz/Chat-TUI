package chat

import (
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	chat_view_content := m.chat_view.View()
	chat_input_box_content := m.chat_input.View()

	box_style := lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("#bbbbbb"))

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		box_style.Render(chat_view_content),
		box_style.Render(chat_input_box_content),
	)

	style := lipgloss.NewStyle()

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		style.Render(content),
	)
}
