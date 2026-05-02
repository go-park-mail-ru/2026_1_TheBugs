package pwd

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"
)

const SessionLength = 10
const CodeLength = 5

func GenerateCode() string {
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%05d", seed.Intn(100000))
}

func GenerateSessionID() string {
	saltBytes := make([]byte, SessionLength)
	rand.Read(saltBytes)
	return base64.RawStdEncoding.EncodeToString(saltBytes)
}
