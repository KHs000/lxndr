package httphandler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/KHs000/lxndr/domain"
	"github.com/KHs000/lxndr/domain/mocks"
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

		body := domain.Response{}
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

func TestTest(t *testing.T) {
	t.Run("should return a user list", func(t *testing.T) {
		controlCursor := true
		cursor := mocks.MockCursor{
			NextFn: func(ctx context.Context) bool {
				hasNext := controlCursor
				controlCursor = !controlCursor
				return hasNext
			},
			CloseFn: func(ctx context.Context) error {
				return nil
			},
			DecodeCursorFn: func() (map[string]interface{}, error) {
				fakeResult := map[string]interface{}{"email": "felipe.carbone@dito.com.br"}
				return fakeResult, nil
			},
		}

		collection := mocks.MockCollection{
			FindFn: func(ctx context.Context, i interface{}) (domain.Cursor, error) {
				return mocks.MockCursor{
					NextFn:         cursor.NextFn,
					CloseFn:        cursor.CloseFn,
					DecodeCursorFn: cursor.DecodeCursorFn,
				}, nil
			},
		}

		database := mocks.MockDatabase{
			CollectionFn: func(name string) domain.Entities {
				return mocks.MockCollection{FindFn: collection.FindFn}
			},
		}

		client := mocks.MockClient{
			DatabaseFn: func(name string) domain.DataLayer {
				return mocks.MockDatabase{CollectionFn: database.CollectionFn}
			},
			CtxFn: func() context.Context {
				return context.Background()
			},
		}

		resp := test(client)
		if len(resp.Data) != 1 {
			t.Errorf("Expected resp.Data to have lenght 1, got %v", len(resp.Data))
		}
	})
}
