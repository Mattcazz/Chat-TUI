package inbox

import (
	"github.com/Mattcazz/Chat-TUI/client/types"
	"github.com/charmbracelet/bubbles/list"
)

type Model struct {
	conversationList list.Model

	client *types.InboxClient
	errorMsg string
	width int
	height int
}

func NewInboxModel(baseClient *types.BaseClient) Model {
	empty_list := make([]list.Item, 0)
	conversationList := list.New(empty_list, list.NewDefaultDelegate(), 0, 0)
	conversationList.Title = "Inbox"
	conversationList.SetShowTitle(true)

	return Model{
		conversationList: conversationList,
		client: &types.InboxClient{Client: *baseClient},
	}
}

func (m *Model) SetSize(width int, height int) {
	 m.width = width
	 m.height = height
 }
