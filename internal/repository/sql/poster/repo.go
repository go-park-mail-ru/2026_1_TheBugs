package poster

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	repository "github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/sql"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/ctxLogger"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/geo"
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

func (r *PosterRepo) GetFlatsAll(ctx context.Context, filters dto.PostersFiltersDTO) ([]entity.PosterFlat, error) {
	log := ctxLogger.GetLogger(ctx).WithField("op", "PosterRepo.GetFlatsAll")
	log.Info("start db query")
	query := `
        SELECT p.id, p.price, p.avatar_url,
               b.address, m.station_name, prop.area, f.floor, p.alias, fc.name AS flat_category
        FROM posters p
        JOIN property prop ON prop.id = p.property_id
        JOIN property_categories pc ON pc.id = prop.category_id
        JOIN buildings b ON b.id = prop.building_id
        LEFT JOIN metro_stations m ON b.metro_station_id = m.id
		LEFT JOIN flat f ON f.property_id = prop.id
		LEFT JOIN flat_categories fc ON fc.id = f.category_id
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

	posters, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.PosterFlat])
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	return posters, nil
}

func (r *PosterRepo) CountPosters(ctx context.Context) (int, error) {
	log := ctxLogger.GetLogger(ctx).WithField("op", "PosterRepo.CountPosters")
	log.Info("start db query")

	var count int
	row := r.pool.QueryRow(ctx, "SELECT get_explain_rows('SELECT * FROM posters')")
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
		SELECT p.id, p.alias, p.price, pc.name AS category_name, pc.alias AS category_alias,
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
	poster.Facilities, err = getPosterFacilities(r, ctx, poster.ID)
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
			   fc.name AS flat_category, fc.room_count
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
               b.address, prop.area, p.alias, pc.name as category_name, pc.alias as category_alias
        FROM posters p
        JOIN property prop ON prop.id = p.property_id
        JOIN property_categories pc ON pc.id = prop.category_id
        JOIN buildings b ON b.id = prop.building_id
		WHERE p.user_id = $1;
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

func getPosterFacilities(r *PosterRepo, ctx context.Context, id int) ([]entity.Facility, error) {
	query := `
		SELECT f.id, f.name, f.alias
		FROM facilities f 
		JOIN facility_property fp ON fp.facility_id = f.id
		JOIN property pr ON pr.id = fp.property_id
		JOIN posters pt ON pt.property_id = pr.id
		WHERE pt.id = $1
		ORDER BY f.name
	`

	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	defer rows.Close()

	facilities, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Facility])
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	return facilities, rows.Err()
}

func (r *PosterRepo) GetMetroStationByRadius(ctx context.Context, buildingGeo dto.GeographyDTO, radius entity.Metre) ([]entity.MetroStation, error) {
	query := `
    SELECT m.id, m.station_name, ST_AsText(m.geo) AS metro_geo
    FROM metro_stations m 
    WHERE ST_DWithin(m.geo, ST_GeogFromText($1), $2)
    ORDER BY m.geo <-> ST_GeogFromText($1)
`
	rows, err := r.pool.Query(ctx, query, geo.GeographyPoint{Lat: buildingGeo.Lat, Lon: buildingGeo.Lon}, int(radius))
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	defer rows.Close()

	stations, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.MetroStation])
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	return stations, rows.Err()
}

func (r *PosterRepo) CreateBuilding(ctx context.Context, poster *dto.PosterInput) (int, error) {
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

func (r *PosterRepo) CreateProperty(ctx context.Context, poster *dto.PosterInput, buildingID int) (int, error) {
	log := ctxLogger.GetLogger(ctx).WithField("op", "PosterRepo.CreateProperty")
	log.Infof("alias: %s", poster.CategoryAlias)

	propertyQuery := `
		INSERT INTO property (category_id, building_id, area)
		SELECT id, $2, $3
		FROM property_categories
		WHERE alias = $1
		RETURNING id;
	`

	var propertyID int
	err := r.pool.QueryRow(ctx, propertyQuery,
		poster.CategoryAlias, buildingID, poster.Area,
	).Scan(&propertyID)

	if err != nil {
		return 0, repository.HandelPgErrors(err)
	}

	return propertyID, nil
}

func (r *PosterRepo) InsertFacilities(ctx context.Context, propertyID int, aliases []string) error {
	log := ctxLogger.GetLogger(ctx).WithField("op", "PosterRepo.InsertFacilities")
	log.Info("start db query")

	selectQuery := `
		SELECT id
		FROM facilities
		WHERE alias = ANY($1::text[])
	`

	rows, err := r.pool.Query(ctx, selectQuery, aliases)
	if err != nil {
		return repository.HandelPgErrors(err)
	}
	defer rows.Close()

	facilityIDs, err := pgx.CollectRows(rows, pgx.RowTo[int])
	if err != nil {
		return repository.HandelPgErrors(err)
	}

	insertQuery := `
		INSERT INTO facility_property (property_id, facility_id)
		VALUES ($1, $2)
	`

	for _, facilityID := range facilityIDs {
		_, err = r.pool.Exec(ctx, insertQuery, propertyID, facilityID)
		if err != nil {
			return repository.HandelPgErrors(err)
		}
	}

	return nil
}

func (r *PosterRepo) Create(ctx context.Context, poster *dto.PosterInput, propertyID int) (int, error) {
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

func (r *PosterRepo) InsertFlat(ctx context.Context, flat *dto.FlatInput) error {
	log := ctxLogger.GetLogger(ctx).WithField("op", "PosterRepo.InsertFlat")
	log.Info("start db query")

	query := `
		INSERT INTO flat (property_id,
			floor, number, category_id)
		VALUES ($1, $2, $3, $4)
		RETURNING property_id
	`
	var flatID int
	err := r.pool.QueryRow(ctx, query, flat.PropertyID,
		flat.Floor, flat.Number, flat.CategoryID,
	).Scan(&flatID)
	if err != nil {
		return repository.HandelPgErrors(err)
	}

	return nil
}

func (r *PosterRepo) InsertPhotos(ctx context.Context, posterID int, photos []dto.PhotoInput) error {
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

func (r *PosterRepo) GetUpdateIDsByAlias(ctx context.Context, alias string) (*dto.PosterUpdateIDs, error) {
	log := ctxLogger.GetLogger(ctx).WithField("op", "PosterRepo.GetUpdateIDsByAlias")
	log.Info("start db query")

	query := `
		SELECT p.id, p.user_id, p.property_id, pr.building_id
		FROM posters p
		JOIN property pr ON pr.id = p.property_id
		WHERE p.alias = $1
	`

	var ids dto.PosterUpdateIDs

	err := r.pool.QueryRow(ctx, query, alias).Scan(
		&ids.UserID,
		&ids.PosterID,
		&ids.PropertyID,
		&ids.BuildingID,
	)
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	return &ids, nil
}

func (r *PosterRepo) Update(ctx context.Context, posterID int, poster *dto.PosterInput) error {
	log := ctxLogger.GetLogger(ctx).WithField("op", "PosterRepo.Update")
	log.Info("start db query")

	query := `
		UPDATE posters
		SET price = $1, description = $2
		WHERE id = $3
	`

	_, err := r.pool.Exec(ctx, query, poster.Price,
		poster.Description, posterID,
	)
	if err != nil {
		return repository.HandelPgErrors(err)
	}

	return nil
}

func (r *PosterRepo) UpdateProperty(ctx context.Context, propertyID int, poster *dto.PosterInput) error {
	log := ctxLogger.GetLogger(ctx).WithField("op", "PosterRepo.UpdateProperty")
	log.Info("start db query")

	query := `
		UPDATE property
		SET category_id = (
				SELECT id
				FROM property_categories
				WHERE alias = $1
			),
			area = $2
		WHERE id = $3
	`

	_, err := r.pool.Exec(ctx, query,
		poster.CategoryAlias,
		poster.Area,
		propertyID,
	)

	if err != nil {
		return repository.HandelPgErrors(err)
	}

	return nil
}

func (r *PosterRepo) UpdateBuilding(ctx context.Context, buildingID int, poster *dto.PosterInput) error {
	log := ctxLogger.GetLogger(ctx).WithField("op", "PosterRepo.UpdateBuilding")
	log.Info("start db query")

	query := `
		UPDATE buildings
		SET address = $1, geo = $2,
			city_id = $3, metro_station_id = $4,
			district = $5, floor_count = $6,
			company_id = $7
		WHERE id = $8
	`

	_, err := r.pool.Exec(ctx, query,
		poster.Address, poster.Geo,
		poster.CityID, poster.MetroStationID,
		poster.District, poster.FloorCount,
		poster.CompanyID, buildingID,
	)
	if err != nil {
		return repository.HandelPgErrors(err)
	}

	return nil
}

func (r *PosterRepo) UpdateFlat(ctx context.Context, flat *dto.FlatInput) error {
	log := ctxLogger.GetLogger(ctx).WithField("op", "PosterRepo.UpdateFlat")
	log.Info("start db query")

	query := `
		UPDATE flat
		SET floor = $1, number = $2, category_id = $3
		WHERE property_id = $4
	`

	_, err := r.pool.Exec(ctx, query, flat.Floor,
		flat.Number, flat.CategoryID, flat.PropertyID,
	)
	if err != nil {
		return repository.HandelPgErrors(err)
	}

	return nil
}

func (r *PosterRepo) GetPhotoPathsByPosterID(ctx context.Context, posterID int) ([]string, error) {
	log := ctxLogger.GetLogger(ctx).WithField("op", "PosterRepo.GetPhotoPathsByPosterID")
	log.Info("start db query")

	query := `
		SELECT img_url
		FROM poster_photos
		WHERE poster_id = $1
	`

	rows, err := r.pool.Query(ctx, query, posterID)
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}
	defer rows.Close()

	photoPaths, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	return photoPaths, nil
}

func (r *PosterRepo) DeleteFacilitiesByPropertyID(ctx context.Context, propertyID int) error {
	log := ctxLogger.GetLogger(ctx).WithField("op", "PosterRepo.DeleteFacilitiesByPropertyID")
	log.Info("start db query")

	query := `
		DELETE FROM facility_property
		WHERE property_id = $1
	`

	_, err := r.pool.Exec(ctx, query, propertyID)
	if err != nil {
		return repository.HandelPgErrors(err)
	}

	return nil
}

func (r *PosterRepo) DeletePhotosByPosterID(ctx context.Context, posterID int) error {
	log := ctxLogger.GetLogger(ctx).WithField("op", "PosterRepo.DeletePhotosByPosterID")
	log.Info("start db query")

	deletePhotosQuery := `
		DELETE FROM poster_photos
		WHERE poster_id = $1
	`

	_, err := r.pool.Exec(ctx, deletePhotosQuery, posterID)
	if err != nil {
		return repository.HandelPgErrors(err)
	}

	clearAvatarQuery := `
		UPDATE posters
		SET avatar_url = NULL
		WHERE id = $1
	`

	_, err = r.pool.Exec(ctx, clearAvatarQuery, posterID)
	if err != nil {
		return repository.HandelPgErrors(err)
	}

	return nil
}
