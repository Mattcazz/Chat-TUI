package app

import (
	"net/http"
	"time"

	"github.com/Mattcazz/Chat-TUI/client/internal/commands"
	"github.com/Mattcazz/Chat-TUI/client/internal/config"
	"github.com/Mattcazz/Chat-TUI/client/internal/logger"
	"github.com/Mattcazz/Chat-TUI/client/modules/ui/chat"
	"github.com/Mattcazz/Chat-TUI/client/modules/ui/inbox"
	"github.com/Mattcazz/Chat-TUI/client/modules/ui/login"
	"github.com/Mattcazz/Chat-TUI/client/styles"
	"github.com/Mattcazz/Chat-TUI/client/types"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type App struct {
	state types.SessionState
	loginModel tea.Model
	inboxModel tea.Model
	chatModel tea.Model

	client *types.BaseClient

	username string
	err error
	width int
	height int
}

func New() App {
	appClient := http.Client{
		Timeout: time.Second * 10,
	}
	config.LoadConfig()

	logger.Init()

	client := &types.BaseClient{Client: appClient, Config: config.Configuration} // TODO only pass host and port
	return App{
		state: types.LoginView,
		loginModel: login.NewLoginModel(client),
		inboxModel: inbox.NewInboxModel(client),
		chatModel: chat.New(),
		client: client,
		err: nil,
	}
}

func (a App) Init() tea.Cmd {
	return tea.Batch(
		a.loginModel.Init(),
		a.inboxModel.Init(),
		a.chatModel.Init(),
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
		logger.Log.Printf("[APP] Successfully logged in")
		logger.Log.Printf("[APP] Switching to inbox view")
		m.state = types.InboxView

		return m, commands.NewUpdateInboxCmd()
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		msg.Width -= lipgloss.NormalBorder().GetLeftSize()
		msg.Width -= lipgloss.NormalBorder().GetRightSize()
		msg.Height -= lipgloss.NormalBorder().GetTopSize()
		msg.Height -= lipgloss.NormalBorder().GetBottomSize()

		m.loginModel, _ = m.loginModel.Update(msg)
		m.inboxModel, _ = m.inboxModel.Update(msg)
		m.chatModel, _ = m.chatModel.Update(msg)

		return m, nil
	}

	switch m.state {
	case types.LoginView:
		m.loginModel, cmd = m.loginModel.Update(msg)
		return m, cmd
	case types.InboxView:
		m.inboxModel, cmd = m.inboxModel.Update(msg)
		return m, cmd
	case types.ChatView:
		m.chatModel, cmd = m.chatModel.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m App) View() string {
	var modelContent string

	switch m.state {
	case types.LoginView:
		modelContent = m.loginModel.View()
	case types.InboxView:
		modelContent = m.inboxModel.View()
	case types.ChatView:
		modelContent = m.chatModel.View()
	default:
		modelContent = "Ya broke it"
	}

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		styles.Default.Border.Render(modelContent),
	)
}


