package chat

import (
	"github.com/charmbracelet/bubbles/textinput"
)

type Model struct {
	// db probs //
	chat_input textinput.Model
	err error
	width int
	height int
}

func New() Model {
	ti := textinput.New()
	ti.Focus()
	
	return Model{
		chat_input: ti,
		err: nil,
	}
}

 func (m *Model) SetSize(width int, height int) {
	 m.width = width
	 m.height = height
 }
