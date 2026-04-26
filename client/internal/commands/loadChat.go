package commands

import (
	"github.com/Mattcazz/Chat-TUI/client/types"
	tea "github.com/charmbracelet/bubbletea"
)

type LoadChatMsg struct {
	Messages[] types.Message
}

func NewLoadChatCmd(messages[] types.Message) func() tea.Msg {
	var msg LoadChatMsg
	msg.Messages = messages

	return func() tea.Msg {
		return msg
	}
}
