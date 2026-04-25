package order

import (
	"net/http"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/response"
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

	order, err := h.uc.CreateOrder(r.Context(), &req)
	if err != nil {
		log.Errorf("h.uc.CreateOrder: %s", err)
		utils.HandelError(w, err)
		return
	}

	var response response.CreatedPosterResponse
	response.Poster = poster

	utils.JSONResponse(w, http.StatusCreated, response)
}
