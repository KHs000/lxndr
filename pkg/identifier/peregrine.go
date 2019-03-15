package identifier

import (
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/mongodb/mongo-go-driver/bson"

	"github.com/KHs000/lxndr/pkg/mongodb"
)

// User ..
type User struct {
	Email string // email
	Hash  string // hash email
	Token string // password
}

// ValidateNewUser ..
func ValidateNewUser(email string) bool {
	coll := mongodb.Collection{Database: "lxndr", CollName: "user"}
	hits := mongodb.Search(mongodb.Conn, coll, bson.M{"email": email})

	for hits.Next(mongodb.Conn.Ctx) {
		var row bson.M
		err := hits.Decode(&row)
		if err != nil {
			log.Println("Error decoding line from search.")
			return false
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
