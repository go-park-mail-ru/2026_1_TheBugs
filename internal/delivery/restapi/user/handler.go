package user

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/generated/user"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/utils"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/ctxLogger"
	"google.golang.org/grpc"
)

type UserHandler struct {
	grpcClient user.UserServiceClient
}

func NewUserHandler(grpcConn *grpc.ClientConn) *UserHandler {
	return &UserHandler{
		grpcClient: user.NewUserServiceClient(grpcConn),
	}
}

// @Summary Get user details
// @Description Returns user details
// @Tags users
// @Produce json
// @Security     BearerAuth
// @Success 200 {object} dto.UserDTO
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /user/me [get]
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	op := "UserHandler.GetMe"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		log.WithError(err).Error("failed to get user ID from context")
		utils.WriteError(w, "failed to get user ID from context", http.StatusUnauthorized)
		return
	}

	user, err := h.grpcClient.GetMe(r.Context(), &user.GetMeRequest{
		UserId: int32(userID),
	})
	if err != nil {
		log.WithError(err).Error("h.grpcClient.GetMe: failed to get user by ID")
		utils.HandelGRPCError(w, err)
		return
	}
	utils.JSONResponse(w, http.StatusOK, user)

}

// @Summary Update user details
// @Description Returns user details
// @Tags users
// @Produce json
// @Param first_name formData string false "New Firstname"
// @Param last_name formData string false "New Lastname"
// @Param phone formData string false "New Phone"
// @Param avatar formData file false "Avatar file"
// @Security     BearerAuth
// @Success 200 {object} dto.UserDTO
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /user/me/profile [put]
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	op := "UserHandler.UpdateProfile"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	var data dto.UpdateProfileRequest
	var file *user.UploadFile

	err := utils.ParseMultipartFormData(r, &data)
	if err != nil {
		log.WithError(err).Error("failed to decode update profile data")
		utils.WriteError(w, "invalid body", http.StatusBadRequest)
		return
	}
	files, ok := r.MultipartForm.File["avatar"]
	if ok && len(files) > 0 {
		fileInput, err := utils.ParseFileInput(files[0])
		if err != nil {
			log.WithError(err).Error("failed to parse avatar file")
			utils.WriteError(w, "failed to read avatar", http.StatusBadRequest)
			return
		}

		data.Avatar = fileInput
		raw, err := io.ReadAll(data.Avatar.File)
		if err != nil {
			log.WithError(err).Error("failed to read avatar file")
			utils.WriteError(w, "failed to read avatar", http.StatusBadRequest)
			return
		}
		file = &user.UploadFile{
			Filename:    data.Avatar.Filename,
			ContentType: data.Avatar.ContentType,
			Size:        data.Avatar.Size,
			Avatar:      raw,
		}
	}
	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		log.WithError(err).Error("failed to get user ID from context")
		utils.WriteError(w, "failed to get user ID", http.StatusUnauthorized)
		return
	}
	data.ID = userID

	user, err := h.grpcClient.UpdateProfile(r.Context(), &user.UpdateProfileRequest{
		Id:        int32(userID),
		Firstname: data.FirstName,
		Lastname:  data.LastName,
		Phone:     data.Phone,
		File:      file,
	})
	if err != nil {
		log.WithError(err).Error("h.grpcClient.UpdateProfile: failed to update user profile")
		utils.HandelGRPCError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, user)
}

// @Summary Get roommate user details
// @Description Returns user details with roommate form
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} user.GetRoommateUserResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /user/{id} [get]
func (h *UserHandler) GetRoommateUser(w http.ResponseWriter, r *http.Request) {
	op := "UserHandler.GetRoommateUser"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	userID, err := utils.ParseIDFromRequest(r)
	if err != nil {
		log.WithError(err).Error("failed to parse user ID from request")
		utils.WriteError(w, "invalid user id", http.StatusBadRequest)
		return
	}

	roommateUser, err := h.grpcClient.GetRoommateUser(r.Context(), &user.GetRoommateUserRequest{
		UserId: int64(userID),
	})
	if err != nil {
		log.WithError(err).Error("h.grpcClient.GetRoommateUser: failed to get roommate user")
		utils.HandelGRPCError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, roommateUser)
}

// @Summary Add roommate match
// @Description Adds current authorized user sympathy to target user. Optional poster_alias links match to poster
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "Target user ID"
// @Param input body AddRoommateMatchRequest false "Optional poster alias"
// @Security BearerAuth
// @Success 204
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /user/{id}/match [post]
func (h *UserHandler) AddRoommateMatch(w http.ResponseWriter, r *http.Request) {
	op := "UserHandler.AddRoommateMatch"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	fromUserID, err := utils.GetUserID(r.Context())
	if err != nil {
		log.WithError(err).Error("failed to get user ID from context")
		utils.WriteError(w, "failed to get user ID", http.StatusUnauthorized)
		return
	}

	toUserID, err := utils.ParseIDFromRequest(r)
	if err != nil {
		log.WithError(err).Error("failed to parse target user ID from request")
		utils.WriteError(w, "invalid user id", http.StatusBadRequest)
		return
	}

	var req dto.AddRoommateMatchRequest
	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err != io.EOF {
			log.WithError(err).Error("failed to decode request body")
			utils.WriteError(w, "invalid request body", http.StatusBadRequest)
			return
		}
	}

	_, err = h.grpcClient.AddRoommateMatch(r.Context(), &user.AddRoommateMatchRequest{
		FromUserId:  int64(fromUserID),
		ToUserId:    int64(toUserID),
		PosterAlias: req.PosterAlias,
	})
	if err != nil {
		log.WithError(err).Error("h.grpcClient.AddRoommateMatch: failed to add roommate match")
		utils.HandelGRPCError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Delete roommate match
// @Description Removes target user from current authorized user roommates
// @Tags users
// @Produce json
// @Param id path int true "Target user ID"
// @Security BearerAuth
// @Success 204
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /user/{id}/match [delete]
func (h *UserHandler) DeleteRoommateMatch(w http.ResponseWriter, r *http.Request) {
	op := "UserHandler.DeleteRoommateMatch"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	fromUserID, err := utils.GetUserID(r.Context())
	if err != nil {
		log.WithError(err).Error("failed to get user ID from context")
		utils.WriteError(w, "failed to get user ID", http.StatusUnauthorized)
		return
	}

	toUserID, err := utils.ParseIDFromRequest(r)
	if err != nil {
		log.WithError(err).Error("failed to parse target user ID from request")
		utils.WriteError(w, "invalid user id", http.StatusBadRequest)
		return
	}

	_, err = h.grpcClient.DeleteRoommateMatch(r.Context(), &user.DeleteRoommateMatchRequest{
		FromUserId: int64(fromUserID),
		ToUserId:   int64(toUserID),
	})
	if err != nil {
		log.WithError(err).Error("h.grpcClient.DeleteRoommateMatch")
		utils.HandelGRPCError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Get roommate contacts
// @Description Returns target user contacts if roommate match is mutual
// @Tags users
// @Produce json
// @Param id path int true "Target user ID"
// @Security BearerAuth
// @Success 200 {object} user.GetRoommateContactsResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /user/{id}/contacts [get]
func (h *UserHandler) GetRoommateContacts(w http.ResponseWriter, r *http.Request) {
	op := "UserHandler.GetRoommateContacts"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	fromUserID, err := utils.GetUserID(r.Context())
	if err != nil {
		log.WithError(err).Error("failed to get user ID from context")
		utils.WriteError(w, "failed to get user ID", http.StatusUnauthorized)
		return
	}

	toUserID, err := utils.ParseIDFromRequest(r)
	if err != nil {
		log.WithError(err).Error("failed to parse target user ID from request")
		utils.WriteError(w, "invalid user id", http.StatusBadRequest)
		return
	}

	contacts, err := h.grpcClient.GetRoommateContacts(r.Context(), &user.GetRoommateContactsRequest{
		FromUserId: int64(fromUserID),
		ToUserId:   int64(toUserID),
	})
	if err != nil {
		log.WithError(err).Error("h.grpcClient.GetRoommateContacts: failed to get roommate contacts")
		utils.HandelGRPCError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, contacts)
}

// @Summary Create roommate form
// @Description Creates roommate form for current authorized user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body dto.CreateRoommateFormRequest true "Roommate form"
// @Success 204
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /user/me/roommate-form [post]
func (h *UserHandler) CreateRoommateForm(w http.ResponseWriter, r *http.Request) {
	op := "UserHandler.CreateRoommateForm"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		log.WithError(err).Error("failed to get user ID from context")
		utils.WriteError(w, "failed to get user ID", http.StatusUnauthorized)
		return
	}

	var data dto.CreateRoommateFormRequest
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.WithError(err).Error("failed to decode roommate form")
		utils.WriteError(w, "invalid body", http.StatusBadRequest)
		return
	}

	_, err = h.grpcClient.CreateRoommateForm(r.Context(), &user.CreateRoommateFormRequest{
		UserId:      int64(userID),
		Gender:      data.Gender,
		Birthday:    data.Birthday,
		Description: data.Description,
		Tags:        data.Tags,
	})
	if err != nil {
		log.WithError(err).Error("h.grpcClient.CreateRoommateForm: failed to create roommate form")
		utils.HandelGRPCError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Get roommate form
// @Description Returns current user roommate form
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.RoommateFormDTO
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /user/me/roommate-form [get]
func (h *UserHandler) GetRoommateForm(w http.ResponseWriter, r *http.Request) {
	op := "UserHandler.GetRoommateForm"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		log.WithError(err).Error("failed to get user ID")
		utils.WriteError(w, "failed to get user ID", http.StatusUnauthorized)
		return
	}

	resp, err := h.grpcClient.GetRoommateForm(r.Context(), &user.GetRoommateFormRequest{
		UserId: int64(userID),
	})
	if err != nil {
		log.WithError(err).Error("h.grpcClient.GetRoommateForm")
		utils.HandelGRPCError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, resp)
}

// @Summary Update roommate form
// @Description Updates roommate form for current authorized user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body dto.RoommateFormDTO true "Roommate form"
// @Success 204
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /user/me/roommate-form [put]
func (h *UserHandler) UpdateRoommateForm(w http.ResponseWriter, r *http.Request) {
	op := "UserHandler.UpdateRoommateForm"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		log.WithError(err).Error("failed to get user ID from context")
		utils.WriteError(w, "failed to get user ID", http.StatusUnauthorized)
		return
	}

	var data dto.RoommateFormDTO
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.WithError(err).Error("failed to decode roommate form")
		utils.WriteError(w, "invalid body", http.StatusBadRequest)
		return
	}

	_, err = h.grpcClient.UpdateRoommateForm(r.Context(), &user.UpdateRoommateFormRequest{
		UserId:      int64(userID),
		Gender:      data.Gender,
		Birthday:    data.Birthday,
		Description: data.Description,
		Tags:        data.Tags,
	})
	if err != nil {
		log.WithError(err).Error("h.grpcClient.UpdateRoommateForm: failed to update roommate form")
		utils.HandelGRPCError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Delete roommate form
// @Description Deletes current user roommate form
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 204
// @Failure 401 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /user/me/roommate-form [delete]
func (h *UserHandler) DeleteRoommateForm(w http.ResponseWriter, r *http.Request) {
	op := "UserHandler.DeleteRoommateForm"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		log.WithError(err).Error("failed to get user ID")
		utils.WriteError(w, "failed to get user ID", http.StatusUnauthorized)
		return
	}

	_, err = h.grpcClient.DeleteRoommateForm(r.Context(), &user.DeleteRoommateFormRequest{
		UserId: int64(userID),
	})
	if err != nil {
		log.WithError(err).Error("h.grpcClient.DeleteRoommateForm")
		utils.HandelGRPCError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary Get incoming roommate matches
// @Description Returns users who added current authorized user to roommates
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} user.GetIncomingRoommateMatchesResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /user/me/roommate-matches/incoming [get]
func (h *UserHandler) GetIncomingRoommateMatches(w http.ResponseWriter, r *http.Request) {
	op := "UserHandler.GetIncomingRoommateMatches"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		log.WithError(err).Error("failed to get user ID from context")
		utils.WriteError(w, "failed to get user ID", http.StatusUnauthorized)
		return
	}

	resp, err := h.grpcClient.GetIncomingRoommateMatches(r.Context(), &user.GetIncomingRoommateMatchesRequest{
		UserId: int64(userID),
	})
	if err != nil {
		log.WithError(err).Error("h.grpcClient.GetIncomingRoommateMatches")
		utils.HandelGRPCError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, resp)
}

// @Summary Get matched roommate matches
// @Description Returns users who have mutual roommate match with current authorized user
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} user.GetMatchedRoommateMatchesResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /user/me/roommate-matches/matched [get]
func (h *UserHandler) GetMatchedRoommateMatches(w http.ResponseWriter, r *http.Request) {
	op := "UserHandler.GetMatchedRoommateMatches"
	log := ctxLogger.GetLogger(r.Context()).WithField("op", op)

	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		log.WithError(err).Error("failed to get user ID from context")
		utils.WriteError(w, "failed to get user ID", http.StatusUnauthorized)
		return
	}

	resp, err := h.grpcClient.GetMatchedRoommateMatches(r.Context(), &user.GetMatchedRoommateMatchesRequest{
		UserId: int64(userID),
	})
	if err != nil {
		log.WithError(err).Error("h.grpcClient.GetMatchedRoommateMatches")
		utils.HandelGRPCError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, resp)
}
