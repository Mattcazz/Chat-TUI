package chat

import (
	"clit_client/internal/commands"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

		// Reduce size of border
		msg.Width -= lipgloss.NormalBorder().GetLeftSize() * 2 // Extra spacing
		msg.Width -= lipgloss.NormalBorder().GetRightSize() * 2 // Extra spacing
		msg.Height -= lipgloss.NormalBorder().GetTopSize() * 2 // 2 boxes inside vertically stacked
		msg.Height -= lipgloss.NormalBorder().GetBottomSize() * 2 // 2 boxes inside vertically stacked
		
		// Calculate size of each thing
		// Input gets height of 1, for example
		chat_view_height := msg.Height - 1
		chat_input_height := 1

		msg.Height = chat_view_height
		m.chat_view, _ = m.chat_view.Update(msg)

		msg.Height = chat_input_height
		m.chat_input, _ = m.chat_input.Update(msg)

		return m, nil
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
