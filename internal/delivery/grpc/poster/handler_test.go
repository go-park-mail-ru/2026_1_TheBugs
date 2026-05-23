package poster

import (
	"context"
	"io"
	"testing"

	posterpb "github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/generated/poster"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/mocks"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/golang/mock/gomock"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestPosterServiceServer_SearchPosters(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	validReq := &posterpb.SearchPostersRequest{
		Limit:      12,
		Offset:     0,
		Facilities: []string{"balcony", "parking"},
	}

	expectedResponse := &dto.PostersResponse{
		Len: 2,
		Posters: []dto.PosterCardDTO{
			{ID: 1, Alias: "flat-1", Price: 120000, Address: "street_1"},
			{ID: 2, Alias: "flat-2", Price: 150000, Address: "street_2"},
		},
	}

	tests := []struct {
		name      string
		req       *posterpb.SearchPostersRequest
		setupMock func(mockUC *mocks.MockPostersUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					SearchPostersUseCase(ctx, gomock.Any()).
					Return(expectedResponse, nil)
			},
			wantErr: false,
		},
		{
			name: "invalid limit",
			req: &posterpb.SearchPostersRequest{
				Limit:  0,
				Offset: 0,
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "invalid offset",
			req: &posterpb.SearchPostersRequest{
				Limit:  12,
				Offset: -1,
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "use case error",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					SearchPostersUseCase(ctx, gomock.Any()).
					Return(nil, entity.ServiceError)
			},
			wantErr:  true,
			wantCode: codes.Internal,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockPostersUseCase(ctrl)
			if test.setupMock != nil {
				test.setupMock(mockUC)
			}

			server := NewPosterServiceServer(mockUC)
			resp, err := server.SearchPosters(ctx, test.req)

			if test.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, test.wantCode, st.Code())
				require.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, int32(expectedResponse.Len), resp.Len)
			require.Len(t, resp.Posters, len(expectedResponse.Posters))
		})
	}
}

func TestPosterServiceServer_GetPosterByAlias(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	alias := "kvartira-na-arbate"
	userID := int64(7)

	validReq := &posterpb.GetPosterByAliasRequest{
		PosterAlias: alias,
		UserId:      &userID,
	}

	expectedPoster := &dto.PosterDTO{
		ID:      33,
		Alias:   alias,
		Price:   135000,
		Address: "Arbatskaya 5k2",
	}

	tests := []struct {
		name      string
		req       *posterpb.GetPosterByAliasRequest
		setupMock func(mockUC *mocks.MockPostersUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					GetPosterByAliasUseCase(ctx, alias, gomock.Any()).
					Return(expectedPoster, nil)
			},
			wantErr: false,
		},
		{
			name: "missing alias",
			req: &posterpb.GetPosterByAliasRequest{
				PosterAlias: "",
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "use case error",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					GetPosterByAliasUseCase(ctx, alias, gomock.Any()).
					Return(nil, entity.ServiceError)
			},
			wantErr:  true,
			wantCode: codes.Internal,
		},
		{
			name: "not found",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					GetPosterByAliasUseCase(ctx, alias, gomock.Any()).
					Return(nil, entity.NotFoundError)
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockPostersUseCase(ctrl)
			if test.setupMock != nil {
				test.setupMock(mockUC)
			}

			server := NewPosterServiceServer(mockUC)
			resp, err := server.GetPosterByAlias(ctx, test.req)

			if test.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, test.wantCode, st.Code())
				require.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotNil(t, resp.Poster)
		})
	}
}

func TestPosterServiceServer_GetPostersByUserID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	userID := int64(7)

	validReq := &posterpb.GetPostersByUserIDRequest{
		UserId: userID,
	}

	expectedPosters := []dto.MyPosterDTO{
		{
			ID:      1,
			Alias:   "flat-1",
			Price:   120000,
			Address: "street_1",
			Area:    55,
		},
		{
			ID:      2,
			Alias:   "flat-2",
			Price:   150000,
			Address: "street_2",
			Area:    70,
		},
	}

	tests := []struct {
		name      string
		req       *posterpb.GetPostersByUserIDRequest
		setupMock func(mockUC *mocks.MockPostersUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					GetPosterByUserID(ctx, int(userID)).
					Return(expectedPosters, nil)
			},
			wantErr: false,
		},
		{
			name: "invalid user id",
			req: &posterpb.GetPostersByUserIDRequest{
				UserId: 0,
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "use case error",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					GetPosterByUserID(ctx, int(userID)).
					Return(nil, entity.ServiceError)
			},
			wantErr:  true,
			wantCode: codes.Internal,
		},
		{
			name: "not found",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					GetPosterByUserID(ctx, int(userID)).
					Return(nil, entity.NotFoundError)
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockPostersUseCase(ctrl)
			if test.setupMock != nil {
				test.setupMock(mockUC)
			}

			server := NewPosterServiceServer(mockUC)
			resp, err := server.GetPostersByUserID(ctx, test.req)

			if test.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, test.wantCode, st.Code())
				require.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Len(t, resp.Posters, len(expectedPosters))
		})
	}
}

func TestPosterServiceServer_AddViewPoster(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	alias := "kvartira-na-arbate"
	userID := int64(7)

	validReq := &posterpb.AddViewPosterRequest{
		Alias:  alias,
		UserId: userID,
	}

	tests := []struct {
		name      string
		req       *posterpb.AddViewPosterRequest
		setupMock func(mockUC *mocks.MockPostersUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					AddViewPoster(ctx, alias, int(userID)).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "missing alias",
			req: &posterpb.AddViewPosterRequest{
				Alias:  "",
				UserId: userID,
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "invalid user id",
			req: &posterpb.AddViewPosterRequest{
				Alias:  alias,
				UserId: 0,
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "use case error",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					AddViewPoster(ctx, alias, int(userID)).
					Return(entity.ServiceError)
			},
			wantErr:  true,
			wantCode: codes.Internal,
		},
		{
			name: "not found",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					AddViewPoster(ctx, alias, int(userID)).
					Return(entity.NotFoundError)
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockPostersUseCase(ctrl)
			if test.setupMock != nil {
				test.setupMock(mockUC)
			}

			server := NewPosterServiceServer(mockUC)
			resp, err := server.AddViewPoster(ctx, test.req)

			if test.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, test.wantCode, st.Code())
				require.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
		})
	}
}

func TestPosterServiceServer_GetViewsPoster(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	alias := "kvartira-na-arbate"

	validReq := &posterpb.GetViewsPosterRequest{
		Alias: alias,
	}

	expectedViews := 42

	tests := []struct {
		name      string
		req       *posterpb.GetViewsPosterRequest
		setupMock func(mockUC *mocks.MockPostersUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					GetViewsPoster(ctx, alias).
					Return(expectedViews, nil)
			},
			wantErr: false,
		},
		{
			name: "missing alias",
			req: &posterpb.GetViewsPosterRequest{
				Alias: "",
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "use case error",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					GetViewsPoster(ctx, alias).
					Return(0, entity.ServiceError)
			},
			wantErr:  true,
			wantCode: codes.Internal,
		},
		{
			name: "not found",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					GetViewsPoster(ctx, alias).
					Return(0, entity.NotFoundError)
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockPostersUseCase(ctrl)
			if test.setupMock != nil {
				test.setupMock(mockUC)
			}

			server := NewPosterServiceServer(mockUC)
			resp, err := server.GetViewsPoster(ctx, test.req)

			if test.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, test.wantCode, st.Code())
				require.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, int32(expectedViews), resp.Views)
		})
	}
}

func TestPosterServiceServer_AddFavoritePoster(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	alias := "kvartira-na-arbate"
	userID := int64(7)

	validReq := &posterpb.AddFavoritePosterRequest{
		Alias:  alias,
		UserId: userID,
	}

	tests := []struct {
		name      string
		req       *posterpb.AddFavoritePosterRequest
		setupMock func(mockUC *mocks.MockPostersUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					AddFavoritePoster(ctx, alias, int(userID)).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "missing alias",
			req: &posterpb.AddFavoritePosterRequest{
				Alias:  "",
				UserId: userID,
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "invalid user id",
			req: &posterpb.AddFavoritePosterRequest{
				Alias:  alias,
				UserId: 0,
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "use case error",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					AddFavoritePoster(ctx, alias, int(userID)).
					Return(entity.ServiceError)
			},
			wantErr:  true,
			wantCode: codes.Internal,
		},
		{
			name: "not found",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					AddFavoritePoster(ctx, alias, int(userID)).
					Return(entity.NotFoundError)
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockPostersUseCase(ctrl)
			if test.setupMock != nil {
				test.setupMock(mockUC)
			}

			server := NewPosterServiceServer(mockUC)
			resp, err := server.AddFavoritePoster(ctx, test.req)

			if test.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, test.wantCode, st.Code())
				require.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
		})
	}
}

func TestPosterServiceServer_GetFavoritePosters(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	userID := int64(7)

	validReq := &posterpb.GetFavoritePostersRequest{
		UserId: userID,
	}

	expectedResponse := &dto.PostersResponse{
		Len: 2,
		Posters: []dto.PosterCardDTO{
			{ID: 1, Alias: "flat-1", Price: 120000, Address: "street_1"},
			{ID: 2, Alias: "flat-2", Price: 150000, Address: "street_2"},
		},
	}

	tests := []struct {
		name      string
		req       *posterpb.GetFavoritePostersRequest
		setupMock func(mockUC *mocks.MockPostersUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					GetFavoritesPoster(ctx, int(userID)).
					Return(expectedResponse, nil)
			},
			wantErr: false,
		},
		{
			name: "invalid user id",
			req: &posterpb.GetFavoritePostersRequest{
				UserId: 0,
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "use case error",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					GetFavoritesPoster(ctx, int(userID)).
					Return(nil, entity.ServiceError)
			},
			wantErr:  true,
			wantCode: codes.Internal,
		},
		{
			name: "not found",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					GetFavoritesPoster(ctx, int(userID)).
					Return(nil, entity.NotFoundError)
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockPostersUseCase(ctrl)
			if test.setupMock != nil {
				test.setupMock(mockUC)
			}

			server := NewPosterServiceServer(mockUC)
			resp, err := server.GetFavoritePosters(ctx, test.req)

			if test.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, test.wantCode, st.Code())
				require.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, int32(expectedResponse.Len), resp.Len)
			require.Len(t, resp.Posters, len(expectedResponse.Posters))
		})
	}
}

func TestPosterServiceServer_DeleteFavoritePoster(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	alias := "kvartira-na-arbate"
	userID := int64(7)

	validReq := &posterpb.DeleteFavoritePosterRequest{
		Alias:  alias,
		UserId: userID,
	}

	tests := []struct {
		name      string
		req       *posterpb.DeleteFavoritePosterRequest
		setupMock func(mockUC *mocks.MockPostersUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					DeleteFavoritePoster(ctx, alias, int(userID)).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "missing alias",
			req: &posterpb.DeleteFavoritePosterRequest{
				Alias:  "",
				UserId: userID,
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "invalid user id",
			req: &posterpb.DeleteFavoritePosterRequest{
				Alias:  alias,
				UserId: 0,
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "use case error",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					DeleteFavoritePoster(ctx, alias, int(userID)).
					Return(entity.ServiceError)
			},
			wantErr:  true,
			wantCode: codes.Internal,
		},
		{
			name: "not found",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					DeleteFavoritePoster(ctx, alias, int(userID)).
					Return(entity.NotFoundError)
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockPostersUseCase(ctrl)
			if test.setupMock != nil {
				test.setupMock(mockUC)
			}

			server := NewPosterServiceServer(mockUC)
			resp, err := server.DeleteFavoritePoster(ctx, test.req)

			if test.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, test.wantCode, st.Code())
				require.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
		})
	}
}

func TestPosterServiceServer_GetFavoritesCountPoster(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	alias := "kvartira-na-arbate"
	userID := int64(7)

	validReq := &posterpb.GetFavoritesCountPosterRequest{
		PosterAlias: alias,
		UserId:      &userID,
	}

	expectedCount := 5
	expectedIsFavorite := true

	tests := []struct {
		name      string
		req       *posterpb.GetFavoritesCountPosterRequest
		setupMock func(mockUC *mocks.MockPostersUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					GetFavoritesCountPoster(ctx, alias, gomock.Any()).
					Return(expectedCount, expectedIsFavorite, nil)
			},
			wantErr: false,
		},
		{
			name: "success without user id",
			req: &posterpb.GetFavoritesCountPosterRequest{
				PosterAlias: alias,
			},
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					GetFavoritesCountPoster(ctx, alias, nil).
					Return(expectedCount, false, nil)
			},
			wantErr: false,
		},
		{
			name: "missing alias",
			req: &posterpb.GetFavoritesCountPosterRequest{
				PosterAlias: "",
				UserId:      &userID,
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "use case error",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					GetFavoritesCountPoster(ctx, alias, gomock.Any()).
					Return(0, false, entity.ServiceError)
			},
			wantErr:  true,
			wantCode: codes.Internal,
		},
		{
			name: "not found",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					GetFavoritesCountPoster(ctx, alias, gomock.Any()).
					Return(0, false, entity.NotFoundError)
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockPostersUseCase(ctrl)
			if test.setupMock != nil {
				test.setupMock(mockUC)
			}

			server := NewPosterServiceServer(mockUC)
			resp, err := server.GetFavoritesCountPoster(ctx, test.req)

			if test.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, test.wantCode, st.Code())
				require.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
		})
	}
}

func TestPosterServiceServer_GetPostersByCoords(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	validReq := &posterpb.GetPostersByCoordsRequest{
		Bounds: &posterpb.MapBounds{
			Bbox: &posterpb.BBox{
				SouthWest: &posterpb.Geography{
					Lat: 55.7,
					Lon: 37.6,
				},
				NorthEast: &posterpb.Geography{
					Lat: 55.8,
					Lon: 37.7,
				},
			},
			Zoom: 12,
		},
		Filters: &posterpb.SearchPostersRequest{
			Limit:      12,
			Offset:     0,
			Facilities: []string{"balcony", "parking"},
		},
	}

	expectedResponse := &dto.GeoJSONFeatureResponse{
		Len: 2,
		Posters: []dto.GeoJSONFeature{
			{
				Type: "Feature",
				Geometry: dto.Geometry{
					Type:        "Point",
					Coordinates: []float64{37.6, 55.7},
				},
				Properties: map[string]any{
					"id": float64(1),
				},
			},
			{
				Type: "Feature",
				Geometry: dto.Geometry{
					Type:        "Point",
					Coordinates: []float64{37.7, 55.8},
				},
				Properties: map[string]any{
					"id": float64(2),
				},
			},
		},
	}

	tests := []struct {
		name      string
		req       *posterpb.GetPostersByCoordsRequest
		setupMock func(mockUC *mocks.MockPostersUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					GetPostersByCoords(ctx, gomock.Any(), gomock.Any()).
					Return(expectedResponse, nil)
			},
			wantErr: false,
		},
		{
			name: "missing bounds",
			req: &posterpb.GetPostersByCoordsRequest{
				Bounds:  nil,
				Filters: validReq.Filters,
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "missing bbox",
			req: &posterpb.GetPostersByCoordsRequest{
				Bounds: &posterpb.MapBounds{
					Bbox: nil,
					Zoom: 12,
				},
				Filters: validReq.Filters,
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "missing south west",
			req: &posterpb.GetPostersByCoordsRequest{
				Bounds: &posterpb.MapBounds{
					Bbox: &posterpb.BBox{
						SouthWest: nil,
						NorthEast: validReq.Bounds.Bbox.NorthEast,
					},
					Zoom: 12,
				},
				Filters: validReq.Filters,
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "missing filters",
			req: &posterpb.GetPostersByCoordsRequest{
				Bounds:  validReq.Bounds,
				Filters: nil,
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "invalid limit",
			req: &posterpb.GetPostersByCoordsRequest{
				Bounds: validReq.Bounds,
				Filters: &posterpb.SearchPostersRequest{
					Limit:  0,
					Offset: 0,
				},
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "invalid offset",
			req: &posterpb.GetPostersByCoordsRequest{
				Bounds: validReq.Bounds,
				Filters: &posterpb.SearchPostersRequest{
					Limit:  12,
					Offset: -1,
				},
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "use case error",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					GetPostersByCoords(ctx, gomock.Any(), gomock.Any()).
					Return(nil, entity.ServiceError)
			},
			wantErr:  true,
			wantCode: codes.Internal,
		},
		{
			name: "not found",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					GetPostersByCoords(ctx, gomock.Any(), gomock.Any()).
					Return(nil, entity.NotFoundError)
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockPostersUseCase(ctrl)
			if test.setupMock != nil {
				test.setupMock(mockUC)
			}

			server := NewPosterServiceServer(mockUC)
			resp, err := server.GetPostersByCoords(ctx, test.req)

			if test.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, test.wantCode, st.Code())
				require.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, int32(expectedResponse.Len), resp.Len)
			require.Len(t, resp.Features, len(expectedResponse.Posters))
		})
	}
}

func TestPosterServiceServer_GetPostersByRadius(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	validReq := &posterpb.GetPostersByRadiusRequest{
		Point: &posterpb.Geography{
			Lat: 55.7,
			Lon: 37.6,
		},
	}

	expectedPosters := []dto.MyPosterDTO{
		{ID: 1, Alias: "flat-1", Price: 120000, Address: "street_1", Area: 55},
		{ID: 2, Alias: "flat-2", Price: 150000, Address: "street_2", Area: 70},
	}

	tests := []struct {
		name      string
		req       *posterpb.GetPostersByRadiusRequest
		setupMock func(mockUC *mocks.MockPostersUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					GetPostersByRadius(ctx, gomock.Any()).
					Return(expectedPosters, nil)
			},
			wantErr: false,
		},
		{
			name: "missing point",
			req: &posterpb.GetPostersByRadiusRequest{
				Point: nil,
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "use case error",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					GetPostersByRadius(ctx, gomock.Any()).
					Return(nil, entity.ServiceError)
			},
			wantErr:  true,
			wantCode: codes.Internal,
		},
		{
			name: "not found",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					GetPostersByRadius(ctx, gomock.Any()).
					Return(nil, entity.NotFoundError)
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockPostersUseCase(ctrl)
			if test.setupMock != nil {
				test.setupMock(mockUC)
			}

			server := NewPosterServiceServer(mockUC)
			resp, err := server.GetPostersByRadius(ctx, test.req)

			if test.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, test.wantCode, st.Code())
				require.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Len(t, resp.Posters, len(expectedPosters))
		})
	}
}

func TestPosterServiceServer_GenerateDescription(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	validReq := &posterpb.GenerateDescriptionRequest{
		Category:     "flat",
		Area:         55.5,
		FlatCategory: "2-room",
		City:         "Moscow",
		Features:     []string{"balcony", "parking"},
	}

	expectedDescription := "Отличная квартира с балконом и парковкой."

	tests := []struct {
		name      string
		req       *posterpb.GenerateDescriptionRequest
		setupMock func(mockUC *mocks.MockPostersUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					GenerateDescription(ctx, gomock.Any()).
					Return(expectedDescription, nil)
			},
			wantErr: false,
		},
		{
			name: "missing category",
			req: &posterpb.GenerateDescriptionRequest{
				Category:     "",
				Area:         55.5,
				FlatCategory: "2-room",
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "invalid area",
			req: &posterpb.GenerateDescriptionRequest{
				Category:     "flat",
				Area:         0,
				FlatCategory: "2-room",
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "missing flat category",
			req: &posterpb.GenerateDescriptionRequest{
				Category:     "flat",
				Area:         55.5,
				FlatCategory: "",
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "use case error",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					GenerateDescription(ctx, gomock.Any()).
					Return("", entity.ServiceError)
			},
			wantErr:  true,
			wantCode: codes.Internal,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockPostersUseCase(ctrl)
			if test.setupMock != nil {
				test.setupMock(mockUC)
			}

			server := NewPosterServiceServer(mockUC)
			resp, err := server.GenerateDescription(ctx, test.req)

			if test.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, test.wantCode, st.Code())
				require.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, expectedDescription, resp.Description)
		})
	}
}

func TestPosterServiceServer_GetPriceHistoryPoster(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	alias := "kvartira-na-arbate"

	validReq := &posterpb.GetPriceHistoryPosterRequest{
		PosterAlias: alias,
	}

	expectedHistory := []dto.PriceHistoryDTO{
		{Date: "2024-01-01", Price: 100000},
		{Date: "2024-02-01", Price: 120000},
	}

	tests := []struct {
		name      string
		req       *posterpb.GetPriceHistoryPosterRequest
		setupMock func(mockUC *mocks.MockPostersUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					GetPriceHistoryPoster(ctx, alias).
					Return(expectedHistory, nil)
			},
			wantErr: false,
		},
		{
			name: "missing alias",
			req: &posterpb.GetPriceHistoryPosterRequest{
				PosterAlias: "",
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "use case error",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					GetPriceHistoryPoster(ctx, alias).
					Return(nil, entity.ServiceError)
			},
			wantErr:  true,
			wantCode: codes.Internal,
		},
		{
			name: "not found",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					GetPriceHistoryPoster(ctx, alias).
					Return(nil, entity.NotFoundError)
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockPostersUseCase(ctrl)
			if test.setupMock != nil {
				test.setupMock(mockUC)
			}

			server := NewPosterServiceServer(mockUC)
			resp, err := server.GetPriceHistoryPoster(ctx, test.req)

			if test.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, test.wantCode, st.Code())
				require.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Len(t, resp.History, len(expectedHistory))
		})
	}
}

type mockCreateFlatPosterStream struct {
	grpc.ServerStream

	ctx      context.Context
	requests []*posterpb.CreateFlatPosterRequest
	resp     *posterpb.CreateFlatPosterResponse
	recvIdx  int
}

func (m *mockCreateFlatPosterStream) Context() context.Context {
	return m.ctx
}

func (m *mockCreateFlatPosterStream) Recv() (*posterpb.CreateFlatPosterRequest, error) {
	if m.recvIdx >= len(m.requests) {
		return nil, io.EOF
	}

	req := m.requests[m.recvIdx]
	m.recvIdx++

	return req, nil
}

func (m *mockCreateFlatPosterStream) SendAndClose(resp *posterpb.CreateFlatPosterResponse) error {
	m.resp = resp
	return nil
}

func TestPosterServiceServer_CreateFlatPoster(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	flatNumber := int32(43)
	district := "Arbat"
	companyID := int64(12)

	validRequests := []*posterpb.CreateFlatPosterRequest{
		{
			Payload: &posterpb.CreateFlatPosterRequest_PosterMeta{
				PosterMeta: &posterpb.FlatPosterMeta{
					UserId:         7,
					Price:          135000,
					Description:    "krutoy remont",
					CategoryAlias:  "flat",
					Area:           96.4,
					GeoLat:         55.7558,
					GeoLon:         37.6173,
					FlatCategoryId: 1,
					FlatNumber:     &flatNumber,
					FlatFloor:      3,
					Address:        "Arbatskaya 5k2",
					City:           "Moscow",
					District:       &district,
					FloorCount:     7,
					CompanyId:      &companyID,
					Features:       []string{"balcony", "parking"},
				},
			},
		},
		{
			Payload: &posterpb.CreateFlatPosterRequest_PhotoMeta{
				PhotoMeta: &posterpb.FlatPosterPhotoMeta{
					Order: 1,
					Source: &posterpb.FlatPosterPhotoMeta_File{
						File: &posterpb.FileInput{
							Filename:    "img1.jpg",
							Size:        4,
							ContentType: "image/jpeg",
						},
					},
				},
			},
		},
		{
			Payload: &posterpb.CreateFlatPosterRequest_PhotoChunk{
				PhotoChunk: &posterpb.PhotoChunk{
					Order: 1,
					Data:  []byte("img1"),
				},
			},
		},
	}

	expectedResult := &dto.CreatedPoster{
		ID:    33,
		Alias: "kvartira-na-arbate",
	}

	tests := []struct {
		name      string
		requests  []*posterpb.CreateFlatPosterRequest
		setupMock func(mockUC *mocks.MockPostersUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name:     "success",
			requests: validRequests,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					CreateFlatPoster(ctx, gomock.Any()).
					Return(expectedResult, nil)
			},
			wantErr: false,
		},
		{
			name:     "missing poster meta",
			requests: []*posterpb.CreateFlatPosterRequest{},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name:     "use case error",
			requests: validRequests,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					CreateFlatPoster(ctx, gomock.Any()).
					Return(nil, entity.ServiceError)
			},
			wantErr:  true,
			wantCode: codes.Internal,
		},
		{
			name:     "not found",
			requests: validRequests,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					CreateFlatPoster(ctx, gomock.Any()).
					Return(nil, entity.NotFoundError)
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockPostersUseCase(ctrl)
			if test.setupMock != nil {
				test.setupMock(mockUC)
			}

			server := NewPosterServiceServer(mockUC)

			stream := &mockCreateFlatPosterStream{
				ctx:      ctx,
				requests: test.requests,
			}

			err := server.CreateFlatPoster(stream)

			if test.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, test.wantCode, st.Code())
				require.Nil(t, stream.resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, stream.resp)
			require.Equal(t, int64(expectedResult.ID), stream.resp.Id)
			require.Equal(t, expectedResult.Alias, stream.resp.Alias)
		})
	}
}

type mockUpdateFlatPosterStream struct {
	grpc.ServerStream

	ctx      context.Context
	requests []*posterpb.UpdateFlatPosterRequest
	resp     *posterpb.UpdateFlatPosterResponse
	recvIdx  int
}

func (m *mockUpdateFlatPosterStream) Context() context.Context {
	return m.ctx
}

func (m *mockUpdateFlatPosterStream) Recv() (*posterpb.UpdateFlatPosterRequest, error) {
	if m.recvIdx >= len(m.requests) {
		return nil, io.EOF
	}

	req := m.requests[m.recvIdx]
	m.recvIdx++

	return req, nil
}

func (m *mockUpdateFlatPosterStream) SendAndClose(resp *posterpb.UpdateFlatPosterResponse) error {
	m.resp = resp
	return nil
}

func TestPosterServiceServer_UpdateFlatPoster(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	alias := "kvartira-na-arbate"
	flatNumber := int32(43)
	district := "Arbat"
	companyID := int64(12)

	validRequests := []*posterpb.UpdateFlatPosterRequest{
		{
			Payload: &posterpb.UpdateFlatPosterRequest_PosterMeta{
				PosterMeta: &posterpb.UpdateFlatPosterMeta{
					Alias: alias,
					Poster: &posterpb.FlatPosterMeta{
						UserId:         7,
						Price:          135000,
						Description:    "krutoy remont",
						CategoryAlias:  "flat",
						Area:           96.4,
						GeoLat:         55.7558,
						GeoLon:         37.6173,
						FlatCategoryId: 1,
						FlatNumber:     &flatNumber,
						FlatFloor:      3,
						Address:        "Arbatskaya 5k2",
						City:           "Moscow",
						District:       &district,
						FloorCount:     7,
						CompanyId:      &companyID,
						Features:       []string{"balcony", "parking"},
					},
				},
			},
		},
		{
			Payload: &posterpb.UpdateFlatPosterRequest_PhotoMeta{
				PhotoMeta: &posterpb.FlatPosterPhotoMeta{
					Order: 1,
					Source: &posterpb.FlatPosterPhotoMeta_File{
						File: &posterpb.FileInput{
							Filename:    "img1.jpg",
							Size:        4,
							ContentType: "image/jpeg",
						},
					},
				},
			},
		},
		{
			Payload: &posterpb.UpdateFlatPosterRequest_PhotoChunk{
				PhotoChunk: &posterpb.PhotoChunk{
					Order: 1,
					Data:  []byte("img1"),
				},
			},
		},
	}

	expectedResult := &dto.CreatedPoster{
		ID:    33,
		Alias: alias,
	}

	tests := []struct {
		name      string
		requests  []*posterpb.UpdateFlatPosterRequest
		setupMock func(mockUC *mocks.MockPostersUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name:     "success",
			requests: validRequests,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					UpdateFlatPoster(ctx, alias, gomock.Any()).
					Return(expectedResult, nil)
			},
			wantErr: false,
		},
		{
			name:     "missing poster meta",
			requests: []*posterpb.UpdateFlatPosterRequest{},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "missing alias",
			requests: []*posterpb.UpdateFlatPosterRequest{
				{
					Payload: &posterpb.UpdateFlatPosterRequest_PosterMeta{
						PosterMeta: &posterpb.UpdateFlatPosterMeta{
							Alias:  "",
							Poster: validRequests[0].GetPosterMeta().Poster,
						},
					},
				},
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "missing poster data",
			requests: []*posterpb.UpdateFlatPosterRequest{
				{
					Payload: &posterpb.UpdateFlatPosterRequest_PosterMeta{
						PosterMeta: &posterpb.UpdateFlatPosterMeta{
							Alias:  alias,
							Poster: nil,
						},
					},
				},
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name:     "use case error",
			requests: validRequests,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					UpdateFlatPoster(ctx, alias, gomock.Any()).
					Return(nil, entity.ServiceError)
			},
			wantErr:  true,
			wantCode: codes.Internal,
		},
		{
			name:     "not found",
			requests: validRequests,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					UpdateFlatPoster(ctx, alias, gomock.Any()).
					Return(nil, entity.NotFoundError)
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockPostersUseCase(ctrl)
			if test.setupMock != nil {
				test.setupMock(mockUC)
			}

			server := NewPosterServiceServer(mockUC)

			stream := &mockUpdateFlatPosterStream{
				ctx:      ctx,
				requests: test.requests,
			}

			err := server.UpdateFlatPoster(stream)

			if test.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, test.wantCode, st.Code())
				require.Nil(t, stream.resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, stream.resp)
			require.Equal(t, int64(expectedResult.ID), stream.resp.Id)
			require.Equal(t, expectedResult.Alias, stream.resp.Alias)
		})
	}
}

func TestPosterServiceServer_DeleteFlatPoster(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	alias := "kvartira-na-arbate"
	userID := int64(7)

	validReq := &posterpb.DeleteFlatPosterRequest{
		Alias:  alias,
		UserId: userID,
	}

	expectedResult := &dto.CreatedPoster{
		ID:    33,
		Alias: alias,
	}

	tests := []struct {
		name      string
		req       *posterpb.DeleteFlatPosterRequest
		setupMock func(mockUC *mocks.MockPostersUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					DeleteFlatPoster(ctx, alias, int(userID)).
					Return(expectedResult, nil)
			},
			wantErr: false,
		},
		{
			name: "missing alias",
			req: &posterpb.DeleteFlatPosterRequest{
				Alias:  "",
				UserId: userID,
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "invalid user id",
			req: &posterpb.DeleteFlatPosterRequest{
				Alias:  alias,
				UserId: 0,
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "use case error",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					DeleteFlatPoster(ctx, alias, int(userID)).
					Return(nil, entity.ServiceError)
			},
			wantErr:  true,
			wantCode: codes.Internal,
		},
		{
			name: "not found",
			req:  validReq,
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					DeleteFlatPoster(ctx, alias, int(userID)).
					Return(nil, entity.NotFoundError)
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockPostersUseCase(ctrl)
			if test.setupMock != nil {
				test.setupMock(mockUC)
			}

			server := NewPosterServiceServer(mockUC)
			resp, err := server.DeleteFlatPoster(ctx, test.req)

			if test.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, test.wantCode, st.Code())
				require.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, int64(expectedResult.ID), resp.Id)
			require.Equal(t, expectedResult.Alias, resp.Alias)
		})
	}
}

func TestPosterServiceServer_GetPosterRoommates(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	alias := "flat-1"

	expectedUsers := []dto.RoommateUserDTO{
		{
			ID:        1,
			FirstName: "Ivan",
			LastName:  "Ivanov",
			AvatarURL: lo.ToPtr("avatar.jpg"),
		},
	}

	tests := []struct {
		name      string
		req       *posterpb.GetPosterRoommatesRequest
		setupMock func(mockUC *mocks.MockPostersUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			req: &posterpb.GetPosterRoommatesRequest{
				Alias: alias,
			},
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					GetPosterRoommates(ctx, alias).
					Return(expectedUsers, nil)
			},
			wantErr: false,
		},
		{
			name: "missing alias",
			req: &posterpb.GetPosterRoommatesRequest{
				Alias: "",
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "use case error",
			req: &posterpb.GetPosterRoommatesRequest{
				Alias: alias,
			},
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					GetPosterRoommates(ctx, alias).
					Return(nil, entity.ServiceError)
			},
			wantErr:  true,
			wantCode: codes.Internal,
		},
		{
			name: "not found",
			req: &posterpb.GetPosterRoommatesRequest{
				Alias: alias,
			},
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					GetPosterRoommates(ctx, alias).
					Return(nil, entity.NotFoundError)
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockPostersUseCase(ctrl)
			if test.setupMock != nil {
				test.setupMock(mockUC)
			}

			server := NewPosterServiceServer(mockUC)

			resp, err := server.GetPosterRoommates(ctx, test.req)

			if test.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, test.wantCode, st.Code())
				require.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, int32(len(expectedUsers)), resp.Total)

			require.Len(t, resp.Users, 1)
			require.Equal(t, int64(expectedUsers[0].ID), resp.Users[0].Id)
			require.Equal(t, expectedUsers[0].FirstName, resp.Users[0].FirstName)
			require.Equal(t, expectedUsers[0].LastName, resp.Users[0].LastName)
			require.Equal(t, expectedUsers[0].AvatarURL, resp.Users[0].AvatarUrl)
		})
	}
}

func TestPosterServiceServer_AddPosterRoommate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	alias := "flat-1"
	userID := int64(10)

	tests := []struct {
		name      string
		req       *posterpb.AddPosterRoommateRequest
		setupMock func(mockUC *mocks.MockPostersUseCase)
		wantErr   bool
		wantCode  codes.Code
	}{
		{
			name: "success",
			req: &posterpb.AddPosterRoommateRequest{
				Alias:  alias,
				UserId: userID,
			},
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					AddPosterRoommate(ctx, alias, int(userID)).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name: "missing alias",
			req: &posterpb.AddPosterRoommateRequest{
				Alias:  "",
				UserId: userID,
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "invalid user id",
			req: &posterpb.AddPosterRoommateRequest{
				Alias:  alias,
				UserId: 0,
			},
			wantErr:  true,
			wantCode: codes.InvalidArgument,
		},
		{
			name: "use case error",
			req: &posterpb.AddPosterRoommateRequest{
				Alias:  alias,
				UserId: userID,
			},
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					AddPosterRoommate(ctx, alias, int(userID)).
					Return(entity.ServiceError)
			},
			wantErr:  true,
			wantCode: codes.Internal,
		},
		{
			name: "not found",
			req: &posterpb.AddPosterRoommateRequest{
				Alias:  alias,
				UserId: userID,
			},
			setupMock: func(mockUC *mocks.MockPostersUseCase) {
				mockUC.EXPECT().
					AddPosterRoommate(ctx, alias, int(userID)).
					Return(entity.NotFoundError)
			},
			wantErr:  true,
			wantCode: codes.NotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := mocks.NewMockPostersUseCase(ctrl)
			if test.setupMock != nil {
				test.setupMock(mockUC)
			}

			server := NewPosterServiceServer(mockUC)

			resp, err := server.AddPosterRoommate(ctx, test.req)

			if test.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, test.wantCode, st.Code())
				require.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
		})
	}
}
