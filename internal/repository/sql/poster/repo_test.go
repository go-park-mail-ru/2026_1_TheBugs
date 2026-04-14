package poster

import (
	"context"
	"regexp"
	"testing"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/geo"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestGetPostersRepo(t *testing.T) {
	expectedListPoster := []entity.PosterFlat{
		{ID: 1, Price: 11111, ImgURL: nil, Address: "street_1", Metro: nil, Area: 35.5, Floor: lo.ToPtr(2), Alias: "alias_1", FlatCategory: lo.ToPtr("1-room")},
		{ID: 2, Price: 22222, ImgURL: nil, Address: "street_2", Metro: nil, Area: 40.0, Floor: lo.ToPtr(3), Alias: "alias_2", FlatCategory: lo.ToPtr("2-room")},
		{ID: 3, Price: 33333, ImgURL: nil, Address: "street_3", Metro: nil, Area: 45.2, Floor: lo.ToPtr(4), Alias: "alias_3", FlatCategory: lo.ToPtr("3-room")},
		{ID: 4, Price: 44444, ImgURL: nil, Address: "street_4", Metro: nil, Area: 50.7, Floor: lo.ToPtr(5), Alias: "alias_4", FlatCategory: lo.ToPtr("4-room")},
		{ID: 5, Price: 55555, ImgURL: nil, Address: "street_5", Metro: nil, Area: 60.1, Floor: lo.ToPtr(6), Alias: "alias_5", FlatCategory: lo.ToPtr("5-room")},
	}

	inputParams := dto.PostersFiltersDTO{
		Limit:  12,
		Offset: 0,
	}

	query := regexp.QuoteMeta(`
		SELECT p.id, p.price, p.avatar_url, b.address, m.station_name, prop.area, f.floor, p.alias, fc.name AS flat_category
		FROM posters p
		JOIN property prop ON prop.id = p.property_id
		JOIN property_categories pc ON pc.id = prop.category_id
		JOIN buildings b ON b.id = prop.building_id
		LEFT JOIN metro_stations m ON b.metro_station_id = m.id
		LEFT JOIN flat f ON f.property_id = prop.id
		LEFT JOIN flat_categories fc ON fc.id = f.category_id
		ORDER BY p.created_at DESC
		LIMIT $1 OFFSET $2`)

	tests := []struct {
		name      string
		params    dto.PostersFiltersDTO
		setupMock func(m pgxmock.PgxPoolIface)
		want      []entity.PosterFlat
		wantErr   error
	}{
		{
			name:   "OK",
			params: inputParams,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "price", "avatar_url", "address", "station_name", "area", "floor", "alias", "flat_category",
				}).AddRow(1, 11111.0, nil, "street_1", nil, 35.5, lo.ToPtr(2), "alias_1", lo.ToPtr("1-room")).
					AddRow(2, 22222.0, nil, "street_2", nil, 40.0, lo.ToPtr(3), "alias_2", lo.ToPtr("2-room")).
					AddRow(3, 33333.0, nil, "street_3", nil, 45.2, lo.ToPtr(4), "alias_3", lo.ToPtr("3-room")).
					AddRow(4, 44444.0, nil, "street_4", nil, 50.7, lo.ToPtr(5), "alias_4", lo.ToPtr("4-room")).
					AddRow(5, 55555.0, nil, "street_5", nil, 60.1, lo.ToPtr(6), "alias_5", lo.ToPtr("5-room"))

				m.ExpectQuery(query).WithArgs(12, 0).WillReturnRows(rows)
			},
			want:    expectedListPoster,
			wantErr: nil,
		},
		{
			name:   "empty_result",
			params: inputParams,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "price", "avatar_url", "address", "station_name", "area", "floor", "alias", "flat_category",
				})

				m.ExpectQuery(query).WithArgs(12, 0).WillReturnRows(rows)
			},
			want:    []entity.PosterFlat{},
			wantErr: nil,
		},
		{
			name:   "collect_rows_error",
			params: inputParams,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "price", "avatar_url", "address", "station_name", "area", "floor", "alias", "flat_category",
				}).AddRow("badId", 11111.0, nil, "street_1", nil, 35.5, lo.ToPtr(2), "alias_1", lo.ToPtr("1-room"))

				m.ExpectQuery(query).WithArgs(12, 0).WillReturnRows(rows)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewPosterRepo(mock)

			got, err := repo.GetFlatsAll(context.Background(), test.params)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetPosterByAliasRepo(t *testing.T) {
	expectedPoster := entity.PosterById{
		ID:               4,
		Alias:            "kvartira-na-arbate",
		Price:            135000,
		Category:         "flat",
		CategoryAlias:    "flat",
		Description:      "krutoy remont",
		PropertyID:       10,
		Area:             96.4,
		Address:          "Arbatskaya 5k2",
		District:         lo.ToPtr("Arbat"),
		Metro:            lo.ToPtr("Arbat"),
		Geo:              geo.GeographyPoint{Lon: 45.3966, Lat: 46.3489},
		MetroGeo:         &geo.GeographyPoint{Lon: 45.3966, Lat: 46.3489},
		City:             "Moscow",
		FloorCount:       7,
		SellerFirstName:  "Sanya",
		SellerLastName:   "Sashenykov",
		Phone:            "+79144564312",
		SellerAvatarURL:  lo.ToPtr("http://seller-avatar.com"),
		CompanyName:      lo.ToPtr("PIC"),
		CompanyAvatarURL: lo.ToPtr("http://img.com"),
		CompanyAlias:     lo.ToPtr("dick"),
		CompanyID:        lo.ToPtr(12),
		Images: []entity.PosterImage{
			{ImgURL: "img1.jpg", Order: 1},
			{ImgURL: "img2.jpg", Order: 2},
		},
		Facilities: []entity.Facility{
			{ID: 1, Name: "Балкон", Alias: "balcony"},
			{ID: 2, Name: "Кондиционер", Alias: "conditioner"},
		},
	}

	inputAlias := "kvartira-na-arbate"

	posterQuery := regexp.QuoteMeta(`
		SELECT p.id, p.alias, p.price, pc.name AS category_name, 
		pc.alias AS category_alias, p.description, prop.area, 
		prop.id AS property_id, ST_AsText(b.geo) AS building_geo, 
		b.address, b.district, ms.station_name, ST_AsText(ms.geo) AS metro_geo, 
		c.city_name, b.floor_count, pr.first_name, pr.last_name, pr.phone, 
		pr.avatar_url as seller_avatar_url, uc.company_name, 
		uc.avatar_url AS company_avatar_url, uc.alias AS company_alias, uc.id AS company_id
		FROM posters p
		JOIN property prop ON prop.id = p.property_id
		JOIN property_categories pc ON pc.id = prop.category_id
		JOIN buildings b ON b.id = prop.building_id
		JOIN cities c ON c.id = b.city_id
		LEFT JOIN metro_stations ms ON ms.id = b.metro_station_id
		LEFT JOIN utility_companies uc ON uc.id = b.company_id
		JOIN users u ON u.id = p.user_id
		JOIN profiles pr ON pr.id = u.profile_id
		WHERE p.alias = $1
	`)

	imagesQuery := regexp.QuoteMeta(`
		SELECT im.img_url, im.sequence_order
		FROM poster_photos im
		WHERE im.poster_id = $1
		ORDER BY im.sequence_order
	`)

	facilitiesQuery := regexp.QuoteMeta(`
		SELECT f.id, f.name, f.alias
		FROM facilities f 
		JOIN facility_property fp ON fp.facility_id = f.id
		JOIN property pr ON pr.id = fp.property_id
		JOIN posters pt ON pt.property_id = pr.id
		WHERE pt.id = $1
		ORDER BY f.name
	`)

	tests := []struct {
		name      string
		param     string
		setupMock func(m pgxmock.PgxPoolIface)
		want      *entity.PosterById
		wantErr   error
	}{
		{
			name:  "ok",
			param: inputAlias,
			setupMock: func(m pgxmock.PgxPoolIface) {
				posterRows := pgxmock.NewRows([]string{
					"id", "alias", "price", "category_name", "category_alias",
					"description", "area", "property_id",
					"building_geo", "address", "district",
					"station_name", "metro_geo", "city_name",
					"floor_count", "first_name", "last_name",
					"phone", "seller_avatar_url", "company_name",
					"company_avatar_url", "company_alias", "company_id",
				}).AddRow(
					4, "kvartira-na-arbate", 135000.0,
					"flat", "flat", "krutoy remont", 96.4, 10,
					"POINT(45.3966 46.3489)", "Arbatskaya 5k2",
					lo.ToPtr("Arbat"), lo.ToPtr("Arbat"),
					&geo.GeographyPoint{Lon: 45.3966, Lat: 46.3489}, "Moscow", 7, "Sanya", "Sashenykov",
					"+79144564312", lo.ToPtr("http://seller-avatar.com"),
					lo.ToPtr("PIC"), lo.ToPtr("http://img.com"), lo.ToPtr("dick"), lo.ToPtr(12),
				)

				imageRows := pgxmock.NewRows([]string{
					"img_url",
					"sequence_order",
				}).AddRow("img1.jpg", 1).
					AddRow("img2.jpg", 2)

				m.ExpectQuery(posterQuery).
					WithArgs(inputAlias).
					WillReturnRows(posterRows)

				m.ExpectQuery(imagesQuery).
					WithArgs(4).
					WillReturnRows(imageRows)

				facilityRows := pgxmock.NewRows([]string{
					"id", "name", "alias",
				}).AddRow(1, "Балкон", "balcony").
					AddRow(2, "Кондиционер", "conditioner")

				m.ExpectQuery(facilitiesQuery).
					WithArgs(4).
					WillReturnRows(facilityRows)
			},
			want:    &expectedPoster,
			wantErr: nil,
		},
		{
			name:  "not_found",
			param: inputAlias,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "alias", "price", "category_name", "category_alias",
					"description", "area", "property_id",
					"building_geo", "address", "district",
					"station_name", "metro_geo", "city_name",
					"floor_count", "first_name", "last_name",
					"phone", "seller_avatar_url", "company_name",
					"company_avatar_url", "company_alias", "company_id",
				})

				m.ExpectQuery(posterQuery).
					WithArgs(inputAlias).
					WillReturnRows(rows)
			},
			want:    &entity.PosterById{},
			wantErr: entity.NotFoundError,
		},
		{
			name:  "collect_poster_error",
			param: inputAlias,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "alias", "price", "category_name", "category_alias",
					"description", "area", "property_id",
					"building_geo", "address", "district",
					"station_name", "metro_geo", "city_name",
					"floor_count", "first_name", "last_name",
					"phone", "seller_avatar_url", "company_name",
					"company_avatar_url", "company_alias", "company_id",
				}).AddRow(
					"bad-id", "kvartira-na-arbate", 135000.0,
					"flat", "flat", "krutoy remont", 96.4, 10,
					"POINT(45.3966 46.3489)", "Arbatskaya 5k2",
					lo.ToPtr("Arbat"), lo.ToPtr("Arbat"),
					&geo.GeographyPoint{Lon: 45.3966, Lat: 46.3489}, "Moscow", 7, "Sanya", "Sashenykov",
					"+79144564312", lo.ToPtr("http://seller-avatar.com"),
					lo.ToPtr("PIC"), lo.ToPtr("http://img.com"), lo.ToPtr("dick"), lo.ToPtr(12),
				)

				m.ExpectQuery(posterQuery).
					WithArgs(inputAlias).
					WillReturnRows(rows)
			},
			want:    &entity.PosterById{},
			wantErr: entity.ServiceError,
		},
		{
			name:  "collect_image_error",
			param: inputAlias,
			setupMock: func(m pgxmock.PgxPoolIface) {
				posterRows := pgxmock.NewRows([]string{
					"id", "alias", "price", "category_name", "category_alias",
					"description", "area", "property_id",
					"building_geo", "address", "district",
					"station_name", "metro_geo", "city_name",
					"floor_count", "first_name", "last_name",
					"phone", "seller_avatar_url", "company_name",
					"company_avatar_url", "company_alias", "company_id",
				}).AddRow(
					4, "kvartira-na-arbate", 135000.0,
					"flat", "flat", "krutoy remont", 96.4, 10,
					"POINT(45.3966 46.3489)", "Arbatskaya 5k2",
					lo.ToPtr("Arbat"), lo.ToPtr("Arbat"),
					nil, "Moscow", 7, "Sanya", "Sashenykov", "+79144564312", lo.ToPtr("http://seller-avatar.com"),
					lo.ToPtr("PIC"), lo.ToPtr("http://img.com"), lo.ToPtr("dick"), lo.ToPtr(12),
				)

				imageRows := pgxmock.NewRows([]string{
					"img_url",
					"sequence_order",
				}).AddRow("img1.jpg", "bad_order")

				m.ExpectQuery(posterQuery).
					WithArgs(inputAlias).
					WillReturnRows(posterRows)

				m.ExpectQuery(imagesQuery).
					WithArgs(4).
					WillReturnRows(imageRows)
			},
			want:    &entity.PosterById{},
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewPosterRepo(mock)

			got, err := repo.GetByAlias(context.Background(), test.param, nil)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetFlatByPropetyIDRepo(t *testing.T) {
	expectedFlat := entity.Flat{
		PropertyID:   10,
		Number:       43,
		Floor:        3,
		FlatCategory: "1-room",
		RoomCount:    1,
	}

	inputPropertyID := 10

	query := regexp.QuoteMeta(`
		SELECT f.property_id, f.number, f.floor,
			   fc.name AS flat_category, fc.room_count
		FROM flat f
		LEFT JOIN flat_categories fc ON fc.id = f.category_id
		WHERE f.property_id = $1;
	`)

	tests := []struct {
		name      string
		param     int
		setupMock func(m pgxmock.PgxPoolIface)
		want      *entity.Flat
		wantErr   error
	}{
		{
			name:  "ok",
			param: inputPropertyID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"property_id", "number", "floor", "flat_category", "room_count",
				}).AddRow(
					10, 43, 3, "1-room", 1,
				)

				m.ExpectQuery(query).
					WithArgs(inputPropertyID).
					WillReturnRows(rows)
			},
			want:    &expectedFlat,
			wantErr: nil,
		},
		{
			name:  "not_found",
			param: inputPropertyID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"property_id", "number", "floor", "flat_category", "room_count",
				})

				m.ExpectQuery(query).
					WithArgs(inputPropertyID).
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: entity.NotFoundError,
		},
		{
			name:  "collect_flat_error",
			param: inputPropertyID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"property_id", "number", "floor", "flat_category", "room_count",
				}).AddRow(
					"bad-id", 43, 3, "1-room", 1,
				)

				m.ExpectQuery(query).
					WithArgs(inputPropertyID).
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewPosterRepo(mock)

			got, err := repo.GetFlatByPropetyID(context.Background(), test.param)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetPosterByUserIDRepo(t *testing.T) {
	expectedPosters := []entity.Poster{
		{
			ID:            1,
			Price:         100000,
			AvatarURl:     lo.ToPtr("img1.jpg"),
			Address:       "street_1",
			Area:          35.5,
			Alias:         "alias_1",
			CategoryName:  "Квартира",
			CategoryAlias: "flat",
		},
		{
			ID:            2,
			Price:         200000,
			AvatarURl:     lo.ToPtr("img2.jpg"),
			Address:       "street_2",
			Area:          45.0,
			Alias:         "alias_2",
			CategoryName:  "Квартира",
			CategoryAlias: "flat",
		},
	}

	inputUserID := 7

	query := regexp.QuoteMeta(`
		SELECT p.id, p.price, p.avatar_url,
               b.address, prop.area, p.alias, pc.name as category_name, pc.alias as category_alias
        FROM posters p
        JOIN property prop ON prop.id = p.property_id
        JOIN property_categories pc ON pc.id = prop.category_id
        JOIN buildings b ON b.id = prop.building_id
		WHERE p.user_id = $1;
	`)

	tests := []struct {
		name      string
		param     int
		setupMock func(m pgxmock.PgxPoolIface)
		want      []entity.Poster
		wantErr   error
	}{
		{
			name:  "ok",
			param: inputUserID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "price", "avatar_url", "address", "area", "alias", "category_name", "category_alias",
				}).AddRow(
					1, 100000.0, lo.ToPtr("img1.jpg"), "street_1", 35.5, "alias_1", "Квартира", "flat",
				).AddRow(
					2, 200000.0, lo.ToPtr("img2.jpg"), "street_2", 45.0, "alias_2", "Квартира", "flat",
				)

				m.ExpectQuery(query).
					WithArgs(inputUserID).
					WillReturnRows(rows)
			},
			want:    expectedPosters,
			wantErr: nil,
		},
		{
			name:  "not_found",
			param: inputUserID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "price", "avatar_url", "address", "area", "alias", "category_name", "category_alias",
				})

				m.ExpectQuery(query).
					WithArgs(inputUserID).
					WillReturnRows(rows)
			},
			want:    []entity.Poster{},
			wantErr: nil,
		},
		{
			name:  "collect_rows_error",
			param: inputUserID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "price", "avatar_url", "address", "area", "alias", "category_name", "category_alias",
				}).AddRow(
					"bad-id", 100000.0, lo.ToPtr("img1.jpg"), "street_1", 35.5, "alias_1", "Квартира", "flat",
				)

				m.ExpectQuery(query).
					WithArgs(inputUserID).
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewPosterRepo(mock)

			got, err := repo.GetByUserID(context.Background(), test.param)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetMetroStationByRadiusRepo(t *testing.T) {
	expectedStations := []entity.MetroStation{
		{
			ID:          1,
			StationName: "Арбатская",
			StationGEO:  geo.GeographyPoint{Lon: 37.605, Lat: 55.752},
		},
		{
			ID:          2,
			StationName: "Смоленская",
			StationGEO:  geo.GeographyPoint{Lon: 37.582, Lat: 55.748},
		},
	}

	inputGeo := dto.GeographyDTO{
		Lon: 37.6173,
		Lat: 55.7558,
	}

	inputRadius := entity.Metre(2500)

	query := regexp.QuoteMeta(`
    SELECT m.id, m.station_name, ST_AsText(m.geo) AS metro_geo
    FROM metro_stations m 
    WHERE ST_DWithin(m.geo, ST_GeogFromText($1), $2)
    ORDER BY m.geo <-> ST_GeogFromText($1)
`)

	tests := []struct {
		name        string
		paramGeo    dto.GeographyDTO
		paramRadius entity.Metre
		setupMock   func(m pgxmock.PgxPoolIface)
		want        []entity.MetroStation
		wantErr     error
	}{
		{
			name:        "ok",
			paramGeo:    inputGeo,
			paramRadius: inputRadius,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "station_name", "metro_geo",
				}).AddRow(
					1, "Арбатская", "POINT(37.605 55.752)",
				).AddRow(
					2, "Смоленская", "POINT(37.582 55.748)",
				)

				m.ExpectQuery(query).
					WithArgs(geo.GeographyPoint{Lat: inputGeo.Lat, Lon: inputGeo.Lon}, int(inputRadius)).
					WillReturnRows(rows)
			},
			want:    expectedStations,
			wantErr: nil,
		},
		{
			name:        "empty_result",
			paramGeo:    inputGeo,
			paramRadius: inputRadius,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "station_name", "metro_geo",
				})

				m.ExpectQuery(query).
					WithArgs(geo.GeographyPoint{Lat: inputGeo.Lat, Lon: inputGeo.Lon}, int(inputRadius)).
					WillReturnRows(rows)
			},
			want:    []entity.MetroStation{},
			wantErr: nil,
		},
		{
			name:        "collect_rows_error",
			paramGeo:    inputGeo,
			paramRadius: inputRadius,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "station_name", "metro_geo",
				}).AddRow(
					"bad-id", "Арбатская", "POINT(37.605 55.752)",
				)

				m.ExpectQuery(query).
					WithArgs(geo.GeographyPoint{Lat: inputGeo.Lat, Lon: inputGeo.Lon}, int(inputRadius)).
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewPosterRepo(mock)

			got, err := repo.GetMetroStationByRadius(context.Background(), test.paramGeo, test.paramRadius)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCreateBuildingRepo(t *testing.T) {
	inputPoster := &dto.PosterInput{
		Address:        "Arbatskaya 5k2",
		Geo:            geo.GeographyPoint{Lon: 37.6173, Lat: 55.7558},
		CityID:         1,
		MetroStationID: lo.ToPtr(5),
		District:       lo.ToPtr("Arbat"),
		FloorCount:     7,
		CompanyID:      lo.ToPtr(12),
	}

	query := regexp.QuoteMeta(`
		INSERT INTO buildings (address, geo, city_id,
			metro_station_id, district, floor_count, company_id)
		VALUES ($1, ST_GeogFromText($2), $3, $4, $5, $6, $7)
		RETURNING id
	`)

	tests := []struct {
		name      string
		param     *dto.PosterInput
		setupMock func(m pgxmock.PgxPoolIface)
		want      int
		wantErr   error
	}{
		{
			name:  "ok",
			param: inputPoster,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id"}).AddRow(11)

				m.ExpectQuery(query).
					WithArgs(
						inputPoster.Address,
						inputPoster.Geo.ToGeo(),
						inputPoster.CityID,
						inputPoster.MetroStationID,
						inputPoster.District,
						inputPoster.FloorCount,
						inputPoster.CompanyID,
					).
					WillReturnRows(rows)
			},
			want:    11,
			wantErr: nil,
		},
		{
			name:  "query_error",
			param: inputPoster,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(query).
					WithArgs(
						inputPoster.Address,
						inputPoster.Geo.ToGeo(),
						inputPoster.CityID,
						inputPoster.MetroStationID,
						inputPoster.District,
						inputPoster.FloorCount,
						inputPoster.CompanyID,
					).
					WillReturnError(entity.ServiceError)
			},
			want:    0,
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewPosterRepo(mock)

			got, err := repo.CreateBuilding(context.Background(), test.param)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCreatePropertyRepo(t *testing.T) {
	inputPoster := &dto.PosterInput{
		CategoryAlias: "flat",
		Area:          96.4,
	}

	inputBuildingID := 11

	query := regexp.QuoteMeta(`
		INSERT INTO property (category_id, building_id, area)
		SELECT pc.id, $2, $3
		FROM property_categories AS pc
		WHERE pc.alias = $1
		RETURNING property.id;
	`)

	tests := []struct {
		name       string
		param      *dto.PosterInput
		buildingID int
		setupMock  func(m pgxmock.PgxPoolIface)
		want       int
		wantErr    error
	}{
		{
			name:       "ok",
			param:      inputPoster,
			buildingID: inputBuildingID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id"}).AddRow(22)

				m.ExpectQuery(query).
					WithArgs(
						inputPoster.CategoryAlias,
						inputBuildingID,
						inputPoster.Area,
					).
					WillReturnRows(rows)
			},
			want:    22,
			wantErr: nil,
		},
		{
			name:       "query_error",
			param:      inputPoster,
			buildingID: inputBuildingID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(query).
					WithArgs(
						inputPoster.CategoryAlias,
						inputBuildingID,
						inputPoster.Area,
					).
					WillReturnError(entity.ServiceError)
			},
			want:    0,
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewPosterRepo(mock)

			got, err := repo.CreateProperty(context.Background(), test.param, test.buildingID)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestInsertFacilitiesRepo(t *testing.T) {
	inputPropertyID := 22
	inputAliases := []string{"balcony", "parking"}

	selectQuery := regexp.QuoteMeta(`
		SELECT id
		FROM facilities
		WHERE alias = ANY($1::text[])
	`)

	insertQuery := regexp.QuoteMeta(`
		INSERT INTO facility_property (property_id, facility_id)
		VALUES ($1, $2)
	`)

	tests := []struct {
		name       string
		propertyID int
		aliases    []string
		setupMock  func(m pgxmock.PgxPoolIface)
		wantErr    error
	}{
		{
			name:       "ok",
			propertyID: inputPropertyID,
			aliases:    inputAliases,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id"}).
					AddRow(1).
					AddRow(2)

				m.ExpectQuery(selectQuery).
					WithArgs(inputAliases).
					WillReturnRows(rows)

				m.ExpectExec(insertQuery).
					WithArgs(inputPropertyID, 1).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))

				m.ExpectExec(insertQuery).
					WithArgs(inputPropertyID, 2).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			wantErr: nil,
		},
		{
			name:       "select_error",
			propertyID: inputPropertyID,
			aliases:    inputAliases,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(selectQuery).
					WithArgs(inputAliases).
					WillReturnError(entity.ServiceError)
			},
			wantErr: entity.ServiceError,
		},
		{
			name:       "collect_rows_error",
			propertyID: inputPropertyID,
			aliases:    inputAliases,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id"}).
					AddRow("bad-id")

				m.ExpectQuery(selectQuery).
					WithArgs(inputAliases).
					WillReturnRows(rows)
			},
			wantErr: entity.ServiceError,
		},
		{
			name:       "insert_error",
			propertyID: inputPropertyID,
			aliases:    inputAliases,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id"}).
					AddRow(1)

				m.ExpectQuery(selectQuery).
					WithArgs(inputAliases).
					WillReturnRows(rows)

				m.ExpectExec(insertQuery).
					WithArgs(inputPropertyID, 1).
					WillReturnError(entity.ServiceError)
			},
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewPosterRepo(mock)

			err = repo.InsertFacilities(context.Background(), test.propertyID, test.aliases)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCreatePosterRepo(t *testing.T) {
	inputPoster := &dto.PosterInput{
		Price:       135000,
		Description: "krutoy remont",
		UserID:      7,
		Alias:       "kvartira-na-arbate",
	}

	inputPropertyID := 22

	query := regexp.QuoteMeta(`
		INSERT INTO posters (price, description,
			user_id, property_id, alias)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`)

	tests := []struct {
		name       string
		param      *dto.PosterInput
		propertyID int
		setupMock  func(m pgxmock.PgxPoolIface)
		want       int
		wantErr    error
	}{
		{
			name:       "ok",
			param:      inputPoster,
			propertyID: inputPropertyID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id"}).AddRow(33)

				m.ExpectQuery(query).
					WithArgs(
						inputPoster.Price,
						inputPoster.Description,
						inputPoster.UserID,
						inputPropertyID,
						inputPoster.Alias,
					).
					WillReturnRows(rows)
			},
			want:    33,
			wantErr: nil,
		},
		{
			name:       "query_error",
			param:      inputPoster,
			propertyID: inputPropertyID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(query).
					WithArgs(
						inputPoster.Price,
						inputPoster.Description,
						inputPoster.UserID,
						inputPropertyID,
						inputPoster.Alias,
					).
					WillReturnError(entity.ServiceError)
			},
			want:    0,
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewPosterRepo(mock)

			got, err := repo.Create(context.Background(), test.param, test.propertyID)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestInsertFlatRepo(t *testing.T) {
	number := 43

	inputFlat := &dto.FlatInput{
		PropertyID: 22,
		Floor:      3,
		Number:     &number,
		CategoryID: 1,
	}

	query := regexp.QuoteMeta(`
		INSERT INTO flat (property_id,
			floor, number, category_id)
		VALUES ($1, $2, $3, $4)
		RETURNING property_id
	`)

	tests := []struct {
		name      string
		param     *dto.FlatInput
		setupMock func(m pgxmock.PgxPoolIface)
		wantErr   error
	}{
		{
			name:  "ok",
			param: inputFlat,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"property_id"}).AddRow(22)

				m.ExpectQuery(query).
					WithArgs(
						inputFlat.PropertyID,
						inputFlat.Floor,
						inputFlat.Number,
						inputFlat.CategoryID,
					).
					WillReturnRows(rows)
			},
			wantErr: nil,
		},
		{
			name:  "query_error",
			param: inputFlat,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(query).
					WithArgs(
						inputFlat.PropertyID,
						inputFlat.Floor,
						inputFlat.Number,
						inputFlat.CategoryID,
					).
					WillReturnError(entity.ServiceError)
			},
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewPosterRepo(mock)

			err = repo.InsertFlat(context.Background(), test.param)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestInsertPhotosRepo(t *testing.T) {
	inputPosterID := 33

	inputPhotos := []dto.PhotoInput{
		{
			Path:  "/poster/img/flat/img1.jpg",
			Order: 1,
		},
		{
			Path:  "/poster/img/flat/img2.jpg",
			Order: 2,
		},
	}

	query := regexp.QuoteMeta(`
		INSERT INTO poster_photos (img_url, sequence_order, poster_id)
		VALUES ($1, $2, $3),($4, $5, $6)`)

	tests := []struct {
		name      string
		posterID  int
		photos    []dto.PhotoInput
		setupMock func(m pgxmock.PgxPoolIface)
		wantErr   error
	}{
		{
			name:     "ok",
			posterID: inputPosterID,
			photos:   inputPhotos,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(query).
					WithArgs(
						"/poster/img/flat/img1.jpg", 1, inputPosterID,
						"/poster/img/flat/img2.jpg", 2, inputPosterID,
					).
					WillReturnResult(pgxmock.NewResult("INSERT", 2))
			},
			wantErr: nil,
		},
		{
			name:     "empty_photos",
			posterID: inputPosterID,
			photos:   []dto.PhotoInput{},
			setupMock: func(m pgxmock.PgxPoolIface) {
			},
			wantErr: nil,
		},
		{
			name:     "exec_error",
			posterID: inputPosterID,
			photos:   inputPhotos,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(query).
					WithArgs(
						"/poster/img/flat/img1.jpg", 1, inputPosterID,
						"/poster/img/flat/img2.jpg", 2, inputPosterID,
					).
					WillReturnError(entity.ServiceError)
			},
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewPosterRepo(mock)

			err = repo.InsertPhotos(context.Background(), test.posterID, test.photos)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestInsertMainPhotoRepo(t *testing.T) {
	inputPosterID := 33
	inputURL := "/poster/img/flat/main.jpg"

	query := regexp.QuoteMeta(`
		UPDATE posters
		SET avatar_url = $1
		WHERE id = $2
	`)

	tests := []struct {
		name      string
		posterID  int
		url       string
		setupMock func(m pgxmock.PgxPoolIface)
		wantErr   error
	}{
		{
			name:     "ok",
			posterID: inputPosterID,
			url:      inputURL,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(query).
					WithArgs(inputURL, inputPosterID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			wantErr: nil,
		},
		{
			name:     "exec_error",
			posterID: inputPosterID,
			url:      inputURL,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(query).
					WithArgs(inputURL, inputPosterID).
					WillReturnError(entity.ServiceError)
			},
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewPosterRepo(mock)

			err = repo.InsertMainPhoto(context.Background(), test.posterID, test.url)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetUpdateIDsByAliasRepo(t *testing.T) {
	expectedIDs := dto.PosterUpdateIDs{
		PosterID:   33,
		UserID:     7,
		PropertyID: 22,
		BuildingID: 11,
	}

	inputAlias := "kvartira-na-arbate"

	query := regexp.QuoteMeta(`
		SELECT p.id, p.user_id, p.property_id, pr.building_id
		FROM posters p
		JOIN property pr ON pr.id = p.property_id
		WHERE p.alias = $1
	`)

	tests := []struct {
		name      string
		param     string
		setupMock func(m pgxmock.PgxPoolIface)
		want      *dto.PosterUpdateIDs
		wantErr   error
	}{
		{
			name:  "ok",
			param: inputAlias,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "user_id", "property_id", "building_id",
				}).AddRow(
					33, 7, 22, 11,
				)

				m.ExpectQuery(query).
					WithArgs(inputAlias).
					WillReturnRows(rows)
			},
			want:    &expectedIDs,
			wantErr: nil,
		},
		{
			name:  "not_found",
			param: inputAlias,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "user_id", "property_id", "building_id",
				})

				m.ExpectQuery(query).
					WithArgs(inputAlias).
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: entity.NotFoundError,
		},
		{
			name:  "query_error",
			param: inputAlias,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(query).
					WithArgs(inputAlias).
					WillReturnError(entity.ServiceError)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewPosterRepo(mock)

			got, err := repo.GetUpdateIDsByAlias(context.Background(), test.param)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetCityByNameRepo(t *testing.T) {
	expectedCity := entity.City{
		ID:   1,
		Name: "Moscow",
	}

	inputName := "Moscow"

	query := regexp.QuoteMeta(`
		SELECT c.id, c.city_name FROM cities c WHERE lower(c.city_name) = lower($1)`)

	tests := []struct {
		name      string
		param     string
		setupMock func(m pgxmock.PgxPoolIface)
		want      *entity.City
		wantErr   error
	}{
		{
			name:  "ok",
			param: inputName,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "city_name",
				}).AddRow(
					1, "Moscow",
				)

				m.ExpectQuery(query).
					WithArgs(inputName).
					WillReturnRows(rows)
			},
			want:    &expectedCity,
			wantErr: nil,
		},
		{
			name:  "not_found",
			param: inputName,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "city_name",
				})

				m.ExpectQuery(query).
					WithArgs(inputName).
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: entity.NotFoundError,
		},
		{
			name:  "collect_error",
			param: inputName,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "city_name",
				}).AddRow(
					"bad-id", "Moscow",
				)

				m.ExpectQuery(query).
					WithArgs(inputName).
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
		{
			name:  "query_error",
			param: inputName,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(query).
					WithArgs(inputName).
					WillReturnError(entity.ServiceError)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewPosterRepo(mock)

			got, err := repo.GetCityByName(context.Background(), test.param)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdatePosterRepo(t *testing.T) {
	inputPoster := &dto.PosterInput{
		Price:       150000,
		Description: "updated description",
	}

	inputPosterID := 33

	query := regexp.QuoteMeta(`
		UPDATE posters
		SET price = $1, description = $2
		WHERE id = $3
	`)

	tests := []struct {
		name      string
		posterID  int
		poster    *dto.PosterInput
		setupMock func(m pgxmock.PgxPoolIface)
		wantErr   error
	}{
		{
			name:     "ok",
			posterID: inputPosterID,
			poster:   inputPoster,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(query).
					WithArgs(
						inputPoster.Price,
						inputPoster.Description,
						inputPosterID,
					).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			wantErr: nil,
		},
		{
			name:     "exec_error",
			posterID: inputPosterID,
			poster:   inputPoster,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(query).
					WithArgs(
						inputPoster.Price,
						inputPoster.Description,
						inputPosterID,
					).
					WillReturnError(entity.ServiceError)
			},
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewPosterRepo(mock)

			err = repo.Update(context.Background(), test.posterID, test.poster)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdatePropertyRepo(t *testing.T) {
	inputPoster := &dto.PosterInput{
		CategoryAlias: "flat",
		Area:          96.4,
	}

	inputPropertyID := 22

	query := regexp.QuoteMeta(`
		UPDATE property
		SET category_id = (
				SELECT id
				FROM property_categories
				WHERE alias = $1
			),
			area = $2
		WHERE id = $3
	`)

	tests := []struct {
		name       string
		propertyID int
		poster     *dto.PosterInput
		setupMock  func(m pgxmock.PgxPoolIface)
		wantErr    error
	}{
		{
			name:       "ok",
			propertyID: inputPropertyID,
			poster:     inputPoster,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(query).
					WithArgs(
						inputPoster.CategoryAlias,
						inputPoster.Area,
						inputPropertyID,
					).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			wantErr: nil,
		},
		{
			name:       "exec_error",
			propertyID: inputPropertyID,
			poster:     inputPoster,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(query).
					WithArgs(
						inputPoster.CategoryAlias,
						inputPoster.Area,
						inputPropertyID,
					).
					WillReturnError(entity.ServiceError)
			},
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewPosterRepo(mock)

			err = repo.UpdateProperty(context.Background(), test.propertyID, test.poster)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdateBuildingRepo(t *testing.T) {
	inputPoster := &dto.PosterInput{
		Address:        "Arbatskaya 5k2",
		Geo:            geo.GeographyPoint{Lon: 37.6173, Lat: 55.7558},
		CityID:         1,
		MetroStationID: lo.ToPtr(5),
		District:       lo.ToPtr("Arbat"),
		FloorCount:     7,
		CompanyID:      lo.ToPtr(12),
	}

	inputBuildingID := 11

	query := regexp.QuoteMeta(`
		UPDATE buildings
		SET address = $1, geo = $2,
			city_id = $3, metro_station_id = $4,
			district = $5, floor_count = $6,
			company_id = $7
		WHERE id = $8
	`)

	tests := []struct {
		name       string
		buildingID int
		poster     *dto.PosterInput
		setupMock  func(m pgxmock.PgxPoolIface)
		wantErr    error
	}{
		{
			name:       "ok",
			buildingID: inputBuildingID,
			poster:     inputPoster,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(query).
					WithArgs(
						inputPoster.Address,
						inputPoster.Geo,
						inputPoster.CityID,
						inputPoster.MetroStationID,
						inputPoster.District,
						inputPoster.FloorCount,
						inputPoster.CompanyID,
						inputBuildingID,
					).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			wantErr: nil,
		},
		{
			name:       "exec_error",
			buildingID: inputBuildingID,
			poster:     inputPoster,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(query).
					WithArgs(
						inputPoster.Address,
						inputPoster.Geo,
						inputPoster.CityID,
						inputPoster.MetroStationID,
						inputPoster.District,
						inputPoster.FloorCount,
						inputPoster.CompanyID,
						inputBuildingID,
					).
					WillReturnError(entity.ServiceError)
			},
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewPosterRepo(mock)

			err = repo.UpdateBuilding(context.Background(), test.buildingID, test.poster)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdateFlatRepo(t *testing.T) {
	number := 43

	inputFlat := &dto.FlatInput{
		PropertyID: 22,
		Floor:      3,
		Number:     &number,
		CategoryID: 1,
	}

	query := regexp.QuoteMeta(`
		UPDATE flat
		SET floor = $1, number = $2, category_id = $3
		WHERE property_id = $4
	`)

	tests := []struct {
		name      string
		flat      *dto.FlatInput
		setupMock func(m pgxmock.PgxPoolIface)
		wantErr   error
	}{
		{
			name: "ok",
			flat: inputFlat,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(query).
					WithArgs(
						inputFlat.Floor,
						inputFlat.Number,
						inputFlat.CategoryID,
						inputFlat.PropertyID,
					).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			wantErr: nil,
		},
		{
			name: "exec_error",
			flat: inputFlat,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(query).
					WithArgs(
						inputFlat.Floor,
						inputFlat.Number,
						inputFlat.CategoryID,
						inputFlat.PropertyID,
					).
					WillReturnError(entity.ServiceError)
			},
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewPosterRepo(mock)

			err = repo.UpdateFlat(context.Background(), test.flat)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetPhotoPathsByPosterIDRepo(t *testing.T) {
	expectedPaths := []string{
		"/poster/img/flat/img1.jpg",
		"/poster/img/flat/img2.jpg",
	}

	inputPosterID := 33

	query := regexp.QuoteMeta(`
		SELECT img_url
		FROM poster_photos
		WHERE poster_id = $1
	`)

	tests := []struct {
		name      string
		param     int
		setupMock func(m pgxmock.PgxPoolIface)
		want      []string
		wantErr   error
	}{
		{
			name:  "ok",
			param: inputPosterID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"img_url"}).
					AddRow("/poster/img/flat/img1.jpg").
					AddRow("/poster/img/flat/img2.jpg")

				m.ExpectQuery(query).
					WithArgs(inputPosterID).
					WillReturnRows(rows)
			},
			want:    expectedPaths,
			wantErr: nil,
		},
		{
			name:  "empty_result",
			param: inputPosterID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"img_url"})

				m.ExpectQuery(query).
					WithArgs(inputPosterID).
					WillReturnRows(rows)
			},
			want:    []string{},
			wantErr: nil,
		},
		{
			name:  "collect_error",
			param: inputPosterID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"img_url"}).
					AddRow(123) // invalid type

				m.ExpectQuery(query).
					WithArgs(inputPosterID).
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
		{
			name:  "query_error",
			param: inputPosterID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(query).
					WithArgs(inputPosterID).
					WillReturnError(entity.ServiceError)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewPosterRepo(mock)

			got, err := repo.GetPhotoPathsByPosterID(context.Background(), test.param)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDeleteFacilitiesByPropertyIDRepo(t *testing.T) {
	inputPropertyID := 22

	query := regexp.QuoteMeta(`
		DELETE FROM facility_property
		WHERE property_id = $1
	`)

	tests := []struct {
		name       string
		propertyID int
		setupMock  func(m pgxmock.PgxPoolIface)
		wantErr    error
	}{
		{
			name:       "ok",
			propertyID: inputPropertyID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(query).
					WithArgs(inputPropertyID).
					WillReturnResult(pgxmock.NewResult("DELETE", 2))
			},
			wantErr: nil,
		},
		{
			name:       "exec_error",
			propertyID: inputPropertyID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(query).
					WithArgs(inputPropertyID).
					WillReturnError(entity.ServiceError)
			},
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewPosterRepo(mock)

			err = repo.DeleteFacilitiesByPropertyID(context.Background(), test.propertyID)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDeletePhotosByPosterIDRepo(t *testing.T) {
	inputPosterID := 33

	deletePhotosQuery := regexp.QuoteMeta(`
		DELETE FROM poster_photos
		WHERE poster_id = $1
	`)

	clearAvatarQuery := regexp.QuoteMeta(`
		UPDATE posters
		SET avatar_url = NULL
		WHERE id = $1
	`)

	tests := []struct {
		name      string
		posterID  int
		setupMock func(m pgxmock.PgxPoolIface)
		wantErr   error
	}{
		{
			name:     "ok",
			posterID: inputPosterID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(deletePhotosQuery).
					WithArgs(inputPosterID).
					WillReturnResult(pgxmock.NewResult("DELETE", 2))

				m.ExpectExec(clearAvatarQuery).
					WithArgs(inputPosterID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			wantErr: nil,
		},
		{
			name:     "delete_photos_error",
			posterID: inputPosterID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(deletePhotosQuery).
					WithArgs(inputPosterID).
					WillReturnError(entity.ServiceError)
			},
			wantErr: entity.ServiceError,
		},
		{
			name:     "clear_avatar_error",
			posterID: inputPosterID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(deletePhotosQuery).
					WithArgs(inputPosterID).
					WillReturnResult(pgxmock.NewResult("DELETE", 2))

				m.ExpectExec(clearAvatarQuery).
					WithArgs(inputPosterID).
					WillReturnError(entity.ServiceError)
			},
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewPosterRepo(mock)

			err = repo.DeletePhotosByPosterID(context.Background(), test.posterID)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCreateCityRepo(t *testing.T) {
	expectedCity := entity.City{
		ID:   1,
		Name: "Moscow",
	}

	inputName := "Moscow"

	query := regexp.QuoteMeta(`
		INSERT INTO cities (city_name) VALUES ($1) RETURNING id, city_name`)

	tests := []struct {
		name      string
		param     string
		setupMock func(m pgxmock.PgxPoolIface)
		want      *entity.City
		wantErr   error
	}{
		{
			name:  "ok",
			param: inputName,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "city_name",
				}).AddRow(
					1, "Moscow",
				)

				m.ExpectQuery(query).
					WithArgs(inputName).
					WillReturnRows(rows)
			},
			want:    &expectedCity,
			wantErr: nil,
		},
		{
			name:  "query_error",
			param: inputName,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(query).
					WithArgs(inputName).
					WillReturnError(entity.ServiceError)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
		{
			name:  "collect_error",
			param: inputName,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "city_name",
				}).AddRow(
					"bad-id", "Moscow",
				)

				m.ExpectQuery(query).
					WithArgs(inputName).
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
		{
			name:  "not_found",
			param: inputName,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "city_name",
				})

				m.ExpectQuery(query).
					WithArgs(inputName).
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: entity.NotFoundError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewPosterRepo(mock)

			got, err := repo.CreateCity(context.Background(), test.param)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDeletePosterRepo(t *testing.T) {
	inputPosterID := 33

	query := regexp.QuoteMeta(`
        UPDATE posters SET deleted_at = NOW()
        WHERE id = $1
    `)

	tests := []struct {
		name      string
		posterID  int
		setupMock func(m pgxmock.PgxPoolIface)
		wantErr   error
	}{
		{
			name:     "ok",
			posterID: inputPosterID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(query).
					WithArgs(inputPosterID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			wantErr: nil,
		},
		{
			name:     "exec_error",
			posterID: inputPosterID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(query).
					WithArgs(inputPosterID).
					WillReturnError(entity.ServiceError)
			},
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewPosterRepo(mock)

			err = repo.Delete(context.Background(), test.posterID)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDeleteFlatRepo(t *testing.T) {
	inputPropertyID := 22

	query := regexp.QuoteMeta(`
        DELETE FROM flat 
        WHERE property_id = $1
    `)

	tests := []struct {
		name       string
		propertyID int
		setupMock  func(m pgxmock.PgxPoolIface)
		wantErr    error
	}{
		{
			name:       "ok",
			propertyID: inputPropertyID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(query).
					WithArgs(inputPropertyID).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
			wantErr: nil,
		},
		{
			name:       "exec_error",
			propertyID: inputPropertyID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(query).
					WithArgs(inputPropertyID).
					WillReturnError(entity.ServiceError)
			},
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewPosterRepo(mock)

			err = repo.DeleteFlat(context.Background(), test.propertyID)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDeletePropertyRepo(t *testing.T) {
	inputPropertyID := 22

	query := regexp.QuoteMeta(`
        DELETE FROM property 
        WHERE id = $1
    `)

	tests := []struct {
		name       string
		propertyID int
		setupMock  func(m pgxmock.PgxPoolIface)
		wantErr    error
	}{
		{
			name:       "ok",
			propertyID: inputPropertyID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(query).
					WithArgs(inputPropertyID).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
			wantErr: nil,
		},
		{
			name:       "exec_error",
			propertyID: inputPropertyID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(query).
					WithArgs(inputPropertyID).
					WillReturnError(entity.ServiceError)
			},
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewPosterRepo(mock)

			err = repo.DeleteProperty(context.Background(), test.propertyID)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDeleteBuildingRepo(t *testing.T) {
	inputBuildingID := 11

	query := regexp.QuoteMeta(`
        DELETE FROM buildings 
        WHERE id = $1
    `)

	tests := []struct {
		name       string
		buildingID int
		setupMock  func(m pgxmock.PgxPoolIface)
		wantErr    error
	}{
		{
			name:       "ok",
			buildingID: inputBuildingID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(query).
					WithArgs(inputBuildingID).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
			wantErr: nil,
		},
		{
			name:       "exec_error",
			buildingID: inputBuildingID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(query).
					WithArgs(inputBuildingID).
					WillReturnError(entity.ServiceError)
			},
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewPosterRepo(mock)

			err = repo.DeleteBuilding(context.Background(), test.buildingID)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
