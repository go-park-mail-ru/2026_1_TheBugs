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

func (h *PosterHandler) GetPosters(w http.ResponseWriter, r *http.Request) {
	var request request.PostersRequest
	var err error

	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	if limit == "" {
		limit = "12"
	}

	if offset == "" {
		offset = "0"
	}

	request.Limit, err = strconv.Atoi(limit)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	request.Offset, err = strconv.Atoi(offset)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	posters, err := h.uc.GetPostersUseCase(request.Limit, request.Offset)
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
