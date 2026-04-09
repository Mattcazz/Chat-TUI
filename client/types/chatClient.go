package types

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Mattcazz/Chat-TUI/client/internal/logger"
	"github.com/Mattcazz/Chat-TUI/pkg"
)

type ChatClient struct {
	Client BaseClient
}

func (c *ChatClient) GetChat(conversationId int64) (pkg.ConversationResponse, error) {
	logger.Log.Printf("[ChatClient] Getting inbox")

	var conversationResponse pkg.ConversationResponse
	requestPath := fmt.Sprintf("conversation/%d", conversationId)
	resp, err := c.Client.doRequest("GET", requestPath, nil, &conversationResponse)
	if err != nil {
		logger.Log.Panic("Failed to get conversation: " + err.Error())
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return conversationResponse, nil
	case http.StatusUnauthorized:
		return pkg.ConversationResponse{}, errors.New("Unauthorized")
	default:
		return pkg.ConversationResponse{}, errors.New("Received an unsupported response status: " + resp.Status);
	}
}

func (c *ChatClient) SendMessage(conversationId int64, message string) error {
	logger.Log.Printf("[ChatClient] Sending message '%s' to conversation with ID %d", message, conversationId)

	requestPath := fmt.Sprintf("conversation/%d/message", conversationId)
	body := map[string]any{
		"content": message,
	}

	_, err := c.Client.doRequest("POST", requestPath, body, nil)

	return err
}
