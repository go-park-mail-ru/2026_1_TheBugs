package order

import (
	"net/http"
	"strings"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/utils"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/order"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/ctxLogger"
)

type OrderHandler struct {
	uc *order.OrderUseCase
}

func NewOrderHandler(uc *order.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		uc: uc,
	}
}

// @Summary Create handling
// @Description Creates a user handling (order) with optional photos
// @Tags handlings
// @Produce json
// @Security     BearerAuth
// @Security     CSRFToken
// @Param category_id formData integer true "Handling category ID"
// @Param description formData string true "Handling description"
// @Param photos.0.file formData file false "First photo file"
// @Param photos.0.order formData integer false "First photo order"
// @Param photos.1.file formData file false "Second photo file"
// @Param photos.1.order formData integer false "Second photo order"
// @Success 201
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /support/orders [post]
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	op := "OrderHandler.CreateOrder"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		log.Errorf("utils.GetUserID: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	var req dto.OrderDTO

	req.UserID = userID
	err = utils.ParseMultipartFormData(r, &req)
	if err != nil {
		log.Errorf("utils.ParseFormData: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	req.Images, err = utils.ParsePhotos(r)
	if err != nil {
		log.Errorf("utils.ParsePhotos: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	err = h.uc.CreateOrder(r.Context(), &req)
	if err != nil {
		log.Errorf("h.uc.CreateOrder: %s", err)
		utils.HandelError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// @Summary Get support agent response
// @Description Sends user prompt to support agent and returns AI response
// @Tags support
// @Accept mpfd
// @Produce json
// @Security BearerAuth
// @Security CSRFToken
// @Param user_prompt formData string true "User prompt"
// @Success 200 {object} dto.ChatResult
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /support/agent-response [post]
func (h *OrderHandler) GetSupportAgentResponse(w http.ResponseWriter, r *http.Request) {
	op := "OrderHandler.GetSupportAgentResponse"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		if err := r.ParseForm(); err != nil {
			log.Errorf("r.ParseForm: %s", err)
			utils.HandelError(w, entity.InvalidInput)
			return
		}
	}

	userPrompt := strings.TrimSpace(r.FormValue("user_prompt"))
	if userPrompt == "" {
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	resp, err := h.uc.GetSupportAgentResponse(r.Context(), userPrompt)
	if err != nil {
		log.Errorf("h.uc.GetSupportAgentResponse: %s", err)
		utils.HandelError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, resp)
}

// @Summary Get user orders
// @Description Returns all support orders of authorized user
// @Tags handlings
// @Produce json
// @Security     BearerAuth
// @Security     CSRFToken
// @Success 200 {object} dto.OrdersResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /support/orders [get]
func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	op := "OrderHandler.GetOrders"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		log.Errorf("utils.GetUserID: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	orders, err := h.uc.GetUserOrders(r.Context(), userID)
	if err != nil {
		log.Errorf("h.uc.GetUserOrders: %s", err)
		utils.HandelError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, orders)
}

// @Summary Get order by id
// @Description Returns full support order info by id
// @Tags handlings
// @Produce json
// @Security     BearerAuth
// @Security     CSRFToken
// @Param id path integer true "Order ID"
// @Success 200 {object} dto.OrderFullDTO
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /support/orders/{id} [get]
func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	op := "OrderHandler.GetOrderByID"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		log.Errorf("utils.GetUserID: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	orderID, err := utils.ParseIDFromRequest(r)
	if err != nil {
		log.Errorf("strconv.Atoi: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	order, err := h.uc.GetOrderByID(r.Context(), userID, orderID)
	if err != nil {
		log.Errorf("utils.ParseIDFromRequest: %s", err)
		utils.HandelError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, order)
}

// @Summary Answer order
// @Description Sends support answer to user email
// @Tags support
// @Accept multipart/form-data
// @Produce json
// @Security     BearerAuth
// @Security     CSRFToken
// @Param id path integer true "Order ID"
// @Param answer formData string true "Support answer"
// @Success 200
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /support/orders/{id}/answer [post]
func (h *OrderHandler) AnswerOrder(w http.ResponseWriter, r *http.Request) {
	op := "OrderHandler.AnswerOrder"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	adminID, err := utils.GetUserID(r.Context())
	if err != nil {
		log.Errorf("utils.GetUserID: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	orderID, err := utils.ParseIDFromRequest(r)
	if err != nil {
		log.Errorf("utils.ParseIDFromRequest: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	var req dto.OrderAnswerDTO

	err = utils.ParseMultipartFormData(r, &req)
	if err != nil {
		log.Errorf("utils.ParseMultipartFormData: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	err = h.uc.AnswerOrder(r.Context(), adminID, orderID, req.Answer)
	if err != nil {
		log.Errorf("h.uc.AnswerOrder: %s", err)
		utils.HandelError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
