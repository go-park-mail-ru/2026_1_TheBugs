package user

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/validator"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/photo"
)

type UserUseCase struct {
	uow  usecase.UnitOfWork
	file usecase.FileRepo
}

func NewUserUseCase(uow usecase.UnitOfWork, file usecase.FileRepo) *UserUseCase {
	return &UserUseCase{
		uow:  uow,
		file: file,
	}
}

func (uc *UserUseCase) GetByID(ctx context.Context, userID int) (*dto.UserDTO, error) {
	return uc.uow.Users().GetByID(ctx, userID)
}

func (uc *UserUseCase) UpdateProfile(ctx context.Context, data dto.UpdateProfileRequest) (*dto.UserDTO, error) {
	if data.Phone != nil {
		if ok := validator.ValidatePhone(*data.Phone); !ok {
			return nil, entity.NewValidationError("phone")
		}
		*data.Phone = validator.NormolizePhoneNumber(*data.Phone)
	}
	if data.FirstName != nil {
		if ok := validator.ValidateName(*data.FirstName); !ok {
			return nil, entity.NewValidationError("first_name")
		}
	}
	if data.LastName != nil {
		if ok := validator.ValidateName(*data.LastName); !ok {
			return nil, entity.NewValidationError("last_name")
		}
	}
	updateDTO := dto.UpdateProfileDTO{ID: data.ID, FirstName: data.FirstName, LastName: data.LastName, Phone: data.Phone}
	if data.Avatar != nil {
		file, err := data.Avatar.Open()
		if err != nil {
			return nil, fmt.Errorf("data.Avatar.Open(): %w", err)
		}
		defer file.Close()

		if !validator.ValidatePhoto(data.Avatar) {
			return nil, entity.NewValidationError("photo")
		}
		path := dto.GenerateAvatarPathForUser(data.ID)
		key := photo.GetKeyFromPath(path)
		size := data.Avatar.Size
		contentType := data.Avatar.Header.Get("Content-Type")

		if err := uc.file.Upload(ctx, key, file, size, contentType); err != nil {
			return nil, fmt.Errorf("uc.file.Upload: %w", err)
		}
		updateDTO.AvatarPath = &path
	}
	return uc.uow.Users().UpdateProfile(ctx, updateDTO)
}
