package login

import (
	"github.com/charmbracelet/bubbles/textinput"
)

type Model struct {
	text_input textinput.Model
	err error
}

func New() Model {
	ti := textinput.New()
	ti.Placeholder = "Username"
	ti.Focus()
	ti.CharLimit = 20
	ti.Width = 20

	return Model{
		text_input: ti,
		err: nil,
	}
}
