package chat_view

import (
	"time"
)

type Message struct {
	author string
	message string
	timestamp time.Time
}

type Model struct {
	messages[] Message

	err error
	width int
	height int
}

func New() Model {
	return Model{
		messages: make([]Message, 0),
		err: nil,
	}
}

 func (m *Model) SetSize(width int, height int) {
	 m.width = width
	 m.height = height
 }
