package authutil

import (
	"os"

	"github.com/golang-jwt/jwt"
)

func ParseClaims(token string) (jwt.MapClaims, error) {

	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {

		if token.Method.Alg() != "HS256" {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	return claims, nil
}

func GenerateToken(id string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": id,
	})

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return "", err
	}

	return t, nil
}
