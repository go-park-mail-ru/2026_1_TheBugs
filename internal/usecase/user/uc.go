package user

import (
	"context"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
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
