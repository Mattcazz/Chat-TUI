package commands

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type NewMessageMsg struct {
	Author string
	Message string
	Timestamp time.Time
}

func NewNewMessageCmd(author string, message string, timestamp time.Time) func() tea.Msg {
	var msg NewMessageMsg
	msg.Author = author
	msg.Message = message
	msg.Timestamp = timestamp

	return func() tea.Msg {
		return msg
	}
}
