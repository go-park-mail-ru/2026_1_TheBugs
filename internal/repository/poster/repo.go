package poster

import (
	"context"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PosterRepo struct {
	pool repository.DB
}

func NewPosterRepo(pool repository.DB) *PosterRepo {
	return &PosterRepo{
		pool: pool,
	}
}

func (r *PosterRepo) GetPosters(ctx context.Context, filters dto.PostersFiltersDTO) ([]entity.Poster, error) {
	query := `SELECT id, price, image_url, address, metro, area, floor, type, rating, beds
		FROM posters
		ORDER BY id
		LIMIT $1 OFFSET $2`

	rows, err := r.pool.Query(ctx, query, filters.Limit, filters.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posters, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Poster])
	if err != nil {
		return nil, entity.CollectPostersErr
	}

	return posters, rows.Err()
}

func (r *PosterRepo) CountPosters(ctx context.Context) (int, error) {
	var count int
	row := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM posters")
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
