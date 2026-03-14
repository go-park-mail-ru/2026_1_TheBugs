package restapi

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/auth"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/middleware"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/poster"
	"github.com/rs/cors"

	"github.com/gorilla/mux"

	_ "github.com/go-park-mail-ru/2026_1_TheBugs/internal/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           DomDeli API
// @version         1.0
// @description     Created by TheBugs in 2026
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT

// @host      localhost:8000
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func RegisterHandlers(app *mux.Router, auth *auth.AuthHandler, post *poster.PosterHandler) {

	c := cors.New(cors.Options{
		AllowedOrigins:   config.Config.CORS.AllowedHosts,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		ExposedHeaders:   []string{"Set-Cookie"},
	})

	app.Use(middleware.LoggingMiddleware)
	app.Use(c.Handler)

	// Routers
	apiGroup := app.PathPrefix("/api").Subrouter()
	apiGroup.Use(mux.CORSMethodMiddleware(apiGroup))
	{
		apiGroup.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
		}).Methods(http.MethodGet)

		apiGroup.HandleFunc("/auth/reg", auth.RegisterUser).Methods(http.MethodPost)
		apiGroup.HandleFunc("/auth/login", auth.LoginUser).Methods(http.MethodPost)
		apiGroup.HandleFunc("/auth/logout", auth.Logout).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.HandleFunc("/auth/refresh", auth.RefreshToken).Methods(http.MethodPost)
		apiGroup.HandleFunc("/auth/vkid", auth.VKLogin).Methods(http.MethodPost)

		apiGroup.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)

		apiGroup.HandleFunc("/posters", post.GetPosters).Methods(http.MethodGet)
	}
}
