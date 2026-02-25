package tokens

import (
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Sub  string `json:"sub"`
	Type string `json:"type"`

	//Email  string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWT(sub string, tokenType string, exp time.Duration) (string, error) {
	secondsFloat := exp.Seconds()
	claims := Claims{
		Sub:  sub,
		Type: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(secondsFloat * float64(time.Second)))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(config.JWTKeys.PrivateKey)
}

func ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return config.JWTKeys.PublicKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}
