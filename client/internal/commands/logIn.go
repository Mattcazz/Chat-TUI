package commands

import (
	tea "github.com/charmbracelet/bubbletea"
)

type LogInMsg struct {
	Username string
}

func NewLogInCmd(username string) func() tea.Msg {
	var msg LogInMsg
	msg.Username = username

	return func() tea.Msg {
		return msg
	}
}
