package slack

import (
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
	originalMessages := []string{"one", "two", "three"}

	err := sendDataSynchronously(&sender, []string{"one", "two", "three"})

	if err != nil {
		t.Error("Expected error to be nil")
	}

	expectedMessages := fmt.Sprintf("%v", originalMessages)
	actualMessages := fmt.Sprintf("%v", sender.messages)

	if expectedMessages != actualMessages {
		t.Errorf("Expected messages should be %s got %s", expectedMessages, actualMessages)
	}
}
