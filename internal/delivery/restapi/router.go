package restapi

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/auth"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/middleware"

	"github.com/gorilla/mux"
)

func RegisterHandlers(app *mux.Router, auth *auth.AuthHandler) {
	app.Use(middleware.LoggingMiddleware)

	// Routers
	apiGroup := app.PathPrefix("/api").Subrouter()
	{
		apiGroup.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(map[string]bool{"ok": true})
		}).Methods(http.MethodGet)

		apiGroup.HandleFunc("/auth/reg", auth.RegisterUser)
		apiGroup.HandleFunc("/auth/login", auth.LoginUser)
		apiGroup.HandleFunc("/auth/refresh", auth.RefreshToken)
	}
}
