package poster

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/request"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/response"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/poster"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/utils"
)

const defaultLimit = "12"
const defaultOffset = "0"

type PosterHandler struct {
	uc *poster.PosterUseCase
}

func NewPosterHandler(uc *poster.PosterUseCase) *PosterHandler {
	return &PosterHandler{
		uc: uc,
	}
}

// @Summary      Получить список объявлений
// @Description  Возвращает количество полученных объявлений и их список
// @Tags         posters
// @Produce      json
// @Param        limit   query     int  false  "Количество объявлений" default(12) minimum(1)
// @Param        offset  query     int  false  "Сдвиг для пагинации"   default(0)  minimum(0)
// @Success      200     {object}  response.PostersResponse
// @Failure 	 400 {object} response.ErrorResponse
// @Failure 	 401 {object} response.ErrorResponse
// @Failure 	 404 {object} response.ErrorResponse
// @Failure 	 500 {object} response.ErrorResponse
// @Router       /posters [get]
func (h *PosterHandler) GetPosters(w http.ResponseWriter, r *http.Request) {
	var request request.PostersRequest

	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	if limit == "" {
		limit = defaultLimit
	}

	if offset == "" {
		offset = defaultOffset
	}

	reqLimit, err := strconv.Atoi(limit)
	if err != nil {
		log.Printf("Atoi: %s", err)
		utils.HandelError(w, err)
		return
	}

	request.Limit = reqLimit

	reqOffset, err := strconv.Atoi(offset)
	if err != nil {
		log.Printf("Atoi: %s", err)
		utils.HandelError(w, err)
		return
	}

	request.Offset = reqOffset

	posters, err := h.uc.GetPostersUseCase(request)
	if err != nil {
		log.Printf("h.uc.GetPostersUseCase: %s", err)
		utils.HandelError(w, err)
		return
	}

	postersLen := len(posters)

	var response response.PostersResponse
	response.Len = postersLen
	response.Posters = posters

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
