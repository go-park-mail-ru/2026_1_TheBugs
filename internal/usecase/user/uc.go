package user

import (
	"context"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/validator"
)

type UserUseCase struct {
	uow usecase.UnitOfWork
}

func NewUserUseCase(uow usecase.UnitOfWork) *UserUseCase {
	return &UserUseCase{
		uow: uow,
	}
}

func (uc *UserUseCase) GetByID(ctx context.Context, userID int) (*dto.UserDTO, error) {
	return uc.uow.Users().GetByID(ctx, userID)
}

func (uc *UserUseCase) UpdateProfile(ctx context.Context, data dto.UpdateProfileDTO) (*dto.UserDTO, error) {
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
	return uc.uow.Users().UpdateProfile(ctx, data)
}
