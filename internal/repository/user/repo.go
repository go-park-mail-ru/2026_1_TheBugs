package user

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository"
	"github.com/jackc/pgx/v5"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
)

type UserRepo struct {
	pool repository.DB
}

func NewUserRepo(pool repository.DB) *UserRepo {
	return &UserRepo{
		pool: pool,
	}
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	sql := `SELECT id, email, salt, hashed_password, provider FROM users WHERE email=$1 AND provider IS NULL`
	row, err := r.pool.Query(ctx, sql, email)

	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	defer row.Close()

	user, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[entity.User])
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	return &user, nil

}

func (r *UserRepo) CreateUser(ctx context.Context, dto dto.CreateUserDTO) (*entity.User, error) {
	var profileID int

	profileSql := `INSERT INTO profiles (phone, first_name, last_name) VALUES ($1, $2, $3) RETURNING id`

	fmt.Println(dto.Phone)

	err := r.pool.QueryRow(ctx, profileSql, dto.Phone, dto.FirstName, dto.LastName).Scan(&profileID)
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}
	sql := `INSERT INTO users (email, hashed_password, salt, profile_id) VALUES ($1, $2, $3, $4) 
			RETURNING id, email, hashed_password, salt`

	row, err := r.pool.Query(ctx, sql, dto.Email, *dto.HashedPassword, *dto.Salt, profileID)
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}
	defer row.Close()

	user, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[entity.User])
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}
	return &user, nil
}

func (r *UserRepo) CreateUserByProvider(ctx context.Context, dto dto.CreateUserByProviderDTO) (*entity.User, error) {
	sql := `INSERT INTO users (email, provider) VALUES ($1, $2) 
			RETURNING id, email, hashed_password, salt, provider`
	row, err := r.pool.Query(ctx, sql, dto.Email, dto.Provider)
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}
	defer row.Close()

	user, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[entity.User])
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}
	return &user, nil
}

func (r *UserRepo) GetUserByProvider(ctx context.Context, provider string, email string) (*entity.User, error) {
	sql := `SELECT id, email, salt, hashed_password, provider FROM users WHERE email=$1 AND provider=$2`
	row, err := r.pool.Query(ctx, sql, email, provider)

	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	defer row.Close()

	user, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[entity.User])
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	return &user, nil
}
