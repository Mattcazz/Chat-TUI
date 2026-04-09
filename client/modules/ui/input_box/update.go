package input_box

import (
	"time"

	"github.com/Mattcazz/Chat-TUI/client/internal/commands"
	"github.com/Mattcazz/Chat-TUI/client/internal/logger"
	"github.com/Mattcazz/Chat-TUI/client/internal/user"

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
			logger.Log.Printf("Hit Enter on input box with global username '%s'", user.UserData.UserName)
			cmd = commands.NewNewMessageCmd(user.UserData.UserName, m.chatInput.Value(), time.Now())
			m.chatInput.Reset()
			return m, cmd
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	m.chatInput, cmd = m.chatInput.Update(msg)

	return m, cmd
}
