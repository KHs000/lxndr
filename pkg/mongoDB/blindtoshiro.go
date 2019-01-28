package mongoDB

import (
	"context"
	"log"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// Connect ..
func Connect() (context.Context, *mongo.Client) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, "mongodb://admin:admin123456@felipe-rabelo-shard-00-00-r4yae.gcp.mongodb.net:27017,felipe-rabelo-shard-00-01-r4yae.gcp.mongodb.net:27017,felipe-rabelo-shard-00-02-r4yae.gcp.mongodb.net:27017/test?ssl=true&replicaSet=felipe-rabelo-shard-0&authSource=admin&retryWrites=true")
	if err != nil {
		log.Fatal(err)
	}

	return ctx, client
}

// Find ..
func Find(ctx context.Context, client *mongo.Client) mongo.Cursor {
	filter := bson.M{}

	collection := client.Database("lxndr").Collection("lxndr-quest")
	res, err := collection.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Close(ctx)

	return res
}
