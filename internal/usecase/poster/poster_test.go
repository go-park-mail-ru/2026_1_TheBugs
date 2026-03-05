package poster

import (
	"testing"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetPostersUseCase(t *testing.T) {
	existingListPoster := []*entity.Poster{
		{Id: 1, Price: 11111, Address: "street_1"},
		{Id: 2, Price: 22222, Address: "street_2"},
		{Id: 3, Price: 33333, Address: "street_3"},
		{Id: 4, Price: 44444, Address: "street_4"},
		{Id: 5, Price: 55555, Address: "street_5"},
	}

	expectedDTO := []dto.PosterDTO{
		{Price: 11111, Address: "street_1"},
		{Price: 22222, Address: "street_2"},
		{Price: 33333, Address: "street_3"},
		{Price: 44444, Address: "street_4"},
		{Price: 55555, Address: "street_5"},
	}

	tests := []struct {
		name      string
		params    dto.PostersFiltersDTO
		setupMock func(m *mocks.MockPosterRepo)
		want      []dto.PosterDTO
		wantErr   error
	}{
		{
			name: "OK",
			params: dto.PostersFiltersDTO{
				Limit:  12,
				Offset: 0,
			},
			setupMock: func(m *mocks.MockPosterRepo) {
				m.EXPECT().CountPosters().Return(5).Times(1)
				m.EXPECT().GetPosters(12, 0).Return(existingListPoster, nil).Times(1)
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
				Limit:  1000,
				Offset: 0,
			},
			setupMock: func(m *mocks.MockPosterRepo) {
				m.EXPECT().CountPosters().Return(5).Times(1)
				m.EXPECT().GetPosters(MaxPostersLimit, 0).Return(existingListPoster, nil).Times(1)
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
			name: "OffsetOutOfRange",
			params: dto.PostersFiltersDTO{
				Limit:  12,
				Offset: 10,
			},
			setupMock: func(m *mocks.MockPosterRepo) {
				m.EXPECT().CountPosters().Return(5).Times(1)
			},
			want:    nil,
			wantErr: entity.OffsetOutOfRange,
		},
		{
			name: "ErrorFromRepo",
			params: dto.PostersFiltersDTO{
				Limit:  12,
				Offset: 0,
			},
			setupMock: func(m *mocks.MockPosterRepo) {
				m.EXPECT().CountPosters().Return(5).Times(1)
				m.EXPECT().GetPosters(12, 0).Return(nil, entity.NotFoundError).Times(1)
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
			if test.setupMock != nil {
				test.setupMock(mockRepo)
			}

			uc := NewPosterUseCase(mockRepo)

			got, err := uc.GetPostersUseCase(test.params)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.Equal(t, test.want, got)
		})
	}
}
