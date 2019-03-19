package mocks

import (
	"context"

	"github.com/KHs000/lxndr/domain"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type (
	// MockClient ..
	MockClient struct {
		*mongo.Client
		Context context.Context

		DatabaseFn func(name string) domain.DataLayer
		CtxFn      func() context.Context
	}

	// MockDatabase ..
	MockDatabase struct {
		*mongo.Database

		CollectionFn func(name string) domain.Entities
	}

	// MockCollection ..
	MockCollection struct {
		*mongo.Collection

		FindFn func(ctx context.Context, i interface{}) (domain.Cursor, error)
	}

	// MockCursor ..
	MockCursor struct {
		mongo.Cursor

		NextFn         func(ctx context.Context) bool
		CloseFn        func(ctx context.Context) error
		DecodeFn       func(i interface{}) error
		DecodeCursorFn func() (map[string]interface{}, error)
	}
)

// Database ..
func (c MockClient) Database(name string) domain.DataLayer {
	return c.DatabaseFn(name)
}

// Ctx ..
func (c MockClient) Ctx() context.Context {
	return c.CtxFn()
}

// Collection ..
func (c MockDatabase) Collection(name string) domain.Entities {
	return c.CollectionFn(name)
}

// Find ..
func (c MockCollection) Find(ctx context.Context, i interface{}) (domain.Cursor, error) {
	cursor, err := c.FindFn(ctx, i)
	if err != nil {
		return MockCursor{}, err
	}
	return cursor, nil
}

// Next ..
func (c MockCursor) Next(ctx context.Context) bool {
	return c.NextFn(ctx)
}

// Close ..
func (c MockCursor) Close(ctx context.Context) error {
	return c.CloseFn(ctx)
}

// Decode ..
func (c MockCursor) Decode(i interface{}) error {
	return c.DecodeFn(i)
}

// DecodeCursor ..
func (c MockCursor) DecodeCursor() (map[string]interface{}, error) {
	return c.DecodeCursorFn()
}
