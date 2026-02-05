package input_box

import (
	"clit_client/internal/commands"
	"time"

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
		case tea.KeyEnter:
			cmd = commands.NewNewMessageCmd("nucieda", m.chat_input.Value(), time.Now())
			m.chat_input.Reset()
			return m, cmd
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	m.chat_input, cmd = m.chat_input.Update(msg)

	return m, cmd
}
