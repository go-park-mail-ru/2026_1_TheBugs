package poster

import (
	"log"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase"
)

var MaxPostersLimit = 12

type PosterUseCase struct {
	repo usecase.PosterRepo
}

func NewPosterUseCase(repo usecase.PosterRepo) *PosterUseCase {
	return &PosterUseCase{
		repo: repo,
	}
}

func (uc *PosterUseCase) GetPostersUseCase(limit, offset int) ([]dto.PosterDTO, error) {
	if limit <= 0 {
		return nil, entity.InvalidInput
	}

	if offset < 0 {
		return nil, entity.InvalidInput
	}

	if limit > MaxPostersLimit {
		limit = MaxPostersLimit
	}

	posters, err := uc.repo.GetPosters(limit, offset)
	if err != nil {
		log.Printf("uc.repo.GetPosters: %s", err)
		return nil, err
	}

	listPosters := make([]dto.PosterDTO, 0, limit)
	for _, poster := range posters {
		posterDTO := dto.PosterDTO{
			Price:   poster.Price,
			ImgURL:  poster.ImgURL,
			Address: poster.Address,
			Metro:   poster.Metro,
			Area:    poster.Area,
			Floor:   poster.Floor,
			Type:    poster.Type,
		}

		listPosters = append(listPosters, posterDTO)
	}

	return listPosters, nil
}
