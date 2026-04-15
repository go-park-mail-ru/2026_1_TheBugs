package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
)

func GetAccessToken(r *http.Request) (string, error) {
	const BearerPrefix = "Bearer "
	cred := r.Header.Get("Authorization")

	if cred == "" {
		return "", errors.New("authorization header missing")
	}

	if !strings.HasPrefix(cred, BearerPrefix) {
		return "", errors.New("invalid authorization header: missing Bearer prefix")
	}
	token := strings.TrimPrefix(cred, BearerPrefix)
	if token == "" {
		return "", errors.New("empty token")
	}

	return token, nil
}

func GetUserID(ctx context.Context) (int, error) {
	val := ctx.Value(entity.UserID{})
	userID, ok := val.(int)
	fmt.Println(val)
	if !ok {
		return 0, fmt.Errorf("wrong userID type %T", userID)
	}
	return userID, nil
}

func SetUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, entity.UserID{}, userID)
}
