package chat

import (
	"clit_client/modules/ui/chat_view"
	"clit_client/modules/ui/input_box"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	chat_view tea.Model
	chat_input tea.Model

	err error
	width int
	height int
}

func New() Model {
	return Model{
		chat_view: chat_view.New(),
		chat_input: input_box.New(),
		err: nil,
	}
}

 func (m *Model) SetSize(width int, height int) {
	 m.width = width
	 m.height = height
 }
