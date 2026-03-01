package restapi

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/auth"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/middleware"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/poster"

	"github.com/gorilla/mux"

	_ "github.com/go-park-mail-ru/2026_1_TheBugs/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

func RegisterHandlers(app *mux.Router, auth *auth.AuthHandler, poster *poster.PosterHandler) {
	app.Use(middleware.LoggingMiddleware)

	// Routers
	apiGroup := app.PathPrefix("/api").Subrouter()
	{
		apiGroup.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)

		apiGroup.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(map[string]bool{"ok": true})
		}).Methods(http.MethodGet)

		apiGroup.HandleFunc("/auth/reg", auth.RegisterUser)
		apiGroup.HandleFunc("/auth/login", auth.LoginUser)

		apiGroup.HandleFunc("/posters", poster.GetPosters)
	}
}
