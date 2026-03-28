package poster

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"
	repository "github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/sql"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/ctxLogger"
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

func (r *PosterRepo) GetAll(ctx context.Context, filters dto.PostersFiltersDTO) ([]entity.Poster, error) {
	log := ctxLogger.GetLogger(ctx).WithField("op", "PosterRepo.GetAll")
	log.Info("start db query")
	query := `
        SELECT p.id, p.price, p.avatar_url,
               b.address, m.station_name, prop.area, f.floor, p.alias
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

	return posters, nil
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

func (r *PosterRepo) GetByAlias(ctx context.Context, posterAlias string) (*entity.PosterById, error) {
	query := `
		SELECT p.id, p.alias, p.price, pc.name AS category,
			   p.description, prop.area, prop.id AS property_id, 
			   ST_AsText(b.geo) AS building_geo, b.address, b.district, 
			   ms.station_name, ST_AsText(ms.geo) AS metro_geo, c.city_name, 
			   b.floor_count, pr.first_name, pr.last_name, pr.phone, pr.avatar_url as seller_avatar_url,
			   uc.company_name, uc.avatar_url AS company_avatar_url, uc.alias AS company_alias, uc.id AS company_id
		FROM posters p
		JOIN property prop ON prop.id = p.property_id
		JOIN property_categories pc ON pc.id = prop.category_id
		JOIN buildings b ON b.id = prop.building_id
		JOIN cities c ON c.id = b.city_id
		LEFT JOIN metro_stations ms ON ms.id = b.metro_station_id
		LEFT JOIN utility_companies uc ON uc.id = b.company_id
		JOIN users u ON u.id = p.user_id
		JOIN profiles pr ON pr.id = u.profile_id
		WHERE p.alias = $1;
	`

	rows, err := r.pool.Query(ctx, query, posterAlias)
	if err != nil {
		return &entity.PosterById{}, repository.HandelPgErrors(err)
	}

	defer rows.Close()

	poster, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.PosterById])
	if err != nil {
		return &entity.PosterById{}, repository.HandelPgErrors(err)
	}

	poster.Images, err = getPosterImages(r, ctx, poster.ID)
	if err != nil {
		return &poster, repository.HandelPgErrors(err)
	}

	return &poster, nil
}

func (r *PosterRepo) GetFlatByPropetyID(ctx context.Context, propertyID int) (*entity.Flat, error) {
	query := `
		SELECT f.property_id, f.number, f.floor,
			   fc.name AS flat_category
		FROM flat f
		LEFT JOIN flat_categories fc ON fc.id = f.category_id
		WHERE f.property_id = $1;
	`

	rows, err := r.pool.Query(ctx, query, propertyID)
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	defer rows.Close()

	flat, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Flat])
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	return &flat, nil
}

func (r *PosterRepo) GetByUserID(ctx context.Context, userID int) ([]entity.Poster, error) {
	sql := `
		SELECT p.id, p.price, p.avatar_url,
               b.address, m.station_name, prop.area, f.floor, p.alias
        FROM posters p
        JOIN property prop ON prop.id = p.property_id
        JOIN flat f ON f.property_id = p.id
        JOIN property_categories pc ON pc.id = prop.category_id
        JOIN buildings b ON b.id = prop.building_id
        JOIN metro_stations m ON b.metro_station_id = m.id
		WHERE p.user_id = $1
		ORDER BY p.created_at DESC
	`
	rows, err := r.pool.Query(ctx, sql, userID)
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}
	posters, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Poster])
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}
	return posters, nil

}

func getPosterImages(r *PosterRepo, ctx context.Context, id int) ([]entity.PosterImage, error) {
	query := `
		SELECT im.img_url, im.sequence_order
		FROM poster_photos im
		WHERE im.poster_id = $1
		ORDER BY im.sequence_order
	`

	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	defer rows.Close()

	images, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.PosterImage])
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	return images, rows.Err()
}
