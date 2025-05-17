package encryption

import (
	"golang.org/x/crypto/bcrypt"
)

func Hash(value string) (string, error) {
	hashedValue, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedValue), nil
}

func VerifyHashedValue(hashedValue, value string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedValue), []byte(value))
}
