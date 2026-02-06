package chat_view

import (
	"clit_client/internal/commands"

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
		var newMsg Message
		newMsg.author = msg.Author
		newMsg.message = msg.Message
		newMsg.timestamp = msg.Timestamp
		m.messages = append(m.messages, newMsg)
	}

	return m, cmd
}
