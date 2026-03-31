package user

import (
	"net/http"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/utils"
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
