package main

import (
	"log"

	"github.com/mongodb/mongo-go-driver/bson/primitive"

	"github.com/KHs000/lxndr/pkg/identifier"
	"github.com/KHs000/lxndr/pkg/mongo"

	"github.com/KHs000/lxndr/pkg/rndtoken"
)

func test() {
	log.Println("This is a test.")
}

func createUser(email string) {
	tkn, hash := rndtoken.SendToken(email)

	newUser := identifier.User{Hash: hash, Token: tkn}

	coll := mongo.Collection{Database: "lxndr", CollName: "user"}
	res := mongo.Insert(mongo.Conn, coll, newUser)

	if id, ok := res.InsertedID.(primitive.ObjectID); ok {
		log.Printf("User Created. ID: %v", primitive.ObjectID.String(id))
	}
}
