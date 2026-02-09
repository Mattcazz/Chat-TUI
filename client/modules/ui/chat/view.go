package chat

import (
	"clit_client/styles"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	chat_view_content := m.chat_view.View()
	chat_input_box_content := m.chat_input.View()

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		styles.Default.Border.Render(chat_view_content),
		styles.Default.Border.Render(chat_input_box_content),
	)

	style := lipgloss.NewStyle()

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		style.Render(content),
	)
}
