package auth

import (
	"encoding/json"
	"net/http"

	"log"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/auth"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/parse"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/utils"
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
		utils.HandelError(w, err)
		return
	}
	err = h.uc.RegisterUseCase(cred.Email, cred.Password)
	if err != nil {
		log.Printf("h.uc.RegisterUseCase: %s", err)
		utils.HandelError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var cred FormDataCredential
	err := parse.ParseFormData(r, &cred)
	log.Println(cred)
	if err != nil {
		log.Printf("parse.ParseFormData: %s", err)
		utils.HandelError(w, err)
		return
	}
	accessCred, err := h.uc.LoginUseCase(cred.Email, cred.Password)
	if err != nil {
		log.Printf("h.uc.LoginUseCase: %s", err)
		utils.HandelError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessCred)
}

func (h AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		log.Printf("r.Cookie: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}
	refreshToken := cookie.Value
	accessCred, err := h.uc.RefreshTokenUseCase(refreshToken)
	if err != nil {
		log.Printf("h.uc.RefreshTokenUseCase: %s", err)
		utils.HandelError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessCred)
}
