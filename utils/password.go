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

	//hash the passsword

	hashedPassword := argon2.IDKey([]byte(password), []byte(salt), 1, 64*1024, 4, 32)

	encodePassword := fmt.Sprintf("%s.%s", base64.RawStdEncoding.EncodeToString(salt))

	base64.RawStdEncoding.EncodeToString(hashedPassword)

	return encodePassword, nil

}

func VerifyPassword(encodedPassword string, password string) bool {
	encodedSaltPassword := password

	parts := strings.Split(encodedSaltPassword, ".")

	decodedHashPassword, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}

	decodedSalt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false
	}

	hashedPassword := argon2.IDKey([]byte(encodedPassword), decodedSalt, 1, 64*1024, 4, 32)

	if bytes.Equal(hashedPassword, decodedHashPassword) {
		return true
	}

	return false
}
