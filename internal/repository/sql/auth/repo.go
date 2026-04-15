package auth

import (
	"context"
	"log"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	repository "github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/sql"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/jackc/pgx/v5"
)

type AuthRepo struct {
	pool repository.DB
}

func NewAuthRepo(pool repository.DB) *AuthRepo {
	return &AuthRepo{
		pool: pool,
	}
}

func (r *AuthRepo) CreateToken(ctx context.Context, token dto.CreateRefreshTokenDTO) error {
	sql := `INSERT INTO refresh_tokens (token_id, user_id, expires_at) VALUES ($1, $2, $3) RETURNING id`
	var tokenID int
	err := r.pool.QueryRow(ctx, sql, token.TokenID, token.UserID, token.ExpiresAt).Scan(&tokenID)
	if err != nil {
		log.Printf("%s", err.Error())
		return repository.HandelPgErrors(err)
	}
	return nil
}

func (r *AuthRepo) GetToken(ctx context.Context, tokenID string, userID int) (*entity.RefreshToken, error) {
	sql := `SELECT id, token_id, user_id, expires_at FROM refresh_tokens WHERE token_id=$1 AND user_id=$2`
	row, err := r.pool.Query(ctx, sql, tokenID, userID)

	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	defer row.Close()

	token, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[entity.RefreshToken])
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	return &token, nil
}
func (r *AuthRepo) DeleteToken(ctx context.Context, tokenID string, userID int) error {
	sql := `DELETE FROM refresh_tokens WHERE token_id=$1 AND user_id=$2 RETURNING id`
	var deletedID int
	err := r.pool.QueryRow(ctx, sql, tokenID, userID).Scan(&deletedID)
	if err != nil {
		return repository.HandelPgErrors(err)
	}
	return nil
}
