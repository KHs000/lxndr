package mongodb

import (
	"context"
	"log"

	"github.com/KHs000/lxndr/domain"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// Conn ..
var Conn *domain.Connection

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
func Search(conn *domain.Connection, coll domain.Collection,
	filter bson.M) mongo.Cursor {
	collection := conn.Client.Database(coll.Database).Collection(coll.CollName)
	res, err := collection.Find(conn.Ctx, filter)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Close(conn.Ctx)
	return res
}

// Insert ..
func Insert(conn *domain.Connection, coll domain.Collection,
	data interface{}) *mongo.InsertOneResult {
	collection := conn.Client.Database(coll.Database).Collection(coll.CollName)
	res, err := collection.InsertOne(conn.Ctx, data)
	if err != nil {
		log.Fatal(err)
	}

	return res
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
