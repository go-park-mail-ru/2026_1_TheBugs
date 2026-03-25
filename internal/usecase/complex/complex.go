package complex

import (
	"context"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
)

type UtilityCompanyUseCase struct {
	repo usecase.UtilityCompanyRepo
}

func NewUtilityCompanyUseCase(repo usecase.UtilityCompanyRepo) *UtilityCompanyUseCase {
	return &UtilityCompanyUseCase{repo: repo}
}

func (uc *UtilityCompanyUseCase) GetUtilityCompany(ctx context.Context, alias string) (*dto.UtilityCompanyDTO, error) {
	return uc.repo.GetByAlias(ctx, alias)
}
