package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"log"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/utils"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"
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

type LoginResponse struct {
	AccessToken    string `json:"access_token"`
	AccessTokenExp int    `json:"expire_at"`
}

func NewAuthHandler(uc *auth.AuthUseCase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

// RegisterUser
// @Summary       Register new user
// @Description   Register new user with email and password
// @Tags          Auth
// @Accept        x-www-form-urlencoded
// @Param         email formData string true "User email"
// @Param         password formData string true "User password"
// @Success       204
// @Failure       400 {string} string "Bad Request"
// @Failure       404 {string} string "Not Found"
// @Failure       500 {string} string "Internal error"
// @Router        /auth/reg [post]
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

// LoginUser
// @Summary       Login user
// @Description   Authenticate user and return access token + set refresh token cookie
// @Tags          Auth
// @Accept        x-www-form-urlencoded
// @Param         email formData string true "User email"
// @Param         password formData string true "User password"
// @Success       200 {object} LoginResponse "Successful login, returns access token"
// @Header        200 {string} Set-Cookie "refresh_token=<NEW_REFRESH_TOKEN>; HttpOnly; Path=/api/auth/refresh; Max-Age=..."
// @Failure       400 {string} string "Bad Request"
// @Failure       404 {string} string "Not Found"
// @Failure       500 {string} string "Internal error"
// @Router        /auth/login [post]
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
	resp := LoginResponse{
		AccessToken:    accessCred.AccessToken,
		AccessTokenExp: accessCred.AccessTokenExp,
	}

	setRefreshCookie(w, accessCred)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(&resp)
}

// RefreshToken
// @Summary       Refresh access token
// @Description   Obtain new access token using refresh token stored in cookie (refresh_token cookie required)
// @Tags          Auth
// @Accept        json
// @Produce       json
// @Success       200 {object} LoginResponse "New access token, also updates refresh token cookie"
// @Header        200 {string} Set-Cookie "new refresh_token=...; HttpOnly; Path=/api/auth/refresh; Max-Age=..."
// @Failure       400 {string} string "Bad Request"
// @Failure       401 {string} string "Unauthorized - invalid or missing refresh token"
// @Failure       500 {string} string "Internal error"
// @Router        /auth/refresh [post]
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
	resp := LoginResponse{
		AccessToken:    accessCred.AccessToken,
		AccessTokenExp: accessCred.AccessTokenExp,
	}

	setRefreshCookie(w, accessCred)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func setRefreshCookie(w http.ResponseWriter, accessCred *dto.UserAccessCredDTO) {
	http.SetCookie(
		w,
		&http.Cookie{
			Name:     "refresh_token",
			Value:    accessCred.RefreshToken,
			Path:     "/api/auth/refresh",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Expires:  time.Now().Add(time.Duration(accessCred.RefreshTokenExp * int(time.Second))),
			MaxAge:   accessCred.RefreshTokenExp,
		},
	)
}
