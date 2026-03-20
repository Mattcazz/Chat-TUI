package login

import (
	"github.com/Mattcazz/Chat-TUI/client/internal/config"
	"github.com/Mattcazz/Chat-TUI/client/types"
	"github.com/charmbracelet/bubbles/textinput"
)

type Model struct {
	usernameInput textinput.Model
	passwordInput textinput.Model
	pk []byte
	sk []byte
	nonce []byte
	signature []byte

	client *types.LoginClient
	config *config.Config

	state types.LoginModelState
	errorMsg string
	width int
	height int
}

func NewLoginModel(baseClient *types.BaseClient) Model {
	usernameTi := textinput.New()
	usernameTi.Placeholder = "Username"
	usernameTi.Focus()
	usernameTi.CharLimit = 25
	usernameTi.Width = 28

	passwordTi := textinput.New()
	passwordTi.Placeholder = "Password"
	passwordTi.EchoMode = textinput.EchoNone
	passwordTi.CharLimit = 0 // inf
	passwordTi.Width = 0

	return Model{
		usernameInput: usernameTi,
		passwordInput: passwordTi,

		pk: nil,
		sk: nil,
		nonce: nil,
		signature: nil,

		client: &types.LoginClient{Client: *baseClient},

		state: types.Normal,
	}
}

func (m *Model) SetSize(width int, height int) {
	 m.width = width
	 m.height = height
 }
