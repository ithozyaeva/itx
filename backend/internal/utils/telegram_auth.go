package utils

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"time"
)

func CheckExpirationDate(date time.Time) bool {
	return time.Now().After(date)
}
func generateRandomSalt(saltSize int) []byte {
	var salt = make([]byte, saltSize)

	_, err := rand.Read(salt[:])

	if err != nil {
		panic(err)
	}

	return salt
}

func HashToken(token string) string {
	var passwordBytes = []byte(token)

	var sha512Hasher = sha512.New()

	passwordBytes = append(passwordBytes, generateRandomSalt(16)...)

	sha512Hasher.Write(passwordBytes)

	var hashedPasswordBytes = sha512Hasher.Sum(nil)

	var hashedPasswordHex = hex.EncodeToString(hashedPasswordBytes)

	return hashedPasswordHex
}
