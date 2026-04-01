package types

import (
	"errors"
	"net/http"

	"github.com/Mattcazz/Chat-TUI/client/internal/logger"
	"github.com/Mattcazz/Chat-TUI/pkg"
)

type InboxClient struct {
	Client BaseClient
}

func (c *InboxClient) GetInbox() (pkg.InboxResponse, error) {
	logger.Log.Printf("[InboxClient] Getting inbox")

	var inboxResponse pkg.InboxResponse
	resp, err := c.Client.doRequest("GET", "inbox", nil, &inboxResponse)
	if err != nil {
		logger.Log.Panic("Failed to get inbox: " + err.Error())
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return inboxResponse, nil
	case http.StatusUnauthorized:
		return pkg.InboxResponse{}, errors.New("Unauthorized")
	default:
		return pkg.InboxResponse{}, errors.New("Received an unsupported response status: " + resp.Status);
	}
}
