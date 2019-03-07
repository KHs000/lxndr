package httphandler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestProcessRequestBody ..
func TestProcessRequestBody(t *testing.T) {
	t.Run("process body should work", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/",
			strings.NewReader(`{"email": "felipe.carbone@dito.com.br"}`))

		type testBody struct {
			email string
		}
		i := testBody{}
		b, err := processRequestBody(r, i)
		if err != nil {
			t.Errorf("Expect error to be nil, got %v", err)
		}

		if b["email"] != "felipe.carbone@dito.com.br" {
			t.Errorf("Expected email to be 'felipe.carbone@dito.com.br', got %v", b["email"])
		}
	})
}

// TestDefaultRoute ..
func TestDefaultRoute(t *testing.T) {
	t.Run("default route should work", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		defaultRoute(w, r)

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
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected StatusCode to be 200, got %v", resp.StatusCode)
		}
		if body.Message != "It works..." {
			t.Errorf("Expected body to be 'It works...', got %v", body.Message)
		}
		if body.Data != nil {
			t.Errorf("Expected data to be nil, got %v", body.Data)
		}
	})

	t.Run("default route should be of GET method", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/", nil)
		w := httptest.NewRecorder()
		defaultRoute(w, r)

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
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected StatusCode to be 400, got %v", resp.StatusCode)
		}
		if body.Message != "Method not allowed." {
			t.Errorf("Expected body to be 'Method not allowed.', got %v", body.Message)
		}
		if body.Data != nil {
			t.Errorf("Expected data to be nil, got %v", body.Data)
		}
	})
}
