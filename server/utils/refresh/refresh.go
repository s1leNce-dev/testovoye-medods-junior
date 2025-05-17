package refresh

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

var MaxAgeRT = 30 * (time.Hour * 24)

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
