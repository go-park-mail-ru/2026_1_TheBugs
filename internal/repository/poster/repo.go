package poster

import (
	"context"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PosterRepo struct {
	pool *pgxpool.Pool
}

func NewPosterRepo(pool *pgxpool.Pool) *PosterRepo {
	return &PosterRepo{
		pool: pool,
	}
}

func (r *PosterRepo) GetPosters(ctx context.Context, filters dto.PostersFiltersDTO) ([]entity.Poster, error) {
	query := `SELECT id, price, image_url, address, metro, area, floor, type
		FROM posters
		ORDER BY id
		LIMIT $1 OFFSET $2`

	rows, err := r.pool.Query(ctx, query, filters.Limit, filters.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posters, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Poster])

	return posters, rows.Err()
}

func (r *PosterRepo) CountPosters(ctx context.Context) int {
	var count int
	row := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM posters")
	row.Scan(&count)
	return count
}
