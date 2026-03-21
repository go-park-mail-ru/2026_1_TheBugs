package auth

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/middleware"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/utils"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/auth"
)

type AuthHandler struct {
	uc *auth.AuthUseCase
}

type FormDataCredential struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

type FormDataCreateUser struct {
	Email     string `schema:"email"`
	Password  string `schema:"password"`
	Phone     string `schema:"phone"`
	FirstName string `schema:"firstname"`
	LastName  string `schema:"lastname"`
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
// @Param         phone formData string true "User phone"
// @Param         firstname formData string true "User firstname"
// @Param         lastname formData string true "User lastname"
// @Success       204
// @Failure       400 {object} response.ValidationErrorResponse
// @Failure       404 {object} response.ErrorResponse
// @Failure       500 {object} response.ErrorResponse
// @Router        /auth/reg [post]
func (h AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	op := "AuthHandler.RegisterUser"
	log := middleware.GetLogger(r.Context()).WithField("op", op)

	var cred FormDataCreateUser
	err := utils.ParseFormData(r, &cred)
	log.Println(cred)
	if err != nil {
		log.Errorf("parse.ParseFormData: %s", err)
		utils.HandelError(w, err)
		return
	}
	err = h.uc.RegisterUseCase(r.Context(), dto.CreateUserDTO{
		Email:     cred.Email,
		Password:  cred.Password,
		Phone:     cred.Phone,
		LastName:  cred.LastName,
		FirstName: cred.FirstName,
	})
	if err != nil {
		log.Errorf("h.uc.RegisterUseCase: %s", err)
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
// @Failure       400 {object} response.ValidationErrorResponse
// @Failure       404 {object} response.ErrorResponse
// @Failure       500 {object} response.ErrorResponse
// @Router        /auth/login [post]
func (h AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	op := "AuthHandler.LoginUser"
	log := middleware.GetLogger(r.Context()).WithField("op", op)

	var cred FormDataCredential

	err := utils.ParseFormData(r, &cred)
	log.Println(cred)
	if err != nil {
		log.Errorf("parse.ParseFormData: %s", err)
		utils.HandelError(w, err)
		return
	}
	accessCred, err := h.uc.LoginUseCase(r.Context(), cred.Email, cred.Password)
	if err != nil {
		log.Errorf("h.uc.LoginUseCase: %s", err)
		utils.HandelError(w, err)
		return
	}
	resp := LoginResponse{
		AccessToken:    accessCred.AccessToken,
		AccessTokenExp: accessCred.AccessTokenExp,
	}

	utils.SetRefreshCookie(w, accessCred)

	utils.JSONResponse(w, http.StatusOK, &resp)

}

// RefreshToken
// @Summary       Refresh access token
// @Description   Obtain new access token using refresh token stored in cookie (refresh_token cookie required)
// @Tags          Auth
// @Accept        json
// @Produce       json
// @Success       200 {object} LoginResponse "New access token, also updates refresh token cookie"
// @Header        200 {string} Set-Cookie "new refresh_token=...; HttpOnly; Path=/api/auth/refresh; Max-Age=..."
// @Failure       400 {object} response.ValidationErrorResponse
// @Failure       401 {object} response.ErrorResponse
// @Failure       500 {object} response.ErrorResponse
// @Router        /auth/refresh [post]
func (h AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	op := "AuthHandler.RefreshToken"
	log := middleware.GetLogger(r.Context()).WithField("op", op)

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		log.Errorf("r.Cookie: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}
	refreshToken := cookie.Value
	accessCred, err := h.uc.RefreshTokenUseCase(r.Context(), refreshToken)
	if err != nil {
		log.Errorf("h.uc.RefreshTokenUseCase: %s", err)
		utils.HandelError(w, err)
		return
	}
	resp := LoginResponse{
		AccessToken:    accessCred.AccessToken,
		AccessTokenExp: accessCred.AccessTokenExp,
	}

	utils.SetRefreshCookie(w, accessCred)

	utils.JSONResponse(w, http.StatusOK, &resp)
}

// Logout godoc
// @Summary      User logout
// @Description  Blacklist access/refresh tokens and delete refresh cookie
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      204 "Successfully logged out (no content)"
// @Failure      400 {object} response.ValidationErrorResponse "Missing tokens"
// @Failure      401 {object} response.ErrorResponse "Invalid tokens"
// @Failure      500 {object} response.ErrorResponse "Blacklist error"
// @Router       /auth/logout [post]
func (h AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	op := "AuthHandler.Logout"
	log := middleware.GetLogger(r.Context()).WithField("op", op)

	accessToken, err := utils.GetAccessToken(r)
	if err != nil {
		log.Errorf("utils.GetAccessToken: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		log.Errorf("r.Cookie: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}
	refreshToken := cookie.Value
	err = h.uc.LogoutUseCase(r.Context(), dto.LogoutDTO{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	})
	if err != nil {
		log.Printf("h.uc.RefreshTokenUseCase: %s", err)
		utils.HandelError(w, err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/api/auth",
		Domain:   config.Config.CORS.CookieHost,
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})
	w.WriteHeader(http.StatusNoContent)
}

// POST /api/auth/vk/login
func (h AuthHandler) VKLogin(w http.ResponseWriter, r *http.Request) {
	op := "AuthHandler.VKLogin"
	log := middleware.GetLogger(r.Context()).WithField("op", op)

	var flow dto.OAuthCodeFlow
	if err := json.NewDecoder(r.Body).Decode(&flow); err != nil {
		log.Errorf("json.NewDecoder(r.Body).Decode(&flow): %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	if flow.Code == "" || flow.DeviceID == nil || flow.State == nil {
		log.Errorf("flow.Code || flow.DeviceID || flow.State empty ")
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	accessCred, err := h.uc.LoginUserFromVKUseCase(r.Context(), flow)
	if err != nil {
		log.Errorf("login vk error: %v", err)
		utils.HandelError(w, err)
		return
	}

	resp := LoginResponse{
		AccessToken:    accessCred.AccessToken,
		AccessTokenExp: accessCred.AccessTokenExp,
	}

	utils.SetRefreshCookie(w, accessCred)
	utils.JSONResponse(w, http.StatusOK, &resp)
}

// LoginYandexUser
// @Summary       Login user from Yandex
// @Description   Authenticate user and return access token + set refresh token cookie
// @Tags          Auth
// @Accept        json
// @Param  flow body dto.OAuthCodeFlow true "OAuth user cred by Authorization code flow"
// @Success       200 {object} LoginResponse "Successful login, returns access token"
// @Header        200 {string} Set-Cookie "refresh_token=<NEW_REFRESH_TOKEN>; HttpOnly; Path=/api/auth/refresh; Max-Age=..."
// @Failure       400 {object} response.ValidationErrorResponse
// @Failure       404 {object} response.ErrorResponse
// @Failure       500 {object} response.ErrorResponse
// @Router        /auth/yandex [post]
func (h AuthHandler) YandexLogin(w http.ResponseWriter, r *http.Request) {
	op := "AuthHandler.YandexLogin"
	log := middleware.GetLogger(r.Context()).WithField("op", op)

	var flow dto.OAuthCodeFlow
	if err := json.NewDecoder(r.Body).Decode(&flow); err != nil {
		log.Errorf("json.NewDecoder(r.Body).Decode(&flow): %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	if flow.Code == "" {
		log.Errorf("flow.Code empty ")
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	accessCred, err := h.uc.LoginUserFromYandexUseCase(r.Context(), flow)
	if err != nil {
		log.Errorf("login yandex error: %v", err)
		utils.HandelError(w, err)
		return
	}

	resp := LoginResponse{
		AccessToken:    accessCred.AccessToken,
		AccessTokenExp: accessCred.AccessTokenExp,
	}

	utils.SetRefreshCookie(w, accessCred)
	utils.JSONResponse(w, http.StatusOK, &resp)
}
