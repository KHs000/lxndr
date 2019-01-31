package identifier

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"

	"github.com/KHs000/lxndr/pkg/mongo"
)

// User ..
type User struct {
	Email string // email
	Hash  string // hash email
	Token string // password
}

// ValidateNewUser ..
func ValidateNewUser(email string) bool {
	coll := mongo.Collection{Database: "lxndr", CollName: "user"}
	hits := mongo.Search(mongo.Conn, coll, bson.M{"email": email})

	for hits.Next(mongo.Conn.Ctx) {
		var row bson.M
		err := hits.Decode(&row)
		if err != nil {
			log.Fatal(err)
		}

		if row["_id"] != nil {
			return false
		}
	}
	return true
}

// IdentityCheck ..
func IdentityCheck(email, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(email))
	if err != nil {
		log.Println("The email doesn't match the hash")
		return false
	}

	log.Println("The email match the hash")
	return true
}
