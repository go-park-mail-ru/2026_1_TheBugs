package user

import (
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
