package httphandler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateUser(t *testing.T) {
	t.Run("should create a new user", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/createUser",
			strings.NewReader(`{"email": "felipe.carbone@dito.com.br"}`))
		w := httptest.NewRecorder()
		createUser(w, r)

		resp := w.Result()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error("Could not read request body")
		}

		body := Response{}
		err = json.Unmarshal(b, &body)
		if err != nil {
			t.Error("Could not parse request body")
		}
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expeceted status code to be %v, got %v",
				http.StatusCreated, resp.StatusCode)
		}
		if body.Message == "" {
			t.Error("Expeceted body to contain ObjectID, got nothing")
		}
	})
}
