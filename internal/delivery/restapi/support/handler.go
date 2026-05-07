package support

import (
	"net/http"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/utils"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/support"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/ctxLogger"
)

type SupportHandler struct {
	uc *support.SupportUseCase
}

func NewSupportHandler(uc *support.SupportUseCase) *SupportHandler {
	return &SupportHandler{
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
func (h *SupportHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	op := "SupportHandler.CreateOrder"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	var req dto.OrderDTO

	err := utils.ParseMultipartFormData(r, &req)
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
		log.Errorf("h.uc.CreateSupport: %s", err)
		utils.HandelError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
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
func (h *SupportHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	op := "SupportHandler.GetOrders"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		log.Errorf("utils.GetUserID: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	supports, err := h.uc.GetUserOrders(r.Context(), userID)
	if err != nil {
		log.Errorf("h.uc.GetUserOrders: %s", err)
		utils.HandelError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, supports)
}

// @Summary Get order by id
// @Description Returns full order info by id
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
func (h *SupportHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	op := "SupportHandler.GetSupportByID"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		log.Errorf("utils.GetUserID: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	supportID, err := utils.ParseIDFromRequest(r)
	if err != nil {
		log.Errorf("strconv.Atoi: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	support, err := h.uc.GetOrderByID(r.Context(), userID, supportID)
	if err != nil {
		log.Errorf("utils.ParseIDFromRequest: %s", err)
		utils.HandelError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, support)
}

// @Summary Answer order
// @Description Sends order answer to user email
// @Tags order
// @Accept multipart/form-data
// @Produce json
// @Security     BearerAuth
// @Security     CSRFToken
// @Param id path integer true "Order ID"
// @Param answer formData string true "Order answer"
// @Success 200
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /support/orders/{id}/answer [post]
func (h *SupportHandler) AnswerOrder(w http.ResponseWriter, r *http.Request) {
	op := "SupportHandler.AnswerOrder"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	adminID, err := utils.GetUserID(r.Context())
	if err != nil {
		log.Errorf("utils.GetUserID: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	supportID, err := utils.ParseIDFromRequest(r)
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

	err = h.uc.AnswerOrder(r.Context(), adminID, supportID, req.Answer)
	if err != nil {
		log.Errorf("h.uc.AnswerSupport: %s", err)
		utils.HandelError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
