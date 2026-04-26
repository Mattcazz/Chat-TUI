package chat

import (
	"github.com/Mattcazz/Chat-TUI/client/modules/ui/chat_view"
	"github.com/Mattcazz/Chat-TUI/client/modules/ui/input_box"
	"github.com/Mattcazz/Chat-TUI/client/types"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	chatView tea.Model
	chatInput tea.Model

	conversationId int64

	client *types.ChatClient

	err error
	width int
	height int
}

func NewChatModel(baseClient *types.BaseClient) Model {
	return Model{
		chatView: chat_view.New(),
		chatInput: input_box.New(),
		client : &types.ChatClient{Client: *baseClient},
	}
}

 func (m *Model) SetSize(width int, height int) {
	 m.width = width
	 m.height = height
 }
