package input_box

import (
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	content := m.chatInput.View()

	style := lipgloss.NewStyle()

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Left, lipgloss.Bottom,
		style.Render(content),
	)
}
