package util

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func GetHashedValue(bytes []byte) string {
    hash, err := bcrypt.GenerateFromPassword(bytes, bcrypt.MinCost)
    if err != nil {
        log.Println(err)
    }
    return string(hash)
}

func HashedValueMatches(hash string, bytes []byte) bool {
    byteHash := []byte(hash)
    err := bcrypt.CompareHashAndPassword(byteHash, bytes)
    if err != nil {
        log.Println(err)
        return false
    }
    
    return true
}