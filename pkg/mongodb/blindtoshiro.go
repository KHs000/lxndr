package mongodb

import (
	"context"
	"log"

	"github.com/KHs000/lxndr/domain"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// Conn ..
var Conn *domain.Connection

// Client ..
var Client domain.MongoClient

// NewClient ..
func NewClient(connStr string) domain.Client {
	ctx := context.Background()

	client, err := mongo.Connect(ctx, connStr)
	if err != nil {
		log.Fatal(err)
	}

	return domain.MongoClient{Client: client, Context: ctx}
}

// Connect ..
func Connect(connStr string) {
	ctx := context.Background()

	client, err := mongo.Connect(ctx, connStr)
	if err != nil {
		log.Fatal(err)
	}

	Conn = &domain.Connection{Ctx: ctx, Client: client}

	if ctx.Err() != nil {
		log.Println(ctx.Err())
	}
}

// Test ..
func Test(client domain.Client, coll domain.Collection, filter interface{}) domain.Cursor {
	collection := client.Database(coll.Database).Collection(coll.CollName)
	res, err := collection.Find(client.Ctx(), filter)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Close(client.Ctx())
	return res
}

// Search ..
func Search(client domain.Client, coll domain.Collection,
	filter bson.M) domain.Cursor {
	collection := client.Database(coll.Database).Collection(coll.CollName)
	res, err := collection.Find(client.Ctx(), filter)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Close(client.Ctx())
	return res
}

// Insert ..
func Insert(client domain.Client, coll domain.Collection,
	data interface{}) (primitive.ObjectID, error) {
	collection := client.Database(coll.Database).Collection(coll.CollName)
	res, err := collection.InsertOne(client.Ctx(), data)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	return res.ObjectID, nil
}

// Update ..
func Update(conn *domain.Connection, coll domain.Collection, filter bson.M,
	data interface{}) *mongo.UpdateResult {

	collection := conn.Client.Database(coll.Database).Collection(coll.CollName)
	res, err := collection.UpdateOne(conn.Ctx, filter, data)
	if err != nil {
		log.Fatal(err)
	}

	return res
}

// Delete ..
func Delete(conn *domain.Connection, coll domain.Collection,
	filter bson.M) *mongo.DeleteResult {
	collection := conn.Client.Database(coll.Database).Collection(coll.CollName)
	res, err := collection.DeleteOne(conn.Ctx, filter)
	if err != nil {
		log.Fatal(err)
	}

	return res
}
