package poster

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/middleware"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/response"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/utils"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
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
// @Router /posters/flats [get]
func (h *PosterHandler) GetFlatsAll(w http.ResponseWriter, r *http.Request) {
	op := "PosterHandler.GetFlatsAll"
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
	response := response.MyPostersResponse{
		Posters: posters,
		Len:     len(posters),
	}
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

// @Summary Create flat poster
// @Description Creates a flat poster with photos
// @Tags posters
// @Produce json
// @Param price formData number true "Poster price"
// @Param description formData string true "Poster description"
// @Param category_id formData integer true "Property category ID"
// @Param area formData number true "Property area"
// @Param address formData string true "Building address"
// @Param city_id formData integer true "City ID"
// @Param metro_station_id formData integer false "Metro station ID"
// @Param district formData string false "District"
// @Param floor_count formData integer true "Building floor count"
// @Param company_id formData integer false "Company ID"
// @Param geo.lat formData number true "Latitude"
// @Param geo.lon formData number true "Longitude"
// @Param flat.category_id formData integer true "Flat category ID"
// @Param flat.flat_number formData integer true "Flat number"
// @Param flat.floor formData integer true "Flat floor"
// @Param photos.0.file formData file false "First photo file"
// @Param photos.0.order formData integer false "First photo order"
// @Success 201 {object} response.CreatedPosterResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /posters/flat [post]
func (h *PosterHandler) CreateFlatPoster(w http.ResponseWriter, r *http.Request) {
	op := "PosterHandler.CreateFlatPoster"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		log.Errorf("utils.GetUserID: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	var req dto.PosterInputFlatDTO

	req.UserID = userID
	err = utils.ParseFormData(r, &req)
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

	poster, err := h.uc.CreateFlatPoster(r.Context(), &req)
	if err != nil {
		log.Errorf("h.uc.CreateFlatPoster: %s", err)
		utils.HandelError(w, err)
		return
	}

	var response response.CreatedPosterResponse
	response.Poster = poster

	utils.JSONResponse(w, http.StatusCreated, response)
}
