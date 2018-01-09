package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const slackDefaultBaseURL = "https://hooks.slack.com/services/"

// URI interface
type URI interface {
	GetURL() string
}

// MessageSender interface
type MessageSender interface {
	SendMessage(message string) error
}

type asyncResult struct {
	success bool
	err     error
}

type jsonMessage struct {
	Data string `json:"text,omitempty"`
}

// API is used to post messages to a Slack webhook
type API struct {
	baseURL  string
	endpoint string
}

// NewAPI returns an instance of API prefilled with the Slack webhook base URL
func NewAPI(endpoint string) *API {
	return &API{
		slackDefaultBaseURL,
		endpoint,
	}
}

// GetURL returns the full API URL
func (api *API) GetURL() string {
	return api.baseURL + api.endpoint
}

// SendMessage sends a single message to a Slack service via HTTP
func (api *API) SendMessage(message string) error {
	return sendMessage(api, message)
}

// SendDataSynchronously sends an array of strings to Slack synchronously
func (api *API) SendDataSynchronously(data []string) error {
	return sendDataSynchronously(api, data)
}

// SendDataConcurrently sends an array of strings to Slack concurrently
func (api *API) SendDataConcurrently(data []string) []error {
	return sendDataConcurrently(api, data)
}

func sendMessage(uri URI, message string) error {
	body, err := json.Marshal(jsonMessage{message})

	if err != nil {
		return fmt.Errorf("marshalling failed with %s: %v", message, err)
	}

	resp, err := http.Post(uri.GetURL(), "application/json", bytes.NewReader(body))

	if err != nil {
		return fmt.Errorf("%s request failed: %v", uri.GetURL(), err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("%s did not respond with a 200 OK", uri.GetURL())
	}

	defer resp.Body.Close()

	return nil
}

func sendDataSynchronously(sender MessageSender, data []string) error {
	for _, text := range data {
		err := sender.SendMessage(text)

		if err != nil {
			return err
		}
	}

	return nil
}

func sendDataConcurrently(sender MessageSender, data []string) []error {
	errors := make([]error, 0)

	ch := make(chan asyncResult, len(data))
	defer close(ch)

	for _, text := range data {
		go func(text string) {
			err := sender.SendMessage(text)

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
