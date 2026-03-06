package input_box

import (
	"github.com/charmbracelet/bubbles/textinput"
)

type Model struct {
	// db probs //
	chatInput textinput.Model
	err error
	width int
	height int
}

func New() Model {
	ti := textinput.New()
	ti.Focus()
	
	return Model{
		chatInput: ti,
		err: nil,
	}
}

 func (m *Model) SetSize(width int, height int) {
	 m.width = width
	 m.height = height
 }
