package httphandler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/KHs000/lxndr/domain"
)

// TestProcessRequestBody ..
func TestProcessRequestBody(t *testing.T) {
	t.Run("should work", func(t *testing.T) {
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
			t.Errorf("Expected email to be 'felipe.carbone@dito.com.br', got %v",
				b["email"])
		}
	})

	t.Run("should return an error for invalid json", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/",
			strings.NewReader(`{"email: "felipe.carbone@dito.com.br"}`))

		type testBody struct {
			email string
		}
		i := testBody{}
		b, err := processRequestBody(r, i)
		if b != nil {
			t.Errorf("Expected body to be nil, got %v", b)
		}
		if err.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code to be %v, got %v",
				http.StatusInternalServerError, err.Code)
		}
		if err.Error != "Oh-uh, something's not quite right." {
			t.Errorf("Expected error to be 'Oh-uh, something's not quite right.', got %v",
				err.Error)
		}
	})
}

// TestDefaultRoute ..
func TestDefaultRoute(t *testing.T) {
	t.Run("should work", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		defaultRoute(w, r)

		resp := w.Result()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error("Could not read request body")
		}

		body := domain.Response{}
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

	t.Run("should be of GET method", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/", nil)
		w := httptest.NewRecorder()
		defaultRoute(w, r)

		resp := w.Result()
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error("Could not read request body")
		}

		body := domain.Response{}
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
