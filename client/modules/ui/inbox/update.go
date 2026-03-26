package inbox

import (
	"github.com/Mattcazz/Chat-TUI/client/internal/commands"
	"github.com/Mattcazz/Chat-TUI/client/internal/logger"
	"github.com/Mattcazz/Chat-TUI/client/internal/user"
	"github.com/Mattcazz/Chat-TUI/client/types"
	"github.com/Mattcazz/Chat-TUI/pkg"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func getConversationItemList(conversations []pkg.InboxConversationResponse) []list.Item {
	conversationList := make([]types.InboxConversation, len(conversations))

	for _, conversation := range conversations {
		var conversationItem types.InboxConversation
		conversationItem.UserName = conversation.UserName
		conversationItem.ID = conversation.ID
		conversationItem.LastMsg = conversation.LastMsg
		conversationItem.LastMsgAt = conversation.LastMsgAt

		conversationList = append(conversationList, conversationItem)
	}

	// Convert to Item list
	itemList := make([]list.Item, len(conversationList))
	for i, _ := range conversationList {
		itemList[i] = &conversationList[i]
	}

	return itemList
}

func (m Model) updateConversationList() tea.Cmd {
	inboxResponse, err := m.client.GetInbox()
	if err != nil {
		logger.Log.Panicln("Failed to get inbox: " + err.Error())
	}

	conversationItemList := getConversationItemList(inboxResponse.Conversations)
	cmd := m.conversationList.SetItems(conversationItemList)

	// Now we update "global" vars
	user.UserData.UserID = inboxResponse.User.ID
	user.UserData.UserName = inboxResponse.User.Username

	return cmd
}

func (m Model) Init() tea.Cmd {
	return commands.NewUpdateInboxCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		return m, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case commands.UpdateInboxMsg:
		return m, m.updateConversationList()
	}

	m.conversationList, cmd = m.conversationList.Update(msg)

	return m, cmd
}
