package cache

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

func GenerateCacheKey(obj any) (string, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)

	return fmt.Sprintf("%x", hash), nil
}
