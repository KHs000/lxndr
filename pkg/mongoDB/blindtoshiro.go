package mongoDB

import (
	"context"
	"log"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// Conn ..
type Conn struct {
	Ctx    context.Context
	Client *mongo.Client
}

// Collection ..
type Collection struct {
	Database string
	CollName string
}

// Connect ..
func Connect() *Conn {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, "mongodb://admin:admin123456@felipe-rabelo-shard-00-00-r4yae.gcp.mongodb.net:27017,felipe-rabelo-shard-00-01-r4yae.gcp.mongodb.net:27017,felipe-rabelo-shard-00-02-r4yae.gcp.mongodb.net:27017/test?ssl=true&replicaSet=felipe-rabelo-shard-0&authSource=admin&retryWrites=true")
	if err != nil {
		log.Fatal(err)
	}

	return &Conn{Ctx: ctx, Client: client}
}

// Find ..
func Find(conn *Conn, coll Collection, filter bson.M) mongo.Cursor {
	collection := conn.Client.Database(coll.Database).Collection(coll.CollName)
	res, err := collection.Find(conn.Ctx, filter)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Close(conn.Ctx)
	return res
}
