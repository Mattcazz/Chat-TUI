package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Mattcazz/Chat-TUI/client/internal/config"
)

type BaseClient struct {
	Client http.Client
	Config config.Config
}

// doRequest is the shared helper that all sub-clients use.
// It handles JSON encoding/decoding and setting the Auth Header.
func (c *BaseClient) doRequest(method string, path string, body interface{}, target interface{}) (*http.Response, error) {
	var bodyReader io.Reader

	// 1. Encode Body (if any)
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBytes)
	}

	// 2. Create Request
	url := fmt.Sprintf("http://%s:%s/%s", c.Config.Network.ServerHost, c.Config.Network.ServerPort, path)
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 3. Set Headers (Content-Type & Auth)
	req.Header.Set("Content-Type", "application/json")
	if c.Config.Jwt != "" {
		req.Header.Set("Authorization", c.Config.Jwt)
	}

	// 4. Execute
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// 5. Decode Response only on success
	if target != nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return resp, nil
}
