package chat

import (
	"github.com/Mattcazz/Chat-TUI/client/internal/commands"
	"github.com/Mattcazz/Chat-TUI/client/internal/logger"
	"github.com/Mattcazz/Chat-TUI/client/types"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) fetchConversationMessages(conversationId int64) ([]types.Message, error) {
	conversationResponse, err := m.client.GetChat(conversationId)
	if err != nil {
		return nil, err
	}

	messages := make([]types.Message, 0)
	for _, msg := range conversationResponse.Messages {
		messages = append(messages, types.Message{
			Author: msg.UserName,
			Message: msg.Content,
			Timestamp: msg.CreatedAt,
		})
	}

	return messages, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Reduce size of border
		msg.Width -= lipgloss.NormalBorder().GetLeftSize() * 2 // Extra spacing
		msg.Width -= lipgloss.NormalBorder().GetRightSize() * 2 // Extra spacing
		msg.Height -= lipgloss.NormalBorder().GetTopSize() * 2 // 2 boxes inside vertically stacked
		msg.Height -= lipgloss.NormalBorder().GetBottomSize() * 2 // 2 boxes inside vertically stacked
		
		// Calculate size of each thing
		// Input gets height of 1, for example
		chatViewHeight := msg.Height - 1
		chatInputHeight := 1

		msg.Height = chatViewHeight
		m.chatView, _ = m.chatView.Update(msg)

		msg.Height = chatInputHeight
		m.chatInput, _ = m.chatInput.Update(msg)

		return m, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case commands.NewMessageMsg:
		err := m.client.SendMessage(m.conversationId, msg.Message)
		if err != nil {
			logger.Log.Panicf("[CHAT] Got an error trying to send a message: %s", err.Error())
		}
		m.chatView, cmd = m.chatView.Update(msg)
		return m, cmd
	case commands.OpenChatMsg:
		m.conversationId = msg.ConversationID
		messages, err := m.fetchConversationMessages(m.conversationId)
		if err != nil {
			logger.Log.Panicf("[CHAT] Error fetching chat from server: %s", err.Error())
		}

		loadChatCmd := commands.NewLoadChatCmd(messages)
		return m, loadChatCmd
	case commands.LoadChatMsg:
		m.chatView, cmd = m.chatView.Update(msg)
		return m, cmd
	case commands.LogInMsg:
		m.username = msg.Username
		return m, nil
	}

	m.chatInput, cmd = m.chatInput.Update(msg)

	return m, cmd
}
