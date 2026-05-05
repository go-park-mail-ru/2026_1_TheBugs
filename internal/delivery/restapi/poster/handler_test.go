package poster

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/generated/poster"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/restapi/utils"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/mocks/grpc_client"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

func strPtr(v string) *string {
	return &v
}

func int32Ptr(v int32) *int32 {
	return &v
}

func int64Ptr(v int64) *int64 {
	return &v
}

func TestPosterHandler_GetFlatsAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockPosterServiceClient(ctrl)

	handler := PosterHandler{grpcClient: mockClient}

	tests := []struct {
		name           string
		query          string
		setupMock      func()
		expectedStatus int
		checkBody      bool
	}{
		{
			name:  "success default params",
			query: "/posters/flats",
			setupMock: func() {
				mockClient.EXPECT().
					SearchPosters(gomock.Any(), &poster.SearchPostersRequest{
						Limit:  12,
						Offset: 0,
					}).
					Return(&poster.SearchPostersResponse{
						Len: 1,
						Posters: []*poster.PosterCard{
							{
								Id:       1,
								Alias:    "flat-1",
								Price:    100000,
								Address:  "Moscow",
								Area:     42,
								ImageUrl: strPtr("image.png"),
							},
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			checkBody:      true,
		},
		{
			name: "success with filters",
			query: "/posters/flats?limit=10&offset=5&search_query=test&utility_company=uk&category=flat&room_count=2" +
				"&min_price=1000&max_price=5000&facilities=wifi,parking&min_square=30&max_square=80" +
				"&min_flat_floor=2&max_flat_floor=9&min_building_floor=5&max_building_floor=20" +
				"&not_first_floor=true&not_last_floor=true",
			setupMock: func() {
				mockClient.EXPECT().
					SearchPosters(gomock.Any(), &poster.SearchPostersRequest{
						Limit:            10,
						Offset:           5,
						SearchQuery:      strPtr("test"),
						UtilityCompany:   strPtr("uk"),
						Category:         strPtr("flat"),
						RoomCount:        int32Ptr(2),
						MinPrice:         int32Ptr(1000),
						MaxPrice:         int32Ptr(5000),
						Facilities:       []string{"wifi", "parking"},
						MinSquare:        int32Ptr(30),
						MaxSquare:        int32Ptr(80),
						MinFlatFloor:     int32Ptr(2),
						MaxFlatFloor:     int32Ptr(9),
						MinBuildingFloor: int32Ptr(5),
						MaxBuildingFloor: int32Ptr(20),
						IsNotFirstFloor:  true,
						IsNotLastFloor:   true,
					}).
					Return(&poster.SearchPostersResponse{
						Len:     0,
						Posters: []*poster.PosterCard{},
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:  "grpc error invalid argument",
			query: "/posters/flats?limit=10",
			setupMock: func() {
				mockClient.EXPECT().
					SearchPosters(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.InvalidArgument, "bad filters"))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:  "grpc error internal",
			query: "/posters/flats?limit=10",
			setupMock: func() {
				mockClient.EXPECT().
					SearchPosters(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			req := httptest.NewRequest(http.MethodGet, tt.query, nil)
			rec := httptest.NewRecorder()

			handler.GetFlatsAll(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)

			if tt.checkBody {
				var resp map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)
				require.Equal(t, float64(1), resp["len"])
			}
		})
	}
}

func TestPosterHandler_GetPoster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockPosterServiceClient(ctrl)

	handler := PosterHandler{grpcClient: mockClient}

	tests := []struct {
		name           string
		alias          string
		setupMock      func()
		expectedStatus int
		checkBody      bool
	}{
		{
			name:  "success",
			alias: "flat-1",
			setupMock: func() {
				mockClient.EXPECT().
					GetPosterByAlias(gomock.Any(), &poster.GetPosterByAliasRequest{
						PosterAlias: "flat-1",
					}).
					Return(&poster.GetPosterByAliasResponse{
						Poster: &poster.Poster{
							Id:    1,
							Alias: "flat-1",
							Price: 100000,
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			checkBody:      true,
		},
		{
			name:           "invalid alias",
			alias:          "",
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:  "grpc not found",
			alias: "unknown",
			setupMock: func() {
				mockClient.EXPECT().
					GetPosterByAlias(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.NotFound, "not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:  "grpc internal error",
			alias: "flat-1",
			setupMock: func() {
				mockClient.EXPECT().
					GetPosterByAlias(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.setupMock != nil {
				test.setupMock()
			}

			req := httptest.NewRequest(http.MethodGet, "/posters/by-alias/"+test.alias, nil)
			req = mux.SetURLVars(req, map[string]string{
				"alias": test.alias,
			})
			rec := httptest.NewRecorder()

			handler.GetPoster(rec, req)

			require.Equal(t, test.expectedStatus, rec.Code)

			if test.checkBody {
				var resp map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)

				posterObj := resp["poster"].(map[string]interface{})
				require.Equal(t, "flat-1", posterObj["alias"])
			}
		})
	}
}

func TestPosterHandler_GetPostersByUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockPosterServiceClient(ctrl)

	handler := PosterHandler{grpcClient: mockClient}

	tests := []struct {
		name           string
		userID         int
		setupMock      func()
		expectedStatus int
		checkBody      bool
	}{
		{
			name:   "success",
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					GetPostersByUserID(gomock.Any(), &poster.GetPostersByUserIDRequest{
						UserId: 10,
					}).
					Return(&poster.GetPostersByUserIDResponse{
						Posters: []*poster.MyPoster{
							{
								Id:      1,
								Alias:   "flat-1",
								Address: "Moscow",
								Area:    42,
								Price:   100000,
								Category: &poster.Category{
									Alias: "flat",
									Name:  "Квартира",
								},
							},
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			checkBody:      true,
		},
		{
			name:           "missing user id",
			userID:         0,
			setupMock:      nil,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "grpc not found",
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					GetPostersByUserID(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.NotFound, "not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:   "grpc internal error",
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					GetPostersByUserID(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.setupMock != nil {
				test.setupMock()
			}

			req := httptest.NewRequest(http.MethodGet, "/posters/me", nil)
			if test.userID != 0 {
				req = req.WithContext(utils.SetUserID(req.Context(), test.userID))
			}
			rec := httptest.NewRecorder()

			handler.GetPostersByUser(rec, req)

			require.Equal(t, test.expectedStatus, rec.Code)

			if test.checkBody {
				var resp map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)

				require.Equal(t, float64(1), resp["len"])

				posters := resp["posters"].([]interface{})
				posterObj := posters[0].(map[string]interface{})
				require.Equal(t, "flat-1", posterObj["alias"])
			}
		})
	}
}

func TestPosterHandler_GetPostersByUserByAlias(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockPosterServiceClient(ctrl)

	handler := PosterHandler{grpcClient: mockClient}

	tests := []struct {
		name           string
		alias          string
		userID         int
		setupMock      func()
		expectedStatus int
		checkBody      bool
	}{
		{
			name:   "success",
			alias:  "flat-1",
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					GetPosterByAlias(gomock.Any(), &poster.GetPosterByAliasRequest{
						PosterAlias: "flat-1",
						UserId:      int64Ptr(10),
					}).
					Return(&poster.GetPosterByAliasResponse{
						Poster: &poster.Poster{
							Id:    1,
							Alias: "flat-1",
							Price: 100000,
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			checkBody:      true,
		},
		{
			name:           "invalid alias",
			alias:          "",
			userID:         10,
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing user id",
			alias:          "flat-1",
			userID:         0,
			setupMock:      nil,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:   "grpc not found",
			alias:  "unknown",
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					GetPosterByAlias(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.NotFound, "not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:   "grpc internal error",
			alias:  "flat-1",
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					GetPosterByAlias(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.setupMock != nil {
				test.setupMock()
			}

			req := httptest.NewRequest(http.MethodGet, "/posters/me/"+test.alias, nil)
			req = mux.SetURLVars(req, map[string]string{
				"alias": test.alias,
			})
			if test.userID != 0 {
				req = req.WithContext(utils.SetUserID(req.Context(), test.userID))
			}
			rec := httptest.NewRecorder()

			handler.GetPostersByUserByAlias(rec, req)

			require.Equal(t, test.expectedStatus, rec.Code)

			if test.checkBody {
				var resp map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)

				posterObj := resp["poster"].(map[string]interface{})
				require.Equal(t, "flat-1", posterObj["alias"])
			}
		})
	}
}

func TestPosterHandler_DeleteFlatPoster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockPosterServiceClient(ctrl)

	handler := PosterHandler{grpcClient: mockClient}

	tests := []struct {
		name           string
		alias          string
		userID         int
		setupMock      func()
		expectedStatus int
		checkBody      bool
	}{
		{
			name:   "success",
			alias:  "flat-1",
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					DeleteFlatPoster(gomock.Any(), &poster.DeleteFlatPosterRequest{
						Alias:  "flat-1",
						UserId: 10,
					}).
					Return(&poster.DeleteFlatPosterResponse{
						Id:    1,
						Alias: "flat-1",
					}, nil)
			},
			expectedStatus: http.StatusOK,
			checkBody:      true,
		},
		{
			name:           "missing user id",
			alias:          "flat-1",
			userID:         0,
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid alias",
			alias:          "",
			userID:         10,
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "grpc not found",
			alias:  "unknown",
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					DeleteFlatPoster(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.NotFound, "not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:   "grpc internal error",
			alias:  "flat-1",
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					DeleteFlatPoster(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.setupMock != nil {
				test.setupMock()
			}

			req := httptest.NewRequest(http.MethodDelete, "/posters/flat/"+test.alias, nil)
			req = mux.SetURLVars(req, map[string]string{
				"alias": test.alias,
			})
			if test.userID != 0 {
				req = req.WithContext(utils.SetUserID(req.Context(), test.userID))
			}
			rec := httptest.NewRecorder()

			handler.DeleteFlatPoster(rec, req)

			require.Equal(t, test.expectedStatus, rec.Code)

			if test.checkBody {
				var resp map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)

				posterObj := resp["poster"].(map[string]interface{})
				require.Equal(t, float64(1), posterObj["id"])
				require.Equal(t, "flat-1", posterObj["alias"])
			}
		})
	}
}

func TestPosterHandler_AddViewPoster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockPosterServiceClient(ctrl)

	handler := PosterHandler{grpcClient: mockClient}

	tests := []struct {
		name           string
		alias          string
		userID         int
		setupMock      func()
		expectedStatus int
	}{
		{
			name:   "success",
			alias:  "flat-1",
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					AddViewPoster(gomock.Any(), &poster.AddViewPosterRequest{
						Alias:  "flat-1",
						UserId: 10,
					}).
					Return(&poster.AddViewPosterResponse{}, nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "missing user id",
			alias:          "flat-1",
			userID:         0,
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid alias",
			alias:          "",
			userID:         10,
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "grpc not found",
			alias:  "unknown",
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					AddViewPoster(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.NotFound, "not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:   "grpc internal error",
			alias:  "flat-1",
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					AddViewPoster(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.setupMock != nil {
				test.setupMock()
			}

			req := httptest.NewRequest(http.MethodPost, "/posters/"+test.alias+"/views", nil)
			req = mux.SetURLVars(req, map[string]string{
				"alias": test.alias,
			})
			if test.userID != 0 {
				req = req.WithContext(utils.SetUserID(req.Context(), test.userID))
			}
			rec := httptest.NewRecorder()

			handler.AddViewPoster(rec, req)

			require.Equal(t, test.expectedStatus, rec.Code)
		})
	}
}

func TestPosterHandler_GetViewsPoster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockPosterServiceClient(ctrl)

	handler := PosterHandler{grpcClient: mockClient}

	tests := []struct {
		name           string
		alias          string
		setupMock      func()
		expectedStatus int
		checkBody      bool
	}{
		{
			name:  "success",
			alias: "flat-1",
			setupMock: func() {
				mockClient.EXPECT().
					GetViewsPoster(gomock.Any(), &poster.GetViewsPosterRequest{
						Alias: "flat-1",
					}).
					Return(&poster.GetViewsPosterResponse{
						Views: 7,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			checkBody:      true,
		},
		{
			name:           "invalid alias",
			alias:          "",
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:  "grpc not found",
			alias: "unknown",
			setupMock: func() {
				mockClient.EXPECT().
					GetViewsPoster(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.NotFound, "not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:  "grpc internal error",
			alias: "flat-1",
			setupMock: func() {
				mockClient.EXPECT().
					GetViewsPoster(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.setupMock != nil {
				test.setupMock()
			}

			req := httptest.NewRequest(http.MethodGet, "/posters/"+test.alias+"/views", nil)
			req = mux.SetURLVars(req, map[string]string{
				"alias": test.alias,
			})
			rec := httptest.NewRecorder()

			handler.GetViewsPoster(rec, req)

			require.Equal(t, test.expectedStatus, rec.Code)

			if test.checkBody {
				var resp map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)

				require.Equal(t, float64(7), resp["views"])
			}
		})
	}
}

func TestPosterHandler_AddFavoritePoster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockPosterServiceClient(ctrl)

	handler := PosterHandler{grpcClient: mockClient}

	tests := []struct {
		name           string
		alias          string
		userID         int
		setupMock      func()
		expectedStatus int
	}{
		{
			name:   "success",
			alias:  "flat-1",
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					AddFavoritePoster(gomock.Any(), &poster.AddFavoritePosterRequest{
						Alias:  "flat-1",
						UserId: 10,
					}).
					Return(&poster.AddFavoritePosterResponse{}, nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "missing user id",
			alias:          "flat-1",
			userID:         0,
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid alias",
			alias:          "",
			userID:         10,
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "grpc not found",
			alias:  "unknown",
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					AddFavoritePoster(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.NotFound, "not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:   "grpc internal error",
			alias:  "flat-1",
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					AddFavoritePoster(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.setupMock != nil {
				test.setupMock()
			}

			req := httptest.NewRequest(http.MethodPost, "/posters/"+test.alias+"/favorites", nil)
			req = mux.SetURLVars(req, map[string]string{
				"alias": test.alias,
			})
			if test.userID != 0 {
				req = req.WithContext(utils.SetUserID(req.Context(), test.userID))
			}
			rec := httptest.NewRecorder()

			handler.AddFavoritePoster(rec, req)

			require.Equal(t, test.expectedStatus, rec.Code)
		})
	}
}

func TestPosterHandler_GetFavoritesPoster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockPosterServiceClient(ctrl)

	handler := PosterHandler{grpcClient: mockClient}

	tests := []struct {
		name           string
		userID         int
		setupMock      func()
		expectedStatus int
		checkBody      bool
	}{
		{
			name:   "success",
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					GetFavoritePosters(gomock.Any(), &poster.GetFavoritePostersRequest{
						UserId: 10,
					}).
					Return(&poster.GetFavoritePostersResponse{
						Len: 1,
						Posters: []*poster.PosterCard{
							{
								Id:       1,
								Alias:    "flat-1",
								Price:    100000,
								Address:  "Moscow",
								Area:     42,
								ImageUrl: strPtr("image.png"),
							},
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			checkBody:      true,
		},
		{
			name:           "missing user id",
			userID:         0,
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "grpc internal error",
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					GetFavoritePosters(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.setupMock != nil {
				test.setupMock()
			}

			req := httptest.NewRequest(http.MethodGet, "/posters/favorites", nil)
			if test.userID != 0 {
				req = req.WithContext(utils.SetUserID(req.Context(), test.userID))
			}
			rec := httptest.NewRecorder()

			handler.GetFavoritesPoster(rec, req)

			require.Equal(t, test.expectedStatus, rec.Code)

			if test.checkBody {
				var resp map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)

				require.Equal(t, float64(1), resp["len"])

				posters := resp["posters"].([]interface{})
				posterObj := posters[0].(map[string]interface{})
				require.Equal(t, "flat-1", posterObj["alias"])
			}
		})
	}
}

func TestPosterHandler_DeleteFavoritePoster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockPosterServiceClient(ctrl)

	handler := PosterHandler{grpcClient: mockClient}

	tests := []struct {
		name           string
		alias          string
		userID         int
		setupMock      func()
		expectedStatus int
	}{
		{
			name:   "success",
			alias:  "flat-1",
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					DeleteFavoritePoster(gomock.Any(), &poster.DeleteFavoritePosterRequest{
						Alias:  "flat-1",
						UserId: 10,
					}).
					Return(&poster.DeleteFavoritePosterResponse{}, nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "missing user id",
			alias:          "flat-1",
			userID:         0,
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid alias",
			alias:          "",
			userID:         10,
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "grpc not found",
			alias:  "unknown",
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					DeleteFavoritePoster(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.NotFound, "not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:   "grpc internal error",
			alias:  "flat-1",
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					DeleteFavoritePoster(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.setupMock != nil {
				test.setupMock()
			}

			req := httptest.NewRequest(http.MethodDelete, "/posters/"+test.alias+"/favorites", nil)
			req = mux.SetURLVars(req, map[string]string{
				"alias": test.alias,
			})
			if test.userID != 0 {
				req = req.WithContext(utils.SetUserID(req.Context(), test.userID))
			}
			rec := httptest.NewRecorder()

			handler.DeleteFavoritePoster(rec, req)

			require.Equal(t, test.expectedStatus, rec.Code)
		})
	}
}

func TestPosterHandler_GetFavoritesCountPoster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockPosterServiceClient(ctrl)

	handler := PosterHandler{grpcClient: mockClient}

	tests := []struct {
		name           string
		alias          string
		userID         int
		setupMock      func()
		expectedStatus int
		checkBody      bool
	}{
		{
			name:   "success with user id",
			alias:  "flat-1",
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					GetFavoritesCountPoster(gomock.Any(), &poster.GetFavoritesCountPosterRequest{
						PosterAlias: "flat-1",
						UserId:      int64Ptr(10),
					}).
					Return(&poster.GetFavoritesCountPosterResponse{
						Count:      5,
						IsFavorite: true,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			checkBody:      true,
		},
		{
			name:   "success without user id",
			alias:  "flat-1",
			userID: 0,
			setupMock: func() {
				mockClient.EXPECT().
					GetFavoritesCountPoster(gomock.Any(), &poster.GetFavoritesCountPosterRequest{
						PosterAlias: "flat-1",
						UserId:      nil,
					}).
					Return(&poster.GetFavoritesCountPosterResponse{
						Count:      5,
						IsFavorite: false,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			checkBody:      true,
		},
		{
			name:           "invalid alias",
			alias:          "",
			userID:         10,
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "grpc not found",
			alias:  "unknown",
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					GetFavoritesCountPoster(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.NotFound, "not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:   "grpc internal error",
			alias:  "flat-1",
			userID: 10,
			setupMock: func() {
				mockClient.EXPECT().
					GetFavoritesCountPoster(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.setupMock != nil {
				test.setupMock()
			}

			req := httptest.NewRequest(http.MethodGet, "/posters/"+test.alias+"/favorites", nil)
			req = mux.SetURLVars(req, map[string]string{
				"alias": test.alias,
			})
			if test.userID != 0 {
				req = req.WithContext(utils.SetUserID(req.Context(), test.userID))
			}
			rec := httptest.NewRecorder()

			handler.GetFavoritesCountPoster(rec, req)

			require.Equal(t, test.expectedStatus, rec.Code)

			if test.checkBody {
				require.NotEmpty(t, rec.Body.String())
			}
		})
	}
}

func TestPosterHandler_GetFlatsMapAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockPosterServiceClient(ctrl)

	handler := PosterHandler{grpcClient: mockClient}

	tests := []struct {
		name           string
		query          string
		setupMock      func()
		expectedStatus int
		checkBody      bool
	}{
		{
			name:  "success",
			query: "/posters/geo?sw_lat=55.1&sw_lon=37.1&ne_lat=56.1&ne_lon=38.1&zoom=10",
			setupMock: func() {
				mockClient.EXPECT().
					GetPostersByCoords(gomock.Any(), gomock.Any()).
					Return(&poster.GetPostersByCoordsResponse{
						Len: 1,
						Features: []*poster.GeoJSONFeature{
							{
								Type: "Feature",
								Properties: map[string]*structpb.Value{
									"alias": structpb.NewStringValue("flat-1"),
									"price": structpb.NewNumberValue(100000),
								},
								Geometry: &poster.Geometry{
									Type:        "Point",
									Coordinates: []float64{37.6, 55.7},
								},
							},
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			checkBody:      true,
		},
		{
			name:           "invalid map filters",
			query:          "/posters/geo?sw_lat=bad&sw_lon=37.1&ne_lat=56.1&ne_lon=38.1&zoom=10",
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing map filters",
			query:          "/posters/geo",
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:  "success with filters",
			query: "/posters/geo?sw_lat=55.1&sw_lon=37.1&ne_lat=56.1&ne_lon=38.1&zoom=10&limit=-1",
			setupMock: func() {
				mockClient.EXPECT().
					GetPostersByCoords(gomock.Any(), gomock.Any()).
					Return(&poster.GetPostersByCoordsResponse{
						Len:      0,
						Features: []*poster.GeoJSONFeature{},
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:  "grpc internal error",
			query: "/posters/geo?sw_lat=55.1&sw_lon=37.1&ne_lat=56.1&ne_lon=38.1&zoom=10",
			setupMock: func() {
				mockClient.EXPECT().
					GetPostersByCoords(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.setupMock != nil {
				test.setupMock()
			}

			req := httptest.NewRequest(http.MethodGet, test.query, nil)
			rec := httptest.NewRecorder()

			handler.GetFlatsMapAll(rec, req)

			require.Equal(t, test.expectedStatus, rec.Code)

			if test.checkBody {
				var resp map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)

				require.Equal(t, float64(1), resp["len"])

				features := resp["features"].([]interface{})
				require.Len(t, features, 1)

				feature := features[0].(map[string]interface{})
				require.Equal(t, "Feature", feature["type"])
			}
		})
	}
}

func TestPosterHandler_GetPostersByPoint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockPosterServiceClient(ctrl)

	handler := PosterHandler{grpcClient: mockClient}

	tests := []struct {
		name           string
		query          string
		setupMock      func()
		expectedStatus int
		checkBody      bool
	}{
		{
			name:  "success",
			query: "/posters/by-point?lat=55.7&lon=37.6",
			setupMock: func() {
				mockClient.EXPECT().
					GetPostersByRadius(gomock.Any(), &poster.GetPostersByRadiusRequest{
						Point: &poster.Geography{
							Lat: 55.7,
							Lon: 37.6,
						},
					}).
					Return(&poster.GetPostersByRadiusResponse{
						Posters: []*poster.MyPoster{
							{
								Id:      1,
								Alias:   "flat-1",
								Address: "Moscow",
								Area:    42,
								Price:   100000,
								Category: &poster.Category{
									Alias: "flat",
									Name:  "Квартира",
								},
							},
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			checkBody:      true,
		},
		{
			name:           "invalid lat",
			query:          "/posters/by-point?lat=bad&lon=37.6",
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid lon",
			query:          "/posters/by-point?lat=55.7&lon=bad",
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing lat",
			query:          "/posters/by-point?lon=37.6",
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing lon",
			query:          "/posters/by-point?lat=55.7",
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:  "grpc internal error",
			query: "/posters/by-point?lat=55.7&lon=37.6",
			setupMock: func() {
				mockClient.EXPECT().
					GetPostersByRadius(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.setupMock != nil {
				test.setupMock()
			}

			req := httptest.NewRequest(http.MethodGet, test.query, nil)
			rec := httptest.NewRecorder()

			handler.GetPostersByPoint(rec, req)

			require.Equal(t, test.expectedStatus, rec.Code)

			if test.checkBody {
				var resp map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)

				require.Equal(t, float64(1), resp["len"])

				posters := resp["posters"].([]interface{})
				posterObj := posters[0].(map[string]interface{})
				require.Equal(t, "flat-1", posterObj["alias"])
			}
		})
	}
}

func TestPosterHandler_GenerateDescription(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockPosterServiceClient(ctrl)

	handler := PosterHandler{grpcClient: mockClient}

	tests := []struct {
		name           string
		body           string
		setupMock      func()
		expectedStatus int
		checkBody      bool
	}{
		{
			name: "success",
			body: `{
				"category":"flat",
				"area":42,
				"flat_category":"studio",
				"city":"Moscow",
				"features":["wifi","parking"]
			}`,
			setupMock: func() {
				mockClient.EXPECT().
					GenerateDescription(gomock.Any(), &poster.GenerateDescriptionRequest{
						Category:     "flat",
						Area:         42,
						FlatCategory: "studio",
						City:         "Moscow",
						Features:     []string{"wifi", "parking"},
					}).
					Return(&poster.GenerateDescriptionResponse{
						Description: "good flat description",
					}, nil)
			},
			expectedStatus: http.StatusOK,
			checkBody:      true,
		},
		{
			name:           "invalid json",
			body:           `{bad json}`,
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "grpc invalid argument",
			body: `{
				"category":"flat",
				"area":0,
				"flat_category":"studio",
				"city":"Moscow",
				"features":[]
			}`,
			setupMock: func() {
				mockClient.EXPECT().
					GenerateDescription(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.InvalidArgument, "bad request"))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "grpc internal error",
			body: `{
				"category":"flat",
				"area":42,
				"flat_category":"studio",
				"city":"Moscow",
				"features":["wifi"]
			}`,
			setupMock: func() {
				mockClient.EXPECT().
					GenerateDescription(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.setupMock != nil {
				test.setupMock()
			}

			req := httptest.NewRequest(http.MethodPost, "/posters/generate-description", bytes.NewBufferString(test.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.GenerateDescription(rec, req)

			require.Equal(t, test.expectedStatus, rec.Code)

			if test.checkBody {
				var resp map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)

				require.Equal(t, "good flat description", resp["description"])
			}
		})
	}
}

func TestPosterHandler_GetPriceHistoryPoster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockPosterServiceClient(ctrl)

	handler := PosterHandler{grpcClient: mockClient}

	tests := []struct {
		name           string
		alias          string
		setupMock      func()
		expectedStatus int
		checkBody      bool
	}{
		{
			name:  "success",
			alias: "flat-1",
			setupMock: func() {
				mockClient.EXPECT().
					GetPriceHistoryPoster(gomock.Any(), &poster.GetPriceHistoryPosterRequest{
						PosterAlias: "flat-1",
					}).
					Return(&poster.GetPriceHistoryPosterResponse{
						History: []*poster.PriceHistory{
							{
								Date:  "2026-05-01",
								Price: 100000,
							},
							{
								Date:  "2026-05-02",
								Price: 110000,
							},
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			checkBody:      true,
		},
		{
			name:           "invalid alias",
			alias:          "",
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:  "grpc not found",
			alias: "unknown",
			setupMock: func() {
				mockClient.EXPECT().
					GetPriceHistoryPoster(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.NotFound, "not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:  "grpc internal error",
			alias: "flat-1",
			setupMock: func() {
				mockClient.EXPECT().
					GetPriceHistoryPoster(gomock.Any(), gomock.Any()).
					Return(nil, status.Error(codes.Internal, "internal error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.setupMock != nil {
				test.setupMock()
			}

			req := httptest.NewRequest(http.MethodGet, "/posters/"+test.alias+"/price-history", nil)
			req = mux.SetURLVars(req, map[string]string{
				"alias": test.alias,
			})
			rec := httptest.NewRecorder()

			handler.GetPriceHistoryPoster(rec, req)

			require.Equal(t, test.expectedStatus, rec.Code)

			if test.checkBody {
				var resp map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)

				require.Equal(t, float64(2), resp["count"])

				history := resp["history"].([]interface{})
				require.Len(t, history, 2)

				firstItem := history[0].(map[string]interface{})
				require.Equal(t, "2026-05-01", firstItem["date"])
				require.Equal(t, float64(100000), firstItem["price"])
			}
		})
	}
}
