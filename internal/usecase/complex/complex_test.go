package complex

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/mocks"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/geo"
	"github.com/golang/mock/gomock"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestUtilityCompanyUseCase_GetUtilityCompany(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	cases := []struct {
		name      string
		alias     string
		setupMock func(repoMock *mocks.MockUtilityCompanyRepo)
		wantErr   error
	}{
		{
			name:  "OK",
			alias: "test",
			setupMock: func(repoMock *mocks.MockUtilityCompanyRepo) {
				expected := &dto.UtilityCompanyDTO{ID: 1, Alias: "test"}
				repoMock.EXPECT().GetByAlias(ctx, "test").Return(expected, nil)
			},
			wantErr: nil,
		},
		{
			name:  "NotFound",
			alias: "nonexistent",
			setupMock: func(repoMock *mocks.MockUtilityCompanyRepo) {
				repoMock.EXPECT().GetByAlias(ctx, "nonexistent").Return(nil, entity.NotFoundError)
			},
			wantErr: entity.NotFoundError,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repoMock := mocks.NewMockUtilityCompanyRepo(ctrl)
			tc.setupMock(repoMock)

			uc := NewUtilityCompanyUseCase(repoMock)
			_, err := uc.GetUtilityCompany(ctx, tc.alias)

			require.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestUtilityCompanyUseCase_GetAllDevelopers(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	cases := []struct {
		name      string
		setupMock func(repoMock *mocks.MockUtilityCompanyRepo)
		wantErr   error
	}{
		{
			name: "OK",
			setupMock: func(repoMock *mocks.MockUtilityCompanyRepo) {
				expected := []dto.DeveloperDTO{{DeveloperID: 1, DeveloperName: "Dev1"}}
				repoMock.EXPECT().GetAllDevelopers(ctx).Return(expected, nil)
			},
			wantErr: nil,
		},
		{
			name: "RepoError",
			setupMock: func(repoMock *mocks.MockUtilityCompanyRepo) {
				repoMock.EXPECT().GetAllDevelopers(ctx).Return(nil, entity.ServiceError)
			},
			wantErr: entity.ServiceError,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repoMock := mocks.NewMockUtilityCompanyRepo(ctrl)
			tc.setupMock(repoMock)

			uc := NewUtilityCompanyUseCase(repoMock)
			_, err := uc.GetAllDevelopers(ctx)

			require.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestUtilityCompanyUseCase_GetAllByDeveloperID(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	cases := []struct {
		name      string
		companyID int
		setupMock func(repoMock *mocks.MockUtilityCompanyRepo)
		wantErr   error
	}{
		{
			name:      "OK",
			companyID: 1,
			setupMock: func(repoMock *mocks.MockUtilityCompanyRepo) {
				expected := []dto.UtilityCompanyCardDTO{{ID: 1, CompanyName: "Company1"}}
				repoMock.EXPECT().GetAllByDeveloperID(ctx, 1).Return(expected, nil)
			},
			wantErr: nil,
		},
		{
			name:      "NotFound",
			companyID: 999,
			setupMock: func(repoMock *mocks.MockUtilityCompanyRepo) {
				repoMock.EXPECT().GetAllByDeveloperID(ctx, 999).Return(nil, entity.NotFoundError)
			},
			wantErr: entity.NotFoundError,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repoMock := mocks.NewMockUtilityCompanyRepo(ctrl)
			tc.setupMock(repoMock)

			uc := NewUtilityCompanyUseCase(repoMock)
			_, err := uc.GetAllByDeveloperID(ctx, tc.companyID)

			require.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestToUtilityCompanyDTO(t *testing.T) {
	in := entity.UtilityCompany{
		ID:          1,
		Phone:       "phone",
		CompanyName: "company_name",
		GEO:         geo.GeographyPoint{Lat: 10, Lon: 12},
		Address:     "address",
		Alias:       "alias",
	}
	photos := []entity.UtilityCompanyPhoto{{ID: lo.ToPtr(1), UtilityCompanyID: lo.ToPtr(1), ImgURL: lo.ToPtr("s"), Order: lo.ToPtr(1)}}

	exp := UtilityCompanyDTO{
		ID:          1,
		Phone:       "phone",
		CompanyName: "company_name",
		GEO:         dto.GeographyDTO{Lat: 10, Lon: 12},
		Address:     "address",
		Alias:       "alias",
		Photos: []dto.PhotoDTO{
			{
				ImgURL: "s",
				Order:  1,
			},
		},
	}
	res := ToUtilityCompanyDTO(&in, photos)
	require.Equal(t, *res, exp)
}

func TestPosterToUtilityCompanyCardDTO_OK(t *testing.T) {
	companyID := 10
	companyName := "Test Company"
	companyAlias := "test-company"
	avatar := "https://img.com/avatar.png"

	poster := &entity.PosterById{
		CompanyID:        &companyID,
		CompanyName:      &companyName,
		CompanyAlias:     &companyAlias,
		CompanyAvatarURL: &avatar,
	}

	dto := PosterToUtilityCompanyCardDTO(poster)

	require.NotNil(t, dto)
	require.Equal(t, 10, dto.ID)
	require.Equal(t, "Test Company", dto.CompanyName)
	require.Equal(t, "test-company", dto.Alias)
	require.Equal(t, &avatar, dto.AvatarURL)
}

func TestPosterToUtilityCompanyCardDTO_NilCompany(t *testing.T) {
	poster := &entity.PosterById{
		CompanyID: nil,
	}

	dto := PosterToUtilityCompanyCardDTO(poster)

	require.Nil(t, dto)
}
