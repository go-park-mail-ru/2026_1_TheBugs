package poster

import (
	"testing"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetPostersUseCase_OK(t *testing.T) {
	testLimit := 12
	testOffset := 0
	existingListPoster := []entity.Poster{
		{
			Id:      1,
			Price:   11111,
			Address: "street_1",
		},
		{
			Id:      2,
			Price:   22222,
			Address: "street_2",
		},
		{
			Id:      3,
			Price:   33333,
			Address: "street_3",
		},
		{
			Id:      4,
			Price:   44444,
			Address: "street_4",
		},
		{
			Id:      5,
			Price:   55555,
			Address: "street_5",
		},
	}

	expected := []dto.PosterDTO{
		{
			Price:   11111,
			Address: "street_1",
		},
		{
			Price:   22222,
			Address: "street_2",
		},
		{
			Price:   33333,
			Address: "street_3",
		},
		{
			Price:   44444,
			Address: "street_4",
		},
		{
			Price:   55555,
			Address: "street_5",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mocks.NewMockPosterRepo(ctrl)

	mock.EXPECT().GetPosters(testLimit, testOffset).Return(existingListPoster, nil).Times(1)

	uc := NewPosterUseCase(mock)

	got, err := uc.GetPostersUseCase(testLimit, testOffset)
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestGetPostersUseCase_InvalidLimit(t *testing.T) {
	testLimit := -2
	testOffset := 0

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mocks.NewMockPosterRepo(ctrl)

	uc := NewPosterUseCase(mock)

	_, err := uc.GetPostersUseCase(testLimit, testOffset)
	require.ErrorIs(t, err, entity.InvalidInput)
}

func TestGetPostersUseCase_InvalidOffset(t *testing.T) {
	testLimit := 12
	testOffset := -3

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mocks.NewMockPosterRepo(ctrl)

	uc := NewPosterUseCase(mock)

	_, err := uc.GetPostersUseCase(testLimit, testOffset)
	require.ErrorIs(t, err, entity.InvalidInput)
}

func TestGetPostersUseCase_ErrorFromRepo(t *testing.T) {
	testLimit := 12
	testOffset := 0

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := mocks.NewMockPosterRepo(ctrl)

	mock.EXPECT().GetPosters(testLimit, testOffset).Return(nil, entity.NotFoundError).Times(1)

	uc := NewPosterUseCase(mock)

	_, err := uc.GetPostersUseCase(testLimit, testOffset)
	require.ErrorIs(t, err, entity.NotFoundError)
}
