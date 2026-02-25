package pwd

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/argon2"
)

func GenerateSalt() (string, error) {
	saltBytes := make([]byte, 16)
	_, err := rand.Read(saltBytes)
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(saltBytes), nil
}

func HashPassword(password string, salt []byte) string {
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return base64.RawStdEncoding.EncodeToString(hash)
}

func VerifyPassword(password string, salt []byte, hash string) bool {
	computedHash := HashPassword(password, salt)
	return computedHash == hash
}
