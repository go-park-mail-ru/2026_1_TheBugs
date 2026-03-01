package poster

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/poster/request"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/poster/response"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/poster"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/utils"
)

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
// @Router       /posters [get]
func (h *PosterHandler) GetPosters(w http.ResponseWriter, r *http.Request) {
	var request request.PostersRequest

	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	if limit == "" {
		limit = "12"
	}

	if offset == "" {
		offset = "0"
	}

	reqLimit, err := strconv.Atoi(limit)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	request.Limit = reqLimit

	reqOffset, err := strconv.Atoi(offset)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
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
