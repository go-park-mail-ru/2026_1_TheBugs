package poster

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/generated/poster"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/response"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/utils"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/ctxLogger"
	"github.com/samber/lo"
	"google.golang.org/grpc"
)

const defaultLimit = "12"
const defaultOffset = "0"

//grpc.ClientStreamingClient[poster.UpdateFlatPosterRequest, poster.UpdateFlatPosterResponse]
//grpc.ClientStreamingClient[poster.CreateFlatPosterRequest, poster.CreateFlatPosterResponse]

type PosterHandler struct {
	grpcClient poster.PosterServiceClient
}

func NewPosterHandler(grpcConn *grpc.ClientConn) *PosterHandler {
	return &PosterHandler{
		grpcClient: poster.NewPosterServiceClient(grpcConn),
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

	req := utils.MapFiltersDTOToProto(params)

	resp, err := h.grpcClient.SearchPosters(r.Context(), req)
	if err != nil {
		log.Errorf("h.grpcClient.SearchPosters: %s", err)
		utils.HandelGRPCError(w, err)
		return
	}

	result := utils.MapSearchPostersResponseToDTO(resp)

	utils.JSONResponse(w, http.StatusOK, result)
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

	resp, err := h.grpcClient.GetPosterByAlias(r.Context(), &poster.GetPosterByAliasRequest{
		PosterAlias: alias,
	})
	if err != nil {
		log.Errorf("h.client.GetPosterByAlias: %s", err)
		utils.HandelGRPCError(w, err)
		return
	}

	var response response.PosterResponse
	response.Poster = utils.MapProtoPosterToDTO(resp.Poster)

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
		log.Errorf("utils.GetUserID: %s", err)
		utils.HandelError(w, err)
		return
	}

	resp, err := h.grpcClient.GetPostersByUserID(r.Context(), &poster.GetPostersByUserIDRequest{
		UserId: int64(userID),
	})
	if err != nil {
		log.Errorf("h.grpcClient.GetPostersByUserID: %s", err)
		utils.HandelGRPCError(w, err)
		return
	}

	posters := utils.MapProtoMyPostersToDTO(resp.Posters)

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
		log.Errorf("utils.GetUserID: %s", err)
		utils.HandelError(w, err)
		return
	}

	resp, err := h.grpcClient.GetPosterByAlias(r.Context(), &poster.GetPosterByAliasRequest{
		PosterAlias: alias,
		UserId:      lo.ToPtr(int64(userID)),
	})
	if err != nil {
		log.Errorf("h.grpcClient.GetPosterByAlias: %s", err)
		utils.HandelGRPCError(w, err)
		return
	}

	var response response.PosterResponse
	response.Poster = utils.MapProtoPosterToDTO(resp.Poster)

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

	stream, err := h.grpcClient.CreateFlatPoster(r.Context())
	if err != nil {
		log.Errorf("h.grpcClient.CreateFlatPoster: %s", err)
		utils.HandelGRPCError(w, err)
		return
	}

	err = utils.SendFlatPosterMeta(stream, &req)
	if err != nil {
		log.Errorf("sendFlatPosterMeta: %s", err)
		utils.HandelError(w, err)
		return
	}

	err = utils.SendFlatPosterPhotos(stream, req.Images)
	if err != nil {
		log.Errorf("sendFlatPosterPhotos: %s", err)
		utils.HandelError(w, err)
		return
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Errorf("stream.CloseAndRecv: %s", err)
		utils.HandelError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusCreated, response.CreatedPosterResponse{
		Poster: &dto.CreatedPoster{
			ID:    int(resp.Id),
			Alias: resp.Alias,
		},
	})
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
		log.Errorf("utils.ParseAliasFromRequest: %s", err)
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

	stream, err := h.grpcClient.UpdateFlatPoster(r.Context())
	if err != nil {
		log.Errorf("h.grpcClient.UpdateFlatPoster: %s", err)
		utils.HandelError(w, err)
		return
	}

	err = utils.SendUpdateFlatPosterMeta(stream, alias, &req)
	if err != nil {
		log.Errorf("utils.SendUpdateFlatPosterMeta: %s", err)
		utils.HandelError(w, err)
		return
	}

	err = utils.SendUpdateFlatPosterPhotos(stream, req.Images)
	if err != nil {
		log.Errorf("utils.SendUpdateFlatPosterPhotos: %s", err)
		utils.HandelError(w, err)
		return
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Errorf("stream.CloseAndRecv: %s", err)
		utils.HandelGRPCError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, response.CreatedPosterResponse{
		Poster: &dto.CreatedPoster{
			ID:    int(resp.Id),
			Alias: resp.Alias,
		},
	})
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

	resp, err := h.grpcClient.DeleteFlatPoster(r.Context(), &poster.DeleteFlatPosterRequest{
		Alias:  alias,
		UserId: int64(userID),
	})
	if err != nil {
		log.Errorf("h.grpcClient.DeleteFlatPoster: %s", err)
		utils.HandelGRPCError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, response.CreatedPosterResponse{
		Poster: &dto.CreatedPoster{
			ID:    int(resp.Id),
			Alias: resp.Alias,
		},
	})
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

	_, err = h.grpcClient.AddViewPoster(r.Context(), &poster.AddViewPosterRequest{
		Alias:  alias,
		UserId: int64(userID),
	})
	if err != nil {
		log.Errorf("h.grpcClient.AddViewPoster: %s", err)
		utils.HandelGRPCError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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

	resp, err := h.grpcClient.GetViewsPoster(r.Context(), &poster.GetViewsPosterRequest{
		Alias: alias,
	})
	if err != nil {
		log.Errorf("h.grpcClient.GetViewsPoster: %s", err)
		utils.HandelGRPCError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, response.PosterViewsResponse{
		Views: int(resp.Views),
	})
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

	_, err = h.grpcClient.AddFavoritePoster(r.Context(), &poster.AddFavoritePosterRequest{
		Alias:  alias,
		UserId: int64(userID),
	})
	if err != nil {
		log.Errorf("h.grpcClient.AddFavoritePoster: %s", err)
		utils.HandelGRPCError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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

	resp, err := h.grpcClient.GetFavoritePosters(r.Context(), &poster.GetFavoritePostersRequest{
		UserId: int64(userID),
	})
	if err != nil {
		log.Errorf("h.grpcClient.GetFavoritePosters: %s", err)
		utils.HandelGRPCError(w, err)
		return
	}

	result := utils.MapSearchPostersResponseToDTO(&poster.SearchPostersResponse{
		Len:     resp.Len,
		Posters: resp.Posters,
	})

	utils.JSONResponse(w, http.StatusOK, result)
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

	_, err = h.grpcClient.DeleteFavoritePoster(r.Context(), &poster.DeleteFavoritePosterRequest{
		Alias:  alias,
		UserId: int64(userID),
	})
	if err != nil {
		log.Errorf("h.grpcClient.DeleteFavoritePoster: %s", err)
		utils.HandelGRPCError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Get poster favorites count
// @Description Returns number of users who added poster to favorites
// @Tags posters
// @Produce json
// @Security     BearerAuth
// @Security     CSRFToken
// @Param alias path string true "Poster alias"
// @Success 200 {object} map[string]int
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /posters/{alias}/favorites [get]
func (h *PosterHandler) GetFavoritesCountPoster(w http.ResponseWriter, r *http.Request) {
	op := "PosterHandler.GetFavoritesCountPoster"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	var userID *int64

	id, err := utils.GetUserID(r.Context())
	if err != nil {
		log.Infof("utils.GetUserID: %s", err)
	} else {
		userID = lo.ToPtr(int64(id))
	}

	alias, err := utils.ParseAliasFromRequest(r)
	if err != nil {
		log.Errorf("utils.ParseAliasFromRequest: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	resp, err := h.grpcClient.GetFavoritesCountPoster(r.Context(), &poster.GetFavoritesCountPosterRequest{
		PosterAlias: alias,
		UserId:      userID,
	})
	if err != nil {
		log.Errorf("h.grpcClient.GetFavoritesCountPoster: %s", err)
		utils.HandelGRPCError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, response.PosterFavoritesCountResponse{
		Count:      int(resp.Count),
		IsFavorite: resp.IsFavorite,
	})
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

	resp, err := h.grpcClient.GetPostersByCoords(r.Context(), &poster.GetPostersByCoordsRequest{
		Bounds:  utils.MapMapBoundsDTOToProto(bbox),
		Filters: utils.MapFiltersDTOToProto(params),
	})
	if err != nil {
		log.Errorf("h.grpcClient.GetPostersByCoords: %s", err)
		utils.HandelGRPCError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, utils.MapGeoJSONResponseToDTO(resp))
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
	op := "PosterHandler.GetPostersByPoint"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	q := r.URL.Query()

	lat, err := strconv.ParseFloat(q.Get("lat"), 64)
	if err != nil {
		log.Errorf("strconv.ParseFloat(q.Get(lat)): %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	lon, err := strconv.ParseFloat(q.Get("lon"), 64)
	if err != nil {
		log.Errorf("strconv.ParseFloat(q.Get(lon)): %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	resp, err := h.grpcClient.GetPostersByRadius(r.Context(), &poster.GetPostersByRadiusRequest{
		Point: &poster.Geography{
			Lat: lat,
			Lon: lon,
		},
	})
	if err != nil {
		log.Errorf("h.grpcClient.GetPostersByRadius: %s", err)
		utils.HandelGRPCError(w, err)
		return
	}

	posters := utils.MapProtoMyPostersToDTO(resp.Posters)

	utils.JSONResponse(w, http.StatusOK, response.MyPostersResponse{
		Posters: posters,
		Len:     len(posters),
	})
}

// @Summary Generate poster description
// @Description Returns description for poster by given parameters
// @Tags posters
// @Produce json
// @Security     BearerAuth
// @Security     CSRFToken
// @Param input body dto.GenerateDescriptionDTO true "Poster parameters"
// @Success 200 {object} response.GenerateDescriptionResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /posters/generate-description [post]
func (h *PosterHandler) GenerateDescription(w http.ResponseWriter, r *http.Request) {
	op := "PosterHandler.GenerateDescription"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	var req dto.GenerateDescriptionDTO
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Errorf("json.NewDecoder.Decode: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	resp, err := h.grpcClient.GenerateDescription(r.Context(), &poster.GenerateDescriptionRequest{
		Category:     req.Category,
		Area:         req.Area,
		FlatCategory: req.FlatCategory,
		City:         req.City,
		Features:     req.Features,
	})
	if err != nil {
		log.Errorf("h.grpcClient.GenerateDescription: %s", err)
		utils.HandelGRPCError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, response.GenerateDescriptionResponse{
		Description: resp.Description,
	})
}

// @Summary Get poster price history
// @Description Returns poster price history
// @Tags posters
// @Produce json
// @Param alias path string true "Poster alias"
// @Success 200 {object} response.PriceHistoryResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /posters/{alias}/price-history [get]
func (h *PosterHandler) GetPriceHistoryPoster(w http.ResponseWriter, r *http.Request) {
	op := "PosterHandler.GetPriceHistoryPoster"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	alias, err := utils.ParseAliasFromRequest(r)
	if err != nil {
		log.Errorf("utils.ParseAliasFromRequest: %s", err)
		utils.HandelError(w, entity.InvalidInput)
		return
	}

	resp, err := h.grpcClient.GetPriceHistoryPoster(r.Context(), &poster.GetPriceHistoryPosterRequest{
		PosterAlias: alias,
	})
	if err != nil {
		log.Errorf("h.grpcClient.GetPriceHistoryPoster: %s", err)
		utils.HandelGRPCError(w, err)
		return
	}

	history := utils.MapProtoPriceHistoryToDTO(resp.History)

	utils.JSONResponse(w, http.StatusOK, response.PriceHistoryResponse{
		History: history,
		Count:   len(history),
	})
}

// @Summary Get poster roommates
// @Description Returns users who want to become roommates for poster
// @Tags posters
// @Produce json
// @Param alias path string true "Alias of poster"
// @Success 200 {object} dto.PosterRoommatesResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /posters/{alias}/roommates [get]
func (h *PosterHandler) GetPosterRoommates(w http.ResponseWriter, r *http.Request) {
	op := "PosterHandler.GetPosterRoommates"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	alias, err := utils.ParseAliasFromRequest(r)
	if err != nil {
		utils.JSONResponse(w, http.StatusBadRequest, response.ErrorResponse{Error: "invalid alias"})
		return
	}

	resp, err := h.grpcClient.GetPosterRoommates(r.Context(), &poster.GetPosterRoommatesRequest{
		Alias: alias,
	})
	if err != nil {
		log.Errorf("h.grpcClient.GetPosterRoommates: %s", err)
		utils.HandelGRPCError(w, err)
		return
	}

	users := make([]dto.RoommateUserDTO, 0, len(resp.Users))
	for _, user := range resp.Users {
		users = append(users, dto.RoommateUserDTO{
			ID:        int(user.Id),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			AvatarURL: user.AvatarUrl,
		})
	}

	utils.JSONResponse(w, http.StatusOK, dto.PosterRoommatesResponse{
		Users: users,
		Len:   int(resp.Total),
	})
}

// @Summary Add poster roommate
// @Description Adds current authorized user to poster roommates
// @Tags posters
// @Produce json
// @Param alias path string true "Alias of poster"
// @Success 204
// @Security     BearerAuth
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /posters/{alias}/roommates [post]
func (h *PosterHandler) AddPosterRoommate(w http.ResponseWriter, r *http.Request) {
	op := "PosterHandler.AddPosterRoommate"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		log.Errorf("utils.GetUserID: %s", err)
		utils.JSONResponse(w, http.StatusUnauthorized, response.ErrorResponse{Error: "unauthorized"})
		return
	}

	alias, err := utils.ParseAliasFromRequest(r)
	if err != nil {
		log.Errorf("utils.ParseAliasFromRequest: %s", err)
		utils.JSONResponse(w, http.StatusBadRequest, response.ErrorResponse{Error: "invalid alias"})
		return
	}

	_, err = h.grpcClient.AddPosterRoommate(r.Context(), &poster.AddPosterRoommateRequest{
		Alias:  alias,
		UserId: int64(userID),
	})
	if err != nil {
		log.Errorf("h.grpcClient.AddPosterRoommate: %s", err)
		utils.HandelGRPCError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Delete poster roommate
// @Description Removes current authorized user from poster roommates
// @Tags posters
// @Produce json
// @Param alias path string true "Alias of poster"
// @Success 204
// @Security BearerAuth
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /posters/{alias}/roommates [delete]
func (h *PosterHandler) DeletePosterRoommate(w http.ResponseWriter, r *http.Request) {
	op := "PosterHandler.DeletePosterRoommate"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		log.Errorf("utils.GetUserID: %s", err)
		utils.JSONResponse(w, http.StatusUnauthorized, response.ErrorResponse{Error: "unauthorized"})
		return
	}

	alias, err := utils.ParseAliasFromRequest(r)
	if err != nil {
		log.Errorf("utils.ParseAliasFromRequest: %s", err)
		utils.JSONResponse(w, http.StatusBadRequest, response.ErrorResponse{Error: "invalid alias"})
		return
	}

	_, err = h.grpcClient.DeletePosterRoommate(r.Context(), &poster.DeletePosterRoommateRequest{
		Alias:  alias,
		UserId: int64(userID),
	})
	if err != nil {
		log.Errorf("h.grpcClient.DeletePosterRoommate: %s", err)
		utils.HandelGRPCError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
