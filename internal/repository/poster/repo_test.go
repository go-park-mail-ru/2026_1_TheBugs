package poster

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"
)

func TestGetPostersRepo(t *testing.T) {
	expectedListPoster := []entity.Poster{
		{Id: 1, Price: 11111, Address: "street_1", Area: 35.5, Floor: 2, Type: "flat"},
		{Id: 2, Price: 22222, Address: "street_2", Area: 40.0, Floor: 3, Type: "flat"},
		{Id: 3, Price: 33333, Address: "street_3", Area: 45.2, Floor: 4, Type: "studio"},
		{Id: 4, Price: 44444, Address: "street_4", Area: 50.7, Floor: 5, Type: "flat"},
		{Id: 5, Price: 55555, Address: "street_5", Area: 60.1, Floor: 6, Type: "apartment"},
	}

	tests := []struct {
		name      string
		params    dto.PostersFiltersDTO
		setupMock func(m pgxmock.PgxPoolIface)
		want      []entity.Poster
		wantErr   error
	}{
		{
			name: "OK",
			params: dto.PostersFiltersDTO{
				Limit:  12,
				Offset: 0,
			},
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "price", "image_url", "address", "metro", "area", "floor", "type", "rating", "beds",
				}).AddRow(1, 11111.0, nil, "street_1", nil, 35.5, 2, "flat", nil, nil).
					AddRow(2, 22222.0, nil, "street_2", nil, 40.0, 3, "flat", nil, nil).
					AddRow(3, 33333.0, nil, "street_3", nil, 45.2, 4, "studio", nil, nil).
					AddRow(4, 44444.0, nil, "street_4", nil, 50.7, 5, "flat", nil, nil).
					AddRow(5, 55555.0, nil, "street_5", nil, 60.1, 6, "apartment", nil, nil)

				m.ExpectQuery(
					`SELECT id, price, image_url, address, metro, area, floor, type, rating, beds
					FROM posters
					ORDER BY id
					LIMIT \$1 OFFSET \$2`,
				).WithArgs(12, 0).WillReturnRows(rows)
			},
			want:    expectedListPoster,
			wantErr: nil,
		},
		{
			name: "empty_result",
			params: dto.PostersFiltersDTO{
				Limit:  12,
				Offset: 0,
			},
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "price", "image_url", "address", "metro", "area", "floor", "type", "rating", "beds",
				})

				m.ExpectQuery(
					`SELECT id, price, image_url, address, metro, area, floor, type, rating, beds
					FROM posters
					ORDER BY id
					LIMIT \$1 OFFSET \$2`,
				).WithArgs(12, 0).WillReturnRows(rows)
			},
			want:    []entity.Poster{},
			wantErr: nil,
		},
		{
			name: "collect_rows_error",
			params: dto.PostersFiltersDTO{
				Limit:  12,
				Offset: 0,
			},
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "price", "image_url", "address", "metro", "area", "floor", "type", "rating", "beds",
				}).AddRow(1, "bad_price", nil, "street_1", nil, 35.5, 2, "flat", nil, nil)

				m.ExpectQuery(
					`SELECT id, price, image_url, address, metro, area, floor, type, rating, beds
					FROM posters
					ORDER BY id
					LIMIT \$1 OFFSET \$2`,
				).WithArgs(12, 0).WillReturnRows(rows)
			},
			want:    nil,
			wantErr: entity.CollectPostersErr,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatalf("new pool: %v", err)
			}
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
