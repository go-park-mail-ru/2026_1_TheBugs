package poster

import (
	"net/http"
	"strconv"

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
// @Description Returns filtered list of apartment posters
// @Tags posters
// @Produce json
// @Param limit query int false "Posters per page" default(12) minimum(1) maximum(100)
// @Param offset query int false "Pagination offset" default(0) minimum(0)
// @Param search_query query string false "Full-text search"
// @Param utility_company query string false "Utility company alias"
// @Param category query string false "Category alias"
// @Param room_count query int false "Exact room count"
// @Param min_price query int false "Min price"
// @Param max_price query int false "Max price"
// @Param facilities query []string false "Facilities aliases"
// @Param min_square query int false "Min area, sq.m"
// @Param max_square query int false "Max area, sq.m"
// @Param min_flat_floor query int false "Min floor"
// @Param max_flat_floor query int false "Max floor"
// @Param min_building_floor query int false "Min building floors"
// @Param max_building_floor query int false "Max building floors"
// @Param not_first_floor query boolean false "Exclude 1st floor" default(false)
// @Param not_last_floor query boolean false "Exclude last floor" default(false)
// @Success 200 {object} dto.PostersResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /posters/flats [get]
func (h *PosterHandler) GetFlatsAll(w http.ResponseWriter, r *http.Request) {
	op := "PosterHandler.GetFlatsAll"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	params, err := utils.ParsePostersFilters(r)
	if err != nil {
		log.Errorf("utils.ParsePostersFilters: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	posters, err := h.uc.SearchPostersUseCase(r.Context(), params)
	if err != nil {
		log.Errorf("h.uc.GetPostersUseCase: %s", err)
		utils.HandelError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, posters)

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

	poster, err := h.uc.GetPosterByAliasUseCase(r.Context(), alias, nil)
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
// @Router /posters/me/{alias} [get]
func (h *PosterHandler) GetPostersByUserByAlias(w http.ResponseWriter, r *http.Request) {
	op := "PosterHandler.GetPostersByUserByAlias"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	alias, err := utils.ParseAliasFromRequest(r)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, response.ErrorResponse{Error: "invalid alias"})
		return
	}

	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		log.Errorf("h.uc.GetPosterByAliasUseCase: %s", err)
		utils.HandelError(w, err)
		return
	}

	poster, err := h.uc.GetPosterByAliasUseCase(r.Context(), alias, &userID)
	if err != nil {
		log.Errorf("h.uc.GetPosterByAliasUseCase: %s", err)
		utils.HandelError(w, err)
		return
	}

	var response response.PosterResponse
	response.Poster = poster

	utils.JSONResponse(w, http.StatusOK, response)

}

// @Summary Create flat poster
// @Description Creates a flat poster with photos
// @Tags posters
// @Produce json
// @Security     BearerAuth
// @Security     CSRFToken
// @Param price formData number true "Poster price"
// @Param description formData string true "Poster description"
// @Param category_alias formData string true "Property category alias"
// @Param area formData number true "Property area"
// @Param address formData string true "Building address"
// @Param city formData string true "City Name"
// @Param district formData string false "District"
// @Param floor_count formData integer true "Building floor count"
// @Param company_id formData integer false "Company ID"
// @Param geo_lat formData number true "Latitude"
// @Param geo_lon formData number true "Longitude"
// @Param flat_category_id formData integer true "Flat category ID"
// @Param flat_number formData integer true "Flat number"
// @Param flat_floor formData integer true "Flat floor"
// @Param features formData []string false "Facilities aliases"
// @Param photos.0.file formData file false "First photo file"
// @Param photos.0.order formData integer false "First photo order"
// @Param photos.1.file formData file false "Second photo file"
// @Param photos.1.order formData integer false "Second photo order"
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

// @Summary Update flat poster
// @Description Updates a flat poster with photos
// @Tags posters
// @Produce json
// @Security     BearerAuth
// @Security     CSRFToken
// @Param alias path string true "Poster alias"
// @Param price formData number false "Poster price"
// @Param description formData string false "Poster description"
// @Param category_alias formData string false "Property category alias"
// @Param area formData number false "Property area"
// @Param address formData string false "Building address"
// @Param city formData string false "City"
// @Param district formData string false "District"
// @Param floor_count formData integer false "Building floor count"
// @Param company_id formData integer false "Company ID"
// @Param geo_lat formData number false "Latitude"
// @Param geo_lon formData number false "Longitude"
// @Param flat_category_id formData integer false "Flat category ID"
// @Param flat_number formData integer false "Flat number"
// @Param flat_floor formData integer false "Flat floor"
// @Param features formData []string false "Facilities aliases"
// @Param photos.0.file formData file false "First photo file"
// @Param photos.0.order formData integer false "First photo order"
// @Param photos.1.file formData file false "Second photo file"
// @Param photos.1.order formData integer false "Second photo order"
// @Success 200 {object} response.CreatedPosterResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /posters/flat/{alias} [put]
func (h *PosterHandler) UpdateFlatPoster(w http.ResponseWriter, r *http.Request) {
	op := "PosterHandler.UpdateFlatPoster"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		log.Errorf("utils.GetUserID: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	alias, err := utils.ParseAliasFromRequest(r)
	if err != nil {
		log.Errorf("utils.ParseIDFromRequest: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	var req dto.PosterInputFlatDTO
	req.UserID = userID

	err = utils.ParseMultipartFormData(r, &req)
	if err != nil {
		log.Errorf("utils.ParseMultipartFormData: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	req.Images, err = utils.ParsePhotos(r)
	if err != nil {
		log.Errorf("utils.ParsePhotos: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}
	log.Info(req.Images)

	poster, err := h.uc.UpdateFlatPoster(r.Context(), alias, &req)
	if err != nil {
		log.Errorf("h.uc.UpdateFlatPoster: %s", err)
		utils.HandelError(w, err)
		return
	}

	var response response.CreatedPosterResponse
	response.Poster = poster

	utils.JSONResponse(w, http.StatusCreated, response)
}

// @Summary Delete flat poster
// @Description Deletes a flat poster with all related data (photos, property, building)
// @Tags posters
// @Produce json
// @Security     BearerAuth
// @Security     CSRFToken
// @Param alias path string true "Poster alias"
// @Success 200 {object} response.CreatedPosterResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /posters/flat/{alias} [delete]
func (h *PosterHandler) DeleteFlatPoster(w http.ResponseWriter, r *http.Request) {
	op := "PosterHandler.DeleteFlatPoster"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		log.Errorf("utils.GetUserID: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}
	alias, err := utils.ParseAliasFromRequest(r)
	if err != nil {
		log.Errorf("utils.ParseAliasFromRequest: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	deletedPoster, err := h.uc.DeleteFlatPoster(r.Context(), alias, userID)
	if err != nil {
		log.Errorf("h.uc.DeleteFlatPoster: %s", err)
		utils.HandelError(w, err)
		return
	}

	var response response.CreatedPosterResponse
	response.Poster = deletedPoster

	utils.JSONResponse(w, http.StatusOK, response)
}

// @Summary Add poster to favorites
// @Description Adds a poster to user's favorites
// @Tags posters
// @Produce json
// @Security     BearerAuth
// @Security     CSRFToken
// @Param alias path string true "Poster alias"
// @Success 200
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /posters/{alias}/favorites [post]
func (h *PosterHandler) AddFavoritePoster(w http.ResponseWriter, r *http.Request) {
	op := "PosterHandler.AddFavoritePoster"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		log.Errorf("utils.GetUserID: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}
	alias, err := utils.ParseAliasFromRequest(r)
	if err != nil {
		log.Errorf("utils.ParseAliasFromRequest: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	err = h.uc.AddFavoritePoster(r.Context(), alias, userID)
	if err != nil {
		log.Errorf("h.uc.AddFavoritePoster: %s", err)
		utils.HandelError(w, err)
		return
	}
}

// @Summary Add poster view
// @Description Adds a view for a poster
// @Tags posters
// @Produce json
// @Security     BearerAuth
// @Security     CSRFToken
// @Param alias path string true "Poster alias"
// @Success 200
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /posters/{alias}/views [post]
func (h *PosterHandler) AddViewPoster(w http.ResponseWriter, r *http.Request) {
	op := "PosterHandler.AddViewPoster"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		log.Errorf("utils.GetUserID: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	alias, err := utils.ParseAliasFromRequest(r)
	if err != nil {
		log.Errorf("utils.ParseAliasFromRequest: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	err = h.uc.AddViewPoster(r.Context(), alias, userID)
	if err != nil {
		log.Errorf("h.uc.AddViewPoster: %s", err)
		utils.HandelError(w, err)
		return
	}
}

// @Summary Get poster views
// @Description Returns poster views count
// @Tags posters
// @Produce json
// @Param alias path string true "Poster alias"
// @Success 200 {object} response.PosterViewsResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /posters/{alias}/views [get]
func (h *PosterHandler) GetViewsPoster(w http.ResponseWriter, r *http.Request) {
	op := "PosterHandler.GetViewsPoster"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	alias, err := utils.ParseAliasFromRequest(r)
	if err != nil {
		log.Errorf("utils.ParseAliasFromRequest: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	views, err := h.uc.GetViewsPoster(r.Context(), alias)
	if err != nil {
		log.Errorf("h.uc.GetViewsPoster: %s", err)
		utils.HandelError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, views)
}

// @Summary Get favorite posters
// @Description Returns all favorite posters of the user
// @Tags posters
// @Produce json
// @Security     BearerAuth
// @Security     CSRFToken
// @Success 200 {object} dto.PostersResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /posters/favorites [get]
func (h *PosterHandler) GetFavoritesPoster(w http.ResponseWriter, r *http.Request) {
	op := "PosterHandler.GetFavoritesPoster"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		log.Errorf("utils.GetUserID: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	posters, err := h.uc.GetFavoritesPoster(r.Context(), userID)
	if err != nil {
		log.Errorf("h.uc.GetFavoritesPoster: %s", err)
		utils.HandelError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, posters)
}

// @Summary Remove poster from favorites
// @Description Deletes a poster from user's favorites
// @Tags posters
// @Produce json
// @Security     BearerAuth
// @Security     CSRFToken
// @Param alias path string true "Poster alias"
// @Success 200
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /posters/{alias}/favorites [delete]
func (h *PosterHandler) DeleteFavoritePoster(w http.ResponseWriter, r *http.Request) {
	op := "PosterHandler.DeleteFavoritePoster"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		log.Errorf("utils.GetUserID: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}
	alias, err := utils.ParseAliasFromRequest(r)
	if err != nil {
		log.Errorf("utils.ParseAliasFromRequest: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	err = h.uc.DeleteFavoritePoster(r.Context(), alias, userID)
	if err != nil {
		log.Errorf("h.uc.DeleteFavoritePoster: %s", err)
		utils.HandelError(w, err)
		return
	}
}

// @Summary Get list of posters JSONGeo
// @Description Returns filtered list of apartment posters on in JSONGeo notation
// @Tags posters
// @Produce json
// @Param sw_lat query float32 true "South West Lat"
// @Param sw_lon query float32 true "South West Lon"
// @Param ne_lat query float32 true "North East Lat"
// @Param ne_lon query float32 true "North East Lon"
// @Param zoom query int true "Map Zoom"
// @Param search_query query string false "Full-text search"
// @Param utility_company query string false "Utility company alias"
// @Param category query string false "Category alias"
// @Param room_count query int false "Exact room count"
// @Param min_price query int false "Min price"
// @Param max_price query int false "Max price"
// @Param facilities query []string false "Facilities aliases"
// @Param min_square query int false "Min area, sq.m"
// @Param max_square query int false "Max area, sq.m"
// @Param min_flat_floor query int false "Min floor"
// @Param max_flat_floor query int false "Max floor"
// @Param min_building_floor query int false "Min building floors"
// @Param max_building_floor query int false "Max building floors"
// @Param not_first_floor query boolean false "Exclude 1st floor" default(false)
// @Param not_last_floor query boolean false "Exclude last floor" default(false)
// @Success 200 {object} dto.GeoJSONFeatureResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /posters/geo [get]
func (h *PosterHandler) GetFlatsMapAll(w http.ResponseWriter, r *http.Request) {
	op := "PosterHandler.GetFlatsMapAll"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	bbox, err := utils.ParseMapFilters(r)
	if err != nil {
		log.Errorf("utils.ParseMapFilters: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}
	params, err := utils.ParsePostersFilters(r)
	if err != nil {
		log.Errorf("utils.ParsePostersFilters: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	posters, err := h.uc.GetPostersByCoords(r.Context(), bbox, params)
	if err != nil {
		log.Errorf("h.uc.GetPostersByCoord: %s", err)
		utils.HandelError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, posters)
}

// @Summary Get list of posters by point
// @Description Returns posters by point
// @Tags posters
// @Produce json
// @Param lat query float32 true "Lat"
// @Param lon query float32 true "Lon"
// @Success 200 {object} response.MyPostersResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /posters/by-point [get]
func (h *PosterHandler) GetPostersByPoint(w http.ResponseWriter, r *http.Request) {
	op := "PosterHandler.GetFlatsMapAll"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)
	q := r.URL.Query()
	lat, err := strconv.ParseFloat(q.Get("lat"), 32)
	if err != nil {
		log.Errorf("strconv.ParseFloat(q.Get(lat): %s", err)
		utils.HandelError(w, entity.InvalidInput)
	}
	lon, err := strconv.ParseFloat(q.Get("lon"), 32)
	if err != nil {
		log.Errorf("strconv.ParseFloat(q.Get(lon): %s", err)
		utils.HandelError(w, entity.InvalidInput)
	}

	posters, err := h.uc.GetPostersByRadius(r.Context(), dto.GeographyDTO{Lat: lat, Lon: lon})
	if err != nil {
		log.Errorf("h.uc.GetPostersByRadius: %s", err)
		utils.HandelError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, response.MyPostersResponse{Posters: posters, Len: len(posters)})
}
