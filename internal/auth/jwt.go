package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func MakeJWT(userID int32, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "synlabs",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   string(userID),
	})

	signedString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", nil
	}

	return signedString, nil
}
