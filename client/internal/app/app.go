package app

type ViewType int
const (
	LOGIN ViewType = iota
	CHAT
)

type App struct {
    currentView ViewType

    login	*login.Model
    chat	*chat.Model

}

func New(apiURL string) App {
	return App {
		currentView: LOGIN,

		login: login_model.New(),
		chat: chat_model.New(),
	}
}

func (a App) Init() tea.Cmd { ... }
func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) { ... }
func (a App) View() string { ... }
