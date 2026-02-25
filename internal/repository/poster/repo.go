package poster

import "github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"

type PosterRepo struct {
	listPoster []entity.Poster // сюда засугем список объявлений
}

func NewPosterRepo() *PosterRepo {
	return &PosterRepo{
		listPoster: []entity.Poster{},
	}
}

func (r *PosterRepo) GetPosters(limit, offset int) ([]entity.Poster, error) {
	total := len(r.listPoster)

	if offset >= total {
		return []entity.Poster{}, nil
	}

	end := offset + limit
	if end > total {
		end = total
	}

	posters := make([]entity.Poster, 0, limit)
	for _, poster := range r.listPoster[offset:end] {
		posters = append(posters, poster)
	}

	return posters, nil
}

// []entity.Poster хз отдавать слайс указателей или нет
// если я вроде не собираюсь изменять само объявление внутри usecase
