package poster

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/mocks"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/geo"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/photo"
	"github.com/golang/mock/gomock"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestGetPostersUseCase(t *testing.T) {
	ctx := context.Background()
	existingListPoster := []entity.PosterFlat{
		{ID: 1, Price: 11111, Address: "street_1"},
		{ID: 2, Price: 22222, Address: "street_2"},
		{ID: 3, Price: 33333, Address: "street_3"},
		{ID: 4, Price: 44444, Address: "street_4"},
		{ID: 5, Price: 55555, Address: "street_5"},
	}

	expectedDTO := []dto.PosterCardDTO{
		{ID: 1, Price: 11111, Address: "street_1"},
		{ID: 2, Price: 22222, Address: "street_2"},
		{ID: 3, Price: 33333, Address: "street_3"},
		{ID: 4, Price: 44444, Address: "street_4"},
		{ID: 5, Price: 55555, Address: "street_5"},
	}

	tests := []struct {
		name      string
		params    dto.PostersFiltersDTO
		setupMock func(m *mocks.MockPosterRepo)
		want      []dto.PosterCardDTO
		wantErr   error
	}{
		{
			name: "OK",
			params: dto.PostersFiltersDTO{
				Limit:  12,
				Offset: 0,
			},
			setupMock: func(m *mocks.MockPosterRepo) {
				m.EXPECT().GetFlatsAll(ctx, dto.PostersFiltersDTO{
					Limit:  12,
					Offset: 0,
				}).Return(existingListPoster, nil).Times(1)
			},
			want:    expectedDTO,
			wantErr: nil,
		},
		{
			name: "InvalidLimit",
			params: dto.PostersFiltersDTO{
				Limit:  -2,
				Offset: 0,
			},
			setupMock: nil,
			want:      nil,
			wantErr:   entity.InvalidInput,
		},
		{
			name: "LimitTooLarge",
			params: dto.PostersFiltersDTO{
				Limit:  MaxPostersLimit + 1,
				Offset: 0,
			},
			setupMock: func(m *mocks.MockPosterRepo) {
				m.EXPECT().GetFlatsAll(ctx, dto.PostersFiltersDTO{
					Limit:  MaxPostersLimit,
					Offset: 0,
				}).Return(existingListPoster, nil).Times(1)
			},
			want:    expectedDTO,
			wantErr: nil,
		},
		{
			name: "InvalidOffset",
			params: dto.PostersFiltersDTO{
				Limit:  12,
				Offset: -3,
			},
			setupMock: nil,
			want:      nil,
			wantErr:   entity.InvalidInput,
		},
		{
			name: "ErrorFromRepo",
			params: dto.PostersFiltersDTO{
				Limit:  12,
				Offset: 0,
			},
			setupMock: func(m *mocks.MockPosterRepo) {
				m.EXPECT().GetFlatsAll(ctx, dto.PostersFiltersDTO{
					Limit:  12,
					Offset: 0,
				}).Return(nil, entity.NotFoundError).Times(1)
			},
			want:    nil,
			wantErr: entity.NotFoundError,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockPosterRepo(ctrl)
			mockUOW := mocks.NewMockUnitOfWork(ctrl)
			mockFile := mocks.NewMockFileRepo(ctrl)

			mockUOW.EXPECT().
				Posters().
				Return(mockRepo).
				AnyTimes()

			if test.setupMock != nil {
				test.setupMock(mockRepo)
			}

			uc := NewPosterUseCase(mockUOW, mockFile, nil)

			got, err := uc.GetPostersUseCase(ctx, test.params)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}
			require.Equal(t, len(test.want), len(got.Posters))
			for i, p := range test.want {
				require.Equal(t, p.ID, got.Posters[i].ID)
				require.Equal(t, p.Address, got.Posters[i].Address)
			}

		})
	}
}

func TestGetPosterByAliasUseCase(t *testing.T) {
	ctx := context.Background()
	alias := "kvartira-na-arbate"
	host := "http://testhost:8000"
	bucket := "test_bucket"
	config.Config.PublicHost = host
	config.Config.Bucket = bucket

	existingPoster := &entity.PosterById{
		ID:              4,
		Alias:           alias,
		Price:           135000,
		Category:        PropertyFlat,
		Description:     "krutoy remont",
		PropertyID:      10,
		Area:            96.4,
		Geo:             geo.GeographyPoint{Lon: 45.3966, Lat: 46.3489},
		Address:         "Arbatskaya 5k2",
		District:        lo.ToPtr("Arbat"),
		Metro:           lo.ToPtr("Arbat"),
		MetroGeo:        &geo.GeographyPoint{Lon: 45.4000, Lat: 46.3519},
		City:            "Moscow",
		FloorCount:      7,
		SellerFirstName: "Sanya",
		SellerLastName:  "Sashenykov",
		Phone:           "+79144564312",
		Images: []entity.PosterImage{
			{ImgURL: photo.MakeUrlFromPath("img1.jpg", host, bucket), Order: 1},
			{ImgURL: photo.MakeUrlFromPath("img2.jpg", host, bucket), Order: 2},
		},
	}

	existingFlat := &entity.Flat{
		PropertyID:   10,
		FlatCategory: "1-room",
		Number:       43,
		Floor:        3,
	}

	expectedPoster := &dto.PosterDTO{
		ID:          4,
		Alias:       alias,
		Price:       135000,
		Category:    dto.CategoryDTO{Name: PropertyFlat, Alias: "flat"},
		Description: "krutoy remont",
		Area:        96.4,
		Geo:         dto.GeographyDTO{Lon: 45.3966, Lat: 46.3489},
		Address:     "Arbatskaya 5k2",
		District:    lo.ToPtr("Arbat"),
		Metro:       lo.ToPtr("Arbat"),
		MetroGeo:    &dto.GeographyDTO{Lon: 45.4000, Lat: 46.3519},
		City:        "Moscow",
		FloorCount:  7,
		Images: []dto.PhotoDTO{
			{ImgURL: photo.MakeUrlFromPath("img1.jpg", host, bucket), Order: 1},
			{ImgURL: photo.MakeUrlFromPath("img2.jpg", host, bucket), Order: 2},
		},
		Seller: dto.PosterSellerDTO{
			SellerFirstName: "Sanya",
			SellerLastName:  "Sashenykov",
			SellerPhone:     "+79144564312",
		},
		Flat: &dto.FlatDTO{
			FlatCategory: "1-room",
			Number:       43,
			Floor:        3,
		},
	}

	tests := []struct {
		name      string
		param     string
		setupMock func(m *mocks.MockPosterRepo)
		want      *dto.PosterDTO
		wantErr   error
	}{
		{
			name:  "ok_flat",
			param: alias,
			setupMock: func(m *mocks.MockPosterRepo) {
				m.EXPECT().
					GetByAlias(ctx, alias, nil).
					Return(existingPoster, nil).
					Times(1)

				m.EXPECT().
					GetFlatByPropetyID(ctx, 10).
					Return(existingFlat, nil).
					Times(1)
			},
			want:    expectedPoster,
			wantErr: nil,
		},
		{
			name:  "repo_error",
			param: alias,
			setupMock: func(m *mocks.MockPosterRepo) {
				m.EXPECT().
					GetByAlias(ctx, alias, nil).
					Return(&entity.PosterById{}, entity.NotFoundError).
					Times(1)
			},
			want:    &dto.PosterDTO{},
			wantErr: entity.NotFoundError,
		},
		{
			name:  "flat_repo_error",
			param: alias,
			setupMock: func(m *mocks.MockPosterRepo) {
				m.EXPECT().
					GetByAlias(ctx, alias, nil).
					Return(existingPoster, nil).
					Times(1)

				m.EXPECT().
					GetFlatByPropetyID(ctx, 10).
					Return(nil, entity.ServiceError).
					Times(1)
			},
			want:    &dto.PosterDTO{},
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockPosterRepo(ctrl)
			mockUOW := mocks.NewMockUnitOfWork(ctrl)
			mockFile := mocks.NewMockFileRepo(ctrl)

			mockUOW.EXPECT().
				Posters().
				Return(mockRepo).
				AnyTimes()

			if test.setupMock != nil {
				test.setupMock(mockRepo)
			}

			uc := NewPosterUseCase(mockUOW, mockFile, nil)

			res, err := uc.GetPosterByAliasUseCase(ctx, alias, nil)

			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)

			require.Equal(t, test.want.ID, res.ID)
			require.Equal(t, test.want.Alias, res.Alias)
			require.Equal(t, test.want.Price, res.Price)
			require.Equal(t, test.want.Category, res.Category)
			require.Equal(t, test.want.Description, res.Description)
			require.Equal(t, test.want.Address, res.Address)
			require.Equal(t, test.want.City, res.City)
			require.Equal(t, test.want.FloorCount, res.FloorCount)
			require.Equal(t, test.want.District, res.District)
			require.Equal(t, test.want.Metro, res.Metro)
			require.Equal(t, test.want.Seller, res.Seller)

			require.Equal(t, len(test.want.Images), len(res.Images))
			for i := range test.want.Images {
				require.Equal(t, test.want.Images[i], res.Images[i])
			}

			if test.want.Flat != nil {
				require.Equal(t, test.want.Flat.FlatCategory, res.Flat.FlatCategory)
				require.Equal(t, test.want.Flat.Number, res.Flat.Number)
				require.Equal(t, test.want.Flat.Floor, res.Flat.Floor)
			}

			require.Equal(t, test.want.Geo.Lat, res.Geo.Lat)
			require.Equal(t, test.want.Geo.Lon, res.Geo.Lon)

			if test.want.MetroGeo != nil {
				require.NotNil(t, res.MetroGeo)
				require.Equal(t, test.want.MetroGeo.Lat, res.MetroGeo.Lat)
				require.Equal(t, test.want.MetroGeo.Lon, res.MetroGeo.Lon)
			}
		})
	}
}

func TestGetPosterByUserID(t *testing.T) {
	t.Parallel()
	userID := 0
	ctx := context.Background()
	tests := []struct {
		name      string
		want      []dto.MyPosterDTO
		setup     func(mock *mocks.MockPosterRepo)
		wantError error
	}{
		{
			name: "OK",
			want: []dto.MyPosterDTO{
				{
					ID:        1,
					Price:     100,
					AvatarURl: lo.ToPtr("img_url"),
					Address:   "addres",
					Area:      10,
					Alias:     "alias",
				},
			},
			setup: func(mock *mocks.MockPosterRepo) {
				posters := []entity.Poster{
					{
						ID:        1,
						Price:     100,
						AvatarURl: lo.ToPtr("img_url"),
						Address:   "addres",
						Area:      10,
						Alias:     "alias",
					},
				}
				mock.EXPECT().GetByUserID(ctx, userID).Return(posters, nil)
			},
			wantError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uow := mocks.NewMockUnitOfWork(ctrl)
			posterRepo := mocks.NewMockPosterRepo(ctrl)
			fileRepo := mocks.NewMockFileRepo(ctrl)

			uow.EXPECT().Posters().Return(posterRepo).AnyTimes()

			tt.setup(posterRepo)

			uc := NewPosterUseCase(uow, fileRepo, nil)
			res, err := uc.GetPosterByUserID(ctx, userID)

			require.NoError(t, err)

			require.Equal(t, len(tt.want), len(res))

			for i, p := range res {
				require.Equal(t, tt.want[i].ID, p.ID)
				require.Equal(t, tt.want[i].Address, p.Address)
				require.Equal(t, tt.want[i].Alias, p.Alias)
				require.Equal(t, tt.want[i].Price, p.Price)
			}
		})

	}
}
