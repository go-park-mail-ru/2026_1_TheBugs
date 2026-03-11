package tokens

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	Sub  string `json:"sub"`
	Type string `json:"type"`

	//Email  string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWT(claims Claims, exp time.Duration) (string, error) {
	secondsFloat := exp.Seconds()
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Duration(secondsFloat * float64(time.Second))))
	claims.IssuedAt = jwt.NewNumericDate(time.Now())
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

func GenerateRefreshToken(tokenID string, userID int, exp time.Duration) (string, error) {
	claims := Claims{
		Sub:  strconv.Itoa(userID),
		Type: entity.RefreshTokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ID: tokenID,
		},
	}
	return GenerateJWT(claims, exp)
}

func GenerateAccessToken(userID int, exp time.Duration) (string, error) {
	ID := uuid.NewString()
	claims := Claims{
		Sub:  strconv.Itoa(userID),
		Type: entity.AccessTokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ID: ID,
		},
	}
	return GenerateJWT(claims, exp)
}
