package httphandler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"

	"github.com/KHs000/lxndr/pkg/identifier"
	"github.com/KHs000/lxndr/pkg/mongo"
	"github.com/KHs000/lxndr/pkg/rndtoken"
)

func createUser(w http.ResponseWriter, r *http.Request) {
	logAccess(r)

	type request struct {
		email string
	}
	b, e := processRequestBody(r, request{})
	if e != nil {
		writeResponse(w, e.Code, e.Error)
		return
	}

	email := b["email"].(string)
	isNew := identifier.ValidateNewUser(email)
	if !isNew {
		writeResponse(w, http.StatusConflict,
			"This email has already been registred.")
		return
	}

	tkn, hash := rndtoken.SendToken(email)
	newUser := identifier.User{Email: email, Hash: hash, Token: tkn}

	coll := mongo.Collection{Database: "lxndr", CollName: "user"}
	res := mongo.Insert(mongo.Conn, coll, newUser)

	if id, ok := res.InsertedID.(primitive.ObjectID); ok {
		message := fmt.Sprintf(`User Created. ID: %v.`, primitive.ObjectID.String(id))
		log.Println(message)
		writeResponse(w, http.StatusCreated, message)
		return
	}

	message := "There might be an error, please retry your operation."
	log.Println(message)
	writeResponse(w, http.StatusInternalServerError, message)
	return
}

func editUser(w http.ResponseWriter, r *http.Request) {
	logAccess(r)

	type request struct {
		email string
	}
	b, e := processRequestBody(r, request{})
	if e != nil {
		writeResponse(w, e.Code, e.Error)
		return
	}

	email := b["email"].(string)
	isNew := identifier.ValidateNewUser(email)
	if isNew {
		log.Println("This email is not registred.")
		writeResponse(w, http.StatusNotFound, "This email is not registred.")
		return
	}

	coll := mongo.Collection{Database: "lxndr", CollName: "user"}
	res := mongo.Update(
		mongo.Conn, coll, bson.M{"email": email}, bson.M{"$set": b})

	if res.MatchedCount != 1 {
		log.Println("This email matched no registry.")
		writeResponse(w, http.StatusNotFound, "This email matched no registry.")
		return
	}

	message := fmt.Sprintf(`User '%v' updated.`, email)
	log.Println(message)
	writeResponse(w, http.StatusOK, message)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	logAccess(r)

	type request struct {
		email string
	}
	b, e := processRequestBody(r, request{})
	if e != nil {
		writeResponse(w, e.Code, e.Error)
		return
	}

	email := b["email"].(string)
	isNew := identifier.ValidateNewUser(email)
	if isNew {
		log.Println("This email is not registred.")
		writeResponse(w, http.StatusNotFound, "This email is not registred.")
		return
	}

	coll := mongo.Collection{Database: "lxndr", CollName: "user"}
	res := mongo.Delete(mongo.Conn, coll, bson.M{"email": email})

	if res.DeletedCount != 1 {
		log.Println("This email matched no registry.")
		writeResponse(w, http.StatusNotFound, "This email matched no registry.")
		return
	}

	message := fmt.Sprintf(`User '%v' deleted.`, email)
	log.Println(message)
	writeResponse(w, http.StatusOK, message)
}
