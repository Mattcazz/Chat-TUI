package commands

import (
	"github.com/Mattcazz/Chat-TUI/client/internal/config"
	"github.com/Mattcazz/Chat-TUI/client/types"
	tea "github.com/charmbracelet/bubbletea"
)

type SetConfigMsg struct {
	Config *config.Config
	Client *types.BaseClient
}

func NewSetConfigCmd(config *config.Config, client *types.BaseClient) func() tea.Msg {
	var msg SetConfigMsg
	msg.Config = config
	msg.Client = client

	return func() tea.Msg {
		return msg
	}
}
