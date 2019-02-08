package main

import (
	"fmt"
	"log"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"

	"github.com/KHs000/lxndr/pkg/httphandler"
	"github.com/KHs000/lxndr/pkg/identifier"
	"github.com/KHs000/lxndr/pkg/mongo"
	"github.com/KHs000/lxndr/pkg/rndtoken"
)

func exportRoutes() {
	httphandler.Handler("/", []string{}, defaultRoute)
	httphandler.Handler("/newUser", []string{"email"}, createUser)
	httphandler.Handler("/editUser", []string{"email"}, editUser)
	httphandler.Handler("/deleteUser", []string{"email"}, deleteUser)
}

func defaultRoute(b map[string]string) httphandler.Res {
	resp := httphandler.Res{}
	resp.S.Code = 200
	resp.S.Message = "It works..."
	return resp
}

func createUser(b map[string]string) httphandler.Res {
	email := b["email"]
	resp := httphandler.Res{}

	isNew := identifier.ValidateNewUser(email)
	if !isNew {
		resp.E.Code = 409
		resp.E.Error = "This email has already been registred."
		return resp
	}

	tkn, hash := rndtoken.SendToken(email)
	newUser := identifier.User{Email: email, Hash: hash, Token: tkn}

	coll := mongo.Collection{Database: "lxndr", CollName: "user"}
	res := mongo.Insert(mongo.Conn, coll, newUser)

	if id, ok := res.InsertedID.(primitive.ObjectID); ok {
		resp.S.Code = 201
		resp.S.Message = fmt.Sprintf(
			"User Created. ID: %v.", primitive.ObjectID.String(id))
		return resp
	}

	resp.E.Code = 500
	resp.E.Error = "There might be an error, please retry your operation."
	return resp
}

func editUser(b map[string]string) httphandler.Res {
	email := b["email"]
	resp := httphandler.Res{}

	isNew := identifier.ValidateNewUser(email)
	if isNew {
		log.Println("This email is not registred.")
		resp.E.Code = 404
		resp.E.Error = "This email is not registred."
		return resp
	}

	coll := mongo.Collection{Database: "lxndr", CollName: "user"}
	res := mongo.Update(
		mongo.Conn, coll, bson.M{"email": email}, bson.M{"$set": b})

	if res.MatchedCount != 1 {
		log.Println("This email matched no registry.")
		resp.E.Code = 404
		resp.E.Error = "This email matched no registry."
		return resp
	}

	resp.S.Code = 200
	resp.S.Message = fmt.Sprintf("%q user updated.", email)
	return resp
}

func deleteUser(b map[string]string) httphandler.Res {
	email := b["email"]
	resp := httphandler.Res{}

	isNew := identifier.ValidateNewUser(email)
	if isNew {
		log.Println("This email is not registred.")
		resp.E.Code = 404
		resp.E.Error = "This email is not registred."
		return resp
	}

	coll := mongo.Collection{Database: "lxndr", CollName: "user"}
	res := mongo.Delete(mongo.Conn, coll, bson.M{"email": email})

	if res.DeletedCount != 1 {
		log.Println("This email matched no registry.")
		resp.E.Code = 404
		resp.E.Error = "This email matched no registry."
		return resp
	}

	resp.S.Code = 200
	resp.S.Message = fmt.Sprintf("%q user deleted.", email)
	return resp
}
