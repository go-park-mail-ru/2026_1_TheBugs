package utils

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
)

func SetRefreshCookie(w http.ResponseWriter, accessCred *dto.UserAccessCredDTO) {
	http.SetCookie(
		w,
		&http.Cookie{
			Name:     "refresh_token",
			Value:    accessCred.RefreshToken,
			Path:     "/api/auth",
			HttpOnly: true,
			Domain:   config.Config.CORS.CookieHost,
			SameSite: http.SameSiteLaxMode,
			Expires:  time.Now().Add(time.Duration(accessCred.RefreshTokenExp * int(time.Second))),
			MaxAge:   accessCred.RefreshTokenExp,
		},
	)
}
