package login

import (
	"github.com/charmbracelet/bubbles/textinput"
)

type Model struct {
	text_input textinput.Model
	err error
	width int
	height int
}

func New() Model {
	ti := textinput.New()
	ti.Placeholder = "Username"
	ti.Focus()
	ti.CharLimit = 25
	ti.Width = 28

	return Model{
		text_input: ti,
		err: nil,
	}
}

 func (m *Model) SetSize(width int, height int) {
	 m.width = width
	 m.height = height
 }
