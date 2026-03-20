package poster

import (
	"context"
	"fmt"
	"log"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase"
)

const (
	MaxPostersLimit = 12
	PropertyFlat    = "flat"
	PropertyHouse   = "house"
)

type PosterUseCase struct {
	repo usecase.PosterRepo
}

func NewPosterUseCase(repo usecase.PosterRepo) *PosterUseCase {
	return &PosterUseCase{
		repo: repo,
	}
}

func (uc *PosterUseCase) GetPostersUseCase(ctx context.Context, filters dto.PostersFiltersDTO) ([]dto.PosterCardDTO, error) {
	if filters.Limit <= 0 || filters.Offset < 0 {
		return nil, entity.InvalidInput
	}

	if filters.Limit > MaxPostersLimit {
		filters.Limit = MaxPostersLimit
	}

	// total := uc.repo.CountPosters()

	// if params.Offset >= total {
	// 	return nil, entity.OffsetOutOfRange
	// }

	posters, err := uc.repo.GetPosters(ctx, filters)
	if err != nil {
		log.Printf("uc.repo.GetPosters: %s", err)
		return nil, err
	}

	return dto.PostersToPostersDTO(posters), nil
}

func (uc *PosterUseCase) GetPosterByAliasUseCase(ctx context.Context, posterAlias string) (*dto.PosterDTO, error) {
	var posterDTO *dto.PosterDTO

	poster, err := uc.repo.GetPosterByAlias(ctx, posterAlias)
	if err != nil {
		log.Printf("uc.repo.GetPosterByAlias: %s", err)
		return nil, err
	}

	posterDTO = dto.PosterToPosterDTO(poster)

	switch poster.Category {
	case PropertyFlat:
		flat, err := uc.repo.GetFlatByPropetyID(ctx, poster.PropertyID)
		if err != nil {
			log.Printf("uc.repo.GetPosters: %s", err)
			return nil, err
		}

		flatDTO := dto.FlatToFlatFlatDTO(flat)
		posterDTO.Flat = flatDTO

	case PropertyHouse:
	default:
		return nil, fmt.Errorf("no such category %w", entity.ServiceError)
	}

	return posterDTO, nil
}
