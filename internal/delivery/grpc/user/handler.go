package user

import (
	"bytes"
	"context"
	"io"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/generated/user"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/utils"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServiceServer struct {
	user.UnimplementedUserServiceServer
	uc delivery.UserUseCase
}

func NewUserServiceServer(uc delivery.UserUseCase) *UserServiceServer {
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

func (s *UserServiceServer) GetRoommateUser(ctx context.Context, req *user.GetRoommateUserRequest) (*user.GetRoommateUserResponse, error) {
	roommateUser, err := s.uc.GetRoommateUser(ctx, int(req.UserId))
	if err != nil {
		return nil, utils.TranslateDomainsError(err)
	}

	tags := make([]*user.RoommateTag, 0, len(roommateUser.Tags))
	for _, tag := range roommateUser.Tags {
		tags = append(tags, &user.RoommateTag{
			Name:  tag.Name,
			Alias: tag.Alias,
		})
	}

	return &user.GetRoommateUserResponse{
		FirstName:   roommateUser.FirstName,
		LastName:    roommateUser.LastName,
		AvatarUrl:   roommateUser.AvatarURL,
		Gender:      roommateUser.Gender,
		Birthday:    roommateUser.Birthday,
		Description: roommateUser.Description,
		Tags:        tags,
	}, nil
}

func (s *UserServiceServer) AddRoommateMatch(ctx context.Context, req *user.AddRoommateMatchRequest) (*user.AddRoommateMatchResponse, error) {
	if req.FromUserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid from_user_id")
	}

	if req.ToUserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid to_user_id")
	}

	err := s.uc.AddRoommateMatch(ctx, int(req.FromUserId), int(req.ToUserId))
	if err != nil {
		return nil, utils.TranslateDomainsError(err)
	}

	return &user.AddRoommateMatchResponse{}, nil
}

func (s *UserServiceServer) GetRoommateContacts(ctx context.Context, req *user.GetRoommateContactsRequest) (*user.GetRoommateContactsResponse, error) {
	if req.FromUserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid from_user_id")
	}

	if req.ToUserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid to_user_id")
	}

	contacts, err := s.uc.GetRoommateContacts(ctx, int(req.FromUserId), int(req.ToUserId))
	if err != nil {
		return nil, utils.TranslateDomainsError(err)
	}

	return &user.GetRoommateContactsResponse{
		Email: contacts.Email,
		Phone: contacts.Phone,
	}, nil
}

func (s *UserServiceServer) CreateRoommateForm(ctx context.Context, req *user.CreateRoommateFormRequest) (*user.CreateRoommateFormResponse, error) {
	if req.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	if req.Gender == "" {
		return nil, status.Error(codes.InvalidArgument, "gender required")
	}

	if req.Birthday == "" {
		return nil, status.Error(codes.InvalidArgument, "birthday required")
	}

	if req.Description == "" {
		return nil, status.Error(codes.InvalidArgument, "description required")
	}

	err := s.uc.CreateRoommateForm(ctx, dto.CreateRoommateFormRequest{
		UserID:      int(req.UserId),
		Gender:      req.Gender,
		Birthday:    req.Birthday,
		Description: req.Description,
		Tags:        req.Tags,
	})
	if err != nil {
		return nil, utils.TranslateDomainsError(err)
	}

	return &user.CreateRoommateFormResponse{}, nil
}

func (s *UserServiceServer) GetRoommateForm(ctx context.Context, req *user.GetRoommateFormRequest) (*user.GetRoommateFormResponse, error) {

	if req.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	form, err := s.uc.GetRoommateForm(ctx, int(req.UserId))
	if err != nil {
		return nil, utils.TranslateDomainsError(err)
	}

	return &user.GetRoommateFormResponse{
		Gender:      form.Gender,
		Birthday:    form.Birthday,
		Description: form.Description,
		Tags:        form.Tags,
	}, nil
}

func (s *UserServiceServer) UpdateRoommateForm(ctx context.Context, req *user.UpdateRoommateFormRequest) (*user.UpdateRoommateFormResponse, error) {
	if req.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	if req.Gender == "" {
		return nil, status.Error(codes.InvalidArgument, "gender required")
	}

	if req.Birthday == "" {
		return nil, status.Error(codes.InvalidArgument, "birthday required")
	}

	if req.Description == "" {
		return nil, status.Error(codes.InvalidArgument, "description required")
	}

	err := s.uc.UpdateRoommateForm(ctx, dto.CreateRoommateFormRequest{
		UserID:      int(req.UserId),
		Gender:      req.Gender,
		Birthday:    req.Birthday,
		Description: req.Description,
		Tags:        req.Tags,
	})
	if err != nil {
		return nil, utils.TranslateDomainsError(err)
	}

	return &user.UpdateRoommateFormResponse{}, nil
}

func (s *UserServiceServer) GetIncomingRoommateMatches(
	ctx context.Context,
	req *user.GetIncomingRoommateMatchesRequest,
) (*user.GetIncomingRoommateMatchesResponse, error) {
	if req.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	resp, err := s.uc.GetIncomingRoommateMatches(ctx, int(req.UserId))
	if err != nil {
		return nil, utils.TranslateDomainsError(err)
	}

	users := make([]*user.RoommateUser, 0, len(resp.Users))
	for _, u := range resp.Users {
		users = append(users, &user.RoommateUser{
			Id:        int64(u.ID),
			FirstName: u.FirstName,
			LastName:  u.LastName,
			AvatarUrl: u.AvatarURL,
		})
	}

	return &user.GetIncomingRoommateMatchesResponse{
		Users: users,
		Len:   int64(resp.Len),
	}, nil
}

func (s *UserServiceServer) GetMatchedRoommateMatches(
	ctx context.Context,
	req *user.GetMatchedRoommateMatchesRequest,
) (*user.GetMatchedRoommateMatchesResponse, error) {
	if req.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	resp, err := s.uc.GetMatchedRoommateMatches(ctx, int(req.UserId))
	if err != nil {
		return nil, utils.TranslateDomainsError(err)
	}

	users := make([]*user.RoommateUser, 0, len(resp.Users))
	for _, u := range resp.Users {
		users = append(users, &user.RoommateUser{
			Id:        int64(u.ID),
			FirstName: u.FirstName,
			LastName:  u.LastName,
			AvatarUrl: u.AvatarURL,
		})
	}

	return &user.GetMatchedRoommateMatchesResponse{
		Users: users,
		Len:   int64(resp.Len),
	}, nil
}
