package rndtoken

import (
	"crypto/md5"
	"encoding/hex"
	"log"

    "golang.org/x/crypto/bcrypt"
)

func GenerateToken (email string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(email), bcrypt.DefaultCost)
    if err != nil {
        log.Fatal(err)
    }

    hasher := md5.New()
    hasher.Write(hash)
    return hex.EncodeToString(hasher.Sum(nil))
}

func RecoverToken (hashKey string) string {
	hasher := md5.New()
    hasher.Write([]byte(hashKey))
    return hex.EncodeToString(hasher.Sum(nil))
} // WIṔ
