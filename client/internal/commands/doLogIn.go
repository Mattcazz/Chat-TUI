package commands

import tea "github.com/charmbracelet/bubbletea"

type DoLogInMsg struct {
}

func NewDoLoginCmd() func() tea.Msg {
	var msg DoLogInMsg

	return func() tea.Msg {
		return msg
	}
}
