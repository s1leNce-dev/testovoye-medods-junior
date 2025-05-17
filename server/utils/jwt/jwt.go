package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT_UserClaims struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	jwt.RegisteredClaims
}

func AccessTokenGen(userID, sessionID string) (string, error) {
	accessClaims := JWT_UserClaims{
		UserID:    userID,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, accessClaims)

	return accessToken.SignedString([]byte(os.Getenv("JWT_ACCESS_TOKEN")))
}

func ValidateAccessToken(tokenString string) (*JWT_UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWT_UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_ACCESS_TOKEN")), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWT_UserClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func ValidateAccessTokenForRefreshToken(tokenString string) (*JWT_UserClaims, error) {
	claims := &JWT_UserClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.SigningMethodHS512.Alg() {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			secret := os.Getenv("JWT_ACCESS_TOKEN")
			return []byte(secret), nil
		},
		jwt.WithoutClaimsValidation(),
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
