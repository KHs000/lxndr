package domain

import (
	"context"

	"github.com/mongodb/mongo-go-driver/mongo"
)

type (
	// Error ..
	Error struct {
		Code  int
		Error string
	}

	// Response ..
	Response struct {
		Message string   `json:"message"`
		Data    []string `json:"data"`
	}

	// User ..
	User struct {
		Email string
		Hash  string
		Token string
	}

	// Connection ..
	Connection struct {
		Ctx    context.Context
		Client *mongo.Client
	}

	// Collection ..
	Collection struct {
		Database string
		CollName string
	}

	// Client ..
	Client interface {
		Database(name string) DataLayer
		Close(ctx context.Context)
	}

	// DataLayer ..
	DataLayer interface {
		Collection(name string) Entities
	}

	// Entities ..
	Entities interface {
		Find(ctx context.Context, i interface{}) (Cursor, error)
	}

	// Cursor ..
	Cursor interface {
		Next(ctx context.Context) bool
		Decode(i interface{}) error
	}

	// MongoClient ..
	MongoClient struct {
		mongo.Client
	}

	// MongoDatabase ..
	MongoDatabase struct {
		*mongo.Database
	}

	// MongoCollection ..
	MongoCollection struct {
		*mongo.Collection
	}

	// MongoCursor ..
	MongoCursor struct {
		mongo.Cursor
	}
)

// Database ..
func (c MongoClient) Database(name string) DataLayer {
	return MongoDatabase{Database: c.Client.Database(name)}
}

// Collection ..
func (d MongoDatabase) Collection(name string) Entities {
	return MongoCollection{Collection: d.Database.Collection(name)}
}

// Find ..
func (c MongoCollection) Find(ctx context.Context, i interface{}) (Cursor, error) {
	cursor, err := c.Collection.Find(ctx, i)
	if err != nil {
		return MongoCursor{}, err
	}
	return MongoCursor{Cursor: cursor}, nil
}
