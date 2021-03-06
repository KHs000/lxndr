package identifier

import (
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/mongodb/mongo-go-driver/bson"

	"github.com/KHs000/lxndr/domain"
	"github.com/KHs000/lxndr/pkg/mongodb"
)

// ValidateNewUser ..
func ValidateNewUser(client domain.Client, email string) bool {
	coll := domain.Collection{Database: "lxndr", CollName: "user"}
	hits := mongodb.Search(client, coll, bson.M{"email": email})

	for hits.Next(client.Ctx()) {
		var row bson.M
		row, err := hits.DecodeCursor()
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
