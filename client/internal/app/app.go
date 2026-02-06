package app

import (
	"clit_client/internal/commands"
	"clit_client/modules/ui/chat"
	"clit_client/modules/ui/login"
	"clit_client/types"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type App struct {
	state types.SessionState
	login_model tea.Model
	chat_model tea.Model

	username string
	err error
	width int
	height int
}

func New() App {
	return App{
		state: types.LoginView,
		login_model: login.New(),
		chat_model: chat.New(),
		err: nil,
	}
}

func (a App) Init() tea.Cmd {
	return nil
}

func (m App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case commands.ChangeStateMsg:
		m.state = msg.State
		return m, nil
	case commands.LogInMsg:
		m.username = msg.Username
		// Check that log in was successful eventually, but
		m.chat_model, _ = m.chat_model.Update(msg)
		return m, commands.NewChangeStateCmd(types.ChatView)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		msg.Width -= lipgloss.NormalBorder().GetLeftSize()
		msg.Width -= lipgloss.NormalBorder().GetRightSize()
		msg.Height -= lipgloss.NormalBorder().GetTopSize()
		msg.Height -= lipgloss.NormalBorder().GetBottomSize()

		m.login_model, _ = m.login_model.Update(msg)
		m.chat_model, _ = m.chat_model.Update(msg)

		return m, nil
	}

	switch m.state {
	case types.LoginView:
		m.login_model, cmd = m.login_model.Update(msg)
		return m, cmd
	case types.ChatView:
		m.chat_model, cmd = m.chat_model.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m App) View() string {
	var model_content string

	switch m.state {
	case types.LoginView:
		model_content = m.login_model.View()
	case types.ChatView:
		model_content = m.chat_model.View()
	default:
		model_content = "Ya broke it"
	}

	style := lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("#bbbbbb"))

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		style.Render(model_content),
	)
}


