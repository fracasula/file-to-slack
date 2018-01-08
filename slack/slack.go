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

type asyncResult struct {
	success bool
	err     error
}

// SendMessage sends a single message to a Slack service via HTTP
func (slack *API) SendMessage(message string) error {
	if slack.baseURL == "" {
		slack.baseURL = slackDefaultBaseURL
	}

	body, err := json.Marshal(jsonMessage{message})

	if err != nil {
		return fmt.Errorf("marshalling failed with %s: %v", message, err)
	}

	resp, err := http.Post(slack.baseURL+slack.URL, "application/json", bytes.NewReader(body))

	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Slack API did not respond with a 200 OK")
	}

	defer resp.Body.Close()

	return nil
}

// SendDataSynchronously sends an array of strings to Slack synchronously
func (slack *API) SendDataSynchronously(data []string) error {
	for _, text := range data {
		err := slack.SendMessage(text)

		if err != nil {
			return err
		}
	}

	return nil
}

// SendDataConcurrently sends an array of strings to Slack concurrently
func (slack *API) SendDataConcurrently(data []string) []error {
	errors := make([]error, 0)

	ch := make(chan asyncResult, len(data))
	defer close(ch)

	for _, text := range data {
		go func(text string) {
			err := slack.SendMessage(text)

			if err != nil {
				ch <- asyncResult{false, err}
			} else {
				ch <- asyncResult{true, nil}
			}
		}(text)
	}

	for i := 0; i < len(data); i++ {
		result, ok := <-ch

		if !result.success && ok {
			errors = append(errors, result.err)
		}

		if !ok {
			break
		}
	}

	return errors
}
