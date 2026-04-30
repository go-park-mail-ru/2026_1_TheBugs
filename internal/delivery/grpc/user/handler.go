package user

import (
	"bytes"
	"context"
	"io"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/generated/user"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/utils"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	userUC "github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/user"
)

type UserServiceServer struct {
	user.UnimplementedUserServiceServer
	uc *userUC.UserUseCase
}

func NewUserServiceServer(uc *userUC.UserUseCase) *UserServiceServer {
	return &UserServiceServer{
		uc: uc,
	}
}

func (s *UserServiceServer) GetMe(ctx context.Context, req *user.GetMeRequest) (*user.GetMeResponse, error) {
	userDTO, err := s.uc.GetByID(ctx, int(req.UserId))
	if err != nil {
		return nil, utils.TranslateDomainsError(err)
	}

	return &user.GetMeResponse{
		Id:        int32(userDTO.ID),
		Email:     userDTO.Email,
		Phone:     userDTO.Phone,
		AvatarUrl: userDTO.AvatarURL,
		Firstname: userDTO.FirstName,
		Lastname:  userDTO.LastName,
	}, nil
}

func (s *UserServiceServer) UpdateProfile(
	ctx context.Context,
	req *user.UpdateProfileRequest,
) (*user.GetMeResponse, error) {

	var fileInput *dto.FileInput
	if req.File != nil {
		fileInput = &dto.FileInput{
			Filename:    req.File.Filename,
			File:        io.NopCloser(bytes.NewReader(req.File.Avatar)),
			Size:        req.File.Size,
			ContentType: req.File.ContentType,
		}
	}

	data, err := s.uc.UpdateProfile(ctx, dto.UpdateProfileRequest{
		ID:        int(req.Id),
		FirstName: req.Firstname,
		LastName:  req.Lastname,
		Phone:     req.Phone,
		Avatar:    fileInput,
	})
	if err != nil {
		return nil, utils.TranslateDomainsError(err)
	}

	return &user.GetMeResponse{
		Id:        int32(data.ID),
		Email:     data.Email,
		Phone:     data.Phone,
		AvatarUrl: data.AvatarURL,
		Firstname: data.FirstName,
		Lastname:  data.LastName,
	}, nil
}
