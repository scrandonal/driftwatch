package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Payload represents the data sent to the webhook endpoint.
type Payload struct {
	Event     string    `json:"event"`
	FilePath  string    `json:"file_path"`
	Timestamp time.Time `json:"timestamp"`
	OldHash   string    `json:"old_hash,omitempty"`
	NewHash   string    `json:"new_hash,omitempty"`
}

// Client sends webhook notifications to a configured URL.
type Client struct {
	URL        string
	HTTPClient *http.Client
}

// New creates a new webhook Client with the given target URL.
func New(url string) *Client {
	return &Client{
		URL: url,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Send encodes the payload as JSON and POSTs it to the webhook URL.
func (c *Client) Send(p Payload) error {
	if p.Timestamp.IsZero() {
		p.Timestamp = time.Now().UTC()
	}

	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}

	resp, err := c.HTTPClient.Post(c.URL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: post request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status code %d", resp.StatusCode)
	}

	return nil
}
