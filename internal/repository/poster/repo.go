package poster

import (
	"context"
	"fmt"

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
	query := `
        SELECT p.id, p.price, p.avatar_url, 
               b.address, m.station_name, prop.area, f.floor
        FROM posters p
        JOIN property prop ON prop.id = p.property_id
		JOIN flat f ON f.property_id = p.id
		JOIN property_categories pc ON pc.id = prop.category_id
        JOIN buildings b ON b.id = prop.building_id
        JOIN metro_stations m ON b.metro_station_id = m.id
	`
	args := []any{filters.Limit, filters.Offset}
	argIndex := 3
	if filters.UtilityCompany != nil {
		query += ` JOIN utility_companies uc ON b.company_id = uc.id 
		WHERE uc.alias = $` + fmt.Sprintf("%d", argIndex)
		args = append(args, *filters.UtilityCompany)
		argIndex++
	}

	query += ` ORDER BY p.created_at DESC LIMIT $1 OFFSET $2`

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	defer rows.Close()

	posters, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Poster])
	if err != nil {
		return nil, repository.HandelPgErrors(err)
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
