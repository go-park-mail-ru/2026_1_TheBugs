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

// @securityDefinitions.apikey   CSRFToken
// @in                           header
// @name                         X-CSRF-Token
// @description                  CSRF токен в заголовке X-CSRF-Token для защищённых

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func RegisterHandlers(app *mux.Router, logger *logrus.Logger, auth *auth.AuthHandler, post *poster.PosterHandler, UtilityCompany *complex.UtilityCompanyHandler, user *user.UserHandler) {

	c := cors.New(cors.Options{
		AllowedOrigins:   config.Config.CORS.AllowedHosts,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		ExposedHeaders:   []string{"Set-Cookie", "X-CSRF-TOKEN"},
	})

	app.Use(middleware.LoggingMiddleware(logger))
	app.Use(c.Handler)

	app.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}).Methods(http.MethodGet)

	// API Routers
	apiGroup := app.PathPrefix("/api").Subrouter()
	apiGroup.Use(middleware.CSRFMiddleware)
	apiGroup.Use(middleware.SecurityMiddleware)
	apiGroup.Use(mux.CORSMethodMiddleware(apiGroup))

	{
		apiGroup.HandleFunc("/csrf-token", auth.GetCSRFToken).Methods(http.MethodGet)

		AuthMiddlewary := auth.GetAuthMiddlewary()
		apiGroup.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
		}).Methods(http.MethodGet)
		apiGroup.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)

		apiGroup.Handle("/user/me", AuthMiddlewary(http.HandlerFunc(user.GetMe))).Methods(http.MethodGet, http.MethodOptions)
		apiGroup.Handle("/user/me/profile", AuthMiddlewary(http.HandlerFunc(user.UpdateProfile))).Methods(http.MethodPut, http.MethodOptions)

		apiGroup.HandleFunc("/auth/login", auth.LoginUser).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.HandleFunc("/auth/reg", auth.RegisterUser).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.HandleFunc("/auth/logout", auth.Logout).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.HandleFunc("/auth/refresh", auth.RefreshToken).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.HandleFunc("/auth/vkid", auth.VKLogin).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.HandleFunc("/auth/yandex", auth.YandexLogin).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.HandleFunc("/auth/recover", auth.SendCodeOnEmail).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.HandleFunc("/auth/recover/verify", auth.VerifyRecoveryCode).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.HandleFunc("/auth/recover/reset", auth.UpdatePassword).Methods(http.MethodPost, http.MethodOptions)

		apiGroup.HandleFunc("/utility-companies/by-alias/{alias}", UtilityCompany.GetUtilityCompany).Methods(http.MethodGet)
		apiGroup.HandleFunc("/utility-companies/developers", UtilityCompany.GetAllDevelopers).Methods(http.MethodGet)
		apiGroup.HandleFunc("/utility-companies/", UtilityCompany.GetUtilityCompaniesByDeveloper).Methods(http.MethodGet)

		apiGroup.HandleFunc("/posters/flats", post.GetFlatsAll).Methods(http.MethodGet)
		apiGroup.HandleFunc("/posters/geo", post.GetFlatsMapAll).Methods(http.MethodGet, http.MethodOptions)
		apiGroup.HandleFunc("/posters/by-point", post.GetPostersByPoint).Methods(http.MethodGet, http.MethodOptions)
		apiGroup.HandleFunc("/posters/by-alias/{alias}", post.GetPoster).Methods(http.MethodGet)

		apiGroup.Handle("/posters/me", AuthMiddlewary(http.HandlerFunc(post.GetPostersByUser))).Methods(http.MethodGet, http.MethodOptions)
		apiGroup.Handle("/posters/me/{alias}", AuthMiddlewary(http.HandlerFunc(post.GetPostersByUserByAlias))).Methods(http.MethodGet, http.MethodOptions)
		apiGroup.Handle("/posters/flat", AuthMiddlewary(http.HandlerFunc(post.CreateFlatPoster))).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.Handle("/posters/flat/{alias}", AuthMiddlewary(http.HandlerFunc(post.UpdateFlatPoster))).Methods(http.MethodPut, http.MethodOptions)
		apiGroup.Handle("/posters/flat/{alias}", AuthMiddlewary(http.HandlerFunc(post.DeleteFlatPoster))).Methods(http.MethodDelete, http.MethodOptions)

		apiGroup.Handle("/posters/{alias}/favorites", AuthMiddlewary(http.HandlerFunc(post.AddFavoritePoster))).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.Handle("/posters/favorites", AuthMiddlewary(http.HandlerFunc(post.GetFavoritesPoster))).Methods(http.MethodGet, http.MethodOptions)
		apiGroup.Handle("/posters/{alias}/favorites", AuthMiddlewary(http.HandlerFunc(post.DeleteFavoritePoster))).Methods(http.MethodDelete, http.MethodOptions)
		apiGroup.Handle("/posters/generate-description", AuthMiddlewary(http.HandlerFunc(post.GenerateDescription))).Methods(http.MethodPost, http.MethodOptions)

		apiGroup.Handle("/posters/{alias}/views", AuthMiddlewary(http.HandlerFunc(post.AddViewPoster))).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.Handle("/posters/{alias}/views", http.HandlerFunc(post.GetViewsPoster)).Methods(http.MethodGet, http.MethodOptions)
	} //alias
}
