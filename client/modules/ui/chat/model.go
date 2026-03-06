package chat

import (
	"github.com/Mattcazz/Chat-TUI/client/modules/ui/chat_view"
	"github.com/Mattcazz/Chat-TUI/client/modules/ui/input_box"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	chatView tea.Model
	chatInput tea.Model

	username string

	err error
	width int
	height int
}

func New() Model {
	return Model{
		chatView: chat_view.New(),
		chatInput: input_box.New(),
		err: nil,
	}
}

 func (m *Model) SetSize(width int, height int) {
	 m.width = width
	 m.height = height
 }
