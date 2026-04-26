package login

import (
	"fmt"

	"github.com/Mattcazz/Chat-TUI/client/internal/config"
	"github.com/Mattcazz/Chat-TUI/client/styles"
	"github.com/Mattcazz/Chat-TUI/client/types"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	var content string
	switch m.state {
	case types.Normal:
		content = lipgloss.JoinVertical(
			lipgloss.Center,
			"Logging in to CLIt",
		)
	case types.NeedsUsername:
		content = lipgloss.JoinVertical(
			lipgloss.Center,
			"Register to CLIt",
			"",
			m.usernameInput.View(),
		)
	case types.NeedsSSHPassword:
		content = lipgloss.JoinVertical(
			lipgloss.Center,
			fmt.Sprintf("SSH key %s requires password", config.Configuration.SSHKeyName),
			"",
			m.passwordInput.View(),
		)
	}

	// append error
	if len(m.errorMsg) > 0 {
		content = lipgloss.JoinVertical(
			lipgloss.Center,
			content,
			m.errorMsg,
		)
	}

	style := styles.Default.Border
	if m.config != nil {
		style.BorderForeground(lipgloss.Color(m.config.Colors.Border))
	}

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		style.Render(content),
	)
}
