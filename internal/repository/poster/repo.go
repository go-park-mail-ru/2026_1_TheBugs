package poster

import "github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"

type PosterRepo struct {
	listPoster []entity.Poster
}

func NewPosterRepo() *PosterRepo {
	return &PosterRepo{
		listPoster: []entity.Poster{},
	}
}

func (r *PosterRepo) GetPosters(limit, offset, end int) ([]*entity.Poster, error) {
	posters := make([]*entity.Poster, 0, end)
	for i := offset; i < end; i++ {
		posters = append(posters, &r.listPoster[i])
	}

	return posters, nil
}

func (r *PosterRepo) CountPosters() (int, error) {
	return len(r.listPoster), nil
}
