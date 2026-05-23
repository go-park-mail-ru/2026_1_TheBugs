package restapi

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/auth"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/complex"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/middleware"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/poster"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/promotion"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/support"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/user"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	prom "github.com/go-park-mail-ru/2026_1_TheBugs/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
func RegisterHandlers(app *mux.Router, logger *logrus.Logger, auth *auth.AuthHandler, post *poster.PosterHandler, UtilityCompany *complex.UtilityCompanyHandler, user *user.UserHandler, support *support.SupportHandler, payment *promotion.PromotionHandler, rps delivery.RateLimitUseCase) {

	c := cors.New(cors.Options{
		AllowedOrigins:   config.Config.CORS.AllowedHosts,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		ExposedHeaders:   []string{"Set-Cookie", "X-CSRF-TOKEN"},
	})
	metrics := prom.NewMetricsMiddleware()
	metrics.Register(entity.ServiceType(config.Config.ServiceName))

	app.Handle("/metrics", promhttp.Handler())
	app.Use(middleware.LoggingMiddleware(logger))
	app.Use(c.Handler)

	app.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}).Methods(http.MethodGet)
	apiGroup := app.PathPrefix("/api").Subrouter()

	apiGroup.HandleFunc("/webhooks/yookassa", payment.YooKassaWebhook).Methods(http.MethodPost, http.MethodOptions)

	restAPI := apiGroup.PathPrefix("/").Subrouter()
	restAPI.Use(metrics.MetricsHTTPMiddleware)
	restAPI.Use(middleware.SecurityMiddleware)
	restAPI.Use(mux.CORSMethodMiddleware(restAPI))

	{
		apiGroup.HandleFunc("/csrf-token", auth.GetCSRFToken).Methods(http.MethodGet)

		AuthMiddlewary := auth.GetAuthMiddlewary()
		RateLimitMiddlewary := middleware.RateLimitMiddleware(rps)
		UserIDMiddleware := auth.GetUserIDMiddlewary()
		apiGroup.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
		}).Methods(http.MethodGet)
		apiGroup.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)

		apiGroup.Handle("/user/me", AuthMiddlewary(http.HandlerFunc(user.GetMe))).Methods(http.MethodGet, http.MethodOptions)
		apiGroup.Handle("/user/me/profile", AuthMiddlewary(http.HandlerFunc(user.UpdateProfile))).Methods(http.MethodPut, http.MethodOptions)

		apiGroup.Handle("/user/me/roommate-form", AuthMiddlewary(http.HandlerFunc(user.CreateRoommateForm))).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.Handle("/user/me/roommate-form", AuthMiddlewary(http.HandlerFunc(user.GetRoommateForm))).Methods(http.MethodGet, http.MethodOptions)
		apiGroup.Handle("/user/me/roommate-form", AuthMiddlewary(http.HandlerFunc(user.UpdateRoommateForm))).Methods(http.MethodPut, http.MethodOptions)
		apiGroup.HandleFunc("/user/{id}", user.GetRoommateUser).Methods(http.MethodGet, http.MethodOptions)
		apiGroup.Handle("/user/{id}/match", AuthMiddlewary(http.HandlerFunc(user.AddRoommateMatch))).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.Handle("/user/{id}/contacts", AuthMiddlewary(http.HandlerFunc(user.GetRoommateContacts))).Methods(http.MethodGet, http.MethodOptions)
		apiGroup.Handle("/user/me/roommate-matches/incoming", AuthMiddlewary(http.HandlerFunc(user.GetIncomingRoommateMatches))).Methods(http.MethodGet, http.MethodOptions)
		apiGroup.Handle("/user/me/roommate-matches/matched", AuthMiddlewary(http.HandlerFunc(user.GetMatchedRoommateMatches))).Methods(http.MethodGet, http.MethodOptions)

		apiGroup.HandleFunc("/auth/login", auth.LoginUser).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.HandleFunc("/auth/admin/login", auth.LoginAdminUser).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.HandleFunc("/auth/reg", auth.RegisterUser).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.HandleFunc("/auth/logout", auth.Logout).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.HandleFunc("/auth/refresh", auth.RefreshToken).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.HandleFunc("/auth/vkid", auth.VKLogin).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.HandleFunc("/auth/yandex", auth.YandexLogin).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.HandleFunc("/auth/recover", auth.SendCodeOnEmail).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.HandleFunc("/auth/recover/verify", auth.VerifyRecoveryCode).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.HandleFunc("/auth/recover/reset", auth.UpdatePassword).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.HandleFunc("/auth/email/verify", auth.VerifyUserEmail).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.HandleFunc("/auth/email", auth.SendVerifyCodeOnEmail).Methods(http.MethodPost, http.MethodOptions)

		apiGroup.HandleFunc("/utility-companies/by-alias/{alias}", UtilityCompany.GetUtilityCompany).Methods(http.MethodGet)
		apiGroup.HandleFunc("/utility-companies/developers", UtilityCompany.GetAllDevelopers).Methods(http.MethodGet)
		apiGroup.HandleFunc("/utility-companies/", UtilityCompany.GetUtilityCompaniesByDeveloper).Methods(http.MethodGet)

		apiGroup.HandleFunc("/posters/flats", post.GetFlatsAll).Methods(http.MethodGet)
		apiGroup.HandleFunc("/posters/geo", post.GetFlatsMapAll).Methods(http.MethodGet, http.MethodOptions)
		apiGroup.HandleFunc("/posters/by-point", post.GetPostersByPoint).Methods(http.MethodGet, http.MethodOptions)
		apiGroup.HandleFunc("/posters/by-alias/{alias}", post.GetPoster).Methods(http.MethodGet)

		apiGroup.Handle("/posters/me", AuthMiddlewary(http.HandlerFunc(post.GetPostersByUser))).Methods(http.MethodGet, http.MethodOptions)
		apiGroup.Handle("/posters/me/{alias}", AuthMiddlewary(http.HandlerFunc(post.GetPostersByUserByAlias))).Methods(http.MethodGet, http.MethodOptions)
		apiGroup.Handle("/posters/flat", RateLimitMiddlewary(10)(AuthMiddlewary(http.HandlerFunc(post.CreateFlatPoster)))).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.Handle("/posters/flat/{alias}", RateLimitMiddlewary(10)(AuthMiddlewary(http.HandlerFunc(post.UpdateFlatPoster)))).Methods(http.MethodPut, http.MethodOptions)
		apiGroup.Handle("/posters/flat/{alias}", AuthMiddlewary(http.HandlerFunc(post.DeleteFlatPoster))).Methods(http.MethodDelete, http.MethodOptions)

		apiGroup.Handle("/posters/{alias}/favorites", AuthMiddlewary(http.HandlerFunc(post.AddFavoritePoster))).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.Handle("/posters/favorites", AuthMiddlewary(http.HandlerFunc(post.GetFavoritesPoster))).Methods(http.MethodGet, http.MethodOptions)
		apiGroup.Handle("/posters/{alias}/favorites", AuthMiddlewary(http.HandlerFunc(post.DeleteFavoritePoster))).Methods(http.MethodDelete, http.MethodOptions)
		apiGroup.Handle("/posters/{alias}/favorites", UserIDMiddleware(http.HandlerFunc(post.GetFavoritesCountPoster))).Methods(http.MethodGet, http.MethodOptions)
		apiGroup.Handle("/posters/generate-description", RateLimitMiddlewary(10)(AuthMiddlewary(http.HandlerFunc(post.GenerateDescription)))).Methods(http.MethodPost, http.MethodOptions)

		apiGroup.Handle("/posters/{alias}/views", AuthMiddlewary(http.HandlerFunc(post.AddViewPoster))).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.Handle("/posters/{alias}/views", http.HandlerFunc(post.GetViewsPoster)).Methods(http.MethodGet, http.MethodOptions)
		apiGroup.Handle("/posters/{alias}/price-history", http.HandlerFunc(post.GetPriceHistoryPoster)).Methods(http.MethodGet, http.MethodOptions)

		apiGroup.Handle("/posters/{alias}/roommates", http.HandlerFunc(post.GetPosterRoommates)).Methods(http.MethodGet, http.MethodOptions)
		apiGroup.Handle("/posters/{alias}/roommates", AuthMiddlewary(http.HandlerFunc(post.AddPosterRoommate))).Methods(http.MethodPost, http.MethodOptions)

		apiGroup.Handle("/support/orders", http.HandlerFunc(support.CreateOrder)).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.Handle("/support/orders", AuthMiddlewary(http.HandlerFunc(support.GetOrders))).Methods(http.MethodGet, http.MethodOptions)
		apiGroup.Handle("/support/orders/{id}", AuthMiddlewary(http.HandlerFunc(support.GetOrderByID))).Methods(http.MethodGet, http.MethodOptions)
		apiGroup.Handle("/support/orders/{id}/answer", AuthMiddlewary(http.HandlerFunc(support.AnswerOrder))).Methods(http.MethodPost, http.MethodOptions)

		apiGroup.Handle("/promotions/payment", AuthMiddlewary(http.HandlerFunc(payment.CreatePayment))).Methods(http.MethodPost, http.MethodOptions)
		//apiGroup.Handle("/promotions/webhooks/yookassa", middleware.IPFilterMiddleware(http.HandlerFunc(payment.YooKassaWebhook), middleware.AllowYookassaIPs)).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.Handle("/promotions/status", AuthMiddlewary(http.HandlerFunc(payment.CheckPaymentStatus))).Methods(http.MethodPost, http.MethodOptions)
		apiGroup.Handle("/promotions/me", AuthMiddlewary(http.HandlerFunc(payment.GetUserPromotions))).Methods(http.MethodGet, http.MethodOptions)
	} //alias
}
