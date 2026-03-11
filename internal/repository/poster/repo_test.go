package poster

import (
	"context"
	"regexp"
	"testing"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"
)

func TestGetPostersRepo(t *testing.T) {
	expectedListPoster := []entity.Poster{
		{Id: 1, Price: 11111, ImgURL: nil, Address: "street_1", Metro: nil, Area: 35.5, Floor: 2},
		{Id: 2, Price: 22222, ImgURL: nil, Address: "street_2", Metro: nil, Area: 40.0, Floor: 3},
		{Id: 3, Price: 33333, ImgURL: nil, Address: "street_3", Metro: nil, Area: 45.2, Floor: 4},
		{Id: 4, Price: 44444, ImgURL: nil, Address: "street_4", Metro: nil, Area: 50.7, Floor: 5},
		{Id: 5, Price: 55555, ImgURL: nil, Address: "street_5", Metro: nil, Area: 60.1, Floor: 6},
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
