package promotion

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/utils"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/ctxLogger"
	yoopayment "github.com/rvinnie/yookassa-sdk-go/yookassa/payment"
	yoowebhook "github.com/rvinnie/yookassa-sdk-go/yookassa/webhook"
)

type PromotionHandler struct {
	usecase delivery.PromotionUseCase
}

func NewPromotionHandler(usecase delivery.PromotionUseCase) *PromotionHandler {
	return &PromotionHandler{
		usecase: usecase,
	}
}

type CreatePaymentRequest struct {
	PromotionCode string `json:"promotion_code"`
	PosterID      int    `json:"poster_id"`
}

type CheckPaymentRequest struct {
	PaymentID string `json:"payment_id"`
}
type CheckPaymentResponse struct {
	Stutus string `json:"payment_status"`
}

// @Summary Create promotion payment
// @Description Creates YooKassa payment and returns confirmation URL
// @Tags promotion
// @Accept json
// @Produce json
// @Param body body CreatePaymentRequest true "Create payment request"
// @Success 200 {object} dto.PaymentDTO
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Security     BearerAuth
// @Security     CSRFToken
// @Router /promotions/payment [post]
func (h *PromotionHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	op := "PromotionHandler.CreatePayment"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	var req CreatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Errorf("json.Decode: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}
	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		log.Errorf("utils.GetUserID: %s", err)
		utils.HandelError(w, err)
		return
	}

	if req.PromotionCode == "" || req.PosterID == 0 {
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	data, err := h.usecase.CreatePaymentOrder(r.Context(), req.PromotionCode, req.PosterID, userID)
	if err != nil {
		log.Errorf("h.usecase.CreatePaymentOrder: %s", err)
		utils.HandelError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, data)
}

// @Summary YooKassa webhook
// @Description Handles YooKassa payment events
// @Tags promotion
// @Accept json
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /promotions/webhooks/yookassa [post]
func (h *PromotionHandler) YooKassaWebhook(w http.ResponseWriter, r *http.Request) {
	op := "PromotionHandler.YooKassaWebhook"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	var req yoowebhook.WebhookEvent[yoopayment.Payment]
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.usecase.ActivatePromotion(r.Context(), req); err != nil {
		log.Errorf("h.usecase.ActivatePromotion: %s", err)
		utils.HandelError(w, entity.ServiceError)
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{"status": "ok"})
}

// @Summary YooKassa check status
// @Description YooKassa check status
// @Tags promotion
// @Accept json
// @Produce json
// @Param body body CheckPaymentRequest true "Create payment request"
// @Success 200 {object} CheckPaymentResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Security     BearerAuth
// @Security     CSRFToken
// @Router /promotions/status [post]
func (h *PromotionHandler) CheckPaymentStatus(w http.ResponseWriter, r *http.Request) {
	op := "PromotionHandler.CheckPaymentStautes"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	var req CheckPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		log.Errorf("utils.GetUserID: %s", err)
		utils.HandelError(w, err)
		return
	}

	status, err := h.usecase.CheckPaymentStatus(r.Context(), userID, req.PaymentID)
	if err != nil {
		log.Errorf("h.usecase.ActivatePromotion: %s", err)
		utils.HandelError(w, entity.ServiceError)
		return
	}

	utils.JSONResponse(w, http.StatusOK, CheckPaymentResponse{Stutus: string(status)})
}

// @Summary Get user`s promotions
// @Description Get user`s promotions
// @Tags promotion
// @Accept json
// @Produce json
// @Success 200 {object} dto.UserPromotionsDTO
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Security     BearerAuth
// @Security     CSRFToken
// @Router /promotions/me [get]
func (h *PromotionHandler) GetUserPromotions(w http.ResponseWriter, r *http.Request) {
	op := "PromotionHandler.GetUserPromotions"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		log.Errorf("utils.GetUserID: %s", err)
		utils.HandelError(w, err)
		return
	}

	resp, err := h.usecase.GetPromotionByUserID(r.Context(), userID)
	if err != nil {
		log.Errorf("h.usecase.ActivatePromotion: %s", err)
		utils.HandelError(w, entity.ServiceError)
		return
	}

	utils.JSONResponse(w, http.StatusOK, resp)
}
