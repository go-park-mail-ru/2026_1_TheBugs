package poster

import (
	"context"
	"io"
	"strings"
	"testing"

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

func TestSearchPostersUseCase(t *testing.T) {
	ctx := context.Background()

	existingSearchResponse := &dto.PostersResponse{
		Posters: []dto.PosterCardDTO{
			{ID: 1},
			{ID: 2},
			{ID: 3},
		},
		Len: 3,
	}

	existingPosterFlats := []entity.PosterFlat{
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
		params    dto.PostersFiltersDTO
		setupMock func(searchMock *mocks.MockSearchRepo, uowMock *mocks.MockUnitOfWork)
		want      *dto.PostersResponse
		wantErr   error
	}{
		{
			name: "OK",
			params: dto.PostersFiltersDTO{
				Limit:  12,
				Offset: 0,
			},
			setupMock: func(searchMock *mocks.MockSearchRepo, uowMock *mocks.MockUnitOfWork) {
				searchMock.EXPECT().SearchPosters(ctx, dto.PostersFiltersDTO{
					Limit:  12,
					Offset: 0,
				}).Return(existingSearchResponse, nil).Times(1)

				mockPosterRepo := mocks.NewMockPosterRepo(gomock.NewController(t))
				uowMock.EXPECT().Posters().Return(mockPosterRepo).Times(1)
				mockPosterRepo.EXPECT().GetFlatsByIDs(ctx, []int{1, 2, 3}).Return(existingPosterFlats, nil).Times(1)
			},
			want: &dto.PostersResponse{
				Posters: expectedDTO,
				Len:     3,
			},
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
			setupMock: func(searchMock *mocks.MockSearchRepo, uowMock *mocks.MockUnitOfWork) {
				searchMock.EXPECT().SearchPosters(ctx, dto.PostersFiltersDTO{
					Limit:  MaxPostersLimit,
					Offset: 0,
				}).Return(existingSearchResponse, nil).Times(1)

				mockPosterRepo := mocks.NewMockPosterRepo(gomock.NewController(t))
				uowMock.EXPECT().Posters().Return(mockPosterRepo).Times(1)
				mockPosterRepo.EXPECT().GetFlatsByIDs(ctx, []int{1, 2, 3}).Return(existingPosterFlats, nil).Times(1)
			},
			want: &dto.PostersResponse{
				Posters: expectedDTO,
				Len:     3,
			},
			wantErr: nil,
		},
		{
			name: "SearchError",
			params: dto.PostersFiltersDTO{
				Limit:  12,
				Offset: 0,
			},
			setupMock: func(searchMock *mocks.MockSearchRepo, uowMock *mocks.MockUnitOfWork) {
				searchMock.EXPECT().SearchPosters(ctx, dto.PostersFiltersDTO{
					Limit:  12,
					Offset: 0,
				}).Return(nil, entity.NotFoundError).Times(1)
			},
			want:    nil,
			wantErr: entity.NotFoundError,
		},
		{
			name: "GetFlatsByIDs error",
			params: dto.PostersFiltersDTO{
				Limit:  12,
				Offset: 0,
			},
			setupMock: func(searchMock *mocks.MockSearchRepo, uowMock *mocks.MockUnitOfWork) {
				searchMock.EXPECT().SearchPosters(ctx, dto.PostersFiltersDTO{
					Limit:  12,
					Offset: 0,
				}).Return(existingSearchResponse, nil).Times(1)

				mockPosterRepo := mocks.NewMockPosterRepo(gomock.NewController(t))
				uowMock.EXPECT().Posters().Return(mockPosterRepo).Times(1)
				mockPosterRepo.EXPECT().GetFlatsByIDs(ctx, []int{1, 2, 3}).Return(nil, entity.NotFoundError).Times(1)
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
			setupMock: func(searchMock *mocks.MockSearchRepo, uowMock *mocks.MockUnitOfWork) {
				emptyResponse := &dto.PostersResponse{Posters: []dto.PosterCardDTO{}, Len: 0}
				searchMock.EXPECT().SearchPosters(ctx, dto.PostersFiltersDTO{
					Limit:  12,
					Offset: 0,
				}).Return(emptyResponse, nil).Times(1)

				mockPosterRepo := mocks.NewMockPosterRepo(gomock.NewController(t))
				uowMock.EXPECT().Posters().Return(mockPosterRepo).Times(1)
				mockPosterRepo.EXPECT().GetFlatsByIDs(ctx, []int{}).Return([]entity.PosterFlat{}, nil).Times(1)
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
				test.setupMock(mockSearch, mockUOW)
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

			uc := NewPosterUseCase(mockUOW, mockFile, nil)

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

			uc := NewPosterUseCase(mockUOW, mockFile, nil)

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

			uc := NewPosterUseCase(mockUOW, mockFile, nil)

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
		setupMock func(repo *mocks.MockPosterRepo, uow *mocks.MockUnitOfWork, file *mocks.MockFileRepo)
		want      *dto.CreatedPoster
		wantErr   error
	}{
		{
			name:   "OK",
			alias:  alias,
			userID: userID,
			setupMock: func(repo *mocks.MockPosterRepo, uow *mocks.MockUnitOfWork, file *mocks.MockFileRepo) {
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
			setupMock: func(repo *mocks.MockPosterRepo, uow *mocks.MockUnitOfWork, file *mocks.MockFileRepo) {
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
			name:   "NotOwner",
			alias:  alias,
			userID: userID,
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
			name:   "GetPhotoPathsError",
			alias:  alias,
			userID: userID,
			setupMock: func(repo *mocks.MockPosterRepo, uow *mocks.MockUnitOfWork, file *mocks.MockFileRepo) {
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
			setupMock: func(repo *mocks.MockPosterRepo, uow *mocks.MockUnitOfWork, file *mocks.MockFileRepo) {
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
			setupMock: func(repo *mocks.MockPosterRepo, uow *mocks.MockUnitOfWork, file *mocks.MockFileRepo) {
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
			setupMock: func(repo *mocks.MockPosterRepo, uow *mocks.MockUnitOfWork, file *mocks.MockFileRepo) {
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

			mockUOW.EXPECT().
				Posters().
				Return(mockRepo).
				AnyTimes()

			if test.setupMock != nil {
				test.setupMock(mockRepo, mockUOW, mockFile)
			}

			uc := NewPosterUseCase(mockUOW, mockFile, nil)

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

			uc := NewPosterUseCase(uow, fileRepo, nil)

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

			uc := NewPosterUseCase(uow, fileRepo, nil)

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
