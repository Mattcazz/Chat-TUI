package login

import (
	"github.com/Mattcazz/Chat-TUI/client/internal/config"
	"github.com/Mattcazz/Chat-TUI/client/types"
	"github.com/charmbracelet/bubbles/textinput"
)

type Model struct {
	username_input textinput.Model
	password_input textinput.Model
	pk []byte
	sk []byte
	nonce []byte
	signature []byte

	client *types.LoginClient
	config *config.Config

	state types.LoginModelState
	err error
	width int
	height int
}

func NewLoginModel(baseClient *types.BaseClient) Model {
	username_ti := textinput.New()
	username_ti.Placeholder = "Username"
	username_ti.Focus()
	username_ti.CharLimit = 25
	username_ti.Width = 28

	password_ti := textinput.New()
	password_ti.Placeholder = "Password"
	password_ti.EchoMode = textinput.EchoNone
	password_ti.CharLimit = 0 // inf
	password_ti.Width = 0

	return Model{
		username_input: username_ti,
		password_input: password_ti,

		pk: nil,
		sk: nil,
		nonce: nil,
		signature: nil,

		client: &types.LoginClient{Client: *baseClient},

		state: types.Normal,
		err: nil,
	}
}

func (m *Model) SetSize(width int, height int) {
	 m.width = width
	 m.height = height
 }
