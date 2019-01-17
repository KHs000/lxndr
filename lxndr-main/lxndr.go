package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/KHs000/lxndr/identifier"
	"github.com/KHs000/lxndr/rndtoken"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

func main() {
	tkn, hash := rndtoken.SendToken("felipe.carbone@dito.com.br")

	fmt.Println("Token")
	fmt.Println(tkn)
	fmt.Println("Hash")
	fmt.Println(hash)

	identifier.IdentityCheck("felipe.carbone@dito.com.br", hash)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, "mongodb://admin:admin123456@felipe-rabelo-shard-00-00-r4yae.gcp.mongodb.net:27017,felipe-rabelo-shard-00-01-r4yae.gcp.mongodb.net:27017,felipe-rabelo-shard-00-02-r4yae.gcp.mongodb.net:27017/test?ssl=true&replicaSet=felipe-rabelo-shard-0&authSource=admin&retryWrites=true")
	if err != nil {
		log.Fatal(err)
	}

	filter := bson.M{"test": "testing around"}

	collection := client.Database("lxndr").Collection("lxndr-quest")
	res, err := collection.FindOne(ctx, filter).DecodeBytes()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Consulta de teste a base do Mongo:")
	fmt.Println(res)
}
