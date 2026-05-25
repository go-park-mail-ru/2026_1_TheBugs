package poster

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/mocks"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase"
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

				m.EXPECT().CountPosters(ctx).
					Return(len(existingListPoster), nil).Times(1)
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

				m.EXPECT().CountPosters(ctx).
					Return(len(existingListPoster), nil).Times(1)
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

			uc := NewPosterUseCase(mockUOW, mockFile, nil, nil, nil, nil)

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

func TestSearchPostersUseCase(t *testing.T) {
	ctx := context.Background()

	existingSearchResponse := &dto.PostersResponse{
		Posters: []dto.PosterCardDTO{
			{ID: 1, Price: 11111, Address: "street_1"},
			{ID: 2, Price: 22222, Address: "street_2"},
			{ID: 3, Price: 33333, Address: "street_3"},
		},
		Len: 3,
	}

	tests := []struct {
		name      string
		params    dto.PostersFiltersDTO
		setupMock func(searchMock *mocks.MockSearchRepo)
		want      *dto.PostersResponse
		wantErr   error
	}{
		{
			name: "OK",
			params: dto.PostersFiltersDTO{
				Limit:  12,
				Offset: 0,
			},
			setupMock: func(searchMock *mocks.MockSearchRepo) {
				searchMock.EXPECT().
					SearchPosters(ctx, dto.PostersFiltersDTO{
						Limit:  12,
						Offset: 0,
					}).
					Return(existingSearchResponse, nil).
					Times(1)
			},
			want:    existingSearchResponse,
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
			name: "LimitTooLarge",
			params: dto.PostersFiltersDTO{
				Limit:  MaxPostersLimit + 1,
				Offset: 0,
			},
			setupMock: func(searchMock *mocks.MockSearchRepo) {
				searchMock.EXPECT().
					SearchPosters(ctx, dto.PostersFiltersDTO{
						Limit:  MaxPostersLimit,
						Offset: 0,
					}).
					Return(existingSearchResponse, nil).
					Times(1)
			},
			want:    existingSearchResponse,
			wantErr: nil,
		},
		{
			name: "SearchError",
			params: dto.PostersFiltersDTO{
				Limit:  12,
				Offset: 0,
			},
			setupMock: func(searchMock *mocks.MockSearchRepo) {
				searchMock.EXPECT().
					SearchPosters(ctx, dto.PostersFiltersDTO{
						Limit:  12,
						Offset: 0,
					}).
					Return(nil, entity.NotFoundError).
					Times(1)
			},
			want:    nil,
			wantErr: entity.NotFoundError,
		},
		{
			name: "Empty search result",
			params: dto.PostersFiltersDTO{
				Limit:  12,
				Offset: 0,
			},
			setupMock: func(searchMock *mocks.MockSearchRepo) {
				emptyResponse := &dto.PostersResponse{
					Posters: []dto.PosterCardDTO{},
					Len:     0,
				}

				searchMock.EXPECT().
					SearchPosters(ctx, dto.PostersFiltersDTO{
						Limit:  12,
						Offset: 0,
					}).
					Return(emptyResponse, nil).
					Times(1)
			},
			want: &dto.PostersResponse{
				Posters: []dto.PosterCardDTO{},
				Len:     0,
			},
			wantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSearch := mocks.NewMockSearchRepo(ctrl)
			mockUOW := mocks.NewMockUnitOfWork(ctrl)
			mockFile := mocks.NewMockFileRepo(ctrl)

			uc := &PosterUseCase{
				search: mockSearch,
				uow:    mockUOW,
				file:   mockFile,
			}

			if test.setupMock != nil {
				test.setupMock(mockSearch)
			}

			got, err := uc.SearchPostersUseCase(ctx, test.params)

			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				require.Nil(t, got)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			require.Equal(t, test.want.Len, got.Len)
			require.Equal(t, len(test.want.Posters), len(got.Posters))

			for i, expectedPoster := range test.want.Posters {
				poster := got.Posters[i]
				require.Equal(t, expectedPoster.ID, poster.ID)
				require.Equal(t, expectedPoster.Price, poster.Price)
				require.Equal(t, expectedPoster.Address, poster.Address)
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
		CategoryAlias:   PropertyFlat,
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

			uc := NewPosterUseCase(mockUOW, mockFile, nil, nil, nil, nil)

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

			uc := NewPosterUseCase(uow, fileRepo, nil, nil, nil, nil)
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

func TestGetMetroStationsByRadius(t *testing.T) {
	ctx := context.Background()

	inputCoords := dto.GeographyDTO{
		Lon: 37.6173,
		Lat: 55.7558,
	}

	existingStations := []entity.MetroStation{
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

	expectedDTO := []dto.MetroStationDTO{
		{
			StationName: "Арбатская",
			StationGeo:  dto.GeographyDTO{Lon: 37.605, Lat: 55.752},
		},
		{
			StationName: "Смоленская",
			StationGeo:  dto.GeographyDTO{Lon: 37.582, Lat: 55.748},
		},
	}

	tests := []struct {
		name      string
		param     dto.GeographyDTO
		setupMock func(m *mocks.MockPosterRepo)
		want      []dto.MetroStationDTO
		wantErr   error
	}{
		{
			name:  "OK",
			param: inputCoords,
			setupMock: func(m *mocks.MockPosterRepo) {
				m.EXPECT().
					GetMetroStationByRadius(ctx, inputCoords, entity.Metre(MetroRadius)).
					Return(existingStations, nil).
					Times(1)
			},
			want:    expectedDTO,
			wantErr: nil,
		},
		{
			name:  "EmptyResult",
			param: inputCoords,
			setupMock: func(m *mocks.MockPosterRepo) {
				m.EXPECT().
					GetMetroStationByRadius(ctx, inputCoords, entity.Metre(MetroRadius)).
					Return([]entity.MetroStation{}, nil).
					Times(1)
			},
			want:    []dto.MetroStationDTO{},
			wantErr: nil,
		},
		{
			name:  "RepoError",
			param: inputCoords,
			setupMock: func(m *mocks.MockPosterRepo) {
				m.EXPECT().
					GetMetroStationByRadius(ctx, inputCoords, entity.Metre(MetroRadius)).
					Return(nil, entity.ServiceError).
					Times(1)
			},
			want:    nil,
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

			uc := NewPosterUseCase(mockUOW, mockFile, nil, nil, nil, nil)

			res, err := uc.GetMetroStationsByRadius(ctx, test.param)

			if test.wantErr != nil {
				require.Error(t, err)
				require.ErrorContains(t, err, test.wantErr.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, len(test.want), len(res))

			for i := range test.want {
				require.Equal(t, test.want[i].StationName, res[i].StationName)
				require.Equal(t, test.want[i].StationGeo.Lat, res[i].StationGeo.Lat)
				require.Equal(t, test.want[i].StationGeo.Lon, res[i].StationGeo.Lon)
			}
		})
	}
}

func TestCreateFlatPoster(t *testing.T) {
	ctx := context.Background()

	district := "Arbat"
	flatNumber := 43
	companyID := 12

	validPoster := &dto.PosterInputFlatDTO{
		UserID:         7,
		Price:          135000,
		Description:    "krutoy remont",
		CategoryAlias:  PropertyFlat,
		Area:           96.4,
		GeoLat:         55.7558,
		GeoLon:         37.6173,
		FlatCategoryID: 1,
		FlatNumber:     &flatNumber,
		FlatFloor:      3,
		Address:        "Arbatskaya 5k2",
		City:           "Moscow",
		District:       &district,
		FloorCount:     7,
		CompanyID:      &companyID,
		Features:       []string{"balcony", "parking"},
		Images: []dto.PhotoInputDTO{
			{
				FileHeader: &dto.FileInput{
					Filename:    "img1.jpg",
					Size:        10,
					ContentType: "image/jpeg",
					File:        io.NopCloser(strings.NewReader("img1")),
				},
				Order: 1,
			},
			{
				FileHeader: &dto.FileInput{
					Filename:    "img2.jpg",
					Size:        10,
					ContentType: "image/jpeg",
					File:        io.NopCloser(strings.NewReader("img2")),
				},
				Order: 2,
			},
		},
	}

	tests := []struct {
		name      string
		poster    *dto.PosterInputFlatDTO
		setupMock func(repo *mocks.MockPosterRepo, uow *mocks.MockUnitOfWork, file *mocks.MockFileRepo)
		want      *dto.CreatedPoster
		wantErr   error
	}{
		{
			name:   "OK",
			poster: validPoster,
			setupMock: func(repo *mocks.MockPosterRepo, uow *mocks.MockUnitOfWork, file *mocks.MockFileRepo) {
				repo.EXPECT().
					GetCityByName(ctx, "Moscow").
					Return(&entity.City{ID: 1, Name: "Moscow"}, nil).
					Times(1)

				repo.EXPECT().
					GetMetroStationByRadius(ctx, dto.GeographyDTO{Lat: 55.7558, Lon: 37.6173}, entity.Metre(MetroRadius)).
					Return([]entity.MetroStation{
						{ID: 5, StationName: "Arbatskaya"},
					}, nil).
					Times(1)

				repo.EXPECT().GetByAlias(ctx, gomock.Any(), gomock.Any()).Return(nil, entity.NotFoundError)

				uow.EXPECT().
					Do(ctx, gomock.Any()).
					DoAndReturn(func(_ context.Context, fn func(r usecase.UnitOfWork) error) error {
						return fn(uow)
					}).
					Times(1)

				repo.EXPECT().
					CreateBuilding(ctx, gomock.Any()).
					Return(11, nil).
					Times(1)

				repo.EXPECT().
					CreateProperty(ctx, gomock.Any(), 11).
					Return(22, nil).
					Times(1)

				repo.EXPECT().
					InsertFacilities(ctx, 22, []string{"balcony", "parking"}).
					Return(nil).
					Times(1)

				repo.EXPECT().
					Create(ctx, gomock.Any(), 22).
					Return(33, nil).
					Times(1)

				repo.EXPECT().
					AddPriceHistory(ctx, 33, float64(135000)).
					Return(nil).
					Times(1)

				repo.EXPECT().
					InsertFlat(ctx, gomock.Any()).
					Return(nil).
					Times(1)

				file.EXPECT().
					Upload(ctx, gomock.Any(), gomock.Any(), int64(10), "image/jpeg").
					Return(nil).
					Times(2)

				repo.EXPECT().
					InsertPhotos(ctx, 33, gomock.Any()).
					Return(nil).
					Times(1)

				repo.EXPECT().
					InsertMainPhoto(ctx, 33, gomock.Any()).
					Return(nil).
					Times(1)
			},
			want: &dto.CreatedPoster{
				ID: 33,
			},
			wantErr: nil,
		},
		{
			name:   "CreateCityIfNotFound",
			poster: validPoster,
			setupMock: func(repo *mocks.MockPosterRepo, uow *mocks.MockUnitOfWork, file *mocks.MockFileRepo) {
				repo.EXPECT().
					GetCityByName(ctx, "Moscow").
					Return(nil, entity.NotFoundError).
					Times(1)

				repo.EXPECT().
					CreateCity(ctx, "Moscow").
					Return(&entity.City{ID: 2, Name: "Moscow"}, nil).
					Times(1)

				repo.EXPECT().
					GetMetroStationByRadius(ctx, dto.GeographyDTO{Lat: 55.7558, Lon: 37.6173}, entity.Metre(MetroRadius)).
					Return([]entity.MetroStation{}, nil).
					Times(1)

				repo.EXPECT().GetByAlias(ctx, gomock.Any(), gomock.Any()).Return(nil, entity.NotFoundError)

				uow.EXPECT().
					Do(ctx, gomock.Any()).
					DoAndReturn(func(_ context.Context, fn func(r usecase.UnitOfWork) error) error {
						return fn(uow)
					}).
					Times(1)

				repo.EXPECT().
					CreateBuilding(ctx, gomock.Any()).
					Return(11, nil).
					Times(1)

				repo.EXPECT().
					CreateProperty(ctx, gomock.Any(), 11).
					Return(22, nil).
					Times(1)

				repo.EXPECT().
					InsertFacilities(ctx, 22, []string{"balcony", "parking"}).
					Return(nil).
					Times(1)

				repo.EXPECT().
					Create(ctx, gomock.Any(), 22).
					Return(33, nil).
					Times(1)

				repo.EXPECT().
					AddPriceHistory(ctx, 33, float64(135000)).
					Return(nil).
					Times(1)

				repo.EXPECT().
					InsertFlat(ctx, gomock.Any()).
					Return(nil).
					Times(1)

				file.EXPECT().
					Upload(ctx, gomock.Any(), gomock.Any(), int64(10), "image/jpeg").
					Return(nil).
					Times(2)

				repo.EXPECT().
					InsertPhotos(ctx, 33, gomock.Any()).
					Return(nil).
					Times(1)

				repo.EXPECT().
					InsertMainPhoto(ctx, 33, gomock.Any()).
					Return(nil).
					Times(1)
			},
			want: &dto.CreatedPoster{
				ID: 33,
			},
			wantErr: nil,
		},
		{
			name:   "GetCityError",
			poster: validPoster,
			setupMock: func(repo *mocks.MockPosterRepo, uow *mocks.MockUnitOfWork, file *mocks.MockFileRepo) {
				repo.EXPECT().
					GetCityByName(ctx, "Moscow").
					Return(nil, entity.ServiceError).
					Times(1)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
		{
			name:   "MetroError",
			poster: validPoster,
			setupMock: func(repo *mocks.MockPosterRepo, uow *mocks.MockUnitOfWork, file *mocks.MockFileRepo) {
				repo.EXPECT().
					GetCityByName(ctx, "Moscow").
					Return(&entity.City{ID: 1, Name: "Moscow"}, nil).
					Times(1)

				repo.EXPECT().
					GetMetroStationByRadius(ctx, dto.GeographyDTO{Lat: 55.7558, Lon: 37.6173}, entity.Metre(MetroRadius)).
					Return(nil, entity.ServiceError).
					Times(1)
			},
			want:    nil,
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
				test.setupMock(mockRepo, mockUOW, mockFile)
			}

			uc := NewPosterUseCase(mockUOW, mockFile, nil, nil, nil, nil)

			res, err := uc.CreateFlatPoster(ctx, test.poster)

			if test.wantErr != nil {
				require.Error(t, err)
				require.ErrorContains(t, err, test.wantErr.Error())
				return
			}

			require.NoError(t, err)
			require.NotNil(t, res)
			require.Equal(t, test.want.ID, res.ID)
			require.NotEmpty(t, res.Alias)
		})
	}
}

func TestUpdateFlatPoster(t *testing.T) {
	ctx := context.Background()
	alias := "kvartira-na-arbate"

	host := "http://testhost:8000"
	bucket := "test_bucket"
	config.Config.PublicHost = host
	config.Config.Bucket = bucket

	district := "Arbat"
	flatNumber := 43
	companyID := 12

	keepPath := "/poster/img/flat/keep.jpg"
	keepURL := host + "/" + bucket + keepPath
	oldPath := "/poster/img/flat/old.jpg"

	validPosterWithPhotos := &dto.PosterInputFlatDTO{
		UserID:         7,
		Price:          135000,
		Description:    "krutoy remont",
		CategoryAlias:  PropertyFlat,
		Area:           96.4,
		GeoLat:         55.7558,
		GeoLon:         37.6173,
		FlatCategoryID: 1,
		FlatNumber:     &flatNumber,
		FlatFloor:      3,
		Address:        "Arbatskaya 5k2",
		City:           "Moscow",
		District:       &district,
		FloorCount:     7,
		CompanyID:      &companyID,
		Features:       []string{"balcony", "parking"},
		Images: []dto.PhotoInputDTO{
			{
				URL:   &keepURL,
				Order: 1,
			},
			{
				FileHeader: &dto.FileInput{
					Filename:    "img2.jpg",
					Size:        10,
					ContentType: "image/jpeg",
					File:        io.NopCloser(strings.NewReader("img2")),
				},
				Order: 2,
			},
		},
	}

	validPosterWithoutNewFiles := &dto.PosterInputFlatDTO{
		UserID:         7,
		Price:          135000,
		Description:    "krutoy remont",
		CategoryAlias:  PropertyFlat,
		Area:           96.4,
		GeoLat:         55.7558,
		GeoLon:         37.6173,
		FlatCategoryID: 1,
		FlatNumber:     &flatNumber,
		FlatFloor:      3,
		Address:        "Arbatskaya 5k2",
		City:           "Moscow",
		District:       &district,
		FloorCount:     7,
		CompanyID:      &companyID,
		Features:       []string{"balcony", "parking"},
		Images: []dto.PhotoInputDTO{
			{
				URL:   &keepURL,
				Order: 1,
			},
			{
				URL:   &keepURL,
				Order: 2,
			},
		},
	}

	ids := &dto.PosterUpdateIDs{
		UserID:     7,
		PosterID:   33,
		PropertyID: 22,
		BuildingID: 11,
	}

	tests := []struct {
		name      string
		alias     string
		poster    *dto.PosterInputFlatDTO
		setupMock func(repo *mocks.MockPosterRepo, uow *mocks.MockUnitOfWork, file *mocks.MockFileRepo)
		want      *dto.CreatedPoster
		wantErr   error
	}{
		{
			name:   "OK",
			alias:  alias,
			poster: validPosterWithoutNewFiles,
			setupMock: func(repo *mocks.MockPosterRepo, uow *mocks.MockUnitOfWork, file *mocks.MockFileRepo) {
				repo.EXPECT().
					GetUpdateIDsByAlias(ctx, alias).
					Return(ids, nil).
					Times(1)

				repo.EXPECT().
					GetCityByName(ctx, "Moscow").
					Return(&entity.City{ID: 1, Name: "Moscow"}, nil).
					Times(1)

				repo.EXPECT().
					GetMetroStationByRadius(ctx, dto.GeographyDTO{Lat: 55.7558, Lon: 37.6173}, entity.Metre(MetroRadius)).
					Return([]entity.MetroStation{
						{ID: 5, StationName: "Arbatskaya"},
					}, nil).
					Times(1)

				repo.EXPECT().
					GetPhotoPathsByPosterID(ctx, 33).
					Return([]string{keepPath, keepPath}, nil).
					Times(1)

				uow.EXPECT().
					Do(ctx, gomock.Any()).
					DoAndReturn(func(_ context.Context, fn func(r usecase.UnitOfWork) error) error {
						return fn(uow)
					}).
					Times(1)

				repo.EXPECT().
					GetLastPriceHistoryByPosterID(ctx, 33).
					Return(&entity.PriceHistory{
						ID:        1,
						Price:     float64(100000),
						ChangedAt: time.Now().Add(-25 * time.Hour),
					}, nil).
					Times(1)

				repo.EXPECT().
					AddPriceHistory(ctx, 33, float64(135000)).
					Return(nil).
					Times(1)

				repo.EXPECT().
					Update(ctx, 33, gomock.Any()).
					Return(nil).
					Times(1)

				repo.EXPECT().
					UpdateProperty(ctx, 22, gomock.Any()).
					Return(nil).
					Times(1)

				repo.EXPECT().
					UpdateBuilding(ctx, 11, gomock.Any()).
					Return(nil).
					Times(1)

				repo.EXPECT().
					DeleteFacilitiesByPropertyID(ctx, 22).
					Return(nil).
					Times(1)

				repo.EXPECT().
					InsertFacilities(ctx, 22, []string{"balcony", "parking"}).
					Return(nil).
					Times(1)

				repo.EXPECT().
					UpdateFlat(ctx, gomock.Any()).
					Return(nil).
					Times(1)

				repo.EXPECT().
					DeletePhotosByPosterID(ctx, 33).
					Return(nil).
					Times(1)

				repo.EXPECT().
					InsertPhotos(ctx, 33, gomock.Any()).
					Return(nil).
					Times(1)

				repo.EXPECT().
					InsertMainPhoto(ctx, 33, gomock.Any()).
					Return(nil).
					Times(1)
			},
			want: &dto.CreatedPoster{
				ID: 33,
			},
			wantErr: nil,
		},
		{
			name:   "CreateCityIfNotFound",
			alias:  alias,
			poster: validPosterWithoutNewFiles,
			setupMock: func(repo *mocks.MockPosterRepo, uow *mocks.MockUnitOfWork, file *mocks.MockFileRepo) {
				repo.EXPECT().
					GetUpdateIDsByAlias(ctx, alias).
					Return(ids, nil).
					Times(1)

				repo.EXPECT().
					GetCityByName(ctx, "Moscow").
					Return(nil, entity.NotFoundError).
					Times(1)

				repo.EXPECT().
					CreateCity(ctx, "Moscow").
					Return(&entity.City{ID: 2, Name: "Moscow"}, nil).
					Times(1)

				repo.EXPECT().
					GetMetroStationByRadius(ctx, dto.GeographyDTO{Lat: 55.7558, Lon: 37.6173}, entity.Metre(MetroRadius)).
					Return([]entity.MetroStation{}, nil).
					Times(1)

				repo.EXPECT().
					GetPhotoPathsByPosterID(ctx, 33).
					Return([]string{keepPath, keepPath}, nil).
					Times(1)

				uow.EXPECT().
					Do(ctx, gomock.Any()).
					DoAndReturn(func(_ context.Context, fn func(r usecase.UnitOfWork) error) error {
						return fn(uow)
					}).
					Times(1)

				repo.EXPECT().
					GetLastPriceHistoryByPosterID(ctx, 33).
					Return(&entity.PriceHistory{
						ID:        1,
						Price:     float64(100000),
						ChangedAt: time.Now().Add(-25 * time.Hour),
					}, nil).
					Times(1)

				repo.EXPECT().
					AddPriceHistory(ctx, 33, float64(135000)).
					Return(nil).
					Times(1)

				repo.EXPECT().Update(ctx, 33, gomock.Any()).Return(nil).Times(1)
				repo.EXPECT().UpdateProperty(ctx, 22, gomock.Any()).Return(nil).Times(1)
				repo.EXPECT().UpdateBuilding(ctx, 11, gomock.Any()).Return(nil).Times(1)
				repo.EXPECT().DeleteFacilitiesByPropertyID(ctx, 22).Return(nil).Times(1)
				repo.EXPECT().InsertFacilities(ctx, 22, []string{"balcony", "parking"}).Return(nil).Times(1)
				repo.EXPECT().UpdateFlat(ctx, gomock.Any()).Return(nil).Times(1)
				repo.EXPECT().DeletePhotosByPosterID(ctx, 33).Return(nil).Times(1)
				repo.EXPECT().InsertPhotos(ctx, 33, gomock.Any()).Return(nil).Times(1)
				repo.EXPECT().InsertMainPhoto(ctx, 33, gomock.Any()).Return(nil).Times(1)
			},
			want: &dto.CreatedPoster{
				ID: 33,
			},
			wantErr: nil,
		},
		{
			name:   "NotOwner",
			alias:  alias,
			poster: validPosterWithoutNewFiles,
			setupMock: func(repo *mocks.MockPosterRepo, uow *mocks.MockUnitOfWork, file *mocks.MockFileRepo) {
				repo.EXPECT().
					GetUpdateIDsByAlias(ctx, alias).
					Return(&dto.PosterUpdateIDs{
						UserID:     999,
						PosterID:   33,
						PropertyID: 22,
						BuildingID: 11,
					}, nil).
					Times(1)
			},
			want:    nil,
			wantErr: entity.NotFoundError,
		},
		{
			name:   "GetUpdateIDsError",
			alias:  alias,
			poster: validPosterWithoutNewFiles,
			setupMock: func(repo *mocks.MockPosterRepo, uow *mocks.MockUnitOfWork, file *mocks.MockFileRepo) {
				repo.EXPECT().
					GetUpdateIDsByAlias(ctx, alias).
					Return(nil, entity.ServiceError).
					Times(1)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
		{
			name:   "MetroError",
			alias:  alias,
			poster: validPosterWithoutNewFiles,
			setupMock: func(repo *mocks.MockPosterRepo, uow *mocks.MockUnitOfWork, file *mocks.MockFileRepo) {
				repo.EXPECT().
					GetUpdateIDsByAlias(ctx, alias).
					Return(ids, nil).
					Times(1)

				repo.EXPECT().
					GetCityByName(ctx, "Moscow").
					Return(&entity.City{ID: 1, Name: "Moscow"}, nil).
					Times(1)

				repo.EXPECT().
					GetMetroStationByRadius(ctx, dto.GeographyDTO{Lat: 55.7558, Lon: 37.6173}, entity.Metre(MetroRadius)).
					Return(nil, entity.ServiceError).
					Times(1)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
		{
			name:   "GetPhotoPathsError",
			alias:  alias,
			poster: validPosterWithPhotos,
			setupMock: func(repo *mocks.MockPosterRepo, uow *mocks.MockUnitOfWork, file *mocks.MockFileRepo) {
				repo.EXPECT().
					GetUpdateIDsByAlias(ctx, alias).
					Return(ids, nil).
					Times(1)

				repo.EXPECT().
					GetCityByName(ctx, "Moscow").
					Return(&entity.City{ID: 1, Name: "Moscow"}, nil).
					Times(1)

				repo.EXPECT().
					GetMetroStationByRadius(ctx, dto.GeographyDTO{Lat: 55.7558, Lon: 37.6173}, entity.Metre(MetroRadius)).
					Return([]entity.MetroStation{}, nil).
					Times(1)

				repo.EXPECT().
					GetPhotoPathsByPosterID(ctx, 33).
					Return(nil, entity.ServiceError).
					Times(1)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
		{
			name:   "OKWithPhotos",
			alias:  alias,
			poster: validPosterWithPhotos,
			setupMock: func(repo *mocks.MockPosterRepo, uow *mocks.MockUnitOfWork, file *mocks.MockFileRepo) {
				repo.EXPECT().
					GetUpdateIDsByAlias(ctx, alias).
					Return(ids, nil).
					Times(1)

				repo.EXPECT().
					GetCityByName(ctx, "Moscow").
					Return(&entity.City{ID: 1, Name: "Moscow"}, nil).
					Times(1)

				repo.EXPECT().
					GetMetroStationByRadius(ctx, dto.GeographyDTO{Lat: 55.7558, Lon: 37.6173}, entity.Metre(MetroRadius)).
					Return([]entity.MetroStation{}, nil).
					Times(1)

				repo.EXPECT().
					GetPhotoPathsByPosterID(ctx, 33).
					Return([]string{
						keepPath,
						oldPath,
					}, nil).
					Times(1)

				file.EXPECT().
					Upload(ctx, gomock.Any(), gomock.Any(), int64(10), "image/jpeg").
					Return(nil).
					Times(1)

				uow.EXPECT().
					Do(ctx, gomock.Any()).
					DoAndReturn(func(_ context.Context, fn func(r usecase.UnitOfWork) error) error {
						return fn(uow)
					}).
					Times(1)

				repo.EXPECT().
					GetLastPriceHistoryByPosterID(ctx, 33).
					Return(&entity.PriceHistory{
						ID:        1,
						Price:     float64(100000),
						ChangedAt: time.Now().Add(-25 * time.Hour),
					}, nil).
					Times(1)

				repo.EXPECT().
					AddPriceHistory(ctx, 33, float64(135000)).
					Return(nil).
					Times(1)

				repo.EXPECT().Update(ctx, 33, gomock.Any()).Return(nil).Times(1)
				repo.EXPECT().UpdateProperty(ctx, 22, gomock.Any()).Return(nil).Times(1)
				repo.EXPECT().UpdateBuilding(ctx, 11, gomock.Any()).Return(nil).Times(1)
				repo.EXPECT().DeleteFacilitiesByPropertyID(ctx, 22).Return(nil).Times(1)
				repo.EXPECT().InsertFacilities(ctx, 22, []string{"balcony", "parking"}).Return(nil).Times(1)
				repo.EXPECT().UpdateFlat(ctx, gomock.Any()).Return(nil).Times(1)
				repo.EXPECT().DeletePhotosByPosterID(ctx, 33).Return(nil).Times(1)
				repo.EXPECT().InsertPhotos(ctx, 33, gomock.Any()).Return(nil).Times(1)
				repo.EXPECT().InsertMainPhoto(ctx, 33, gomock.Any()).Return(nil).Times(1)

				file.EXPECT().
					Delete(ctx, "poster/img/flat/old.jpg").
					Return(nil).
					Times(1)
			},
			want: &dto.CreatedPoster{
				ID: 33,
			},
			wantErr: nil,
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
				test.setupMock(mockRepo, mockUOW, mockFile)
			}

			uc := NewPosterUseCase(mockUOW, mockFile, nil, nil, nil, nil)

			res, err := uc.UpdateFlatPoster(ctx, test.alias, test.poster)

			if test.wantErr != nil {
				require.Error(t, err)
				require.ErrorContains(t, err, test.wantErr.Error())
				return
			}

			require.NoError(t, err)
			require.NotNil(t, res)
			require.Equal(t, test.want.ID, res.ID)
			require.Equal(t, test.alias, res.Alias)
		})
	}
}

func TestDeleteFlatPoster(t *testing.T) {
	ctx := context.Background()
	alias := "kvartira-na-arbate"
	userID := 7

	ids := &dto.PosterUpdateIDs{
		UserID:     7,
		PosterID:   33,
		PropertyID: 22,
		BuildingID: 11,
	}

	oldPaths := []string{
		"/poster/img/flat/old1.jpg",
		"/poster/img/flat/old2.jpg",
	}

	tests := []struct {
		name      string
		alias     string
		userID    int
		setupMock func(repo *mocks.MockPosterRepo, uow *mocks.MockUnitOfWork, file *mocks.MockFileRepo, search *mocks.MockSearchRepo)
		want      *dto.CreatedPoster
		wantErr   error
	}{
		{
			name:   "OK",
			alias:  alias,
			userID: userID,
			setupMock: func(repo *mocks.MockPosterRepo, uow *mocks.MockUnitOfWork, file *mocks.MockFileRepo, search *mocks.MockSearchRepo) {
				repo.EXPECT().
					GetUpdateIDsByAlias(ctx, alias).
					Return(ids, nil).
					Times(1)

				repo.EXPECT().
					GetPhotoPathsByPosterID(ctx, 33).
					Return(oldPaths, nil).
					Times(1)

				uow.EXPECT().
					Do(ctx, gomock.Any()).
					DoAndReturn(func(_ context.Context, fn func(r usecase.UnitOfWork) error) error {
						return fn(uow)
					}).
					Times(1)

				repo.EXPECT().
					DeletePhotosByPosterID(ctx, 33).
					Return(nil).
					Times(1)

				repo.EXPECT().
					DeleteFacilitiesByPropertyID(ctx, 22).
					Return(nil).
					Times(1)

				repo.EXPECT().
					Delete(ctx, 33).
					Return(nil).
					Times(1)

				repo.EXPECT().
					DeleteFlat(ctx, 22).
					Return(nil).
					Times(1)

				repo.EXPECT().
					DeleteProperty(ctx, 22).
					Return(nil).
					Times(1)

				repo.EXPECT().
					DeleteBuilding(ctx, 11).
					Return(nil).
					Times(1)

				file.EXPECT().
					Delete(ctx, "poster/img/flat/old1.jpg").
					Return(nil).
					Times(1)

				file.EXPECT().
					Delete(ctx, "poster/img/flat/old2.jpg").
					Return(nil).
					Times(1)

				search.EXPECT().
					DeletePoster(ctx, 33).
					Return(nil).
					Times(1)
			},
			want: &dto.CreatedPoster{
				ID:    33,
				Alias: alias,
			},
			wantErr: nil,
		},
		{
			name:   "OKWithoutPhotos",
			alias:  alias,
			userID: userID,
			setupMock: func(repo *mocks.MockPosterRepo, uow *mocks.MockUnitOfWork, file *mocks.MockFileRepo, search *mocks.MockSearchRepo) {
				repo.EXPECT().
					GetUpdateIDsByAlias(ctx, alias).
					Return(ids, nil).
					Times(1)

				repo.EXPECT().
					GetPhotoPathsByPosterID(ctx, 33).
					Return([]string{}, nil).
					Times(1)

				uow.EXPECT().
					Do(ctx, gomock.Any()).
					DoAndReturn(func(_ context.Context, fn func(r usecase.UnitOfWork) error) error {
						return fn(uow)
					}).
					Times(1)

				repo.EXPECT().
					DeleteFacilitiesByPropertyID(ctx, 22).
					Return(nil).
					Times(1)

				repo.EXPECT().
					Delete(ctx, 33).
					Return(nil).
					Times(1)

				repo.EXPECT().
					DeleteFlat(ctx, 22).
					Return(nil).
					Times(1)

				repo.EXPECT().
					DeleteProperty(ctx, 22).
					Return(nil).
					Times(1)

				repo.EXPECT().
					DeleteBuilding(ctx, 11).
					Return(nil).
					Times(1)

				search.EXPECT().
					DeletePoster(ctx, 33).
					Return(nil).
					Times(1)
			},
			want: &dto.CreatedPoster{
				ID:    33,
				Alias: alias,
			},
			wantErr: nil,
		},
		{
			name:   "GetUpdateIDsError",
			alias:  alias,
			userID: userID,
			setupMock: func(repo *mocks.MockPosterRepo, uow *mocks.MockUnitOfWork, file *mocks.MockFileRepo, search *mocks.MockSearchRepo) {
				repo.EXPECT().
					GetUpdateIDsByAlias(ctx, alias).
					Return(nil, entity.ServiceError).
					Times(1)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
		{
			name:   "NotOwner",
			alias:  alias,
			userID: userID,
			setupMock: func(repo *mocks.MockPosterRepo, uow *mocks.MockUnitOfWork, file *mocks.MockFileRepo, search *mocks.MockSearchRepo) {
				repo.EXPECT().
					GetUpdateIDsByAlias(ctx, alias).
					Return(&dto.PosterUpdateIDs{
						UserID:     999,
						PosterID:   33,
						PropertyID: 22,
						BuildingID: 11,
					}, nil).
					Times(1)
			},
			want:    nil,
			wantErr: entity.NotFoundError,
		},
		{
			name:   "GetPhotoPathsError",
			alias:  alias,
			userID: userID,
			setupMock: func(repo *mocks.MockPosterRepo, uow *mocks.MockUnitOfWork, file *mocks.MockFileRepo, search *mocks.MockSearchRepo) {
				repo.EXPECT().
					GetUpdateIDsByAlias(ctx, alias).
					Return(ids, nil).
					Times(1)

				repo.EXPECT().
					GetPhotoPathsByPosterID(ctx, 33).
					Return(nil, entity.ServiceError).
					Times(1)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
		{
			name:   "DeletePhotosError",
			alias:  alias,
			userID: userID,
			setupMock: func(repo *mocks.MockPosterRepo, uow *mocks.MockUnitOfWork, file *mocks.MockFileRepo, search *mocks.MockSearchRepo) {
				repo.EXPECT().
					GetUpdateIDsByAlias(ctx, alias).
					Return(ids, nil).
					Times(1)

				repo.EXPECT().
					GetPhotoPathsByPosterID(ctx, 33).
					Return(oldPaths, nil).
					Times(1)

				uow.EXPECT().
					Do(ctx, gomock.Any()).
					DoAndReturn(func(_ context.Context, fn func(r usecase.UnitOfWork) error) error {
						return fn(uow)
					}).
					Times(1)

				repo.EXPECT().
					DeletePhotosByPosterID(ctx, 33).
					Return(entity.ServiceError).
					Times(1)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
		{
			name:   "DeleteFacilitiesError",
			alias:  alias,
			userID: userID,
			setupMock: func(repo *mocks.MockPosterRepo, uow *mocks.MockUnitOfWork, file *mocks.MockFileRepo, search *mocks.MockSearchRepo) {
				repo.EXPECT().
					GetUpdateIDsByAlias(ctx, alias).
					Return(ids, nil).
					Times(1)

				repo.EXPECT().
					GetPhotoPathsByPosterID(ctx, 33).
					Return(oldPaths, nil).
					Times(1)

				uow.EXPECT().
					Do(ctx, gomock.Any()).
					DoAndReturn(func(_ context.Context, fn func(r usecase.UnitOfWork) error) error {
						return fn(uow)
					}).
					Times(1)

				repo.EXPECT().
					DeletePhotosByPosterID(ctx, 33).
					Return(nil).
					Times(1)

				repo.EXPECT().
					DeleteFacilitiesByPropertyID(ctx, 22).
					Return(entity.ServiceError).
					Times(1)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
		{
			name:   "DeletePosterError",
			alias:  alias,
			userID: userID,
			setupMock: func(repo *mocks.MockPosterRepo, uow *mocks.MockUnitOfWork, file *mocks.MockFileRepo, search *mocks.MockSearchRepo) {
				repo.EXPECT().
					GetUpdateIDsByAlias(ctx, alias).
					Return(ids, nil).
					Times(1)

				repo.EXPECT().
					GetPhotoPathsByPosterID(ctx, 33).
					Return(oldPaths, nil).
					Times(1)

				uow.EXPECT().
					Do(ctx, gomock.Any()).
					DoAndReturn(func(_ context.Context, fn func(r usecase.UnitOfWork) error) error {
						return fn(uow)
					}).
					Times(1)

				repo.EXPECT().
					DeletePhotosByPosterID(ctx, 33).
					Return(nil).
					Times(1)

				repo.EXPECT().
					DeleteFacilitiesByPropertyID(ctx, 22).
					Return(nil).
					Times(1)

				repo.EXPECT().
					Delete(ctx, 33).
					Return(entity.ServiceError).
					Times(1)
			},
			want:    nil,
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
			mockSearch := mocks.NewMockSearchRepo(ctrl)

			mockUOW.EXPECT().
				Posters().
				Return(mockRepo).
				AnyTimes()

			if test.setupMock != nil {
				test.setupMock(mockRepo, mockUOW, mockFile, mockSearch)
			}

			uc := NewPosterUseCase(mockUOW, mockFile, mockSearch, nil, nil, nil)

			res, err := uc.DeleteFlatPoster(ctx, test.alias, test.userID)

			if test.wantErr != nil {
				require.Error(t, err)
				require.ErrorContains(t, err, test.wantErr.Error())
				return
			}

			require.NoError(t, err)
			require.NotNil(t, res)
			require.Equal(t, test.want.ID, res.ID)
			require.Equal(t, test.want.Alias, res.Alias)
		})
	}
}

func TestAddViewPoster(t *testing.T) {
	ctx := context.Background()
	alias := "test-alias"
	userID := 1

	tests := []struct {
		name      string
		setup     func(mock *mocks.MockPosterRepo)
		wantError error
	}{
		{
			name: "OK",
			setup: func(mock *mocks.MockPosterRepo) {
				poster := &entity.PosterById{
					ID: 10,
				}

				mock.EXPECT().
					GetByAlias(ctx, alias, nil).
					Return(poster, nil)

				mock.EXPECT().
					AddView(ctx, userID, poster.ID)
			},
			wantError: nil,
		},
		{
			name: "RepoError",
			setup: func(mock *mocks.MockPosterRepo) {
				mock.EXPECT().
					GetByAlias(ctx, alias, nil).
					Return(nil, entity.ServiceError)
			},
			wantError: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uow := mocks.NewMockUnitOfWork(ctrl)
			posterRepo := mocks.NewMockPosterRepo(ctrl)
			fileRepo := mocks.NewMockFileRepo(ctrl)

			uow.EXPECT().
				Posters().
				Return(posterRepo).
				AnyTimes()

			if test.setup != nil {
				test.setup(posterRepo)
			}

			uc := NewPosterUseCase(uow, fileRepo, nil, nil, nil, nil)

			err := uc.AddViewPoster(ctx, alias, userID)

			if test.wantError != nil {
				require.Error(t, err)
				require.ErrorContains(t, err, test.wantError.Error())
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestGetViewsPoster(t *testing.T) {
	ctx := context.Background()
	alias := "test-alias"

	tests := []struct {
		name      string
		setup     func(mock *mocks.MockPosterRepo)
		want      int
		wantError error
	}{
		{
			name: "OK",
			setup: func(mock *mocks.MockPosterRepo) {
				poster := &entity.PosterById{
					ID: 10,
				}

				mock.EXPECT().
					GetByAlias(ctx, alias, nil).
					Return(poster, nil)

				mock.EXPECT().
					GetViewsCount(ctx, poster.ID).
					Return(5, nil)
			},
			want:      5,
			wantError: nil,
		},
		{
			name: "RepoErrorGetByAlias",
			setup: func(mock *mocks.MockPosterRepo) {
				mock.EXPECT().
					GetByAlias(ctx, alias, nil).
					Return(nil, entity.ServiceError)
			},
			want:      0,
			wantError: entity.ServiceError,
		},
		{
			name: "RepoErrorGetViews",
			setup: func(mock *mocks.MockPosterRepo) {
				poster := &entity.PosterById{
					ID: 10,
				}

				mock.EXPECT().
					GetByAlias(ctx, alias, nil).
					Return(poster, nil)

				mock.EXPECT().
					GetViewsCount(ctx, poster.ID).
					Return(0, entity.ServiceError)
			},
			want:      0,
			wantError: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			uow := mocks.NewMockUnitOfWork(ctrl)
			posterRepo := mocks.NewMockPosterRepo(ctrl)
			fileRepo := mocks.NewMockFileRepo(ctrl)

			uow.EXPECT().
				Posters().
				Return(posterRepo).
				AnyTimes()

			if test.setup != nil {
				test.setup(posterRepo)
			}

			uc := NewPosterUseCase(uow, fileRepo, nil, nil, nil, nil)

			res, err := uc.GetViewsPoster(ctx, alias)

			if test.wantError != nil {
				require.Error(t, err)
				require.ErrorContains(t, err, test.wantError.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, res)
		})
	}
}

func TestAddFavoritePoster(t *testing.T) {
	ctx := context.Background()
	alias := "kvartira-na-arbate"
	userID := 1

	existingPoster := &entity.PosterById{
		ID:    4,
		Alias: alias,
	}

	tests := []struct {
		name      string
		param     string
		setupMock func(m *mocks.MockPosterRepo)
		wantErr   error
	}{
		{
			name:  "ok",
			param: alias,
			setupMock: func(m *mocks.MockPosterRepo) {
				m.EXPECT().
					GetByAlias(ctx, alias, nil).
					Return(existingPoster, nil).
					Times(1)

				m.EXPECT().
					AddFavorite(ctx, userID, existingPoster.ID).
					Return(nil).
					Times(1)
			},
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
			wantErr: entity.NotFoundError,
		},
		{
			name:  "add_favorite_error",
			param: alias,
			setupMock: func(m *mocks.MockPosterRepo) {
				m.EXPECT().
					GetByAlias(ctx, alias, nil).
					Return(existingPoster, nil).
					Times(1)

				m.EXPECT().
					AddFavorite(ctx, userID, existingPoster.ID).
					Return(entity.ServiceError).
					Times(1)
			},
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

			uc := NewPosterUseCase(mockUOW, mockFile, nil, nil, nil, nil)

			err := uc.AddFavoritePoster(ctx, alias, userID)

			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestGetFavoritesPoster(t *testing.T) {
	ctx := context.Background()
	userID := 1

	existingListPoster := []entity.PosterFlat{
		{ID: 1, Price: 11111, Address: "street_1"},
		{ID: 2, Price: 22222, Address: "street_2"},
		{ID: 3, Price: 33333, Address: "street_3"},
	}

	expectedDTO := []dto.PosterCardDTO{
		{ID: 1, Price: 11111, Address: "street_1"},
		{ID: 2, Price: 22222, Address: "street_2"},
		{ID: 3, Price: 33333, Address: "street_3"},
	}

	tests := []struct {
		name      string
		setupMock func(m *mocks.MockPosterRepo)
		want      *dto.PostersResponse
		wantErr   error
	}{
		{
			name: "OK",
			setupMock: func(m *mocks.MockPosterRepo) {
				m.EXPECT().
					GetFavoritesFlatsByUserID(ctx, userID).
					Return(existingListPoster, nil).
					Times(1)

				m.EXPECT().
					CountFavoritesByUserID(ctx, userID).
					Return(len(existingListPoster), nil).
					Times(1)
			},
			want: &dto.PostersResponse{
				Posters: expectedDTO,
				Len:     len(existingListPoster),
			},
			wantErr: nil,
		},
		{
			name: "RepoError",
			setupMock: func(m *mocks.MockPosterRepo) {
				m.EXPECT().
					GetFavoritesFlatsByUserID(ctx, userID).
					Return(nil, entity.ServiceError).
					Times(1)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
		{
			name: "CountError",
			setupMock: func(m *mocks.MockPosterRepo) {
				m.EXPECT().
					GetFavoritesFlatsByUserID(ctx, userID).
					Return(existingListPoster, nil).
					Times(1)

				m.EXPECT().
					CountFavoritesByUserID(ctx, userID).
					Return(0, entity.ServiceError).
					Times(1)
			},
			want:    nil,
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

			uc := NewPosterUseCase(mockUOW, mockFile, nil, nil, nil, nil)

			got, err := uc.GetFavoritesPoster(ctx, userID)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want.Len, got.Len)
			require.Equal(t, len(test.want.Posters), len(got.Posters))

			for i, p := range test.want.Posters {
				require.Equal(t, p.ID, got.Posters[i].ID)
				require.Equal(t, p.Address, got.Posters[i].Address)
				require.Equal(t, p.Price, got.Posters[i].Price)
			}
		})
	}
}

func TestDeleteFavoritePoster(t *testing.T) {
	ctx := context.Background()
	alias := "kvartira-na-arbate"
	userID := 1

	existingPoster := &entity.PosterById{
		ID:    4,
		Alias: alias,
	}

	tests := []struct {
		name      string
		param     string
		setupMock func(m *mocks.MockPosterRepo)
		wantErr   error
	}{
		{
			name:  "ok",
			param: alias,
			setupMock: func(m *mocks.MockPosterRepo) {
				m.EXPECT().
					GetByAlias(ctx, alias, nil).
					Return(existingPoster, nil).
					Times(1)

				m.EXPECT().
					DeleteFavorite(ctx, userID, existingPoster.ID).
					Return(nil).
					Times(1)
			},
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
			wantErr: entity.NotFoundError,
		},
		{
			name:  "delete_error",
			param: alias,
			setupMock: func(m *mocks.MockPosterRepo) {
				m.EXPECT().
					GetByAlias(ctx, alias, nil).
					Return(existingPoster, nil).
					Times(1)

				m.EXPECT().
					DeleteFavorite(ctx, userID, existingPoster.ID).
					Return(entity.ServiceError).
					Times(1)
			},
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

			uc := NewPosterUseCase(mockUOW, mockFile, nil, nil, nil, nil)

			err := uc.DeleteFavoritePoster(ctx, alias, userID)

			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestGetPostersByCoords(t *testing.T) {
	ctx := context.Background()

	boundsLowZoom := dto.MapBounds{
		Zoom: 10,
	}

	boundsHighZoom := dto.MapBounds{
		Zoom: 15,
	}

	filters := dto.PostersFiltersDTO{}

	clusters := []entity.ClusterPoint{
		{Count: 2, Lat: 55.7, Lon: 37.6},
		{Count: 3, Lat: 55.8, Lon: 37.7},
	}

	posters := []entity.AnyPoint{
		{
			ID:  1,
			Lat: 55.7,
			Lon: 37.6,
		},
		{
			ID:  2,
			Lat: 55.8,
			Lon: 37.7,
		},
	}

	tests := []struct {
		name      string
		bounds    dto.MapBounds
		setupMock func(search *mocks.MockSearchRepo)
		wantLen   int
		wantErr   error
	}{
		{
			name:   "OK clusters (zoom < 13)",
			bounds: boundsLowZoom,
			setupMock: func(search *mocks.MockSearchRepo) {
				search.EXPECT().
					GetClustersByMapBounds(ctx, boundsLowZoom, gomock.Any()).
					Return(clusters, nil).
					Times(1)
			},
			wantLen: len(clusters),
			wantErr: nil,
		},
		{
			name:   "OK posters (zoom >= 13)",
			bounds: boundsHighZoom,
			setupMock: func(search *mocks.MockSearchRepo) {
				search.EXPECT().
					GetPostersByMapBounds(ctx, boundsHighZoom, gomock.Any()).
					Return(posters, nil).
					Times(1)
			},
			wantLen: len(posters),
			wantErr: nil,
		},
		{
			name:   "Clusters error",
			bounds: boundsLowZoom,
			setupMock: func(search *mocks.MockSearchRepo) {
				search.EXPECT().
					GetClustersByMapBounds(ctx, boundsLowZoom, gomock.Any()).
					Return(nil, entity.ServiceError).
					Times(1)
			},
			wantLen: 0,
			wantErr: entity.ServiceError,
		},
		{
			name:   "Posters error",
			bounds: boundsHighZoom,
			setupMock: func(search *mocks.MockSearchRepo) {
				search.EXPECT().
					GetPostersByMapBounds(ctx, boundsHighZoom, gomock.Any()).
					Return(nil, entity.ServiceError).
					Times(1)
			},
			wantLen: 0,
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSearch := mocks.NewMockSearchRepo(ctrl)
			mockUOW := mocks.NewMockUnitOfWork(ctrl)
			mockFile := mocks.NewMockFileRepo(ctrl)

			if test.setupMock != nil {
				test.setupMock(mockSearch)
			}

			uc := NewPosterUseCase(mockUOW, mockFile, mockSearch, nil, nil, nil)

			got, err := uc.GetPostersByCoords(ctx, test.bounds, filters)

			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			require.Equal(t, test.wantLen, got.Len)
		})
	}
}

func TestGetPostersByRadius(t *testing.T) {
	ctx := context.Background()

	point := dto.GeographyDTO{
		Lat: 55.7558,
		Lon: 37.6173,
	}

	existingPosters := []entity.Poster{
		{
			ID:      1,
			Alias:   "flat-1",
			Address: "street_1",
			Area:    50,
			Price:   11111,
		},
		{
			ID:      2,
			Alias:   "flat-2",
			Address: "street_2",
			Area:    60,
			Price:   22222,
		},
	}

	tests := []struct {
		name      string
		setupMock func(repo *mocks.MockPosterRepo)
		wantLen   int
		wantErr   error
	}{
		{
			name: "OK",
			setupMock: func(repo *mocks.MockPosterRepo) {
				repo.EXPECT().
					GetPostersByRadius(ctx, point, entity.Metre(10)).
					Return(existingPosters, nil).
					Times(1)
			},
			wantLen: len(existingPosters),
			wantErr: nil,
		},
		{
			name: "RepoError",
			setupMock: func(repo *mocks.MockPosterRepo) {
				repo.EXPECT().
					GetPostersByRadius(ctx, point, entity.Metre(10)).
					Return(nil, entity.ServiceError).
					Times(1)
			},
			wantLen: 0,
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

			uc := NewPosterUseCase(mockUOW, mockFile, nil, nil, nil, nil)

			got, err := uc.GetPostersByRadius(ctx, point)

			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.Len(t, got, test.wantLen)

			for i, p := range got {
				require.Equal(t, existingPosters[i].ID, p.ID)
				require.Equal(t, existingPosters[i].Address, p.Address)
				require.Equal(t, existingPosters[i].Price, p.Price)
				require.Equal(t, existingPosters[i].Area, p.Area)
			}
		})
	}
}

func TestGenerateDescription(t *testing.T) {
	ctx := context.Background()

	input := dto.GenerateDescriptionDTO{
		Category:     "квартира",
		FlatCategory: "2-комнатная",
		Area:         50.5,
		Features:     []string{"балкон", "парковка"},
	}

	mockResponse := &dto.ChatResult{
		Content: "Отличная квартира с балконом и парковкой.",
	}

	tests := []struct {
		name      string
		setupMock func(agent *mocks.MockLLMAgent)
		want      string
		wantErr   error
	}{
		{
			name: "OK",
			setupMock: func(agent *mocks.MockLLMAgent) {
				agent.EXPECT().
					Chat(
						ctx,
						gomock.Any(),
						gomock.Any(),
					).
					Return(mockResponse, nil).
					Times(1)
			},
			want:    mockResponse.Content,
			wantErr: nil,
		},
		{
			name: "AgentError",
			setupMock: func(agent *mocks.MockLLMAgent) {
				agent.EXPECT().
					Chat(
						ctx,
						gomock.Any(),
						gomock.Any(),
					).
					Return(nil, entity.ServiceError).
					Times(1)
			},
			want:    "",
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAgent := mocks.NewMockLLMAgent(ctrl)
			mockUOW := mocks.NewMockUnitOfWork(ctrl)
			mockFile := mocks.NewMockFileRepo(ctrl)

			if test.setupMock != nil {
				test.setupMock(mockAgent)
			}

			uc := &PosterUseCase{
				agent: mockAgent,
				uow:   mockUOW,
				file:  mockFile,
			}

			res, err := uc.GenerateDescription(ctx, input)

			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, res)
		})
	}
}

func TestGetPriceHistoryPoster(t *testing.T) {
	ctx := context.Background()
	alias := "kvartira-na-arbate"

	t1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	t3 := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)

	existingHistory := []entity.PriceHistory{
		{ID: 1, Price: 100000, ChangedAt: t1},
		{ID: 2, Price: 110000, ChangedAt: t2},
		{ID: 3, Price: 120000, ChangedAt: t3},
	}

	expectedDTO := []dto.PriceHistoryDTO{
		{Date: t1.Format("2006-01-02"), Price: 100000},
		{Date: t2.Format("2006-01-02"), Price: 110000},
		{Date: t3.Format("2006-01-02"), Price: 120000},
	}

	tests := []struct {
		name      string
		alias     string
		setupMock func(repo *mocks.MockPosterRepo)
		want      []dto.PriceHistoryDTO
		wantErr   error
	}{
		{
			name:  "OK",
			alias: alias,
			setupMock: func(repo *mocks.MockPosterRepo) {
				repo.EXPECT().
					GetPriceHistoryByAlias(ctx, alias).
					Return(existingHistory, nil).
					Times(1)
			},
			want:    expectedDTO,
			wantErr: nil,
		},
		{
			name:  "RepoError",
			alias: alias,
			setupMock: func(repo *mocks.MockPosterRepo) {
				repo.EXPECT().
					GetPriceHistoryByAlias(ctx, alias).
					Return(nil, entity.ServiceError).
					Times(1)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
		{
			name:  "EmptyHistory",
			alias: alias,
			setupMock: func(repo *mocks.MockPosterRepo) {
				repo.EXPECT().
					GetPriceHistoryByAlias(ctx, alias).
					Return([]entity.PriceHistory{}, nil).
					Times(1)
			},
			want:    []dto.PriceHistoryDTO{},
			wantErr: nil,
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

			uc := NewPosterUseCase(mockUOW, mockFile, nil, nil, nil, nil)

			got, err := uc.GetPriceHistoryPoster(ctx, test.alias)

			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, len(test.want), len(got))

			for i := range test.want {
				require.Equal(t, test.want[i].Price, got[i].Price)
				require.Equal(t, test.want[i].Date, got[i].Date)
			}
		})
	}
}

func TestGetPosterRoommates(t *testing.T) {
	ctx := context.Background()

	alias := "flat-1"

	existingRoommates := []dto.RoommateUserDTO{
		{
			ID:        1,
			FirstName: "Ivan",
			LastName:  "Ivanov",
			AvatarURL: lo.ToPtr("avatar.jpg"),
		},
		{
			ID:        2,
			FirstName: "Petr",
			LastName:  "Petrov",
			AvatarURL: nil,
		},
	}

	tests := []struct {
		name      string
		param     string
		setupMock func(m *mocks.MockPosterRepo)
		want      []dto.RoommateUserDTO
		wantErr   error
	}{
		{
			name:  "OK",
			param: alias,
			setupMock: func(m *mocks.MockPosterRepo) {
				m.EXPECT().
					GetPosterRoommates(ctx, alias).
					Return(existingRoommates, nil).
					Times(1)
			},
			want:    existingRoommates,
			wantErr: nil,
		},
		{
			name:  "RepoError",
			param: alias,
			setupMock: func(m *mocks.MockPosterRepo) {
				m.EXPECT().
					GetPosterRoommates(ctx, alias).
					Return(nil, entity.NotFoundError).
					Times(1)
			},
			want:    nil,
			wantErr: entity.NotFoundError,
		},
		{
			name:  "EmptyResult",
			param: alias,
			setupMock: func(m *mocks.MockPosterRepo) {
				m.EXPECT().
					GetPosterRoommates(ctx, alias).
					Return([]dto.RoommateUserDTO{}, nil).
					Times(1)
			},
			want:    []dto.RoommateUserDTO{},
			wantErr: nil,
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

			uc := NewPosterUseCase(mockUOW, mockFile, nil, nil)

			res, err := uc.GetPosterRoommates(ctx, test.param)

			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, len(test.want), len(res))

			for i, user := range test.want {
				require.Equal(t, user.ID, res[i].ID)
				require.Equal(t, user.FirstName, res[i].FirstName)
				require.Equal(t, user.LastName, res[i].LastName)
				require.Equal(t, user.AvatarURL, res[i].AvatarURL)
			}
		})
	}
}

func TestAddPosterRoommate(t *testing.T) {
	ctx := context.Background()

	alias := "flat-1"
	userID := 10

	tests := []struct {
		name      string
		setupMock func(repo *mocks.MockPosterRepo)
		wantErr   error
	}{
		{
			name: "OK",
			setupMock: func(repo *mocks.MockPosterRepo) {
				repo.EXPECT().
					HasRoommateForm(ctx, userID).
					Return(true, nil).
					Times(1)

				repo.EXPECT().
					AddPosterRoommate(ctx, alias, userID).
					Return(nil).
					Times(1)
			},
			wantErr: nil,
		},
		{
			name: "NoRoommateForm",
			setupMock: func(repo *mocks.MockPosterRepo) {
				repo.EXPECT().
					HasRoommateForm(ctx, userID).
					Return(false, nil).
					Times(1)
			},
			wantErr: entity.InvalidInput,
		},
		{
			name: "HasRoommateFormRepoError",
			setupMock: func(repo *mocks.MockPosterRepo) {
				repo.EXPECT().
					HasRoommateForm(ctx, userID).
					Return(false, entity.ServiceError).
					Times(1)
			},
			wantErr: entity.ServiceError,
		},
		{
			name: "AddPosterRoommateRepoError",
			setupMock: func(repo *mocks.MockPosterRepo) {
				repo.EXPECT().
					HasRoommateForm(ctx, userID).
					Return(true, nil).
					Times(1)

				repo.EXPECT().
					AddPosterRoommate(ctx, alias, userID).
					Return(entity.NotFoundError).
					Times(1)
			},
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

			uc := NewPosterUseCase(mockUOW, mockFile, nil, nil)

			err := uc.AddPosterRoommate(ctx, alias, userID)

			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestDeletePosterRoommate(t *testing.T) {
	ctx := context.Background()

	alias := "flat-1"
	userID := 10

	tests := []struct {
		name      string
		setupMock func(repo *mocks.MockPosterRepo)
		wantErr   error
	}{
		{
			name: "OK",
			setupMock: func(repo *mocks.MockPosterRepo) {
				repo.EXPECT().
					DeletePosterRoommate(ctx, alias, userID).
					Return(nil).
					Times(1)
			},
			wantErr: nil,
		},
		{
			name: "NotFound",
			setupMock: func(repo *mocks.MockPosterRepo) {
				repo.EXPECT().
					DeletePosterRoommate(ctx, alias, userID).
					Return(entity.NotFoundError).
					Times(1)
			},
			wantErr: entity.NotFoundError,
		},
		{
			name: "RepoError",
			setupMock: func(repo *mocks.MockPosterRepo) {
				repo.EXPECT().
					DeletePosterRoommate(ctx, alias, userID).
					Return(entity.ServiceError).
					Times(1)
			},
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

			uc := NewPosterUseCase(mockUOW, mockFile, nil, nil)

			err := uc.DeletePosterRoommate(ctx, alias, userID)

			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}
