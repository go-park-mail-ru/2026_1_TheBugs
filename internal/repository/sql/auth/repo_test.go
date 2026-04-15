package auth

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"
)

func TestCreateToken(t *testing.T) {
	query := regexp.QuoteMeta(`INSERT INTO refresh_tokens (token_id, user_id, expires_at) VALUES ($1, $2, $3) RETURNING id`)

	expiresAt := time.Date(2026, 3, 9, 12, 0, 0, 0, time.UTC)
	inputDTO := dto.CreateRefreshTokenDTO{
		TokenID:   "token-123",
		UserID:    42,
		ExpiresAt: expiresAt,
	}

	tests := []struct {
		name      string
		dto       dto.CreateRefreshTokenDTO
		setupMock func(m pgxmock.PgxPoolIface)
		wantErr   error
	}{
		{
			name: "OK",
			dto:  inputDTO,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id"}).AddRow(1)

				m.ExpectQuery(query).WithArgs("token-123", 42, expiresAt).WillReturnRows(rows)
			},
			wantErr: nil,
		},
		{
			name: "not_found",
			dto:  inputDTO,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id"})

				m.ExpectQuery(query).WithArgs("token-123", 42, expiresAt).WillReturnRows(rows)
			},
			wantErr: entity.NotFoundError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewAuthRepo(mock)

			err = repo.CreateToken(context.Background(), test.dto)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetToken(t *testing.T) {
	query := regexp.QuoteMeta(`SELECT id, token_id, user_id, expires_at FROM refresh_tokens WHERE token_id=$1 AND user_id=$2`)

	expiresAt := time.Date(2026, 3, 9, 12, 0, 0, 0, time.UTC)
	expectedToken := &entity.RefreshToken{
		ID:        1,
		TokenID:   "token-123",
		UserID:    42,
		ExpiresAt: expiresAt,
	}

	tests := []struct {
		name      string
		tokenID   string
		userID    int
		setupMock func(m pgxmock.PgxPoolIface)
		want      *entity.RefreshToken
		wantErr   error
	}{
		{
			name:    "OK",
			tokenID: "token-123",
			userID:  42,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "token_id", "user_id", "expires_at",
				}).AddRow(1, "token-123", 42, expiresAt)

				m.ExpectQuery(query).WithArgs("token-123", 42).WillReturnRows(rows)
			},
			want:    expectedToken,
			wantErr: nil,
		},
		{
			name:    "not_found",
			tokenID: "missing-token",
			userID:  42,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "token_id", "user_id", "expires_at",
				})

				m.ExpectQuery(query).WithArgs("missing-token", 42).WillReturnRows(rows)
			},
			want:    nil,
			wantErr: entity.NotFoundError,
		},
		{
			name:    "collect_row_error",
			tokenID: "token-123",
			userID:  42,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "token_id", "user_id", "expires_at",
				}).AddRow("bad_id", "token-123", 42, expiresAt)

				m.ExpectQuery(query).WithArgs("token-123", 42).WillReturnRows(rows)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewAuthRepo(mock)

			got, err := repo.GetToken(context.Background(), test.tokenID, test.userID)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDeleteToken(t *testing.T) {
	query := regexp.QuoteMeta(`DELETE FROM refresh_tokens WHERE token_id=$1 AND user_id=$2 RETURNING id`)

	tests := []struct {
		name      string
		tokenID   string
		userID    int
		setupMock func(m pgxmock.PgxPoolIface)
		wantErr   error
	}{
		{
			name:    "OK",
			tokenID: "token-123",
			userID:  42,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id"}).AddRow(1)

				m.ExpectQuery(query).WithArgs("token-123", 42).WillReturnRows(rows)
			},
			wantErr: nil,
		},
		{
			name:    "not_found",
			tokenID: "missing-token",
			userID:  42,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id"})

				m.ExpectQuery(query).WithArgs("missing-token", 42).WillReturnRows(rows)
			},
			wantErr: entity.NotFoundError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewAuthRepo(mock)

			err = repo.DeleteToken(context.Background(), test.tokenID, test.userID)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
