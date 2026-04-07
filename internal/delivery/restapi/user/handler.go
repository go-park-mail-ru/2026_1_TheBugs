package user

import (
	"net/http"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/utils"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/user"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/ctxLogger"
)

type UserHandler struct {
	uc *user.UserUseCase
}

func NewUserHandler(uc *user.UserUseCase) *UserHandler {
	return &UserHandler{
		uc: uc,
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

	user, err := h.uc.GetByID(r.Context(), userID)
	if err != nil {
		log.WithError(err).Error("h.uc.GetByID: failed to get user by ID")
		utils.HandelError(w, err)
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

	err := utils.ParseMultipartFormData(r, &data)
	if err != nil {
		log.WithError(err).Error("failed to decode update profile data")
		utils.WriteError(w, "invalid body", http.StatusBadRequest)
		return
	}
	files, ok := r.MultipartForm.File["avatar"]
	if ok && len(files) > 0 {
		data.Avatar = files[0]
	}
	userID, err := utils.GetUserID(r.Context())
	if err != nil {
		log.WithError(err).Error("failed to get user ID from context")
		utils.WriteError(w, "failed to get user ID", http.StatusUnauthorized)
		return
	}
	data.ID = userID

	user, err := h.uc.UpdateProfile(r.Context(), data)
	if err != nil {
		log.WithError(err).Error("h.uc.UpdateProfile: failed to update user profile")
		utils.HandelError(w, err)
		return
	}

	utils.JSONResponse(w, http.StatusOK, user)
}
