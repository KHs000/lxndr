package httphandler

import (
	"context"
	"testing"

	"github.com/KHs000/lxndr/domain"
	"github.com/KHs000/lxndr/domain/mocks"
)

func newMockClient() mocks.MockClient {
	return mocks.MockClient{
		DatabaseFn: func(name string) domain.DataLayer {
			return mocks.MockDatabase{
				CollectionFn: func(name string) domain.Entities {
					return mocks.MockCollection{}
				},
			}
		},
		CtxFn: func() context.Context {
			return context.Background()
		},
	}
}

func TestCreateUser(t *testing.T) {
	t.Run("should create a new user", func(t *testing.T) {
		// N tem como reaproveitar (fácil) o client
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

		client := mocks.MockClient{
			DatabaseFn: func(name string) domain.DataLayer {
				return mocks.MockDatabase{
					CollectionFn: func(name string) domain.Entities {
						return mocks.MockCollection{
							FindFn: func(ctx context.Context, i interface{}) (domain.Cursor, error) {
								return mocks.MockCursor{
									NextFn:         cursor.NextFn,
									CloseFn:        cursor.CloseFn,
									DecodeCursorFn: cursor.DecodeCursorFn,
								}, nil
							},
						}
					},
				}
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
