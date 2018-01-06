package slack

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

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

	slack := API{
		testServer.URL,
		Endpoint{serviceEndpoint},
	}
	err := slack.SendMessage("Ciao!")

	if err == nil {
		t.Error("Function did not return error on 500")
	}
}
