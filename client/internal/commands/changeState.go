package commands

import (
	"clit_client/types"

	tea "github.com/charmbracelet/bubbletea"
)

type ChangeStateMsg struct {
	State types.SessionState
}

func NewChangeStateCmd(state types.SessionState) func() tea.Msg {
	var msg ChangeStateMsg
	msg.State = state

	return func() tea.Msg {
		return msg
	}
}
