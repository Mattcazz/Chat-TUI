package types

import "time"

type InboxConversation struct {
	UserName  string
	ID        int64
	LastMsg   string
	LastMsgAt time.Time
}

func (c InboxConversation) Title() string {
	return c.UserName
}

func (c InboxConversation) Description() string {
	if len(c.LastMsg) > 50 {
		return c.LastMsg[:50]
	}
	return c.LastMsg
}

// Assume searching is done exclusively by name and not by anything else
// Perhaps in the future we will add support for fuzzy finding inside a conversation from outside it
func (c InboxConversation) FilterValue() string {
	return c.UserName
}
