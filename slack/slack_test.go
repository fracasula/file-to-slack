package slack

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockedSender struct {
	err      error
	messages []string
}

func (sender *mockedSender) SendMessage(message string) error {
	sender.messages = append(sender.messages, message)

	return sender.err
}

func inSlice(needle string, haystack []string) bool {
	for _, elm := range haystack {
		if elm == needle {
			return true
		}
	}

	return false
}

func TestNewAPI(t *testing.T) {
	endpoint := "/a/nice/endpoint"
	api := NewAPI(endpoint)

	if api.baseURL != slackDefaultBaseURL {
		t.Errorf("Expected baseURL to be %s got %s", slackDefaultBaseURL, api.baseURL)
	}

	if api.endpoint != endpoint {
		t.Errorf("Expected endpoint to be %s got %s", endpoint, api.endpoint)
	}

	if api.GetURL() != slackDefaultBaseURL+endpoint {
		t.Errorf("Expected GetURL() to return %s got %s", slackDefaultBaseURL+endpoint, api.GetURL())
	}
}

func TestSendsMessageReturnsErrorOn500(t *testing.T) {
	const serviceEndpoint = "/this/is/a/test"

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Error("Request method should be POST")
		}

		if r.RequestURI != serviceEndpoint {
			t.Error("Request URI should be " + serviceEndpoint)
		}

		if body, err := ioutil.ReadAll(r.Body); err != nil {
			if string(body) != "{\"text\":\"Ciao!\"}" {
				t.Error("Request body not as expected")
			}
		}

		w.WriteHeader(http.StatusInternalServerError)
	}))

	defer testServer.Close()

	api := &API{
		testServer.URL,
		serviceEndpoint,
	}
	err := api.SendMessage("Ciao!")

	if err == nil {
		t.Error("Function did not return error on 500")
	}
}

func TestSendsAllMessagesSynchronously(t *testing.T) {
	mockedMessages := make([]string, 0)
	sender := mockedSender{nil, mockedMessages}

	messagesToSend := []string{"one", "two", "three"}
	err := sendDataSynchronously(&sender, messagesToSend)

	if err != nil {
		t.Error("Expected error to be nil")
	}

	expectedMessages := fmt.Sprintf("%v", messagesToSend)
	actualMessages := fmt.Sprintf("%v", sender.messages)

	if expectedMessages != actualMessages {
		t.Errorf("Expected messages should be %s got %s", expectedMessages, actualMessages)
	}
}

func TestSendsAllMessagesSynchronouslyWithError(t *testing.T) {
	mockedMessages := make([]string, 0)
	sender := mockedSender{errors.New("Just a random error"), mockedMessages}

	messagesToSend := []string{"one", "two", "three"}
	err := sendDataSynchronously(&sender, messagesToSend)

	if err.Error() != "Just a random error" {
		t.Error("Expected error not to be nil")
	}

	expectedMessages := fmt.Sprintf("%v", []string{"one"})
	actualMessages := fmt.Sprintf("%v", sender.messages)

	if expectedMessages != actualMessages {
		t.Errorf("Expected messages should be %s got %s", expectedMessages, actualMessages)
	}
}

func TestSendsAllMessagesConcurrently(t *testing.T) {
	mockedMessages := make([]string, 0)
	sender := mockedSender{nil, mockedMessages}

	messagesToSend := []string{"one", "two", "three"}
	errors := sendDataConcurrently(&sender, messagesToSend)

	if errors != nil {
		t.Errorf("Expected errors to be nil, got %v", errors)
	}

	for _, actualMessage := range sender.messages {
		if !inSlice(actualMessage, messagesToSend) {
			t.Errorf("Expected message %s not found", actualMessage)
		}
	}
}

func TestSendsAllMessagesConcurrentlyWithError(t *testing.T) {
	err := errors.New("Just a random error")
	mockedMessages := make([]string, 0)
	sender := mockedSender{err, mockedMessages}

	messagesToSend := []string{"one", "two", "three"}
	errors := sendDataConcurrently(&sender, messagesToSend)

	expectedErrors := fmt.Sprintf("%v", []error{err, err, err})
	actualErrors := fmt.Sprintf("%v", errors)

	if expectedErrors != actualErrors {
		t.Errorf("Expected errors to be %s, got %s", expectedErrors, actualErrors)
	}

	for _, actualMessage := range sender.messages {
		if !inSlice(actualMessage, messagesToSend) {
			t.Errorf("Expected message %s not found", actualMessage)
		}
	}
}
