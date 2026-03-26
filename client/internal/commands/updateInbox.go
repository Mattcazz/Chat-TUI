package commands

import tea "github.com/charmbracelet/bubbletea"

type UpdateInboxMsg struct {
}

func NewUpdateInboxCmd() func() tea.Msg {
	var msg UpdateInboxMsg

	return func() tea.Msg {
		return msg
	}
}
