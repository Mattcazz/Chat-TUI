package chat

import (
	"github.com/Mattcazz/Chat-TUI/client/styles"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	chatViewContent := m.chatView.View()
	chatInputBoxContent := m.chatInput.View()

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		styles.Default.Border.Render(chatViewContent),
		styles.Default.Border.Render(chatInputBoxContent),
	)

	style := lipgloss.NewStyle()

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		style.Render(content),
	)
}
