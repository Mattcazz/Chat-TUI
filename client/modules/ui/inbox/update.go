package inbox

import (
	"github.com/Mattcazz/Chat-TUI/client/internal/commands"
	"github.com/Mattcazz/Chat-TUI/client/internal/logger"
	"github.com/Mattcazz/Chat-TUI/client/internal/user"
	"github.com/Mattcazz/Chat-TUI/client/styles"
	"github.com/Mattcazz/Chat-TUI/client/types"
	"github.com/Mattcazz/Chat-TUI/pkg"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func getConversationItemList(conversations []pkg.InboxConversationResponse) []list.Item {
	itemList := make([]list.Item, 0, len(conversations))

	for _, conversation := range conversations {
		itemList = append(itemList, &types.InboxConversation{
			UserName: conversation.UserName,
			ID: conversation.ID,
			LastMsg: conversation.LastMsg,
			LastMsgAt: conversation.LastMsgAt,
		})
		logger.Log.Printf("New conversation with %s logged", conversation.UserName)
	}

	return itemList
}

func (m *Model) updateConversationList() tea.Cmd {
	logger.Log.Printf("[INBOX] Updating conversation list")
	inboxResponse, err := m.client.GetInbox()
	if err != nil {
		logger.Log.Panicf("Failed to get inbox: %s", err.Error())
	}

	conversationItemList := getConversationItemList(inboxResponse.Conversations)
	cmd := m.conversationList.SetItems(conversationItemList)

	// Now we update "global" vars
	user.UserData.UserID = inboxResponse.User.ID
	user.UserData.UserName = inboxResponse.User.Username
	logger.Log.Printf("Setting global username to '%s'", inboxResponse.User.Username)

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
		conversationListWidth := m.width - styles.Default.Border.GetHorizontalFrameSize()
		conversationListHeight := m.height - styles.Default.Border.GetVerticalFrameSize() - 1 // TODO -1 is for "Inbox" title, replace with its own model
		m.conversationList.SetSize(conversationListWidth, conversationListHeight)

		return m, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			selectedItem := m.conversationList.SelectedItem()
			inboxConversation, ok := selectedItem.(*types.InboxConversation)
			if !ok {
				// Invalid item?
				logger.Log.Panicf("Unable to cast item '%s' back to InboxConversation", selectedItem.FilterValue())
			}
			
			openChatCmd := commands.NewOpenChatCmd(inboxConversation.UserName, inboxConversation.ID)
			return m, openChatCmd
		}
	case commands.UpdateInboxMsg:
		return m, m.updateConversationList()
	}

	m.conversationList, cmd = m.conversationList.Update(msg)

	return m, cmd
}
