package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const slackDefaultBaseURL = "https://hooks.slack.com/services/"

// API is used to post messages to a Slack webhook
type API struct {
	baseURL string
	Endpoint
}

// Endpoint is used to store the RequestURI
type Endpoint struct {
	URL string
}

type jsonMessage struct {
	Data string `json:"text,omitempty"`
}

// SendMessage sends a single message to a Slack service via HTTP
func (slack *API) SendMessage(message string) error {
	if slack.baseURL == "" {
		slack.baseURL = slackDefaultBaseURL
	}

	body, err := json.Marshal(jsonMessage{message})

	if err != nil {
		return err
	}

	resp, err := http.Post(slack.baseURL+slack.URL, "application/json", bytes.NewReader(body))

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Something went wrong, Slack did not respond with a 200 OK")
	}

	defer resp.Body.Close()

	return nil
}
