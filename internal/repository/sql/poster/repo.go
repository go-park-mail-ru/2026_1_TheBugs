package poster

import (
	"context"
	"fmt"
	"strings"

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
	log := ctxLogger.GetLogger(ctx).WithField("op", "PosterRepo.CountPosters")
	log.Info("start db query")

	var count int
	row := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM posters")
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *PosterRepo) GetByAlias(ctx context.Context, posterAlias string) (*entity.PosterById, error) {
	log := ctxLogger.GetLogger(ctx).WithField("op", "PosterRepo.GetByAlias")
	log.Info("start db query")

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
	log := ctxLogger.GetLogger(ctx).WithField("op", "PosterRepo.GetFlatByPropetyID")
	log.Info("start db query")

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

func (r *PosterRepo) CreateBuilding(ctx context.Context, poster *entity.PosterInput) (int, error) {
	log := ctxLogger.GetLogger(ctx).WithField("op", "PosterRepo.CreateBuilding")
	log.Info("start db query")

	buildingQuery := `
		INSERT INTO buildings (address, geo, city_id,
			metro_station_id, district, floor_count, company_id)
		VALUES ($1, ST_GeogFromText($2), $3, $4, $5, $6, $7)
		RETURNING id
	`

	var buildingID int
	err := r.pool.QueryRow(ctx, buildingQuery, poster.Address,
		poster.Geo.ToGeo(), poster.CityID, poster.MetroStationID,
		poster.District, poster.FloorCount, poster.CompanyID,
	).Scan(&buildingID)

	if err != nil {
		return 0, repository.HandelPgErrors(err)
	}

	return buildingID, nil
}

func (r *PosterRepo) CreateProperty(ctx context.Context, poster *entity.PosterInput, buildingID int) (int, error) {
	log := ctxLogger.GetLogger(ctx).WithField("op", "PosterRepo.CreateProperty")
	log.Info("start db query")

	propertyQuery := `
		INSERT INTO property (category_id, building_id, area)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var propertyID int
	err := r.pool.QueryRow(ctx, propertyQuery,
		poster.CategoryID, buildingID, poster.Area,
	).Scan(&propertyID)

	if err != nil {
		return 0, repository.HandelPgErrors(err)
	}

	return propertyID, nil
}

func (r *PosterRepo) Create(ctx context.Context, poster *entity.PosterInput, propertyID int) (int, error) {
	log := ctxLogger.GetLogger(ctx).WithField("op", "PosterRepo.Create")
	log.Info("start db query")

	posterQuery := `
		INSERT INTO posters (price, description,
			user_id, property_id, alias)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var posterID int
	err := r.pool.QueryRow(ctx, posterQuery,
		poster.Price, poster.Description,
		poster.UserID, propertyID, poster.Alias,
	).Scan(&posterID)

	if err != nil {
		return 0, repository.HandelPgErrors(err)
	}

	return posterID, nil
}

func (r *PosterRepo) InsertFlat(ctx context.Context, flat *entity.FlatInput) error {
	log := ctxLogger.GetLogger(ctx).WithField("op", "PosterRepo.InsertFlat")
	log.Info("start db query")

	query := `
		INSERT INTO flat (property_id,
			floor, number, category_id)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.pool.Exec(ctx, query, flat.PropertyID,
		flat.Floor, flat.Number, flat.CategoryID,
	)
	if err != nil {
		return repository.HandelPgErrors(err)
	}

	return nil
}

func (r *PosterRepo) InsertPhotos(ctx context.Context, posterID int, photos []entity.PhotoInput) error {
	log := ctxLogger.GetLogger(ctx).WithField("op", "PosterRepo.InsertPhotos")
	log.Info("start db query")

	if len(photos) == 0 {
		return nil
	}

	args := make([]any, 0, len(photos)*3)
	list := make([]string, 0, len(photos))

	for i, photo := range photos {
		base := i * 3
		list = append(list, fmt.Sprintf("($%d, $%d, $%d)", base+1, base+2, base+3))
		args = append(args, photo.Path, photo.Order, posterID)
	}

	query := `
		INSERT INTO poster_photos (img_url, sequence_order, poster_id)
		VALUES ` + strings.Join(list, ",")

	_, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return repository.HandelPgErrors(err)
	}

	return nil
}

func (r *PosterRepo) InsertMainPhoto(ctx context.Context, posterID int, avatarURL string) error {
	log := ctxLogger.GetLogger(ctx).WithField("op", "PosterRepo.InsertMainPhoto")
	log.Info("start db query")

	query := `
		UPDATE posters
		SET avatar_url = $1
		WHERE id = $2
	`

	_, err := r.pool.Exec(ctx, query, avatarURL, posterID)
	if err != nil {
		return repository.HandelPgErrors(err)
	}

	return nil
}
