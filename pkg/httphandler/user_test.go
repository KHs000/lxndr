package httphandler

import (
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

	})
}
