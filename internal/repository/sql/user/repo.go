package user

import (
	"context"
	"fmt"

	repository "github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/sql"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
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

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	sql := `SELECT id, email, salt, hashed_password, provider FROM users WHERE email=$1`
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

func (r *UserRepo) GetByID(ctx context.Context, id int) (*dto.UserDTO, error) {
	sql := `SELECT u.id, u.email, p.first_name, p.last_name, p.avatar_url , p.phone
			FROM users u
			JOIN profiles p ON u.profile_id = p.id
			WHERE u.id=$1`
	row, err := r.pool.Query(ctx, sql, id)

	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	defer row.Close()

	user, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[entity.UserDetails])
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	return dto.UserToDTO(&user), nil

}

func (r *UserRepo) UpdateProfile(ctx context.Context, data dto.UpdateProfileDTO) (*dto.UserDTO, error) {
	sql := `UPDATE profiles p
			SET phone = COALESCE($1, p.phone),
				first_name = COALESCE($2, p.first_name),
				last_name = COALESCE($3, p.last_name),
				avatar_url = COALESCE($4, p.avatar_url)
			FROM users u
			WHERE p.id = u.profile_id AND u.id = $5  -- Join users here
			RETURNING p.id, u.email, p.first_name, p.last_name, p.avatar_url, p.phone;`
	row, err := r.pool.Query(ctx, sql, data.Phone, data.FirstName, data.LastName, data.AvatarURL, data.ID)

	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	defer row.Close()

	user, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[entity.UserDetails])
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	return dto.UserToDTO(&user), nil

}

func (r *UserRepo) Create(ctx context.Context, dto dto.CreateUserDTO) (*entity.User, error) {
	var profileID int

	profileSql := `INSERT INTO profiles (phone, first_name, last_name) VALUES ($1, $2, $3) RETURNING id`

	fmt.Println(dto.Phone)

	err := r.pool.QueryRow(ctx, profileSql, dto.Phone, dto.FirstName, dto.LastName).Scan(&profileID)
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}
	sql := `INSERT INTO users (email, hashed_password, salt, profile_id) VALUES ($1, $2, $3, $4) 
			RETURNING id, email, hashed_password, salt, provider`

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

func (r *UserRepo) CreateByProvider(ctx context.Context, dto dto.CreateUserByProviderDTO) (*entity.User, error) {
	var profileID int

	profileSql := `INSERT INTO profiles (phone, first_name, last_name) VALUES ($1, $2, $3) RETURNING id`

	fmt.Println(dto.Phone)

	err := r.pool.QueryRow(ctx, profileSql, dto.Phone, dto.FirstName, dto.LastName).Scan(&profileID)
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}
	sql := `INSERT INTO users (email, provider, provider_id, profile_id) VALUES ($1, $2, $3, $4) 
			RETURNING id, email, hashed_password, salt, provider`
	row, err := r.pool.Query(ctx, sql, dto.Email, dto.Provider, dto.ProviderID, profileID)
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

func (r *UserRepo) GetByProvider(ctx context.Context, provider string, email string) (*entity.User, error) {
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

func (r *UserRepo) UpdatePwd(ctx context.Context, email string, pwd string, salt string) error {
	sql := `UPDATE users SET hashed_password=$1, salt=$2 WHERE email=$3`
	ct, err := r.pool.Exec(ctx, sql, pwd, salt, email)

	if err != nil {
		return repository.HandelPgErrors(err)
	}

	if ct.RowsAffected() == 0 {
		return repository.HandelPgErrors(pgx.ErrNoRows)
	}

	return nil
}
