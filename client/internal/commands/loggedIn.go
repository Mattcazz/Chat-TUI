package commands

import (
	tea "github.com/charmbracelet/bubbletea"
)

type LoggedInMsg struct {
	Username string
}

func NewLoggedInCmd(username string) func() tea.Msg {
	var msg LoggedInMsg
	msg.Username = username

	return func() tea.Msg {
		return msg
	}
}
