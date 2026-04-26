package chat_view

import (
	"github.com/Mattcazz/Chat-TUI/client/internal/commands"
	"github.com/Mattcazz/Chat-TUI/client/internal/logger"
	"github.com/Mattcazz/Chat-TUI/client/types"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case commands.NewMessageMsg:
		var newMsg types.Message
		newMsg.Author = msg.Author
		newMsg.Message = msg.Message
		newMsg.Timestamp = msg.Timestamp
		logger.Log.Printf("[CHAT VIEW] New message by '%s': %s", newMsg.Author, newMsg.Message)
		m.messages = append(m.messages, newMsg)
	case commands.LoadChatMsg:
		m.messages = make([]types.Message, len(msg.Messages))
		for i, message := range msg.Messages {
			m.messages[len(msg.Messages)-1-i] = message
		}
	}

	return m, cmd
}
