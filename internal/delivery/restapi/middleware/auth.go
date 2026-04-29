package middleware

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/generated/auth"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/utils"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/ctxLogger"
)

func AuthMiddleware(uc auth.AuthServiceClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			op := "AuthMiddleware"
			log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

			accessToken, err := utils.GetAccessToken(r)
			if err != nil {
				log.Errorf("utils.GetAccessToken: %s", err)
				utils.WriteError(w, "invalid access token", http.StatusUnauthorized)
				return
			}
			resp, err := uc.CheckAccessToken(r.Context(), &auth.CheckAccessTokenRequest{AccessToken: accessToken})
			if err != nil {
				log.Errorf("uc.CheckAccessToken: %s", err)
				utils.WriteError(w, "invalid access token", http.StatusUnauthorized)
				return
			}
			log.Info(resp)
			userID, err := strconv.Atoi(resp.UserId)
			if err != nil {
				log.Errorf("strconv.Atoi: %s", err)
				utils.WriteError(w, "invalid user ID", http.StatusUnauthorized)
				return
			}
			log.Info(userID)
			ctx := utils.SetUserID(r.Context(), userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
