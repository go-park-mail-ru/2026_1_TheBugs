package poster

import (
	"log"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/poster/request"
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

func (uc *PosterUseCase) GetPostersUseCase(params request.PostersRequest) ([]dto.PosterDTO, error) {
	if params.Limit <= 0 || params.Offset < 0 {
		return nil, entity.InvalidInput
	}

	if params.Limit > MaxPostersLimit {
		params.Limit = MaxPostersLimit
	}

	total, err := uc.repo.CountPosters()
	if err != nil {
		log.Printf("uc.repo.CountPosters: %s", err)
		return nil, err
	}

	if params.Offset >= total {
		return nil, entity.OffsetOutOfRange
	}

	end := params.Offset + params.Limit
	if end > total {
		end = total
	}

	posters, err := uc.repo.GetPosters(params.Limit, params.Offset, end)
	if err != nil {
		log.Printf("uc.repo.GetPosters: %s", err)
		return nil, err
	}

	listPosters := make([]dto.PosterDTO, 0, params.Limit)
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
