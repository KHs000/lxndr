package domain

import (
	"context"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"gopkg.in/mgo.v2/bson"
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
		Ctx() context.Context
	}

	// DataLayer ..
	DataLayer interface {
		Collection(name string) Entities
	}

	// Entities ..
	Entities interface {
		Find(ctx context.Context, i interface{}) (Cursor, error)
		InsertOne(ctx context.Context, i interface{}) (MongoInsert, error)
	}

	// Cursor ..
	Cursor interface {
		Next(ctx context.Context) bool
		Close(ctx context.Context) error
		Decode(i interface{}) error
		DecodeCursor() (map[string]interface{}, error)
	}

	// MongoClient ..
	MongoClient struct {
		Client  *mongo.Client
		Context context.Context
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

	// MongoInsert ..
	MongoInsert struct {
		ObjectID primitive.ObjectID
	}
)

// Database ..
func (c MongoClient) Database(name string) DataLayer {
	return MongoDatabase{Database: c.Client.Database(name)}
}

// Ctx ..
func (c MongoClient) Ctx() context.Context {
	return c.Context
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

// InsertOne ..
func (c MongoCollection) InsertOne(ctx context.Context, i interface{}) (MongoInsert, error) {
	result, err := c.Collection.InsertOne(ctx, i)
	if err != nil {
		return MongoInsert{}, err
	}
	return MongoInsert{ObjectID: result.InsertedID.(primitive.ObjectID)}, nil
}

// Next ..
func (c MongoCursor) Next(ctx context.Context) bool {
	return c.Cursor.Next(ctx)
}

// Close ..
func (c MongoCursor) Close(ctx context.Context) error {
	return c.Cursor.Close(ctx)
}

// Decode ..
func (c MongoCursor) Decode(i interface{}) error {
	return c.Cursor.Decode(i)
}

// DecodeCursor ..
func (c MongoCursor) DecodeCursor() (map[string]interface{}, error) {
	var row bson.M
	err := c.Cursor.Decode(&row)
	if err != nil {
		return nil, err
	}
	return row, nil
}
