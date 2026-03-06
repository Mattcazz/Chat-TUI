package commands

import (
	tea "github.com/charmbracelet/bubbletea"
)

type LoadSSHKeysMsg struct {
}

func NewLoadSSHKeysCmd() func() tea.Msg {
	var msg LoadSSHKeysMsg

	return func() tea.Msg {
		return msg
	}
}
