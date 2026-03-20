package complex

import (
	"context"
	"regexp"
	"testing"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"
)

func TestGetComplexRepo(t *testing.T) {
	testImg := "test.img"
	orderImg := 1
	photoID := 1
	expectedComplex := dto.UtilityCompanyDTO{
		ID:          1,
		Phone:       "8 985 981 09 77",
		CompanyName: "Рога и Копыта",
		Address:     "street_1",
		AvatarURL:   nil,
		Alias:       "toke_roga",
		Photos: []dto.PhotoDTO{
			{
				ImgURL: testImg,
				Order:  orderImg,
			},
		},
		GEO: dto.GeographyDTO{Lat: 1.112, Lon: 121.12},
	}

	query := regexp.QuoteMeta(`
        SELECT uc.id, uc.phone, uc.company_name, ST_AsText(uc.geo) AS geo, uc.address, uc.avatar_url, uc.alias,
               up.id as photo_id, up.utility_company_id, up.img_url, up.sequence_order
        FROM utility_companies uc
        LEFT JOIN utility_companies_photos up ON uc.id = up.utility_company_id
        WHERE uc.alias = $1`)

	tests := []struct {
		name      string
		alias     string
		setupMock func(m pgxmock.PgxPoolIface)
		want      *dto.UtilityCompanyDTO
		wantErr   error
	}{
		{
			name:  "OK_with_photos",
			alias: "toke_roga",
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id",
					"phone",
					"company_name",
					"geo",
					"address",
					"avatar_url",
					"alias",
					"photo_id",
					"utility_company_id",
					"img_url",
					"sequence_order",
				}).AddRow(
					int(1),
					"8 985 981 09 77",
					"Рога и Копыта",
					"POINT(121.12 1.112)",
					"street_1",
					nil,
					"toke_roga",
					&photoID,
					&photoID,
					&testImg,
					&orderImg,
				)

				m.ExpectQuery(query).WithArgs("toke_roga").WillReturnRows(rows)
			},
			want:    &expectedComplex,
			wantErr: nil,
		},
		{
			name:  "empty_result",
			alias: "unknown_alias",
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id",
					"phone",
					"company_name",
					"geo",
					"address",
					"avatar_url",
					"alias",
					"photo_id",
					"utility_company_id",
					"img_url",
					"sequence_order",
				})

				m.ExpectQuery(query).WithArgs("unknown_alias").WillReturnRows(rows)
			},
			want:    &dto.UtilityCompanyDTO{},
			wantErr: entity.NotFoundError,
		},
		{
			name:  "scan_error",
			alias: "toke_roga",
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id",
					"phone",
					"company_name",
					"geo",
					"address",
					"avatar_url",
					"alias",
					"photo_id",
					"utility_company_id",
					"img_url",
					"sequence_order",
				}).AddRow(
					"bad_id",
					"8 985 981 09 77",
					"Рога и Копыта",
					"POINT(121.12 1.112)",
					"street_1",
					nil,
					"toke_roga",
					&photoID,
					&photoID,
					&testImg,
					&orderImg,
				)

				m.ExpectQuery(query).WithArgs("toke_roga").WillReturnRows(rows)
			},
			want:    &dto.UtilityCompanyDTO{},
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewUtilityCompanyRepo(mock)

			got, err := repo.GetUtilityCompanyByAlias(context.Background(), test.alias)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				require.NoError(t, mock.ExpectationsWereMet())
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
