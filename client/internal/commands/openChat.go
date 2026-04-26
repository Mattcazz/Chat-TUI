package commands

import (
	tea "github.com/charmbracelet/bubbletea"
)

type OpenChatMsg struct {
	Username string
	ConversationID int64
}

func NewOpenChatCmd(username string, conversationId int64) func() tea.Msg {
	var msg OpenChatMsg
	msg.Username = username
	msg.ConversationID = conversationId

	return func() tea.Msg {
		return msg
	}
}
