package restapi

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/auth"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/complex"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/middleware"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/poster"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/user"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"

	_ "github.com/go-park-mail-ru/2026_1_TheBugs/docs"
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
func RegisterHandlers(app *mux.Router, logger *logrus.Logger, auth *auth.AuthHandler, post *poster.PosterHandler, UtilityCompany *complex.UtilityCompanyHandler, user *user.UserHandler) {

	c := cors.New(cors.Options{
		AllowedOrigins:   config.Config.CORS.AllowedHosts,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		ExposedHeaders:   []string{"Set-Cookie"},
	})

	app.Use(middleware.LoggingMiddleware(logger))
	app.Use(c.Handler)

	app.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}).Methods(http.MethodGet)

	// API Routers
	apiGroup := app.PathPrefix("/api").Subrouter()
	apiGroup.Use(mux.CORSMethodMiddleware(apiGroup))
	{
		AuthMiddlewary := auth.GetAuthMiddlewary()
		apiGroup.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
		}).Methods(http.MethodGet)
		apiGroup.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)

		apiGroup.Handle("/user/me", AuthMiddlewary(http.HandlerFunc(user.GetMe))).Methods(http.MethodGet, http.MethodOptions)
		apiGroup.Handle("/user/me/profile", AuthMiddlewary(http.HandlerFunc(user.UpdateProfile))).Methods(http.MethodPut, http.MethodOptions)

		apiGroup.HandleFunc("/auth/login", auth.LoginUser).Methods(http.MethodPost)
		apiGroup.HandleFunc("/auth/reg", auth.RegisterUser).Methods(http.MethodPost)
		apiGroup.HandleFunc("/auth/logout", auth.Logout).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.HandleFunc("/auth/refresh", auth.RefreshToken).Methods(http.MethodPost)
		apiGroup.HandleFunc("/auth/vkid", auth.VKLogin).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.HandleFunc("/auth/yandex", auth.YandexLogin).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.HandleFunc("/auth/recover", auth.SendCodeOnEmail).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.HandleFunc("/auth/recover/verify", auth.VerifyRecoveryCode).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.HandleFunc("/auth/recover/reset", auth.UpdatePassword).Methods(http.MethodPost, http.MethodOptions)

		apiGroup.HandleFunc("/utility-companies/by-alias/{alias}", UtilityCompany.GetUtilityCompany).Methods(http.MethodGet)
		apiGroup.HandleFunc("/utility-companies/developers", UtilityCompany.GetAllDevelopers).Methods(http.MethodGet)
		apiGroup.HandleFunc("/utility-companies/", UtilityCompany.GetUtilityCompaniesByDeveloper).Methods(http.MethodGet)

		apiGroup.HandleFunc("/posters/flats", post.GetFlatsAll).Methods(http.MethodGet)
		apiGroup.HandleFunc("/posters/by-alias/{alias}", post.GetPoster).Methods(http.MethodGet)

		apiGroup.Handle("/posters/me", AuthMiddlewary(http.HandlerFunc(post.GetPostersByUser))).Methods(http.MethodGet, http.MethodOptions)
		apiGroup.Handle("/posters/me/{alias}", AuthMiddlewary(http.HandlerFunc(post.GetPostersByUserByAlias))).Methods(http.MethodGet, http.MethodOptions)
		apiGroup.Handle("/posters/flat", AuthMiddlewary(http.HandlerFunc(post.CreateFlatPoster))).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.Handle("/posters/flat/{alias}", AuthMiddlewary(http.HandlerFunc(post.UpdateFlatPoster))).Methods(http.MethodPut, http.MethodOptions)
		apiGroup.Handle("/posters/flat/{alias}", AuthMiddlewary(http.HandlerFunc(post.DeleteFlatPoster))).Methods(http.MethodDelete, http.MethodOptions)

	} //alias
}
