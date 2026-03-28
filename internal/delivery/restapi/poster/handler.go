package poster

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/response"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/utils"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/poster"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/ctxLogger"
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

// @Summary Get list of posters
// @Description Returns the number of retrieved posters and their list
// @Tags posters
// @Produce json
// @Param limit query int false "Number of posters" default(12) minimum(1)
// @Param offset query int false "Pagination offset" default(0) minimum(0)
// @Param utility_company query string false "Utility company"
// @Success 200 {object} response.PostersResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /posters [get]
func (h *PosterHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	op := "PosterHandler.GetAll"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	var params dto.PostersFiltersDTO

	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")
	utilityCompany := r.URL.Query().Get("utility_company")

	if limit == "" {
		limit = defaultLimit
	}

	if offset == "" {
		offset = defaultOffset
	}

	reqLimit, err := strconv.Atoi(limit)
	if err != nil {
		log.Errorf("Atoi: %s", err)
		utils.HandelError(w, err)
		return
	}

	params.Limit = reqLimit

	reqOffset, err := strconv.Atoi(offset)
	if err != nil {
		log.Errorf("Atoi: %s", err)
		utils.HandelError(w, err)
		return
	}

	params.Offset = reqOffset
	if utilityCompany != "" {
		params.UtilityCompany = &utilityCompany
	}

	posters, err := h.uc.GetPostersUseCase(r.Context(), params)
	if err != nil {
		log.Errorf("h.uc.GetPostersUseCase: %s", err)
		utils.HandelError(w, err)
		return
	}

	response := response.PostersResponse{
		Posters: posters,
		Len:     len(posters),
	}

	utils.JSONResponse(w, http.StatusOK, response)

}

// @Summary Get poster by alias
// @Description Returns the poster with all information about it
// @Tags posters
// @Produce json
// @Param alias path string true "Alias of poster"
// @Success 200 {object} response.PosterResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /posters/by-alias/{alias} [get]
func (h *PosterHandler) GetPoster(w http.ResponseWriter, r *http.Request) {
	op := "PosterHandler.GetPoster"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	alias, err := utils.ParseAliasFromRequest(r)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, response.ErrorResponse{Error: "invalid alias"})
		return
	}

	poster, err := h.uc.GetPosterByAliasUseCase(r.Context(), alias)
	if err != nil {
		log.Errorf("h.uc.GetPosterByAliasUseCase: %s", err)
		utils.HandelError(w, err)
		return
	}

	var response response.PosterResponse
	response.Poster = poster

	utils.JSONResponse(w, http.StatusOK, response)
}

// @Summary Get list of user`s posters
// @Description Returns the number of retrieved posters and their list
// @Tags posters
// @Produce json
// @Security     BearerAuth
// @Success 200 {object} response.PosterResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /posters/me [get]
func (h *PosterHandler) GetPostersByUser(w http.ResponseWriter, r *http.Request) {
	op := "PosterHandler.GetPostersByUser"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		log.Errorf("h.uc.GetPosterByAliasUseCase: %s", err)
		utils.HandelError(w, err)
		return
	}

	posters, err := h.uc.GetPosterByUserID(r.Context(), userID)
	if err != nil {
		log.Errorf("h.uc.GetPosterByUserID: %s", err)
		utils.HandelError(w, err)
		return
	}
	response := response.PostersResponse{
		Posters: posters,
		Len:     len(posters),
	}
	utils.JSONResponse(w, http.StatusOK, response)

}
