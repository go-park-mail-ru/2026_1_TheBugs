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
	sql := `SELECT id, email, salt, hashed_password, provider, is_verified, is_admin FROM users WHERE email=$1`
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

func (r *UserRepo) GetAdminByID(ctx context.Context, userID int) (*entity.User, error) {
	sql := `SELECT id, email, salt, hashed_password, provider, is_admin FROM users WHERE id=$1`
	row, err := r.pool.Query(ctx, sql, userID)
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
	sql := `SELECT u.id, u.email, p.first_name, p.last_name, p.avatar_url, p.phone
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
	row, err := r.pool.Query(ctx, sql, data.Phone, data.FirstName, data.LastName, data.AvatarPath, data.ID)

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
			RETURNING id, email, hashed_password, salt, provider, is_verified, is_admin`

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
	sql := `INSERT INTO users (email, provider, provider_id, profile_id, is_verified) VALUES ($1, $2, $3, $4, TRUE) 
			RETURNING id, email, hashed_password, salt, provider, is_verified, is_admin`
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
	sql := `SELECT id, email, salt, hashed_password, provider, is_verified, is_admin FROM users WHERE email=$1 AND provider=$2`
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

func (r *UserRepo) VerifyEmail(ctx context.Context, email string) error {
	sql := `UPDATE users SET is_verified=$1 WHERE email=$2`
	ct, err := r.pool.Exec(ctx, sql, true, email)

	if err != nil {
		return repository.HandelPgErrors(err)
	}

	if ct.RowsAffected() == 0 {
		return repository.HandelPgErrors(pgx.ErrNoRows)
	}
	return nil
}

func (r *UserRepo) GetRoommateUser(ctx context.Context, userID int) (*entity.RoommateUser, error) {
	sql := `
		SELECT p.first_name, p.last_name, p.avatar_url,
			   rf.gender, rf.birthday::TEXT AS birthday, rf.description
		FROM users u
		JOIN profiles p ON u.profile_id = p.id
		JOIN roommate_forms rf ON rf.user_id = u.id
		WHERE u.id = $1
	`

	rows, err := r.pool.Query(ctx, sql, userID)
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}
	defer rows.Close()

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[entity.RoommateUser])
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	return &user, nil
}

func (r *UserRepo) GetRoommateTags(ctx context.Context, userID int) ([]entity.RoommateTag, error) {
	sql := `
		SELECT rt.name, rt.alias
		FROM roommate_tags rt
		JOIN roommate_form_tags rft ON rft.roommate_tag_id = rt.id
		JOIN roommate_forms rf ON rf.id = rft.roommate_form_id
		WHERE rf.user_id = $1
		ORDER BY rt.name
	`

	rows, err := r.pool.Query(ctx, sql, userID)
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}
	defer rows.Close()

	tags, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.RoommateTag])
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	return tags, nil
}

func (r *UserRepo) AddRoommateMatch(ctx context.Context, fromUserID int, toUserID int, posterAlias *string) error {
	if posterAlias != nil {
		sql := `
			INSERT INTO roommate_matches (from_user_id, to_user_id, poster_id)
			SELECT $1, $2, p.id
			FROM posters p
			WHERE p.alias = $3
			  AND p.deleted_at IS NULL
		`

		ct, err := r.pool.Exec(ctx, sql, fromUserID, toUserID, *posterAlias)
		if err != nil {
			return repository.HandelPgErrors(err)
		}

		if ct.RowsAffected() == 0 {
			return repository.HandelPgErrors(pgx.ErrNoRows)
		}

		return nil
	}

	sql := `
		INSERT INTO roommate_matches (from_user_id, to_user_id, poster_id)
		VALUES (
			$1,
			$2,
			(
				SELECT reverse_rm.poster_id
				FROM roommate_matches reverse_rm
				WHERE reverse_rm.from_user_id = $2
				  AND reverse_rm.to_user_id = $1
			)
		)
	`

	_, err := r.pool.Exec(ctx, sql, fromUserID, toUserID)
	if err != nil {
		return repository.HandelPgErrors(err)
	}

	return nil
}

func (r *UserRepo) IsRoommateMatch(ctx context.Context, fromUserID int, toUserID int) (bool, error) {
	sql := `
		SELECT EXISTS(
			SELECT 1
			FROM roommate_matches rm1
			JOIN roommate_matches rm2
			  ON rm1.from_user_id = rm2.to_user_id
			 AND rm1.to_user_id = rm2.from_user_id
			WHERE rm1.from_user_id = $1
			  AND rm1.to_user_id = $2
		)
	`

	var isMatched bool
	err := r.pool.QueryRow(ctx, sql, fromUserID, toUserID).Scan(&isMatched)
	if err != nil {
		return false, repository.HandelPgErrors(err)
	}

	return isMatched, nil
}

func (r *UserRepo) GetRoommateContacts(ctx context.Context, userID int) (*dto.RoommateContactsDTO, error) {
	sql := `
		SELECT u.email, p.phone
		FROM users u
		JOIN profiles p ON p.id = u.profile_id
		WHERE u.id = $1
	`

	rows, err := r.pool.Query(ctx, sql, userID)
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}
	defer rows.Close()

	contacts, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[dto.RoommateContactsDTO])
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	return &contacts, nil
}

func (r *UserRepo) CreateRoommateForm(ctx context.Context, data dto.CreateRoommateFormRequest) error {
	sql := `
		INSERT INTO roommate_forms (user_id, gender, birthday, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var roommateFormID int
	err := r.pool.QueryRow(ctx, sql, data.UserID, data.Gender, data.Birthday, data.Description).Scan(&roommateFormID)
	if err != nil {
		return repository.HandelPgErrors(err)
	}

	for _, tag := range data.Tags {
		sql = `
			INSERT INTO roommate_form_tags (roommate_form_id, roommate_tag_id)
			SELECT $1, rt.id
			FROM roommate_tags rt
			WHERE rt.alias = $2
		`

		ct, err := r.pool.Exec(ctx, sql, roommateFormID, tag)
		if err != nil {
			return repository.HandelPgErrors(err)
		}

		if ct.RowsAffected() == 0 {
			return repository.HandelPgErrors(pgx.ErrNoRows)
		}
	}

	return nil
}

func (r *UserRepo) GetRoommateForm(ctx context.Context, userID int) (*entity.RoommateForm, error) {
	sql := `
		SELECT gender, birthday::TEXT AS birthday, description
		FROM roommate_forms
		WHERE user_id = $1
	`

	rows, err := r.pool.Query(ctx, sql, userID)
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}
	defer rows.Close()

	form, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[entity.RoommateForm])
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	return &form, nil
}

func (r *UserRepo) GetRoommateFormTags(ctx context.Context, userID int) ([]string, error) {
	sql := `
		SELECT rt.alias
		FROM roommate_tags rt
		JOIN roommate_form_tags rft ON rft.roommate_tag_id = rt.id
		JOIN roommate_forms rf ON rf.id = rft.roommate_form_id
		WHERE rf.user_id = $1
		ORDER BY rt.alias
	`

	rows, err := r.pool.Query(ctx, sql, userID)
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}
	defer rows.Close()

	tags, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	return tags, nil
}

func (r *UserRepo) UpdateRoommateForm(ctx context.Context, data dto.CreateRoommateFormRequest) error {
	sql := `
		UPDATE roommate_forms
		SET gender = $1,
			birthday = $2,
			description = $3,
			updated_at = NOW()
		WHERE user_id = $4
		RETURNING id
	`

	var roommateFormID int
	err := r.pool.QueryRow(ctx, sql, data.Gender, data.Birthday, data.Description, data.UserID).Scan(&roommateFormID)
	if err != nil {
		return repository.HandelPgErrors(err)
	}

	sql = `
		DELETE FROM roommate_form_tags
		WHERE roommate_form_id = $1
	`

	_, err = r.pool.Exec(ctx, sql, roommateFormID)
	if err != nil {
		return repository.HandelPgErrors(err)
	}

	for _, tag := range data.Tags {
		sql = `
			INSERT INTO roommate_form_tags (roommate_form_id, roommate_tag_id)
			SELECT $1, rt.id
			FROM roommate_tags rt
			WHERE rt.alias = $2
		`

		ct, err := r.pool.Exec(ctx, sql, roommateFormID, tag)
		if err != nil {
			return repository.HandelPgErrors(err)
		}

		if ct.RowsAffected() == 0 {
			return repository.HandelPgErrors(pgx.ErrNoRows)
		}
	}

	return nil
}

func (r *UserRepo) GetIncomingRoommateMatches(ctx context.Context, userID int) ([]dto.RoommateUserDTO, error) {
	sql := `
		SELECT u.id, p.first_name, p.last_name, p.avatar_url, poster.alias AS poster_alias
		FROM roommate_matches rm
		JOIN users u ON u.id = rm.from_user_id
		JOIN profiles p ON p.id = u.profile_id
		JOIN posters poster ON poster.id = rm.poster_id
		WHERE rm.to_user_id = $1
		  AND poster.deleted_at IS NULL
		  AND NOT EXISTS (
			  SELECT 1
			  FROM roommate_matches reverse_rm
			  WHERE reverse_rm.from_user_id = rm.to_user_id
			    AND reverse_rm.to_user_id = rm.from_user_id
		  )
		ORDER BY rm.created_at DESC
	`

	rows, err := r.pool.Query(ctx, sql, userID)
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[dto.RoommateUserDTO])
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	return users, nil
}

func (r *UserRepo) GetMatchedRoommateMatches(ctx context.Context, userID int) ([]dto.RoommateUserDTO, error) {
	sql := `
		SELECT u.id, p.first_name, p.last_name, p.avatar_url, poster.alias AS poster_alias
		FROM roommate_matches rm1
		JOIN roommate_matches rm2
		  ON rm1.from_user_id = rm2.to_user_id
		 AND rm1.to_user_id = rm2.from_user_id
		JOIN users u ON u.id = rm1.to_user_id
		JOIN profiles p ON p.id = u.profile_id
		JOIN posters poster ON poster.id = rm1.poster_id
		WHERE rm1.from_user_id = $1
		  AND poster.deleted_at IS NULL
		ORDER BY rm1.created_at DESC
	`

	rows, err := r.pool.Query(ctx, sql, userID)
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[dto.RoommateUserDTO])
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	return users, nil
}
