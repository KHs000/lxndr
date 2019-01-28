package identifier

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func IdentityCheck(email, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(email))

	if err != nil {
		fmt.Println("A conta de email não bate com o hash")
	} else {
		fmt.Println("A conta de email bate com o hash")
	}

	return true
}
