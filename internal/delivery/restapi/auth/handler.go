package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/middleware"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/request"
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

// SendCodeOnEmail
// @Summary       Send recovery code to email
// @Description   Generates recovery session, sends verification code to user's email and sets session_id cookie
// @Tags          Auth
// @Accept        json
// @Produce       json
// @Param         data body request.UserEmail true "User email"
// @Success       200 {object} map[string]string "Status OK"
// @Header        200 {string} Set-Cookie "session_id=<SESSION_ID>; HttpOnly; Path=/api/auth; Max-Age=600"
// @Failure       400 {object} response.ValidationErrorResponse
// @Failure       500 {object} response.ErrorResponse
// @Router        /auth/recover [post]
func (h AuthHandler) SendCodeOnEmail(w http.ResponseWriter, r *http.Request) {
	op := "AuthHandler.SendCodeOnEmail"
	log := middleware.GetLogger(r.Context()).WithField("op", op)

	var data request.UserEmail
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Errorf("decode: %v", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	sessionID, err := h.uc.SendVerificationCode(r.Context(), data.Email)
	if err != nil {
		log.Errorf("SendVerificationCode: %v", err)
		utils.HandelError(w, err)
		return
	}
	http.SetCookie(
		w,
		&http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/api/auth",
			HttpOnly: true,
			Domain:   config.Config.CORS.CookieHost,
			SameSite: http.SameSiteLaxMode,
			Expires:  time.Now().Add(time.Duration(int(config.Config.JWT.RecoverExp) * int(time.Second))),
			MaxAge:   int(config.Config.JWT.RecoverExp),
		},
	)

	utils.JSONResponse(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

// VerifyRecoveryCode
// @Summary       Verify recovery code
// @Description   Verifies code from email using session_id cookie and marks session as verified
// @Tags          Auth
// @Accept        json
// @Produce       json
// @Param         data body request.VerifyCodeDTO true "Verification code"
// @Success       200 {object} map[string]string "Status verified"
// @Failure       400 {object} response.ValidationErrorResponse
// @Failure       401 {object} response.ErrorResponse "Invalid code or session"
// @Failure       500 {object} response.ErrorResponse
// @Router        /auth/recover/verify [post]
func (h AuthHandler) VerifyRecoveryCode(w http.ResponseWriter, r *http.Request) {
	op := "AuthHandler.VerifyRecoveryCode"
	log := middleware.GetLogger(r.Context()).WithField("op", op)

	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Errorf("r.Cookie: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}
	sessionId := cookie.Value
	if sessionId == "" {
		log.Errorf("sessionId is empty")
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	var data request.VerifyCodeDTO
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Errorf("decode: %v", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	err = h.uc.CheckRecoveryCode(r.Context(), sessionId, data.Code)
	if err != nil {
		log.Errorf("CheckRecoveryCode: %v", err)
		utils.HandelError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{
		"status": "verified",
	})
}

// UpdatePassword
// @Summary       Update user password
// @Description   Updates password using verified recovery session (session_id cookie required)
// @Tags          Auth
// @Accept        json
// @Produce       json
// @Param         data body request.UpdatePwdDTO true "New password"
// @Success       200 {object} map[string]string "Password updated"
// @Failure       400 {object} response.ValidationErrorResponse
// @Failure       401 {object} response.ErrorResponse "Session not verified or invalid"
// @Failure       500 {object} response.ErrorResponse
// @Router        /auth/recover/reset [post]
func (h AuthHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	op := "AuthHandler.UpdatePassword"
	log := middleware.GetLogger(r.Context()).WithField("op", op)
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Errorf("r.Cookie: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}
	sessionId := cookie.Value
	if sessionId == "" {
		log.Errorf("sessionId is empty")
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	var data request.UpdatePwdDTO
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Errorf("decode: %v", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	err = h.uc.UpdateUserPassword(r.Context(), sessionId, data.Password)
	if err != nil {
		log.Errorf("UpdateUserPassword: %v", err)
		utils.HandelError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{
		"status": "password_updated",
	})
}
