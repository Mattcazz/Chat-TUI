package login

import (
	"clit_client/styles"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		"Log-In to CLIt",
		"",
		m.text_input.View(),
	)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		styles.Default.Border.Render(content),
	)
}
