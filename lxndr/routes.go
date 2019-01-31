package main

import (
	"log"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"

	"github.com/KHs000/lxndr/pkg/identifier"
	"github.com/KHs000/lxndr/pkg/mongo"
	"github.com/KHs000/lxndr/pkg/rndtoken"
)

func createUser(email string) {
	isNew := identifier.ValidateNewUser(email)
	if !isNew {
		log.Fatal("This email has already been registred")
	}

	tkn, hash := rndtoken.SendToken(email)
	newUser := identifier.User{Email: email, Hash: hash, Token: tkn}

	coll := mongo.Collection{Database: "lxndr", CollName: "user"}
	res := mongo.Insert(mongo.Conn, coll, newUser)

	if id, ok := res.InsertedID.(primitive.ObjectID); ok {
		log.Printf("User Created. ID: %v", primitive.ObjectID.String(id))
	}
}

func editUser(email string, data interface{}) {
	isNew := identifier.ValidateNewUser(email)
	if isNew {
		log.Fatal("This email is not registred")
	}

	coll := mongo.Collection{Database: "lxndr", CollName: "user"}
	res := mongo.Update(mongo.Conn, coll, bson.M{"email": email}, data)

	if id, ok := res.UpsertedID.(primitive.ObjectID); ok {
		log.Printf("User Updated. ID %v", id)
	}
}
