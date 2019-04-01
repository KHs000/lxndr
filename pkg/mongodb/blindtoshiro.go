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
func Update(client domain.Client, coll domain.Collection, filter bson.M,
	data interface{}) (domain.MongoUpdate, error) {

	collection := client.Database(coll.Database).Collection(coll.CollName)
	res, err := collection.UpdateOne(client.Ctx(), filter, data)
	if err != nil {
		return domain.MongoUpdate{}, err
	}

	return res, nil
}

// Delete ..
func Delete(client domain.Client, coll domain.Collection,
	filter bson.M) (domain.MongoDelete, error) {
	collection := client.Database(coll.Database).Collection(coll.CollName)
	res, err := collection.DeleteOne(client.Ctx(), filter)
	if err != nil {
		return domain.MongoDelete{}, err
	}

	return domain.MongoDelete{DeletedCount: int(res.DeletedCount)}, nil
}
