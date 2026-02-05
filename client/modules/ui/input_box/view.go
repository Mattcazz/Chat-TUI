package input_box

import (
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	content := m.chat_input.View()

	style := lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("#bbbbbb"))

	return style.Render(content)
}
