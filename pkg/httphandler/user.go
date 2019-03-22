package httphandler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"

	"github.com/KHs000/lxndr/domain"
	"github.com/KHs000/lxndr/pkg/identifier"
	"github.com/KHs000/lxndr/pkg/mongodb"
	"github.com/KHs000/lxndr/pkg/rndtoken"
)

type request struct {
	email string
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	logAccess(r)
	defer recovery("Method not allowed.")
	validateMethod(w, r, "POST")

	b, e := processRequestBody(r, request{})
	if e != nil {
		writeResponse(w, e.Code, e.Error)
		return
	}

	code, resp := createUser(mongodb.Client, b)
	writeResponse(w, code, resp)
}

func createUser(client domain.Client, body map[string]interface{}) (int, interface{}) {
	resp := domain.Response{}

	email := body["email"].(string)
	isNew := identifier.ValidateNewUser(client, email)
	if !isNew {
		resp.Message = "This email has already been registred."
		log.Println("Email conflict.")
		return http.StatusConflict, resp
	}

	tkn, hash := rndtoken.SendToken(email)
	newUser := domain.User{Email: email, Hash: hash, Token: tkn}

	coll := domain.Collection{Database: "lxndr", CollName: "user"}
	id, err := mongodb.Insert(client, coll, newUser)

	if err != nil {
		message := "There might be an error, please retry your operation."
		resp.Message = message
		log.Println(message)
		return http.StatusInternalServerError, domain.Response{Message: message}
	}

	message := fmt.Sprintf(`%v`, primitive.ObjectID.String(id))
	resp.Message = message
	log.Println(message)
	return http.StatusCreated, domain.Response{Message: message}
}

func editUserHandler(w http.ResponseWriter, r *http.Request) {
	logAccess(r)
	defer recovery("Method not allowed.")
	validateMethod(w, r, "POST")

	// b, e := processRequestBody(r, request{})
	// if e != nil {
	// 	writeResponse(w, e.Code, e.Error)
	// 	return
	// }

	// code, resp := editUser(mongodb.Client, b)
	// writeResponse(w, code, resp)
}

func editUser(w http.ResponseWriter, r *http.Request) {
	logAccess(r)
	defer recovery("Method not allowed.")
	validateMethod(w, r, "POST")
	resp := domain.Response{}

	type request struct {
		email string
	}
	b, e := processRequestBody(r, request{})
	if e != nil {
		writeResponse(w, e.Code, e.Error)
		return
	}

	email := b["email"].(string)
	isNew := identifier.ValidateNewUser(mongodb.Client, email)
	if isNew {
		resp.Message = "This email is not registred."
		log.Println("Email not registred.")
		writeResponse(w, http.StatusNotFound, resp)
		return
	}

	coll := domain.Collection{Database: "lxndr", CollName: "user"}
	res := mongodb.Update(
		mongodb.Conn, coll, bson.M{"email": email}, bson.M{"$set": b})

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
	defer recovery("Method not allowed.")
	validateMethod(w, r, "POST")
	resp := domain.Response{}

	type request struct {
		email string
	}
	b, e := processRequestBody(r, request{})
	if e != nil {
		writeResponse(w, e.Code, e.Error)
		return
	}

	email := b["email"].(string)
	isNew := identifier.ValidateNewUser(mongodb.Client, email)
	if isNew {
		resp.Message = "This email is not registred."
		log.Println("Email not registred.")
		writeResponse(w, http.StatusNotFound, resp)
		return
	}

	coll := domain.Collection{Database: "lxndr", CollName: "user"}
	res := mongodb.Delete(mongodb.Conn, coll, bson.M{"email": email})

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
	defer recovery("Method not allowed.")
	validateMethod(w, r, "GET")
	resp := domain.Response{}

	coll := domain.Collection{Database: "lxndr", CollName: "user"}
	res := mongodb.Search(mongodb.Client, coll, nil)

	var list []string
	for res.Next(mongodb.Conn.Ctx) {
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

func testHandler(w http.ResponseWriter, r *http.Request) {
	logAccess(r)
	defer recovery("Method not allowed.")
	validateMethod(w, r, "GET")

	client := domain.MongoClient{Client: mongodb.Conn.Client, Context: mongodb.Conn.Ctx}
	resp := test(client)

	log.Printf("Found %v users", len(resp.Data))
	writeResponse(w, http.StatusOK, resp)
}

func test(client domain.Client) domain.Response {
	resp := domain.Response{}

	coll := domain.Collection{Database: "lxndr", CollName: "user"}
	res := mongodb.Test(client, coll, bson.M{})

	var list []string
	for res.Next(client.Ctx()) {
		row, err := res.DecodeCursor()
		if err != nil {
			resp.Message = "Internal server error."
			log.Println("Error decoding line from search.")
			return resp
		}
		list = append(list, row["email"].(string))
	}

	resp.Message = ""
	resp.Data = list
	return resp
}
