package utils

import (
	"crypto/sha512"
	"encoding/hex"
)

func HashPassword(password string) string {
	hash := sha512.Sum512([]byte(password))
	return hex.EncodeToString(hash[:])
}

func CheckPasswordHash(password, hash string) bool {
	return HashPassword(password) == hash
}
