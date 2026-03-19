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
func generateRandomSalt(saltSize int) ([]byte, error) {
	var salt = make([]byte, saltSize)

	_, err := rand.Read(salt[:])

	if err != nil {
		return nil, err
	}

	return salt, nil
}

func HashToken(token string) (string, error) {
	var passwordBytes = []byte(token)

	var sha512Hasher = sha512.New()

	salt, err := generateRandomSalt(16)
	if err != nil {
		return "", err
	}

	passwordBytes = append(passwordBytes, salt...)

	sha512Hasher.Write(passwordBytes)

	var hashedPasswordBytes = sha512Hasher.Sum(nil)

	var hashedPasswordHex = hex.EncodeToString(hashedPasswordBytes)

	return hashedPasswordHex, nil
}
