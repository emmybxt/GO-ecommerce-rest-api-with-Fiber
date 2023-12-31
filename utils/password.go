package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) (string, error) {

	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	// hash password with argon2
	hashedPassword := argon2.IDKey([]byte(password), []byte(salt), 1, 64*1024, 4, 32)

	//  encode salt and hashed password to base64
	encodedPassword := fmt.Sprintf("%s.%s", base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hashedPassword))
	return encodedPassword, nil
}

func VerifyPassword(encodedPassword string, password string) bool {
	encodedSaltAndPassword := password
	// this method is less efficient
	parts := strings.Split(encodedSaltAndPassword, ".")
	// this method is more efficient
	//parts := helpers.SplitString(encodedSaltAndPassword, ".")
	decodedHashedPassword, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}
	decodedSalt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}

	// hash the password with the same salt
	hashedPassword := argon2.IDKey([]byte(encodedPassword), decodedSalt, 1, 64*1024, 4, 32)

	// compare the hashedPassword with hash
	if bytes.Equal(hashedPassword, decodedHashedPassword) {
		return true
	}
	return false
}
