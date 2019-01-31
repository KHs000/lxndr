package mongo

import (
	"context"
	"log"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// Connection ..
type Connection struct {
	Ctx    context.Context
	Client *mongo.Client
}

// Collection ..
type Collection struct {
	Database string
	CollName string
}

// Conn ..
var Conn *Connection

// Connect ..
func Connect(connStr string) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctx, connStr)
	if err != nil {
		log.Fatal(err)
	}

	Conn = &Connection{Ctx: ctx, Client: client}
}

// Search ..
func Search(conn *Connection, coll Collection, filter bson.M) mongo.Cursor {
	collection := conn.Client.Database(coll.Database).Collection(coll.CollName)
	res, err := collection.Find(conn.Ctx, filter)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Close(conn.Ctx)
	return res
}

// Insert ..
func Insert(conn *Connection, coll Collection, data interface{}) *mongo.InsertOneResult {
	collection := conn.Client.Database(coll.Database).Collection(coll.CollName)
	res, err := collection.InsertOne(conn.Ctx, data)
	if err != nil {
		log.Fatal(err)
	}

	return res
}

// Update ..
func Update(conn *Connection, coll Collection, filter bson.M,
	data interface{}) *mongo.UpdateResult {

	collection := conn.Client.Database(coll.Database).Collection(coll.CollName)
	res, err := collection.UpdateOne(conn.Ctx, filter, data)
	if err != nil {
		log.Fatal(err)
	}

	return res
}
