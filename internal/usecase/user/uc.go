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
	user, err := uc.uow.Users().GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("uc.uow.Users().GetByID: %w", err)
	}
	user.MakeAvatarPath()
	return user, nil
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
		defer data.Avatar.File.Close()

		if !validator.ValidatePhoto(data.Avatar) {
			return nil, entity.NewValidationError("photo")
		}
		user, err := uc.uow.Users().GetByID(ctx, data.ID)
		if err != nil {
			return nil, fmt.Errorf("uc.uow.Users().GetByID: %w", err)
		}
		path := dto.GenerateAvatarPathForUser(data.ID)
		key := photo.GetKeyFromPath(path)
		size := data.Avatar.Size
		contentType := data.Avatar.ContentType

		if err := uc.file.Upload(ctx, key, data.Avatar.File, size, contentType); err != nil {
			return nil, fmt.Errorf(" uc.file.Upload: %w", err)
		}
		if user.AvatarURL != nil {
			if err := uc.file.Delete(ctx, photo.GetKeyFromPath(*user.AvatarURL)); err != nil {
				return nil, fmt.Errorf(" uc.file.Upload: %w", err)
			}
		}
		updateDTO.AvatarPath = &path
	}
	user, err := uc.uow.Users().UpdateProfile(ctx, updateDTO)
	if err != nil {
		return nil, fmt.Errorf("uc.uow.Users().UpdateProfile: %w", err)
	}
	user.MakeAvatarPath()
	return user, nil
}
