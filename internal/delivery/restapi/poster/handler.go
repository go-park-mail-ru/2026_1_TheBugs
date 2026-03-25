package poster

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/middleware"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/response"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/utils"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/poster"
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
	log := middleware.GetLogger(r.Context()).WithField("op", op)

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
		utils.WriteError(w, "bad limit query param", http.StatusBadRequest)
		return
	}

	params.Limit = reqLimit

	reqOffset, err := strconv.Atoi(offset)
	if err != nil {
		log.Errorf("Atoi: %s", err)
		utils.WriteError(w, "bad offset query param", http.StatusBadRequest)
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

	postersLen := len(posters)

	var response response.PostersResponse
	response.Len = postersLen
	response.Posters = posters

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
	log := middleware.GetLogger(r.Context()).WithField("op", op)

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

// @Summary Get metro stations by geo
// @Description Returns the metro stations
// @Tags posters
// @Produce json
// @Param lat query float32 true "lat"
// @Param lon query float32 true "lon"
// @Success 200 {object} response.MetroResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /posters/metro-stations [get]
func (h *PosterHandler) GetMetrosStation(w http.ResponseWriter, r *http.Request) {
	op := "PosterHandler.GetMetrosStation"
	log := middleware.GetLogger(r.Context()).WithField("op", op)

	latVal := r.URL.Query().Get("lat")
	lonVal := r.URL.Query().Get("lon")
	log.Infof("lat: %s, lon: %s", latVal, lonVal)
	lat, err := strconv.ParseFloat(latVal, 32)
	if err != nil {
		log.Errorf("Atoi: %s", err)
		utils.WriteError(w, "lat query param is requered", http.StatusBadRequest)
		return
	}
	lon, err := strconv.ParseFloat(lonVal, 32)
	if err != nil {
		log.Errorf("Atoi: %s", err)
		utils.WriteError(w, "lon query param is requered", http.StatusBadRequest)
		return
	}

	stations, err := h.uc.GetMetroStationsByRadius(r.Context(), dto.GeographyDTO{Lat: lat, Lon: lon})
	if err != nil {
		log.Errorf("h.uc.GetMetroStationsByRadius: %s", err)
		utils.HandelError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, response.MetroResponse{
		MetroStations: stations,
		Len:           len(stations),
	})
}
