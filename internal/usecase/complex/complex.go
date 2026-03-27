package complex

import (
	"context"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase"
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

func (uc *UtilityCompanyUseCase) GetAllDevelopers(ctx context.Context) ([]dto.DeveloperDTO, error) {
	return uc.repo.GetAllDevelopers(ctx)
}

func (uc *UtilityCompanyUseCase) GetAllByDeveloperID(ctx context.Context, companyID int) ([]dto.UtilityCompanyCardDTO, error) {
	return uc.repo.GetAllByDeveloperID(ctx, companyID)
}
