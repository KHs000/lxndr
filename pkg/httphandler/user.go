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

type response struct {
	Message string   `json:"message"`
	Data    []string `json:"data"`
}

func createUser(w http.ResponseWriter, r *http.Request) {
	logAccess(r)
	resp := response{}

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
		resp.Message = "This email has already been registred."
		log.Println("Email conflict.")
		writeResponse(w, http.StatusConflict, resp)
		return
	}

	tkn, hash := rndtoken.SendToken(email)
	newUser := identifier.User{Email: email, Hash: hash, Token: tkn}

	coll := mongo.Collection{Database: "lxndr", CollName: "user"}
	res := mongo.Insert(mongo.Conn, coll, newUser)

	if id, ok := res.InsertedID.(primitive.ObjectID); ok {
		message := fmt.Sprintf(`%v`, primitive.ObjectID.String(id))
		resp.Message = message
		log.Println(message)
		writeResponse(w, http.StatusCreated, resp)
		return
	}

	message := "There might be an error, please retry your operation."
	resp.Message = message
	log.Println(message)
	writeResponse(w, http.StatusInternalServerError, resp)
	return
}

func editUser(w http.ResponseWriter, r *http.Request) {
	logAccess(r)
	resp := response{}

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
		resp.Message = "This email is not registred."
		log.Println("This email is not registred.")
		writeResponse(w, http.StatusNotFound, resp)
		return
	}

	coll := mongo.Collection{Database: "lxndr", CollName: "user"}
	res := mongo.Update(
		mongo.Conn, coll, bson.M{"email": email}, bson.M{"$set": b})

	if res.MatchedCount != 1 {
		resp.Message = "This email matched no registry."
		log.Println("This email matched no registry.")
		writeResponse(w, http.StatusNotFound, resp)
		return
	}

	message := fmt.Sprintf(`User '%v' updated.`, email)
	resp.Message = message
	log.Println(message)
	writeResponse(w, http.StatusOK, resp)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	logAccess(r)
	resp := response{}

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
		resp.Message = "This email is not registred."
		log.Println("This email is not registred.")
		writeResponse(w, http.StatusNotFound, resp)
		return
	}

	coll := mongo.Collection{Database: "lxndr", CollName: "user"}
	res := mongo.Delete(mongo.Conn, coll, bson.M{"email": email})

	if res.DeletedCount != 1 {
		resp.Message = "This email matched no registry."
		log.Println("This email matched no registry.")
		writeResponse(w, http.StatusNotFound, resp)
		return
	}

	message := fmt.Sprintf(`User '%v' deleted.`, email)
	resp.Message = message
	log.Println(message)
	writeResponse(w, http.StatusOK, resp)
}

func listUsers(w http.ResponseWriter, r *http.Request) {
	logAccess(r)
	resp := response{}

	coll := mongo.Collection{Database: "lxndr", CollName: "user"}
	res := mongo.Search(mongo.Conn, coll, nil)

	var list []string
	for res.Next(mongo.Conn.Ctx) {
		var row bson.M
		err := res.Decode(&row)
		if err != nil {
			resp.Message = "Internal server error."
			log.Println("Error decoding line from search.")
			writeResponse(w, http.StatusInternalServerError, resp)
		}

		list = append(list, row["email"].(string))
	}

	resp.Message = ""
	resp.Data = list
	log.Printf("Found %v users", len(list))
	writeResponse(w, http.StatusOK, resp)
}
