package complex

import (
	"net/http"
	"net/http/httptest"
	"testing"

	complexGRPC "github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/generated/complex"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/mocks/grpc_client"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestUtilityCompanyHandler_GetUtilityCompany(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockComplexServiceClient(ctrl)
	handler := UtilityCompanyHandler{grpcClient: mockClient}

	tests := []struct {
		name           string
		alias          string
		setupMock      func()
		expectedStatus int
	}{
		{
			name:  "success",
			alias: "test-alias",
			setupMock: func() {
				mockClient.EXPECT().
					GetComplex(gomock.Any(), &complexGRPC.GetComplexRequest{Alias: "test-alias"}).
					Return(&complexGRPC.GetComplexResponse{
						Id:          1,
						Alias:       "test-alias",
						CompanyName: "Test Company",
						Description: "Best utility complex",
						Address:     "Moscow, Lenina 1",
						AvatarUrl:   lo.ToPtr("https://example.com/avatar.png"),
						Geo: &complexGRPC.GeoPoint{
							Lat: 55.7558,
							Lon: 37.6173,
						},
						Photos: []*complexGRPC.Photo{
							{Url: "https://example.com/photo1.jpg"},
						},
						Developer: &complexGRPC.Developer{
							Id:            10,
							DeveloperName: "Dev Group",
							AvatarUrl:     lo.ToPtr("https://example.com/dev.png"),
						},
						Phone: "+7-999-111-22-33",
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing alias",
			alias:          "",
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:  "grpc not found",
			alias: "missing",
			setupMock: func() {
				mockClient.EXPECT().
					GetComplex(gomock.Any(), &complexGRPC.GetComplexRequest{Alias: "missing"}).
					Return(nil, status.Error(codes.NotFound, "not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			req := httptest.NewRequest(http.MethodGet, "/utility-companies/by-alias/"+tt.alias, nil)
			if tt.alias != "" {
				req = mux.SetURLVars(req, map[string]string{"alias": tt.alias})
			}
			rec := httptest.NewRecorder()

			handler.GetUtilityCompany(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestUtilityCompanyHandler_GetAllDevelopers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockComplexServiceClient(ctrl)
	handler := UtilityCompanyHandler{grpcClient: mockClient}

	tests := []struct {
		name           string
		setupMock      func()
		expectedStatus int
	}{
		{
			name: "success",
			setupMock: func() {
				mockClient.EXPECT().
					GetAllDevelopers(gomock.Any(), &emptypb.Empty{}).
					Return(&complexGRPC.GetAllDevelopersResponse{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "grpc error",
			setupMock: func() {
				mockClient.EXPECT().
					GetAllDevelopers(gomock.Any(), &emptypb.Empty{}).
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

			req := httptest.NewRequest(http.MethodGet, "/utility-companies/developers", nil)
			rec := httptest.NewRecorder()

			handler.GetAllDevelopers(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}

func TestUtilityCompanyHandler_GetUtilityCompaniesByDeveloper(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := grpc_client.NewMockComplexServiceClient(ctrl)
	handler := UtilityCompanyHandler{grpcClient: mockClient}

	tests := []struct {
		name           string
		url            string
		setupMock      func()
		expectedStatus int
	}{
		{
			name: "success",
			url:  "/utility-companies/?developer_id=12",
			setupMock: func() {
				mockClient.EXPECT().
					GetComplexesByDeveloperID(gomock.Any(), &complexGRPC.GetComplexesByDeveloperIDRequest{DeveloperID: 12}).
					Return(&complexGRPC.GetComplexesByDeveloperIDResponse{
						Complexes: []*complexGRPC.Complex{},
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing developer id",
			url:            "/utility-companies/",
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid developer id",
			url:            "/utility-companies/?developer_id=abc",
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "grpc error",
			url:  "/utility-companies/?developer_id=15",
			setupMock: func() {
				mockClient.EXPECT().
					GetComplexesByDeveloperID(gomock.Any(), &complexGRPC.GetComplexesByDeveloperIDRequest{DeveloperID: 15}).
					Return(nil, status.Error(codes.Internal, "db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			rec := httptest.NewRecorder()

			handler.GetUtilityCompaniesByDeveloper(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}
