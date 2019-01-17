package rndtoken

import (
	"crypto/md5"
	"encoding/hex"
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/KHs000/lxndr/identifier"
)

func SendToken(email string) (string, string) {
	tkn, hash := generateToken(email)

	return tkn, string(hash)
}

func ValidateToken(email, hash string) bool {
	return identifier.IdentityCheck(email, hash)
}

func generateToken(email string) (string, []byte) {
	hash, err := bcrypt.GenerateFromPassword([]byte(email), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	hasher := md5.New()
	hasher.Write(hash)
	return hex.EncodeToString(hasher.Sum(nil)), hash
}
