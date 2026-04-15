package complex

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestGetComplexRepo(t *testing.T) {
	t.Parallel()

	expectedComplex := &dto.UtilityCompanyDTO{
		ID:          1,
		Phone:       "8 985 981 09 77",
		CompanyName: "Рога и Копыта",
		Address:     "street_1",
		AvatarURL:   nil,
		Alias:       "toke_roga",
		Photos: []dto.PhotoDTO{
			{ImgURL: "test.img", Order: 1},
		},
		GEO: dto.GeographyDTO{Lat: 1.112, Lon: 121.12},
		Developer: dto.DeveloperDTO{
			DeveloperID:   42,
			DeveloperName: "Developer Name",
			AvatarURL:     nil,
		},
	}

	query := regexp.QuoteMeta(`
        SELECT uc.id, uc.phone, uc.company_name, uc.description, ST_AsText(uc.geo) AS geo, uc.address, uc.avatar_url, uc.alias,
               up.id as photo_id, up.utility_company_id, up.img_url, up.sequence_order, 
               d.id, d.developer_name, d.avatar_url
        FROM utility_companies uc
        LEFT JOIN developers d ON d.id = uc.developer_id
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
			name:  "OK_with_photos_and_developer",
			alias: "toke_roga",
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "phone", "company_name", "description", "geo", "address", "avatar_url", "alias",
					"photo_id", "utility_company_id", "img_url", "sequence_order",
					"d.id", "d.developer_name", "d.avatar_url",
				}).AddRow(
					1, "8 985 981 09 77", "Рога и Копыта", nil, "POINT(121.12 1.112)", "street_1", nil, "toke_roga",
					lo.ToPtr(1), lo.ToPtr(1), lo.ToPtr("test.img"), lo.ToPtr(1),
					42, "Developer Name", nil,
				)
				m.ExpectQuery(query).WithArgs("toke_roga").WillReturnRows(rows)
			},
			want:    expectedComplex,
			wantErr: nil,
		},
		{
			name:  "OK_no_photos_with_developer",
			alias: "no_photos",
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "phone", "company_name", "description", "geo", "address", "avatar_url", "alias",
					"photo_id", "utility_company_id", "img_url", "sequence_order",
					"d.id", "d.developer_name", "d.avatar_url",
				}).AddRow(
					2, "phone2", "company2", nil, "POINT(22.2 2.2)", "addr2", nil, "no_photos",
					nil, nil, nil, nil,
					43, "Dev2", nil,
				)
				m.ExpectQuery(query).WithArgs("no_photos").WillReturnRows(rows)
			},
			want: &dto.UtilityCompanyDTO{
				ID:          2,
				Phone:       "phone2",
				CompanyName: "company2",
				Address:     "addr2",
				AvatarURL:   nil,
				Alias:       "no_photos",
				Photos:      []dto.PhotoDTO{},
				GEO:         dto.GeographyDTO{Lat: 2.2, Lon: 22.2},
				Developer: dto.DeveloperDTO{
					DeveloperID:   43,
					DeveloperName: "Dev2",
					AvatarURL:     nil,
				},
			},
			wantErr: nil,
		},
		{
			name:  "OK_with_photos_no_developer",
			alias: "no_dev",
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "phone", "company_name", "description", "geo", "address", "avatar_url", "alias",
					"photo_id", "utility_company_id", "img_url", "sequence_order",
					"d.id", "d.developer_name", "d.avatar_url",
				}).AddRow(
					3, "phone3", "company3", "desc3", "POINT(33.3 3.3)", "addr3", nil, "no_dev",
					lo.ToPtr(1), lo.ToPtr(1), lo.ToPtr("photo3.jpg"), lo.ToPtr(1),

					nil, nil, nil,
				)
				m.ExpectQuery(query).WithArgs("no_dev").WillReturnRows(rows)
			},
			want: &dto.UtilityCompanyDTO{
				ID:          3,
				Phone:       "phone3",
				CompanyName: "company3",
				Address:     "addr3",
				AvatarURL:   nil,
				Description: "desc3",
				Alias:       "no_dev",
				Photos:      []dto.PhotoDTO{{ImgURL: "photo3.jpg", Order: 1}},
				GEO:         dto.GeographyDTO{Lat: 3.3, Lon: 33.3},
			},
			wantErr: nil,
		},
		{
			name:  "not_found",
			alias: "unknown",
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(query).WithArgs("unknown").WillReturnRows(
					pgxmock.NewRows([]string{
						"id", "phone", "company_name", "description", "geo", "address", "avatar_url", "alias",
						"photo_id", "utility_company_id", "img_url", "sequence_order",
						"d.id", "d.developer_name", "d.avatar_url",
					}),
				)
			},
			want:    nil,
			wantErr: entity.NotFoundError,
		},
		{
			name:  "db_error",
			alias: "db_error",
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(query).WithArgs("db_error").WillReturnError(errors.New("db failed"))
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
			repo := NewUtilityCompanyRepo(mock)

			got, err := repo.GetByAlias(context.Background(), test.alias)

			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.want, got)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUtilityCompanyPgRepo_GetAllDevelopers(t *testing.T) {
	t.Parallel()

	query := regexp.QuoteMeta(`
        SELECT d.id, d.developer_name, d.avatar_url
        FROM developers d
        ORDER BY d.id`)

	expectedDevelopers := []dto.DeveloperDTO{
		{
			DeveloperID:   1,
			DeveloperName: "Dev1",
			AvatarURL:     lo.ToPtr("avatar1.jpg"),
		},
		{
			DeveloperID:   2,
			DeveloperName: "Dev2",
			AvatarURL:     nil,
		},
	}

	tests := []struct {
		name      string
		setupMock func(m pgxmock.PgxPoolIface)
		want      []dto.DeveloperDTO
		wantErr   error
	}{
		{
			name: "OK_multiple_developers",
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "developer_name", "avatar_url"}).
					AddRow(1, "Dev1", lo.ToPtr("avatar1.jpg")).
					AddRow(2, "Dev2", nil)
				m.ExpectQuery(query).WillReturnRows(rows)
			},
			want:    expectedDevelopers,
			wantErr: nil,
		},
		{
			name: "OK_empty_result",
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(query).WillReturnRows(pgxmock.NewRows([]string{"id", "developer_name", "avatar_url"}))
			},
			want:    []dto.DeveloperDTO{},
			wantErr: nil,
		},
		{
			name: "query_error",
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(query).WillReturnError(errors.New("db error"))
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
		{
			name: "scan_error",
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id"}).AddRow("bad_id")
				m.ExpectQuery(query).WillReturnRows(rows)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			tt.setupMock(mock)
			repo := NewUtilityCompanyRepo(mock)

			got, err := repo.GetAllDevelopers(context.Background())

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUtilityCompanyPgRepo_GetAllByDeveloperID(t *testing.T) {
	t.Parallel()

	queryTemplate := regexp.QuoteMeta(`
		SELECT uc.id, uc.company_name, uc.avatar_url, uc.alias
		FROM utility_companies uc
		WHERE uc.developer_id = $1
		ORDER BY uc.id
	`)

	expectedCompanies := []dto.UtilityCompanyCardDTO{
		{
			ID:          1,
			CompanyName: "Company1",
			AvatarURL:   lo.ToPtr("avatar1.jpg"),
			Alias:       "company1_alias",
		},
		{
			ID:          2,
			CompanyName: "Company2",
			AvatarURL:   nil,
			Alias:       "company2_alias",
		},
	}

	tests := []struct {
		name        string
		developerID int
		setupMock   func(m pgxmock.PgxPoolIface)
		want        []dto.UtilityCompanyCardDTO
		wantErr     error
	}{
		{
			name:        "OK_multiple_companies",
			developerID: 42,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "company_name", "avatar_url", "alias"}).
					AddRow(1, "Company1", lo.ToPtr("avatar1.jpg"), "company1_alias").
					AddRow(2, "Company2", nil, "company2_alias")
				m.ExpectQuery(queryTemplate).WithArgs(42).WillReturnRows(rows)
			},
			want:    expectedCompanies,
			wantErr: nil,
		},
		{
			name:        "OK_no_companies",
			developerID: 999,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(queryTemplate).WithArgs(999).
					WillReturnRows(pgxmock.NewRows([]string{"id", "company_name", "avatar_url", "alias"}))
			},
			want:    []dto.UtilityCompanyCardDTO{},
			wantErr: nil,
		},
		{
			name:        "query_error",
			developerID: 42,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(queryTemplate).WithArgs(42).
					WillReturnError(errors.New("database connection failed"))
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
		{
			name:        "scan_error",
			developerID: 42,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id"}).AddRow("bad_id")
				m.ExpectQuery(queryTemplate).WithArgs(42).WillReturnRows(rows)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			tt.setupMock(mock)
			repo := NewUtilityCompanyRepo(mock)

			got, err := repo.GetAllByDeveloperID(context.Background(), tt.developerID)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
