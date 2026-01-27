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
	Width(m.text_input.CharLimit + 3 /* "> " + 1 because it behaves weirdly if it's exact or the wrong parity */).
	Align(lipgloss.Center, lipgloss.Center).
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("#bbbbbb"))

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		style.Render(content),
	)
}
