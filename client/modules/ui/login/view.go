package login

import (
	"github.com/Mattcazz/Chat-TUI/client/styles"

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
