package poster

import (
	"context"
	"regexp"
	"testing"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/geo"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestGetPostersRepo(t *testing.T) {
	expectedListPoster := []entity.Poster{
		{ID: 1, Price: 11111, ImgURL: nil, Address: "street_1", Metro: nil, Area: 35.5, Floor: 2},
		{ID: 2, Price: 22222, ImgURL: nil, Address: "street_2", Metro: nil, Area: 40.0, Floor: 3},
		{ID: 3, Price: 33333, ImgURL: nil, Address: "street_3", Metro: nil, Area: 45.2, Floor: 4},
		{ID: 4, Price: 44444, ImgURL: nil, Address: "street_4", Metro: nil, Area: 50.7, Floor: 5},
		{ID: 5, Price: 55555, ImgURL: nil, Address: "street_5", Metro: nil, Area: 60.1, Floor: 6},
	}

	inputParams := dto.PostersFiltersDTO{
		Limit:  12,
		Offset: 0,
	}

	query := regexp.QuoteMeta(`
        SELECT p.id, p.price, p.avatar_url, 
               b.address, m.station_name, a.area, a.floor
        FROM posters p
        JOIN apartments a ON a.id = p.apartment_id
        JOIN buildings b ON b.id = a.building_id
        JOIN metro_stations m ON b.metro_station_id = m.id
		ORDER BY p.created_at DESC 
		LIMIT $1 OFFSET $2`)

	tests := []struct {
		name      string
		params    dto.PostersFiltersDTO
		setupMock func(m pgxmock.PgxPoolIface)
		want      []entity.Poster
		wantErr   error
	}{
		{
			name:   "OK",
			params: inputParams,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "price", "avatar_url", "address", "station_name", "area", "floor",
				}).AddRow(1, 11111.0, nil, "street_1", nil, 35.5, 2).
					AddRow(2, 22222.0, nil, "street_2", nil, 40.0, 3).
					AddRow(3, 33333.0, nil, "street_3", nil, 45.2, 4).
					AddRow(4, 44444.0, nil, "street_4", nil, 50.7, 5).
					AddRow(5, 55555.0, nil, "street_5", nil, 60.1, 6)

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
					"id", "price", "avatar_url", "address", "station_name", "area", "floor",
				})

				m.ExpectQuery(query).WithArgs(12, 0).WillReturnRows(rows)
			},
			want:    []entity.Poster{},
			wantErr: nil,
		},
		{
			name:   "collect_rows_error",
			params: inputParams,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "price", "avatar_url", "address", "station_name", "area", "floor",
				}).AddRow(1, "bad_price", nil, "street_1", nil, 35.5, 2)

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

			got, err := repo.GetPosters(context.Background(), test.params)
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
		CompanyName:      lo.ToPtr("PIC"),
		CompanyAvatarURL: lo.ToPtr("http://img.com"),
		CompanyAlias:     lo.ToPtr("dick"),
		CompanyID:        lo.ToPtr(12),
		Images: []entity.PosterImage{
			{ImgURL: "img1.jpg", Order: 1},
			{ImgURL: "img2.jpg", Order: 2},
		},
	}

	inputAlias := "kvartira-na-arbate"

	posterQuery := regexp.QuoteMeta(`
		SELECT p.id, p.alias, p.price, pc.name AS category,
			   p.description, prop.area, prop.id AS property_id, 
			   ST_AsText(b.geo) AS building_geo, b.address, b.district, 
			   ms.station_name, ST_AsText(ms.geo) AS metro_geo, c.city_name, 
			   b.floor_count, pr.first_name, pr.last_name, pr.phone,
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
	`)

	imagesQuery := regexp.QuoteMeta(`
		SELECT im.img_url, im.sequence_order
		FROM poster_photos im
		WHERE im.poster_id = $1
		ORDER BY im.sequence_order
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
					"id", "alias", "price", "category",
					"description", "area", "property_id",
					"building_geo", "address", "district",
					"station_name", "metro_geo", "city_name",
					"floor_count", "first_name", "last_name",
					"phone", "company_name", "company_avatar_url", "company_alias", "company_id",
				}).AddRow(
					4, "kvartira-na-arbate", 135000.0,
					"flat", "krutoy remont", 96.4, 10,
					"POINT(45.3966 46.3489)", "Arbatskaya 5k2",
					lo.ToPtr("Arbat"), lo.ToPtr("Arbat"),
					lo.ToPtr("Arbat"), lo.ToPtr("Arbat"),
					&geo.GeographyPoint{Lon: 45.3966, Lat: 46.3489}, "Moscow", 7, "Sanya", "Sashenykov",
					"+79144564312", lo.ToPtr("PIC"), lo.ToPtr("http://img.com"), lo.ToPtr("dick"), lo.ToPtr(12),
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
			},
			want:    expectedPoster,
			wantErr: nil,
		},
		{
			name:  "not_found",
			param: inputAlias,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "alias", "price", "category",
					"description", "area", "property_id",
					"building_geo", "address", "district",
					"station_name", "metro_geo", "city_name",
					"floor_count", "first_name", "last_name",
					"phone", "company_name", "company_avatar_url", "company_alias", "company_id",
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
					"id", "alias", "price", "category",
					"description", "area", "property_id",
					"building_geo", "address", "district",
					"station_name", "metro_geo", "city_name",
					"floor_count", "first_name", "last_name",
					"phone", "company_name", "company_avatar_url", "company_alias", "company_id",
				}).AddRow(
					"bad-id", "kvartira-na-arbate", 135000.0,
					"flat", "krutoy remont", 96.4, 10,
					"POINT(45.3966 46.3489)", "Arbatskaya 5k2",
					lo.ToPtr("Arbat"), lo.ToPtr("Arbat"),
					lo.ToPtr("Arbat"), lo.ToPtr("Arbat"),
					&geo.GeographyPoint{Lon: 45.3966, Lat: 46.3489}, "Moscow", 7, "Sanya", "Sashenykov", "+79144564312",
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
					"id", "alias", "price", "category",
					"description", "area", "property_id",
					"building_geo", "address", "district",
					"station_name", "metro_geo", "city_name",
					"floor_count", "first_name", "last_name",
					"phone", "company_name", "company_avatar_url", "company_alias", "company_id",
				}).AddRow(
					4, "kvartira-na-arbate", 135000.0,
					"flat", "krutoy remont", 96.4, 10,
					"POINT(45.3966 46.3489)", "Arbatskaya 5k2",
					lo.ToPtr("Arbat"), lo.ToPtr("Arbat"),
					lo.ToPtr("Arbat"), lo.ToPtr("Arbat"),
					nil, "Moscow", 7, "Sanya", "Sashenykov", "+79144564312",
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

			got, err := repo.GetPosterByAlias(context.Background(), test.param)
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
