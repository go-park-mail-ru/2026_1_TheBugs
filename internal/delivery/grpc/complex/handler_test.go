package complex

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/mocks"
	usecase "github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/complex"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	complexpb "github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/generated/complex"
)

func TestComplexServiceServer_GetComplex(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	alias := "sunrise-park"
	companyAvatar := "https://example.com/company.png"
	developerAvatar := "https://example.com/dev.png"
	company := &dto.UtilityCompanyDTO{
		ID:          10,
		Phone:       "+1234567890",
		CompanyName: "Sunrise Park",
		Description: "Modern complex",
		GEO:         dto.GeographyDTO{Lat: 55.751244, Lon: 37.618423},
		Address:     "Main street, 1",
		AvatarURL:   &companyAvatar,
		Alias:       alias,
		Photos: []dto.PhotoDTO{
			{ImgURL: "https://example.com/photo1.png", Order: 1},
			{ImgURL: "https://example.com/photo2.png", Order: 2},
		},
		Developer: dto.DeveloperDTO{
			DeveloperID:   7,
			DeveloperName: "DevCo",
			AvatarURL:     &developerAvatar,
		},
	}

	tests := []struct {
		name      string
		req       *complexpb.GetComplexRequest
		setupMock func(repo *mocks.MockUtilityCompanyRepo)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			req:  &complexpb.GetComplexRequest{Alias: alias},
			setupMock: func(repo *mocks.MockUtilityCompanyRepo) {
				repo.EXPECT().
					GetByAlias(ctx, alias).
					Return(company, nil)
			},
			wantErr: false,
		},
		{
			name:      "missing alias",
			req:       &complexpb.GetComplexRequest{Alias: ""},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "not found",
			req:  &complexpb.GetComplexRequest{Alias: alias},
			setupMock: func(repo *mocks.MockUtilityCompanyRepo) {
				repo.EXPECT().
					GetByAlias(ctx, alias).
					Return(nil, entity.NotFoundError)
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockUtilityCompanyRepo(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(repo)
			}

			uc := usecase.NewUtilityCompanyUseCase(repo)
			server := NewComplexServiceServer(uc)

			resp, err := server.GetComplex(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tt.wantCode, st.Code())
				require.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, int32(company.ID), resp.Id)
			require.Equal(t, company.Alias, resp.Alias)
			require.Equal(t, company.CompanyName, resp.CompanyName)
			require.Equal(t, company.Description, resp.Description)
			require.Equal(t, company.Address, resp.Address)
			require.Equal(t, company.AvatarURL, resp.AvatarUrl)
			require.Equal(t, company.Phone, resp.Phone)
			require.NotNil(t, resp.Geo)
			require.Equal(t, company.GEO.Lat, resp.Geo.Lat)
			require.Equal(t, company.GEO.Lon, resp.Geo.Lon)
			require.Len(t, resp.Photos, len(company.Photos))
			for i, photo := range resp.Photos {
				require.Equal(t, company.Photos[i].ImgURL, photo.Url)
				require.Equal(t, int32(company.Photos[i].Order), photo.Order)
			}
			require.NotNil(t, resp.Developer)
			require.Equal(t, int32(company.Developer.DeveloperID), resp.Developer.Id)
			require.Equal(t, company.Developer.DeveloperName, resp.Developer.DeveloperName)
			require.Equal(t, company.Developer.AvatarURL, resp.Developer.AvatarUrl)
		})
	}
}

func TestComplexServiceServer_GetAllDevelopers(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	avatar := "https://example.com/dev.png"
	devs := []dto.DeveloperDTO{
		{
			DeveloperID:   1,
			DeveloperName: "DevOne",
			AvatarURL:     &avatar,
		},
		{
			DeveloperID:   2,
			DeveloperName: "DevTwo",
			AvatarURL:     nil,
		},
	}

	tests := []struct {
		name      string
		setupMock func(repo *mocks.MockUtilityCompanyRepo)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			setupMock: func(repo *mocks.MockUtilityCompanyRepo) {
				repo.EXPECT().
					GetAllDevelopers(ctx).
					Return(devs, nil)
			},
			wantErr: false,
		},
		{
			name: "not found",
			setupMock: func(repo *mocks.MockUtilityCompanyRepo) {
				repo.EXPECT().
					GetAllDevelopers(ctx).
					Return(nil, entity.NotFoundError)
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockUtilityCompanyRepo(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(repo)
			}

			uc := usecase.NewUtilityCompanyUseCase(repo)
			server := NewComplexServiceServer(uc)

			resp, err := server.GetAllDevelopers(ctx, &emptypb.Empty{})

			if tt.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tt.wantCode, st.Code())
				require.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, int32(len(devs)), resp.Length)
			require.Len(t, resp.Developers, len(devs))
			for i, dev := range resp.Developers {
				require.Equal(t, int32(devs[i].DeveloperID), dev.Id)
				require.Equal(t, devs[i].DeveloperName, dev.DeveloperName)
				require.Equal(t, devs[i].AvatarURL, dev.AvatarUrl)
			}
		})
	}
}

func TestComplexServiceServer_GetComplexesByDeveloperID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	devID := int32(12)
	avatar := "https://example.com/complex.png"
	complexes := []dto.UtilityCompanyCardDTO{
		{ID: 101, Alias: "first", CompanyName: "First Co", AvatarURL: &avatar},
		{ID: 102, Alias: "second", CompanyName: "Second Co", AvatarURL: nil},
	}

	tests := []struct {
		name      string
		req       *complexpb.GetComplexesByDeveloperIDRequest
		setupMock func(repo *mocks.MockUtilityCompanyRepo)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			req:  &complexpb.GetComplexesByDeveloperIDRequest{DeveloperID: devID},
			setupMock: func(repo *mocks.MockUtilityCompanyRepo) {
				repo.EXPECT().
					GetAllByDeveloperID(ctx, int(devID)).
					Return(complexes, nil)
			},
			wantErr: false,
		},
		{
			name:      "invalid developer id",
			req:       &complexpb.GetComplexesByDeveloperIDRequest{DeveloperID: -1},
			setupMock: nil,
			wantErr:   true,
			wantCode:  codes.InvalidArgument,
		},
		{
			name: "not found",
			req:  &complexpb.GetComplexesByDeveloperIDRequest{DeveloperID: devID},
			setupMock: func(repo *mocks.MockUtilityCompanyRepo) {
				repo.EXPECT().
					GetAllByDeveloperID(ctx, int(devID)).
					Return(nil, entity.NotFoundError)
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockUtilityCompanyRepo(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(repo)
			}

			uc := usecase.NewUtilityCompanyUseCase(repo)
			server := NewComplexServiceServer(uc)

			resp, err := server.GetComplexesByDeveloperID(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tt.wantCode, st.Code())
				require.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, int32(len(complexes)), resp.Length)
			require.Len(t, resp.Complexes, len(complexes))
			for i, comp := range resp.Complexes {
				require.Equal(t, int32(complexes[i].ID), comp.Id)
				require.Equal(t, complexes[i].Alias, comp.Alias)
				require.Equal(t, complexes[i].CompanyName, comp.CompanyName)
				require.Equal(t, complexes[i].AvatarURL, comp.AvatarUrl)
			}
		})
	}
}
