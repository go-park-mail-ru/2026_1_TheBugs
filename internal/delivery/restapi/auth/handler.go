package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/generated/auth"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/middleware"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/utils"
	ctxLogger "github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/ctxLogger"
	"google.golang.org/grpc"
)

type AuthHandler struct {
	grpcClient auth.AuthServiceClient
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

func NewAuthHandler(grpcConn *grpc.ClientConn) *AuthHandler {
	return &AuthHandler{
		grpcClient: auth.NewAuthServiceClient(grpcConn),
	}
}

func (h AuthHandler) GetAuthMiddlewary() func(http.Handler) http.Handler {
	return middleware.AuthMiddleware(h.grpcClient)
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
// @Security      CSRFToken
// @Failure       400 {object} response.ValidationErrorResponse
// @Failure       404 {object} response.ErrorResponse
// @Failure       500 {object} response.ErrorResponse
// @Security     CSRFToken
// @Router        /auth/reg [post]
func (h AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	op := "AuthHandler.RegisterUser"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	var cred FormDataCreateUser
	err := utils.ParseFormData(r, &cred)
	log.Println(cred)
	if err != nil {
		log.Errorf("parse.ParseFormData: %s", err)
		utils.WriteError(w, "invalid form data", http.StatusBadRequest)
		return
	}

	resp, err := h.grpcClient.RegisterUser(r.Context(), &auth.RegisterUserRequest{
		Email:     cred.Email,
		Password:  cred.Password,
		Phone:     cred.Phone,
		Firstname: cred.FirstName,
		Lastname:  cred.LastName,
	})
	if err != nil {
		log.Errorf("grpcClient.RegisterUser: %s", err)
		utils.HandelGRPCError(w, err)
		return
	}
	log.Infof("RegisterUser response: %v", resp)
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
// @Security     CSRFToken
// @Router        /auth/login [post]
func (h AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	op := "AuthHandler.LoginUser"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	var cred FormDataCredential
	err := utils.ParseFormData(r, &cred)
	log.Println(cred)
	if err != nil {
		log.Errorf("parse.ParseFormData: %s", err)
		utils.WriteError(w, "invalid form data", http.StatusBadRequest)
		return
	}

	loginResp, err := h.grpcClient.LoginUser(r.Context(), &auth.LoginUserRequest{
		Email:    cred.Email,
		Password: cred.Password,
	})
	if err != nil {
		log.Errorf("grpcClient.LoginUser: %s", err)
		utils.HandelGRPCError(w, err)
		return
	}

	resp := LoginResponse{
		AccessToken:    loginResp.AccessToken,
		AccessTokenExp: int(loginResp.AccessTokenExp),
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    loginResp.RefreshToken,
		Path:     "/api/auth",
		Domain:   config.Config.CORS.CookieHost,
		MaxAge:   int(loginResp.RefreshTokenExp),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

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
// @Security     CSRFToken
// @Router        /auth/refresh [post]
func (h AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	op := "AuthHandler.RefreshToken"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		log.Errorf("r.Cookie: %s", err)
		utils.WriteError(w, "invalid cookie", http.StatusBadRequest)
		return
	}
	refreshToken := cookie.Value

	loginResp, err := h.grpcClient.RefreshToken(r.Context(), &auth.RefreshTokenRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		log.Errorf("grpcClient.RefreshToken: %s", err)
		utils.HandelGRPCError(w, err)
		return
	}

	resp := LoginResponse{
		AccessToken:    loginResp.AccessToken,
		AccessTokenExp: int(loginResp.AccessTokenExp),
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    loginResp.RefreshToken,
		Path:     "/api/auth",
		Domain:   config.Config.CORS.CookieHost,
		MaxAge:   int(loginResp.RefreshTokenExp),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	utils.JSONResponse(w, http.StatusOK, &resp)
}

// Logout godoc
// @Summary      User logout
// @Description  Blacklist access/refresh tokens and delete refresh cookie
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CSRFToken
// @Success      204 "Successfully logged out (no content)"
// @Failure      400 {object} response.ValidationErrorResponse "Missing tokens"
// @Failure      401 {object} response.ErrorResponse "Invalid tokens"
// @Failure      500 {object} response.ErrorResponse "Blacklist error"
// @Router       /auth/logout [post]
func (h AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	op := "AuthHandler.Logout"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	accessToken, err := utils.GetAccessToken(r)
	if err != nil {
		log.Errorf("utils.GetAccessToken: %s", err)
		utils.WriteError(w, "invalid access token", http.StatusUnauthorized)
		return
	}
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		log.Errorf("r.Cookie: %s", err)
		utils.WriteError(w, "invalid cookie", http.StatusBadRequest)
		return
	}
	refreshToken := cookie.Value

	_, err = h.grpcClient.Logout(r.Context(), &auth.LogoutRequest{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
	if err != nil {
		log.Printf("grpcClient.Logout: %s", err)
		utils.HandelGRPCError(w, err)
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
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	type VKLoginRequest struct {
		Code         string `json:"code"`
		DeviceID     string `json:"device_id"`
		State        string `json:"state"`
		CodeVerifier string `json:"code_verifier"`
	}
	var flow VKLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&flow); err != nil {
		log.Errorf("json.NewDecoder(r.Body).Decode(&flow): %s", err)
		utils.WriteError(w, "invalid body", http.StatusBadRequest)
		return
	}

	if flow.Code == "" {
		log.Errorf("flow.Code empty")
		utils.WriteError(w, "invalid body", http.StatusBadRequest)
		return
	}

	loginResp, err := h.grpcClient.VKLogin(r.Context(), &auth.VKLoginRequest{
		Code:         flow.Code,
		DeviceId:     flow.DeviceID,
		State:        flow.State,
		CodeVerifier: flow.CodeVerifier,
	})
	if err != nil {
		log.Errorf("grpcClient.VKLogin: %v", err)
		utils.HandelGRPCError(w, err)
		return
	}

	resp := LoginResponse{
		AccessToken:    loginResp.AccessToken,
		AccessTokenExp: int(loginResp.AccessTokenExp),
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    loginResp.RefreshToken,
		Path:     "/api/auth",
		Domain:   config.Config.CORS.CookieHost,
		MaxAge:   int(loginResp.RefreshTokenExp),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

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
// @Security CSRFToken
// @Router        /auth/yandex [post]
func (h AuthHandler) YandexLogin(w http.ResponseWriter, r *http.Request) {
	op := "AuthHandler.YandexLogin"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	type YandexLoginRequest struct {
		Code         string `json:"code"`
		DeviceID     string `json:"device_id"`
		State        string `json:"state"`
		CodeVerifier string `json:"code_verifier"`
	}
	var flow YandexLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&flow); err != nil {
		log.Errorf("json.NewDecoder(r.Body).Decode(&flow): %s", err)
		utils.WriteError(w, "invalid body", http.StatusBadRequest)
		return
	}

	if flow.Code == "" {
		log.Errorf("flow.Code empty")
		utils.WriteError(w, "invalid body", http.StatusBadRequest)
		return
	}

	loginResp, err := h.grpcClient.YandexLogin(r.Context(), &auth.YandexLoginRequest{
		Code:         flow.Code,
		DeviceId:     flow.DeviceID,
		State:        flow.State,
		CodeVerifier: flow.CodeVerifier,
	})
	if err != nil {
		log.Errorf("grpcClient.YandexLogin: %v", err)
		utils.HandelGRPCError(w, err)
		return
	}

	resp := LoginResponse{
		AccessToken:    loginResp.AccessToken,
		AccessTokenExp: int(loginResp.AccessTokenExp),
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    loginResp.RefreshToken,
		Path:     "/api/auth",
		Domain:   config.Config.CORS.CookieHost,
		MaxAge:   int(loginResp.RefreshTokenExp),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

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
// @Security	  CSRFToken
// @Router        /auth/recover [post]
func (h AuthHandler) SendCodeOnEmail(w http.ResponseWriter, r *http.Request) {
	op := "AuthHandler.SendCodeOnEmail"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	type UserEmail struct {
		Email string `json:"email"`
	}
	var data UserEmail
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Errorf("decode: %v", err)
		utils.WriteError(w, "invalid body", http.StatusBadRequest)
		return
	}

	sessionResp, err := h.grpcClient.SendCodeOnEmail(r.Context(), &auth.SendCodeOnEmailRequest{
		Email: data.Email,
	})
	if err != nil {
		log.Errorf("grpcClient.SendCodeOnEmail: %v", err)
		utils.HandelGRPCError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionResp.SessionId,
		Path:     "/api/auth",
		HttpOnly: true,
		Domain:   config.Config.CORS.CookieHost,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(time.Duration(config.Config.JWT.RecoverExp) * time.Second),
		MaxAge:   int(config.Config.JWT.RecoverExp),
	})

	utils.JSONResponse(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

// SendVerifyCodeOnEmail
// @Summary       Send verify code to email
// @Description   Generates recovery session, sends verification code to user's email and sets session_id cookie
// @Tags          Auth
// @Accept        json
// @Produce       json
// @Param         data body request.UserEmail true "User email"
// @Success       200 {object} map[string]string "Status OK"
// @Header        200 {string} Set-Cookie "session_id=<SESSION_ID>; HttpOnly; Path=/api/auth; Max-Age=600"
// @Failure       400 {object} response.ValidationErrorResponse
// @Failure       500 {object} response.ErrorResponse
// @Security	  CSRFToken
// @Router        /auth/email [post]
func (h AuthHandler) SendVerifyCodeOnEmail(w http.ResponseWriter, r *http.Request) {
	op := "AuthHandler.SendVerifyCodeOnEmail"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	type UserEmail struct {
		Email string `json:"email"`
	}
	var data UserEmail
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Errorf("decode: %v", err)
		utils.WriteError(w, "invalid body", http.StatusBadRequest)
		return
	}

	sessionResp, err := h.grpcClient.SendVerifyCodeOnEmail(r.Context(), &auth.SendVerifyCodeOnEmailRequest{
		Email: data.Email,
	})
	if err != nil {
		log.Errorf("grpcClient.SendVerifyCodeOnEmail: %v", err)
		utils.HandelGRPCError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionResp.SessionId,
		Path:     "/api/auth",
		HttpOnly: true,
		Domain:   config.Config.CORS.CookieHost,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(time.Duration(config.Config.JWT.RecoverExp) * time.Second),
		MaxAge:   int(config.Config.JWT.RecoverExp),
	})

	utils.JSONResponse(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

// VerifyUserEmail
// @Summary       Verify user email from code
// @Description   Verifies code from email using session_id cookie and marks session as verified
// @Tags          Auth
// @Accept        json
// @Produce       json
// @Param         data body request.VerifyCodeDTO true "Verification code"
// @Success       200 {object} map[string]string "Status verified"
// @Failure       400 {object} response.ValidationErrorResponse
// @Failure       401 {object} response.ErrorResponse "Invalid code or session"
// @Failure       500 {object} response.ErrorResponse
// @Security      CSRFToken
// @Router        /auth/email/verify [post]
func (h AuthHandler) VerifyUserEmail(w http.ResponseWriter, r *http.Request) {
	op := "AuthHandler.VerifyUserEmail"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Errorf("r.Cookie: %s", err)
		utils.WriteError(w, "invalid cookie", http.StatusBadRequest)
		return
	}
	sessionId := cookie.Value
	if sessionId == "" {
		log.Errorf("sessionId is empty")
		utils.WriteError(w, "session_id is empty", http.StatusBadRequest)
		return
	}

	type VerifyCodeDTO struct {
		Code string `json:"code"`
	}
	var data VerifyCodeDTO
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Errorf("decode: %v", err)
		utils.WriteError(w, "invalid body", http.StatusBadRequest)
		return
	}

	_, err = h.grpcClient.VerifyCode(r.Context(), &auth.VerifyCodeRequest{
		SessionId: sessionId,
		Code:      data.Code,
	})
	if err != nil {
		log.Errorf("grpcClient.VerifyCode: %v", err)
		utils.HandelGRPCError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{
		"status": "email verified",
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
// @Security      CSRFToken
// @Router        /auth/recover/verify [post]
func (h AuthHandler) VerifyRecoveryCode(w http.ResponseWriter, r *http.Request) {
	op := "AuthHandler.VerifyRecoveryCode"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Errorf("r.Cookie: %s", err)
		utils.WriteError(w, "invalid cookie", http.StatusBadRequest)
		return
	}
	sessionId := cookie.Value
	if sessionId == "" {
		log.Errorf("sessionId is empty")
		utils.WriteError(w, "session_id is empty", http.StatusBadRequest)
		return
	}

	type VerifyCodeDTO struct {
		Code string `json:"code"`
	}
	var data VerifyCodeDTO
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Errorf("decode: %v", err)
		utils.WriteError(w, "invalid body", http.StatusBadRequest)
		return
	}

	_, err = h.grpcClient.VerifyRecoveryCode(r.Context(), &auth.VerifyRecoveryCodeRequest{
		SessionId: sessionId,
		Code:      data.Code,
	})
	if err != nil {
		log.Errorf("grpcClient.VerifyRecoveryCode: %v", err)
		utils.HandelGRPCError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{
		"status": "verified",
	})
}

// UpdatePassword
// @Summary       Verify user email
// @Description   Updates password using verified recovery session (session_id cookie required)
// @Tags          Auth
// @Accept        json
// @Produce       json
// @Param         data body request.UpdatePwdDTO true "New password"
// @Success       200 {object} map[string]string "Password updated"
// @Failure       400 {object} response.ValidationErrorResponse
// @Failure       401 {object} response.ErrorResponse "Session not verified or invalid"
// @Failure       500 {object} response.ErrorResponse
// @Security      CSRFToken
// @Router        /auth/recover/reset [post]
func (h AuthHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	op := "AuthHandler.UpdatePassword"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Errorf("r.Cookie: %s", err)
		utils.WriteError(w, "invalid cookie", http.StatusBadRequest)
		return
	}
	sessionId := cookie.Value
	if sessionId == "" {
		log.Errorf("sessionId is empty")
		utils.WriteError(w, "session_id is empty", http.StatusBadRequest)
		return
	}

	type UpdatePwdDTO struct {
		Password string `json:"password"`
	}
	var data UpdatePwdDTO
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Errorf("decode: %v", err)
		utils.WriteError(w, "invalid body", http.StatusBadRequest)
		return
	}

	_, err = h.grpcClient.UpdatePassword(r.Context(), &auth.UpdatePasswordRequest{
		SessionId: sessionId,
		Password:  data.Password,
	})
	if err != nil {
		log.Errorf("grpcClient.UpdatePassword: %v", err)
		utils.HandelGRPCError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{
		"status": "password_updated",
	})
}

// @Summary       Get CSRF token
// @Description   Create and get csrf token
// @Tags          Auth
// @Accept        json
// @Produce       json
// @Success       200 {object} map[string]string "CSRF token"
// @Failure       403 {object} response.ErrorResponse
// @Failure       500 {object} response.ErrorResponse
// @Router        /csrf-token [get]
func (h AuthHandler) GetCSRFToken(w http.ResponseWriter, r *http.Request) {
	token := generateCSRFToken()
	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   3600,
	})
	log.Println(token)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"csrf_token": token})
}

func generateCSRFToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
