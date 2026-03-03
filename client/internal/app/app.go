package app

import (
	"net/http"
	"time"

	"github.com/Mattcazz/Chat-TUI/client/internal/commands"
	"github.com/Mattcazz/Chat-TUI/client/internal/config"
	"github.com/Mattcazz/Chat-TUI/client/internal/logger"
	"github.com/Mattcazz/Chat-TUI/client/modules/ui/chat"
	"github.com/Mattcazz/Chat-TUI/client/modules/ui/login"
	"github.com/Mattcazz/Chat-TUI/client/styles"
	"github.com/Mattcazz/Chat-TUI/client/types"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type App struct {
	state types.SessionState
	login_model tea.Model
	chat_model tea.Model

	client *types.BaseClient

	username string
	err error
	width int
	height int
}

func New() App {
	app_client := http.Client{
		Timeout: time.Second * 10,
	}
	config.LoadConfig()

	logger.Init()

	client := &types.BaseClient{Client: app_client, Config: config.Configuration} // TODO only pass host and port
	return App{
		state: types.LoginView,
		login_model: login.NewLoginModel(client),
		chat_model: chat.New(),
		client: client,
		err: nil,
	}
}

func (a App) Init() tea.Cmd {
	return tea.Batch(
		a.login_model.Init(),
		a.chat_model.Init(),
	)
}

func (m App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case commands.ChangeStateMsg:
		logger.Log.Printf("[APP] Changing state to: %s", msg.State)
		m.state = msg.State
		return m, nil
	case commands.LogInMsg:
		logger.Log.Printf("[APP] Successfully logged in with username: %s", msg.Username)
		m.username = msg.Username
		m.chat_model, _ = m.chat_model.Update(msg)
		logger.Log.Printf("[APP] Switching to chat view...")
		m.state = types.ChatView

		return m, nil
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

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		styles.Default.Border.Render(model_content),
	)
}


