package restapi

import (
	"2026_1_the_bugs/internal/delivery/restapi/middleware"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(app *mux.Router) {
	app.Use(middleware.LoggingMiddleware)

	// Routers
	apiGroup := app.PathPrefix("/api").Subrouter()
	{
		apiGroup.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(map[string]bool{"ok": true})
		}).Methods(http.MethodGet)
	}
}
