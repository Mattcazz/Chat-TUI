package chat

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
		msg.Author = m.username
		m.chat_view, cmd = m.chat_view.Update(msg)
		return m, cmd
	case commands.LogInMsg:
		m.username = msg.Username
		return m, nil
	}

	m.chat_input, cmd = m.chat_input.Update(msg)

	return m, cmd
}
