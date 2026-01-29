package login

import (
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		"Log-In to CLIt",
		"",
		m.text_input.View(),
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
