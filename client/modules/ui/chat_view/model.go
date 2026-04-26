package chat_view

import (
	"github.com/Mattcazz/Chat-TUI/client/types"
)

type Model struct {
	// First in slice is the newest message
	// Last in slice is the oldest message
	messages[] types.Message

	err error
	width int
	height int
}

func New() Model {
	return Model{
		messages: make([]types.Message, 0),
		err: nil,
	}
}

 func (m *Model) SetSize(width int, height int) {
	 m.width = width
	 m.height = height
 }
