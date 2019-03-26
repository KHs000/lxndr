package httphandler

import (
	"context"
	"net/http"
	"testing"

	"github.com/KHs000/lxndr/domain"
	"github.com/KHs000/lxndr/domain/mocks"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

func TestCreateUser(t *testing.T) {
	fakeEmail := map[string]interface{}{"email": "fake@email.com"}

	t.Run("should create a new user", func(t *testing.T) {
		client := defaultClientFn(
			defaultDatabaseFn(
				defaultCollectionFn(
					defaultFindFn(
						defaultNextFn(false),
						defaultCloseFn(),
						nil,
					),
					defaultInsertFn(primitive.ObjectID{}),
					nil,
				),
			),
		)

		code, _ := createUser(client, fakeEmail)

		if code != http.StatusCreated {
			t.Errorf("Expected status code to be 201, got %v", code)
		}
	})

	t.Run("should not allow repeated emails", func(t *testing.T) {
		client := defaultClientFn(
			defaultDatabaseFn(
				defaultCollectionFn(
					defaultFindFn(
						defaultNextFn(true),
						defaultCloseFn(),
						defaultDecodeCursorFn(map[string]interface{}{"_id": "fakeID"}),
					),
					nil, nil,
				),
			),
		)

		code, _ := createUser(client, fakeEmail)

		if code != http.StatusConflict {
			t.Errorf("Expected status code to be 409, got %v", code)
		}
	})
}

func TestEditUser(t *testing.T) {
	fakeEdit := map[string]interface{}{"email": "fake@email.com", "hash": "fakeHash"}

	t.Run("should edit a user based on a valid email", func(t *testing.T) {
		customNextFn := func(control bool) func(ctx context.Context) bool {
			return func(ctx context.Context) bool {
				hasNext := control
				control = !control
				return hasNext
			}
		}

		client := defaultClientFn(
			defaultDatabaseFn(
				defaultCollectionFn(
					defaultFindFn(
						customNextFn(true),
						defaultCloseFn(),
						defaultDecodeCursorFn(map[string]interface{}{"_id": "fakeID"}),
					),
					nil,
					defaultUpdateFn(1),
				),
			),
		)

		code, _ := editUser(client, fakeEdit)
		if code != http.StatusOK {
			t.Errorf("Expected status code to be 200, got %v", code)
		}
	})

	t.Run("should not allow edit for unregistered emails", func(t *testing.T) {
		client := defaultClientFn(
			defaultDatabaseFn(
				defaultCollectionFn(
					defaultFindFn(
						defaultNextFn(false),
						defaultCloseFn(),
						nil,
					),
					nil, nil,
				),
			),
		)

		code, _ := editUser(client, fakeEdit)
		if code != http.StatusNotFound {
			t.Errorf("Expected status code to be 404, got %v", code)
		}
	})
}

func TestTest(t *testing.T) {
	t.Run("should return a user list", func(t *testing.T) {
		customNextFn := func(control bool) func(ctx context.Context) bool {
			return func(ctx context.Context) bool {
				hasNext := control
				control = !control
				return hasNext
			}
		}

		client := defaultClientFn(
			defaultDatabaseFn(
				defaultCollectionFn(
					defaultFindFn(
						customNextFn(true),
						defaultCloseFn(),
						defaultDecodeCursorFn(
							map[string]interface{}{"email": "test@email.com"},
						),
					),
					nil, nil,
				),
			),
		)

		resp := test(client)
		if len(resp.Data) != 1 {
			t.Errorf("Expected resp.Data to have lenght 1, got %v", len(resp.Data))
		}
	})
}

func defaultClientFn(databaseFn func(name string) domain.DataLayer) domain.Client {
	return mocks.MockClient{
		DatabaseFn: databaseFn,
		CtxFn:      defaultCtxFn(),
	}
}

func defaultDatabaseFn(
	collectionFn func(name string) domain.Entities,
) func(name string) domain.DataLayer {
	return func(name string) domain.DataLayer {
		return mocks.MockDatabase{
			CollectionFn: collectionFn,
		}
	}
}

func defaultCtxFn() func() context.Context {
	return func() context.Context {
		return context.Background()
	}
}

func defaultCollectionFn(
	findFn func(
		context.Context,
		interface{},
	) (domain.Cursor, error),

	insertFn func(
		ctx context.Context,
		i interface{},
	) (domain.MongoInsert, error),

	updateFn func(
		ctx context.Context,
		filter bson.M,
		i interface{},
	) (domain.MongoUpdate, error),
) func(name string) domain.Entities {
	return func(name string) domain.Entities {
		return mocks.MockCollection{
			FindFn:      findFn,
			InsertOneFn: insertFn,
			UpdateOneFn: updateFn,
		}
	}
}

func defaultFindFn(
	nextFn func(context.Context) bool,
	closeFn func(context.Context) error,
	decodeFn func() (map[string]interface{}, error),
) func(context.Context, interface{}) (domain.Cursor, error) {
	return func(ctx context.Context, i interface{}) (domain.Cursor, error) {
		return mocks.MockCursor{
			NextFn:         nextFn,
			CloseFn:        closeFn,
			DecodeCursorFn: decodeFn,
		}, nil
	}
}

func defaultInsertFn(
	pID primitive.ObjectID,
) func(ctx context.Context, i interface{}) (domain.MongoInsert, error) {
	return func(ctx context.Context, i interface{}) (domain.MongoInsert, error) {
		fakeResult := domain.MongoInsert{ObjectID: pID}
		return fakeResult, nil
	}
}

func defaultUpdateFn(count int) func(
	ctx context.Context,
	filter bson.M,
	i interface{},
) (domain.MongoUpdate, error) {
	return func(
		ctx context.Context,
		filter bson.M,
		i interface{},
	) (domain.MongoUpdate, error) {
		fakeResult := domain.MongoUpdate{MatchedCount: count}
		return fakeResult, nil
	}
}

func defaultNextFn(nextReturn bool) func(context.Context) bool {
	return func(ctx context.Context) bool {
		return nextReturn
	}
}

func defaultCloseFn() func(context.Context) error {
	return func(ctx context.Context) error {
		return nil
	}
}

func defaultDecodeCursorFn(
	fakeResult map[string]interface{},
) func() (map[string]interface{}, error) {
	return func() (map[string]interface{}, error) {
		return fakeResult, nil
	}
}
