package auth

import (
	"net/http"

	"log"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/auth"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/parse"
)

type AuthHandler struct {
	uc *auth.AuthUseCase
}

type FormDataCredential struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

func NewAuthHandler(uc *auth.AuthUseCase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

func (h AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var cred FormDataCredential
	err := parse.ParseFormData(r, &cred)
	log.Println(cred)
	if err != nil {
		log.Printf("parse.ParseFormData: %s", err)
		http.Error(w, "bad request", http.StatusBadRequest)
	}
	err = h.uc.RegisterUseCase(cred.Email, cred.Password)
	if err != nil {
		log.Printf("h.uc.RegisterUseCase: %s", err)
		http.Error(w, "bad request", http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusNoContent)
}
