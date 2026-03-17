package poster

import (
	"context"
	"log"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase"
)

const MaxPostersLimit = 12

type PosterUseCase struct {
	repo usecase.PosterRepo
}

func NewPosterUseCase(repo usecase.PosterRepo) *PosterUseCase {
	return &PosterUseCase{
		repo: repo,
	}
}

func (uc *PosterUseCase) GetPostersUseCase(ctx context.Context, filters dto.PostersFiltersDTO) ([]dto.PosterDTO, error) {
	if filters.Limit <= 0 || filters.Offset < 0 {
		return nil, entity.InvalidInput
	}

	if filters.Limit > MaxPostersLimit {
		filters.Limit = MaxPostersLimit
	}

	posters, err := uc.repo.GetPosters(ctx, filters)
	if err != nil {
		log.Printf("uc.repo.GetPosters: %s", err)
		return nil, err
	}

	return dto.PostersToPostersDTO(posters), nil
}
