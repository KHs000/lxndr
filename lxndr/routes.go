package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"

	"github.com/KHs000/lxndr/pkg/identifier"
	"github.com/KHs000/lxndr/pkg/mongo"
	"github.com/KHs000/lxndr/pkg/rndtoken"
)

func exportRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defaultRoute(w, r)
	})

	http.HandleFunc("/newUser", func(w http.ResponseWriter, r *http.Request) {
		rb, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		type validBody struct {
			Email string `json:"email"`
		}
		b := validBody{}
		err = json.Unmarshal([]byte(rb), &b)
		if err != nil {
			log.Println("Invalid body.")
			w.WriteHeader(400)
			w.Write([]byte("400 - Bad Request"))
			return
		}

		res, err := createUser(b.Email)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(res)
	})
}

func defaultRoute(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, "It works...") }

func createUser(email string) (string, error) {
	isNew := identifier.ValidateNewUser(email)
	if !isNew {
		return "This email has already been registred", nil
	}

	tkn, hash := rndtoken.SendToken(email)
	newUser := identifier.User{Email: email, Hash: hash, Token: tkn}

	coll := mongo.Collection{Database: "lxndr", CollName: "user"}
	res := mongo.Insert(mongo.Conn, coll, newUser)

	if id, ok := res.InsertedID.(primitive.ObjectID); ok {
		return fmt.Sprintf("User Created. ID: %v.", primitive.ObjectID.String(id)), nil
	}

	return "There might be an error, please retry your operation.", nil
}

func editUser(email string, data interface{}) {
	isNew := identifier.ValidateNewUser(email)
	if isNew {
		log.Fatal("This email is not registred")
	}

	coll := mongo.Collection{Database: "lxndr", CollName: "user"}
	res := mongo.Update(mongo.Conn, coll, bson.M{"email": email}, bson.M{"$set": data})

	if res.MatchedCount != 1 {
		log.Fatal("There was as error during the update.")
	}

	log.Printf("%v user updated.", res.ModifiedCount)
}
