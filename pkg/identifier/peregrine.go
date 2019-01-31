package identifier

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// User ..
type User struct {
	Hash  string // email
	Token string // password
}

// IdentityCheck ..
func IdentityCheck(email, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(email))
	if err != nil {
		log.Println("A conta de email não bate com o hash")
		return false
	}

	log.Println("A conta de email bate com o hash")
	return true
}
