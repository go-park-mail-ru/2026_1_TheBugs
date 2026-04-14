package tokens

import (
	"context"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/redis/mocks"
	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestTokenRepo_SetBlacklist_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := mocks.NewMockCmdable(ctrl)
	repo := NewTokenRepo(mockRedis)

	ttl := 10 * time.Minute

	cmd := redis.NewStatusCmd(context.Background())
	cmd.SetVal("OK")

	mockRedis.EXPECT().
		Set(gomock.Any(), "blacklist:token", "1", ttl).
		Return(cmd)

	err := repo.SetBlacklist(context.Background(), "token", ttl)
	require.NoError(t, err)
}

func TestTokenRepo_IsBlacklisted_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := mocks.NewMockCmdable(ctrl)
	repo := NewTokenRepo(mockRedis)

	cmd := redis.NewStringCmd(context.Background())
	cmd.SetVal("1")

	mockRedis.EXPECT().
		Get(gomock.Any(), "blacklist:token").
		Return(cmd)

	res, err := repo.IsBlacklisted(context.Background(), "token")
	require.NoError(t, err)
	require.True(t, res)
}

func TestTokenRepo_IsBlacklisted_InBlacklist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := mocks.NewMockCmdable(ctrl)
	repo := NewTokenRepo(mockRedis)

	cmd := redis.NewStringCmd(context.Background())
	cmd.SetVal("0")

	mockRedis.EXPECT().
		Get(gomock.Any(), "blacklist:token").
		Return(cmd)

	res, err := repo.IsBlacklisted(context.Background(), "token")
	require.NoError(t, err)
	require.False(t, res)
}

func TestTokenRepo_CreateRecoverSession_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := mocks.NewMockCmdable(ctrl)
	repo := NewTokenRepo(mockRedis)

	ctx := context.Background()
	sessionID := "session1"
	ttl := 10 * time.Minute

	data := entity.RecoverSession{
		Email:    "test@mail.com",
		Code:     "1234",
		Attempts: 0,
		Verified: true,
	}

	key := "auth:recover:session:session1"

	hsetCmd := redis.NewIntCmd(ctx)
	hsetCmd.SetVal(1)

	expCmd := redis.NewBoolCmd(ctx)
	expCmd.SetVal(true)

	mockRedis.EXPECT().
		HSet(gomock.Any(), key, gomock.Any()).
		Return(hsetCmd)

	mockRedis.EXPECT().
		Expire(gomock.Any(), key, ttl).
		Return(expCmd)

	err := repo.CreateRecoverSession(ctx, sessionID, data, ttl)
	require.NoError(t, err)
}

func TestTokenRepo_GetRecoverSession_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := mocks.NewMockCmdable(ctrl)
	repo := NewTokenRepo(mockRedis)

	ctx := context.Background()
	sessionID := "session1"

	key := "auth:recover:session:session1"

	cmd := redis.NewMapStringStringCmd(ctx)
	cmd.SetVal(map[string]string{
		"email":    "test@mail.com",
		"code":     "1234",
		"attempts": "2",
		"verified": "1",
	})

	mockRedis.EXPECT().
		HGetAll(gomock.Any(), key).
		Return(cmd)

	res, err := repo.GetRecoverSession(ctx, sessionID)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, "test@mail.com", res.Email)
	require.Equal(t, 2, res.Attempts)
	require.True(t, res.Verified)
}

func TestTokenRepo_GetRecoverSession_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := mocks.NewMockCmdable(ctrl)
	repo := NewTokenRepo(mockRedis)

	ctx := context.Background()
	sessionID := "session1"

	key := "auth:recover:session:session1"

	cmd := redis.NewMapStringStringCmd(ctx)
	cmd.SetVal(map[string]string{})

	mockRedis.EXPECT().
		HGetAll(gomock.Any(), key).
		Return(cmd)

	res, err := repo.GetRecoverSession(ctx, sessionID)
	require.NoError(t, err)
	require.Nil(t, res)
}

func TestTokenRepo_DeleteRecoverSession_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := mocks.NewMockCmdable(ctrl)
	repo := NewTokenRepo(mockRedis)

	ctx := context.Background()
	sessionID := "session1"

	key := "auth:recover:session:session1"

	cmd := redis.NewIntCmd(ctx)
	cmd.SetVal(1)

	mockRedis.EXPECT().
		Del(gomock.Any(), key).
		Return(cmd)

	err := repo.DeleteRecoverSession(ctx, sessionID)
	require.NoError(t, err)
}

func TestTokenRepo_SetRecoverVerified_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := mocks.NewMockCmdable(ctrl)
	repo := NewTokenRepo(mockRedis)

	ctx := context.Background()
	sessionID := "session1"

	key := "auth:recover:session:session1"

	cmd := redis.NewIntCmd(ctx)
	cmd.SetVal(1)

	mockRedis.EXPECT().
		HSet(gomock.Any(), key, "verified", true).
		Return(cmd)

	err := repo.SetRecoverVerified(ctx, sessionID, true)
	require.NoError(t, err)
}
