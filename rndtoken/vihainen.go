package rndtoken

import (
	"crypto/md5"
	"encoding/hex"
	"log"

    "golang.org/x/crypto/bcrypt"

    "github.com/KHs000/lxndr/identifier"
)

func SendToken (email string) string {
    tkn, _ := generateToken(email)

    return tkn
}

func ValidateToken (email, key string) bool {
    return identifier.IdentityCheck(email, key)
}

func generateToken (email string) (string, []byte) {
	hash, err := bcrypt.GenerateFromPassword([]byte(email), bcrypt.DefaultCost)
    if err != nil {
        log.Fatal(err)
    }

    hasher := md5.New()
    hasher.Write(hash)
    return hex.EncodeToString(hasher.Sum(nil)), hash
}
