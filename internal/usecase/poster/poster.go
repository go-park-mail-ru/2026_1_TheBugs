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
	PropertyFlat    = "Квартиры" // TODO: делать по id а не по имени а то оно меняется либо поменять в
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

	posters, err := uc.repo.GetAll(ctx, filters)
	if err != nil {
		log.Printf("uc.repo.GetAll: %s", err)
		return nil, err
	}

	return dto.PostersToPostersDTO(posters), nil
}

func (uc *PosterUseCase) GetPosterByAliasUseCase(ctx context.Context, posterAlias string) (*dto.PosterDTO, error) {
	var posterDTO *dto.PosterDTO

	poster, err := uc.repo.GetByAlias(ctx, posterAlias)
	if err != nil {
		log.Printf("uc.repo.GetByAlias: %s", err)
		return nil, err
	}

	posterDTO = dto.PosterToPosterDTO(poster)
	log.Println(poster, posterDTO)

	switch poster.Category {
	case PropertyFlat:
		flat, err := uc.repo.GetFlatByPropetyID(ctx, poster.PropertyID)
		if err != nil {
			log.Printf("uc.repo.GetAll: %s", err)
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

func (uc *PosterUseCase) GetPosterByUserID(ctx context.Context, userID int) ([]dto.PosterCardDTO, error) {
	posters, err := uc.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return dto.PostersToPostersDTO(posters), nil
}
