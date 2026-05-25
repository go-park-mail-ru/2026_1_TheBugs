package user

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestGetUserByEmail(t *testing.T) {
	query := regexp.QuoteMeta(`SELECT id, email, salt, hashed_password, provider, is_verified FROM users WHERE email=$1`)
	expectedUser := &entity.User{
		ID:             1,
		Email:          "test@mail.ru",
		Salt:           lo.ToPtr("salt123"),
		HashedPassword: lo.ToPtr("hash123"),
	}

	tests := []struct {
		name      string
		email     string
		setupMock func(m pgxmock.PgxPoolIface)
		want      *entity.User
		wantErr   error
	}{
		{
			name:  "OK",
			email: "test@mail.ru",
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "email", "salt", "hashed_password", "provider", "is_verified",
				}).AddRow(1, "test@mail.ru", lo.ToPtr("salt123"), lo.ToPtr("hash123"), nil, false)

				m.ExpectQuery(query).WithArgs("test@mail.ru").WillReturnRows(rows)
			},
			want:    expectedUser,
			wantErr: nil,
		},
		{
			name:  "not_found",
			email: "missing@mail.ru",
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "email", "salt", "hashed_password", "provider", "is_verified",
				})

				m.ExpectQuery(query).WithArgs("missing@mail.ru").WillReturnRows(rows)
			},
			want:    nil,
			wantErr: entity.NotFoundError,
		},
		{
			name:  "collect_row_error",
			email: "test@mail.ru",
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "email", "salt", "hashed_password", "provider", "is_verified",
				}).AddRow("bad_id", "test@mail.ru", lo.ToPtr("salt123"), lo.ToPtr("hash123"), nil, false)

				m.ExpectQuery(query).WithArgs("test@mail.ru").WillReturnRows(rows)
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

			repo := NewUserRepo(mock)

			got, err := repo.GetByEmail(context.Background(), test.email)
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

func TestGetUserByID(t *testing.T) {
	query := regexp.QuoteMeta(`SELECT u.id, u.email, p.first_name, p.last_name, p.avatar_url, p.phone
			FROM users u
			JOIN profiles p ON u.profile_id = p.id
			WHERE u.id=$1`)

	expectedUser := &dto.UserDTO{
		ID:        1,
		Email:     "test@mail.ru",
		FirstName: "John",
		LastName:  "Doe",
		AvatarURL: nil,
		Phone:     "123456",
	}

	tests := []struct {
		name      string
		id        int
		setupMock func(m pgxmock.PgxPoolIface)
		want      *dto.UserDTO
		wantErr   error
	}{
		{
			name: "OK",
			id:   1,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "email", "first_name", "last_name", "avatar_url", "phone",
				}).AddRow(1, "test@mail.ru", "John", "Doe", nil, "123456")

				m.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)
			},
			want:    expectedUser,
			wantErr: nil,
		},
		{
			name: "not_found",
			id:   999,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "email", "first_name", "last_name", "avatar_url", "phone",
				})

				m.ExpectQuery(query).WithArgs(999).WillReturnRows(rows)
			},
			want:    nil,
			wantErr: entity.NotFoundError,
		},
		{
			name: "collect_row_error",
			id:   1,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "email", "first_name", "last_name", "avatar_url", "phone",
				}).AddRow("bad_id", "test@mail.ru", "John", "Doe", nil, "123456")

				m.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)
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

			repo := NewUserRepo(mock)

			got, err := repo.GetByID(context.Background(), test.id)
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

func TestUpdateProfile(t *testing.T) {
	query := regexp.QuoteMeta(`UPDATE profiles p
			SET phone = COALESCE($1, p.phone),
				first_name = COALESCE($2, p.first_name),
				last_name = COALESCE($3, p.last_name),
				avatar_url = COALESCE($4, p.avatar_url)
			FROM users u
			WHERE p.id = u.profile_id AND u.id = $5  -- Join users here
			RETURNING p.id, u.email, p.first_name, p.last_name, p.avatar_url, p.phone;`)

	phone := lo.ToPtr("123456")
	firstName := lo.ToPtr("John")
	lastName := lo.ToPtr("Doe")
	var avatarPath *string

	inputDTO := dto.UpdateProfileDTO{
		ID:         1,
		Phone:      phone,
		FirstName:  firstName,
		LastName:   lastName,
		AvatarPath: avatarPath,
	}

	expectedUser := &dto.UserDTO{
		ID:        1,
		Email:     "test@mail.ru",
		FirstName: "John",
		LastName:  "Doe",
		AvatarURL: nil,
		Phone:     "123456",
	}

	tests := []struct {
		name      string
		dto       dto.UpdateProfileDTO
		setupMock func(m pgxmock.PgxPoolIface)
		want      *dto.UserDTO
		wantErr   error
	}{
		{
			name: "OK",
			dto:  inputDTO,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "email", "first_name", "last_name", "avatar_url", "phone",
				}).AddRow(1, "test@mail.ru", "John", "Doe", nil, "123456")

				m.ExpectQuery(query).
					WithArgs(phone, firstName, lastName, avatarPath, 1).
					WillReturnRows(rows)
			},
			want:    expectedUser,
			wantErr: nil,
		},
		{
			name: "not_found",
			dto: dto.UpdateProfileDTO{
				ID:         999,
				Phone:      phone,
				FirstName:  firstName,
				LastName:   lastName,
				AvatarPath: avatarPath,
			},
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "email", "first_name", "last_name", "avatar_url", "phone",
				})

				m.ExpectQuery(query).
					WithArgs(phone, firstName, lastName, avatarPath, 999).
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: entity.NotFoundError,
		},
		{
			name: "collect_row_error",
			dto:  inputDTO,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "email", "first_name", "last_name", "avatar_url", "phone",
				}).AddRow("bad_id", "test@mail.ru", "John", "Doe", nil, "123456")

				m.ExpectQuery(query).
					WithArgs(phone, firstName, lastName, avatarPath, 1).
					WillReturnRows(rows)
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

			repo := NewUserRepo(mock)

			got, err := repo.UpdateProfile(context.Background(), test.dto)
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

func TestCreateUser(t *testing.T) {
	profileQuery := regexp.QuoteMeta(`INSERT INTO profiles (phone, first_name, last_name) VALUES ($1, $2, $3) RETURNING id`)
	userQuery := regexp.QuoteMeta(`INSERT INTO users (email, hashed_password, salt, profile_id) VALUES ($1, $2, $3, $4) 
			RETURNING id, email, hashed_password, salt, provider, is_verified`)

	phone := "123456"
	firstName := "John"
	lastName := "Doe"

	inputDTO := dto.CreateUserDTO{
		Email:          "test@mail.ru",
		Salt:           lo.ToPtr("salt123"),
		HashedPassword: lo.ToPtr("hash123"),
		Phone:          phone,
		FirstName:      firstName,
		LastName:       lastName,
	}

	expectedUser := &entity.User{
		ID:             1,
		Email:          "test@mail.ru",
		Salt:           lo.ToPtr("salt123"),
		HashedPassword: lo.ToPtr("hash123"),
	}

	tests := []struct {
		name      string
		dto       dto.CreateUserDTO
		setupMock func(m pgxmock.PgxPoolIface)
		want      *entity.User
		wantErr   error
	}{
		{
			name: "OK",
			dto:  inputDTO,
			setupMock: func(m pgxmock.PgxPoolIface) {
				profileRows := pgxmock.NewRows([]string{"id"}).AddRow(1)
				m.ExpectQuery(profileQuery).WithArgs("123456", "John", "Doe").WillReturnRows(profileRows)

				rows := pgxmock.NewRows([]string{
					"id", "email", "hashed_password", "salt", "provider", "is_verified",
				}).AddRow(1, "test@mail.ru", lo.ToPtr("hash123"), lo.ToPtr("salt123"), nil, false)

				m.ExpectQuery(userQuery).WithArgs("test@mail.ru", "hash123", "salt123", 1).WillReturnRows(rows)
			},
			want:    expectedUser,
			wantErr: nil,
		},
		{
			name: "collect_row_error",
			dto:  inputDTO,
			setupMock: func(m pgxmock.PgxPoolIface) {
				profileRows := pgxmock.NewRows([]string{"id"}).AddRow(1)
				m.ExpectQuery(profileQuery).WithArgs("123456", "John", "Doe").WillReturnRows(profileRows)

				rows := pgxmock.NewRows([]string{
					"id", "email", "hashed_password", "salt", "provider", "is_verified",
				}).AddRow("bad_id", "test@mail.ru", lo.ToPtr("hash123"), lo.ToPtr("salt123"), nil, false)

				m.ExpectQuery(userQuery).WithArgs("test@mail.ru", "hash123", "salt123", 1).WillReturnRows(rows)
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

			repo := NewUserRepo(mock)

			got, err := repo.Create(context.Background(), test.dto)
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

func TestCreateUserByProvider(t *testing.T) {
	profileQuery := regexp.QuoteMeta(`INSERT INTO profiles (phone, first_name, last_name) VALUES ($1, $2, $3) RETURNING id`)
	userQuery := regexp.QuoteMeta(`INSERT INTO users (email, provider, provider_id, profile_id, is_verified) VALUES ($1, $2, $3, $4, TRUE) 
			RETURNING id, email, hashed_password, salt, provider, is_verified`)

	phone := "123456"
	firstName := "John"
	lastName := "Doe"
	provider := entity.ProviderType("provider")
	providerID := lo.ToPtr("provider-123")

	inputDTO := dto.CreateUserByProviderDTO{
		Email:      "test@mail.ru",
		Provider:   provider,
		ProviderID: providerID,
		Phone:      phone,
		FirstName:  firstName,
		LastName:   lastName,
	}

	expectedUser := &entity.User{
		ID:             1,
		Email:          "test@mail.ru",
		Salt:           nil,
		HashedPassword: nil,
		Provider:       lo.ToPtr("provider"),
	}

	tests := []struct {
		name      string
		dto       dto.CreateUserByProviderDTO
		setupMock func(m pgxmock.PgxPoolIface)
		want      *entity.User
		wantErr   error
	}{
		{
			name: "OK",
			dto:  inputDTO,
			setupMock: func(m pgxmock.PgxPoolIface) {
				profileRows := pgxmock.NewRows([]string{"id"}).AddRow(1)
				m.ExpectQuery(profileQuery).WithArgs("123456", "John", "Doe").WillReturnRows(profileRows)

				rows := pgxmock.NewRows([]string{
					"id", "email", "hashed_password", "salt", "provider", "is_verified",
				}).AddRow(1, "test@mail.ru", nil, nil, lo.ToPtr("provider"), false)

				m.ExpectQuery(userQuery).WithArgs("test@mail.ru", provider, providerID, 1).WillReturnRows(rows)
			},
			want:    expectedUser,
			wantErr: nil,
		},
		{
			name: "collect_row_error",
			dto:  inputDTO,
			setupMock: func(m pgxmock.PgxPoolIface) {
				profileRows := pgxmock.NewRows([]string{"id"}).AddRow(1)
				m.ExpectQuery(profileQuery).WithArgs("123456", "John", "Doe").WillReturnRows(profileRows)

				rows := pgxmock.NewRows([]string{
					"id", "email", "hashed_password", "salt", "provider", "is_verified",
				}).AddRow("bad_id", "test@mail.ru", nil, nil, lo.ToPtr("provider"), false)

				m.ExpectQuery(userQuery).WithArgs("test@mail.ru", provider, providerID, 1).WillReturnRows(rows)
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

			repo := NewUserRepo(mock)

			got, err := repo.CreateByProvider(context.Background(), test.dto)
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

func TestGetUserByProvider(t *testing.T) {
	query := regexp.QuoteMeta(`SELECT id, email, salt, hashed_password, provider, is_verified FROM users WHERE email=$1 AND provider=$2`)
	provider := "provider"

	expectedUser := &entity.User{
		ID:             1,
		Email:          "test@mail.ru",
		Salt:           nil,
		HashedPassword: nil,
		Provider:       lo.ToPtr(provider),
	}

	tests := []struct {
		name      string
		provider  string
		email     string
		setupMock func(m pgxmock.PgxPoolIface)
		want      *entity.User
		wantErr   error
	}{
		{
			name:     "OK",
			provider: provider,
			email:    "test@mail.ru",
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "email", "salt", "hashed_password", "provider", "is_verified",
				}).AddRow(1, "test@mail.ru", nil, nil, lo.ToPtr(provider), false)

				m.ExpectQuery(query).
					WithArgs("test@mail.ru", provider).
					WillReturnRows(rows)
			},
			want:    expectedUser,
			wantErr: nil,
		},
		{
			name:     "not_found",
			provider: provider,
			email:    "missing@mail.ru",
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "email", "salt", "hashed_password", "provider", "is_verified",
				})

				m.ExpectQuery(query).
					WithArgs("missing@mail.ru", provider).
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: entity.NotFoundError,
		},
		{
			name:     "collect_row_error",
			provider: provider,
			email:    "test@mail.ru",
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "email", "salt", "hashed_password", "provider", "is_verified",
				}).AddRow("bad_id", "test@mail.ru", nil, nil, lo.ToPtr(provider), false)

				m.ExpectQuery(query).
					WithArgs("test@mail.ru", provider).
					WillReturnRows(rows)
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

			repo := NewUserRepo(mock)

			got, err := repo.GetByProvider(context.Background(), test.provider, test.email)
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

func TestUpdatePwd(t *testing.T) {
	query := regexp.QuoteMeta(`UPDATE users SET hashed_password=$1, salt=$2 WHERE email=$3`)

	email := "test@mail.ru"
	pwd := "new_hash"
	salt := "new_salt"

	tests := []struct {
		name      string
		email     string
		pwd       string
		salt      string
		setupMock func(m pgxmock.PgxPoolIface)
		wantErr   error
	}{
		{
			name:  "OK",
			email: email,
			pwd:   pwd,
			salt:  salt,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(query).
					WithArgs(pwd, salt, email).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			wantErr: nil,
		},
		{
			name:  "not_found",
			email: email,
			pwd:   pwd,
			salt:  salt,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(query).
					WithArgs(pwd, salt, email).
					WillReturnResult(pgxmock.NewResult("UPDATE", 0))
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

			repo := NewUserRepo(mock)

			err = repo.UpdatePwd(context.Background(), test.email, test.pwd, test.salt)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetRoommateUser(t *testing.T) {
	query := regexp.QuoteMeta(`
		SELECT p.first_name, p.last_name, p.avatar_url,
			   rf.gender, rf.birthday::TEXT AS birthday, rf.description
		FROM users u
		JOIN profiles p ON u.profile_id = p.id
		JOIN roommate_forms rf ON rf.user_id = u.id
		WHERE u.id = $1
	`)

	expectedUser := &entity.RoommateUser{
		FirstName:   "John",
		LastName:    "Doe",
		AvatarURL:   lo.ToPtr("https://example.com/avatar.jpg"),
		Gender:      "male",
		Birthday:    "2000-01-01",
		Description: lo.ToPtr("good roommate"),
	}

	tests := []struct {
		name      string
		id        int
		setupMock func(m pgxmock.PgxPoolIface)
		want      *entity.RoommateUser
		wantErr   error
	}{
		{
			name: "OK",
			id:   42,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"first_name", "last_name", "avatar_url", "gender", "birthday", "description",
				}).AddRow(
					"John",
					"Doe",
					lo.ToPtr("https://example.com/avatar.jpg"),
					"male",
					"2000-01-01",
					lo.ToPtr("good roommate"),
				)

				m.ExpectQuery(query).
					WithArgs(42).
					WillReturnRows(rows)
			},
			want:    expectedUser,
			wantErr: nil,
		},
		{
			name: "not_found",
			id:   999,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"first_name", "last_name", "avatar_url", "gender", "birthday", "description",
				})

				m.ExpectQuery(query).
					WithArgs(999).
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: entity.NotFoundError,
		},
		{
			name: "collect_row_error",
			id:   42,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"first_name", "last_name", "avatar_url", "gender", "birthday", "description",
				}).AddRow(
					"John",
					"Doe",
					lo.ToPtr("https://example.com/avatar.jpg"),
					123,
					"2000-01-01",
					lo.ToPtr("good roommate"),
				)

				m.ExpectQuery(query).
					WithArgs(42).
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
		{
			name: "query_error",
			id:   42,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(query).
					WithArgs(42).
					WillReturnError(errors.New("db error"))
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

			repo := NewUserRepo(mock)

			got, err := repo.GetRoommateUser(context.Background(), test.id)
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

func TestGetRoommateTags(t *testing.T) {
	query := regexp.QuoteMeta(`
		SELECT rt.name, rt.alias
		FROM roommate_tags rt
		JOIN roommate_form_tags rft ON rft.roommate_tag_id = rt.id
		JOIN roommate_forms rf ON rf.id = rft.roommate_form_id
		WHERE rf.user_id = $1
		ORDER BY rt.name
	`)

	expectedTags := []entity.RoommateTag{
		{
			Name:  "Без животных",
			Alias: "no_pets",
		},
		{
			Name:  "Не курю",
			Alias: "no_smoking",
		},
	}

	tests := []struct {
		name      string
		id        int
		setupMock func(m pgxmock.PgxPoolIface)
		want      []entity.RoommateTag
		wantErr   error
	}{
		{
			name: "OK",
			id:   42,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"name", "alias",
				}).AddRow(
					"Без животных",
					"no_pets",
				).AddRow(
					"Не курю",
					"no_smoking",
				)

				m.ExpectQuery(query).
					WithArgs(42).
					WillReturnRows(rows)
			},
			want:    expectedTags,
			wantErr: nil,
		},
		{
			name: "empty",
			id:   42,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"name", "alias",
				})

				m.ExpectQuery(query).
					WithArgs(42).
					WillReturnRows(rows)
			},
			want:    []entity.RoommateTag{},
			wantErr: nil,
		},
		{
			name: "collect_rows_error",
			id:   42,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"name", "alias",
				}).AddRow(
					123,
					"no_smoking",
				)

				m.ExpectQuery(query).
					WithArgs(42).
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
		{
			name: "query_error",
			id:   42,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(query).
					WithArgs(42).
					WillReturnError(errors.New("db error"))
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

			repo := NewUserRepo(mock)

			got, err := repo.GetRoommateTags(context.Background(), test.id)
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

func TestAddRoommateMatch(t *testing.T) {
	insertWithAliasQuery := regexp.QuoteMeta(`
		INSERT INTO roommate_matches (from_user_id, to_user_id, poster_id)
		SELECT $1, $2, p.id
		FROM posters p
		WHERE p.alias = $3
		  AND p.deleted_at IS NULL
	`)

	insertWithoutAliasQuery := regexp.QuoteMeta(`
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
	`)

	fromUserID := 10
	toUserID := 42
	posterAlias := "flat-1"

	tests := []struct {
		name        string
		fromUserID  int
		toUserID    int
		posterAlias *string
		setupMock   func(m pgxmock.PgxPoolIface)
		wantErr     error
	}{
		{
			name:        "OK_with_alias",
			fromUserID:  fromUserID,
			toUserID:    toUserID,
			posterAlias: &posterAlias,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(insertWithAliasQuery).
					WithArgs(fromUserID, toUserID, posterAlias).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			wantErr: nil,
		},
		{
			name:        "OK_without_alias_accept_request",
			fromUserID:  fromUserID,
			toUserID:    toUserID,
			posterAlias: nil,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(insertWithoutAliasQuery).
					WithArgs(fromUserID, toUserID).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			wantErr: nil,
		},
		{
			name:        "poster_not_found",
			fromUserID:  fromUserID,
			toUserID:    toUserID,
			posterAlias: &posterAlias,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(insertWithAliasQuery).
					WithArgs(fromUserID, toUserID, posterAlias).
					WillReturnResult(pgxmock.NewResult("INSERT", 0))
			},
			wantErr: entity.NotFoundError,
		},
		{
			name:        "already_exists_with_alias",
			fromUserID:  fromUserID,
			toUserID:    toUserID,
			posterAlias: &posterAlias,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(insertWithAliasQuery).
					WithArgs(fromUserID, toUserID, posterAlias).
					WillReturnError(&pgconn.PgError{Code: "23505"})
			},
			wantErr: entity.AlredyExitError,
		},
		{
			name:        "already_exists_without_alias",
			fromUserID:  fromUserID,
			toUserID:    toUserID,
			posterAlias: nil,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(insertWithoutAliasQuery).
					WithArgs(fromUserID, toUserID).
					WillReturnError(&pgconn.PgError{Code: "23505"})
			},
			wantErr: entity.AlredyExitError,
		},
		{
			name:        "service_error_with_alias",
			fromUserID:  fromUserID,
			toUserID:    toUserID,
			posterAlias: &posterAlias,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(insertWithAliasQuery).
					WithArgs(fromUserID, toUserID, posterAlias).
					WillReturnError(errors.New("db error"))
			},
			wantErr: entity.ServiceError,
		},
		{
			name:        "service_error_without_alias",
			fromUserID:  fromUserID,
			toUserID:    toUserID,
			posterAlias: nil,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(insertWithoutAliasQuery).
					WithArgs(fromUserID, toUserID).
					WillReturnError(errors.New("db error"))
			},
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewUserRepo(mock)

			err = repo.AddRoommateMatch(context.Background(), test.fromUserID, test.toUserID, test.posterAlias)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetRoommateContacts(t *testing.T) {
	query := regexp.QuoteMeta(`
		SELECT u.email, p.phone
		FROM users u
		JOIN profiles p ON p.id = u.profile_id
		WHERE u.id = $1
	`)

	expectedContacts := &dto.RoommateContactsDTO{
		Email: "target@example.com",
		Phone: "+79991234567",
	}

	tests := []struct {
		name      string
		userID    int
		setupMock func(m pgxmock.PgxPoolIface)
		want      *dto.RoommateContactsDTO
		wantErr   error
	}{
		{
			name:   "OK",
			userID: 42,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"email", "phone",
				}).AddRow(
					"target@example.com",
					"+79991234567",
				)

				m.ExpectQuery(query).
					WithArgs(42).
					WillReturnRows(rows)
			},
			want:    expectedContacts,
			wantErr: nil,
		},
		{
			name:   "not_found",
			userID: 999,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"email", "phone",
				})

				m.ExpectQuery(query).
					WithArgs(999).
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: entity.NotFoundError,
		},
		{
			name:   "collect_row_error",
			userID: 42,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"email", "phone",
				}).AddRow(
					123,
					"+79991234567",
				)

				m.ExpectQuery(query).
					WithArgs(42).
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
		{
			name:   "query_error",
			userID: 42,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(query).
					WithArgs(42).
					WillReturnError(errors.New("db error"))
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

			repo := NewUserRepo(mock)

			got, err := repo.GetRoommateContacts(context.Background(), test.userID)
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

func TestIsRoommateMatch(t *testing.T) {
	query := regexp.QuoteMeta(`
		SELECT EXISTS(
			SELECT 1
			FROM roommate_matches rm1
			JOIN roommate_matches rm2
			  ON rm1.from_user_id = rm2.to_user_id
			 AND rm1.to_user_id = rm2.from_user_id
			WHERE rm1.from_user_id = $1
			  AND rm1.to_user_id = $2
		)
	`)

	fromUserID := 10
	toUserID := 42

	tests := []struct {
		name       string
		fromUserID int
		toUserID   int
		setupMock  func(m pgxmock.PgxPoolIface)
		want       bool
		wantErr    error
	}{
		{
			name:       "OK_true",
			fromUserID: fromUserID,
			toUserID:   toUserID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"exists"}).AddRow(true)

				m.ExpectQuery(query).
					WithArgs(fromUserID, toUserID).
					WillReturnRows(rows)
			},
			want:    true,
			wantErr: nil,
		},
		{
			name:       "OK_false",
			fromUserID: fromUserID,
			toUserID:   toUserID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"exists"}).AddRow(false)

				m.ExpectQuery(query).
					WithArgs(fromUserID, toUserID).
					WillReturnRows(rows)
			},
			want:    false,
			wantErr: nil,
		},
		{
			name:       "query_error",
			fromUserID: fromUserID,
			toUserID:   toUserID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(query).
					WithArgs(fromUserID, toUserID).
					WillReturnError(errors.New("db error"))
			},
			want:    false,
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewUserRepo(mock)

			got, err := repo.IsRoommateMatch(context.Background(), test.fromUserID, test.toUserID)
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

func TestCreateRoommateForm(t *testing.T) {
	formQuery := regexp.QuoteMeta(`
		INSERT INTO roommate_forms (user_id, gender, birthday, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`)

	tagQuery := regexp.QuoteMeta(`
		INSERT INTO roommate_form_tags (roommate_form_id, roommate_tag_id)
		SELECT $1, rt.id
		FROM roommate_tags rt
		WHERE rt.alias = $2
	`)

	input := dto.CreateRoommateFormRequest{
		UserID:      10,
		Gender:      "male",
		Birthday:    "2000-01-01",
		Description: "good roommate",
		Tags:        []string{"no_smoking", "no_pets"},
	}

	tests := []struct {
		name      string
		data      dto.CreateRoommateFormRequest
		setupMock func(m pgxmock.PgxPoolIface)
		wantErr   error
	}{
		{
			name: "OK",
			data: input,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id"}).AddRow(1)

				m.ExpectQuery(formQuery).
					WithArgs(input.UserID, input.Gender, input.Birthday, input.Description).
					WillReturnRows(rows)

				m.ExpectExec(tagQuery).
					WithArgs(1, "no_smoking").
					WillReturnResult(pgxmock.NewResult("INSERT", 1))

				m.ExpectExec(tagQuery).
					WithArgs(1, "no_pets").
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			wantErr: nil,
		},
		{
			name: "OK_empty_tags",
			data: dto.CreateRoommateFormRequest{
				UserID:      10,
				Gender:      "female",
				Birthday:    "2000-01-01",
				Description: "good roommate",
				Tags:        []string{},
			},
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id"}).AddRow(1)

				m.ExpectQuery(formQuery).
					WithArgs(10, "female", "2000-01-01", "good roommate").
					WillReturnRows(rows)
			},
			wantErr: nil,
		},
		{
			name: "form_already_exists",
			data: input,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(formQuery).
					WithArgs(input.UserID, input.Gender, input.Birthday, input.Description).
					WillReturnError(&pgconn.PgError{Code: "23505"})
			},
			wantErr: entity.AlredyExitError,
		},
		{
			name: "form_foreign_key_error",
			data: input,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(formQuery).
					WithArgs(input.UserID, input.Gender, input.Birthday, input.Description).
					WillReturnError(&pgconn.PgError{Code: "23503"})
			},
			wantErr: entity.ServiceError,
		},
		{
			name: "form_service_error",
			data: input,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(formQuery).
					WithArgs(input.UserID, input.Gender, input.Birthday, input.Description).
					WillReturnError(errors.New("db error"))
			},
			wantErr: entity.ServiceError,
		},
		{
			name: "tag_insert_error",
			data: input,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id"}).AddRow(1)

				m.ExpectQuery(formQuery).
					WithArgs(input.UserID, input.Gender, input.Birthday, input.Description).
					WillReturnRows(rows)

				m.ExpectExec(tagQuery).
					WithArgs(1, "no_smoking").
					WillReturnError(errors.New("db error"))
			},
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewUserRepo(mock)

			err = repo.CreateRoommateForm(context.Background(), test.data)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetRoommateForm(t *testing.T) {
	query := regexp.QuoteMeta(`
		SELECT gender, birthday::TEXT AS birthday, description
		FROM roommate_forms
		WHERE user_id = $1
	`)

	expectedForm := &entity.RoommateForm{
		Gender:      "male",
		Birthday:    "2000-01-01",
		Description: "good roommate",
	}

	tests := []struct {
		name      string
		userID    int
		setupMock func(m pgxmock.PgxPoolIface)
		want      *entity.RoommateForm
		wantErr   error
	}{
		{
			name:   "OK",
			userID: 10,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"gender", "birthday", "description",
				}).AddRow(
					"male",
					"2000-01-01",
					"good roommate",
				)

				m.ExpectQuery(query).
					WithArgs(10).
					WillReturnRows(rows)
			},
			want:    expectedForm,
			wantErr: nil,
		},
		{
			name:   "not_found",
			userID: 999,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"gender", "birthday", "description",
				})

				m.ExpectQuery(query).
					WithArgs(999).
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: entity.NotFoundError,
		},
		{
			name:   "collect_row_error",
			userID: 10,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"gender", "birthday", "description",
				}).AddRow(
					123,
					"2000-01-01",
					"good roommate",
				)

				m.ExpectQuery(query).
					WithArgs(10).
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
		{
			name:   "query_error",
			userID: 10,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(query).
					WithArgs(10).
					WillReturnError(errors.New("db error"))
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

			repo := NewUserRepo(mock)

			got, err := repo.GetRoommateForm(context.Background(), test.userID)
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

func TestGetRoommateFormTags(t *testing.T) {
	query := regexp.QuoteMeta(`
		SELECT rt.alias
		FROM roommate_tags rt
		JOIN roommate_form_tags rft ON rft.roommate_tag_id = rt.id
		JOIN roommate_forms rf ON rf.id = rft.roommate_form_id
		WHERE rf.user_id = $1
		ORDER BY rt.alias
	`)

	expectedTags := []string{"no_pets", "no_smoking"}

	tests := []struct {
		name      string
		userID    int
		setupMock func(m pgxmock.PgxPoolIface)
		want      []string
		wantErr   error
	}{
		{
			name:   "OK",
			userID: 10,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"alias"}).
					AddRow("no_pets").
					AddRow("no_smoking")

				m.ExpectQuery(query).
					WithArgs(10).
					WillReturnRows(rows)
			},
			want:    expectedTags,
			wantErr: nil,
		},
		{
			name:   "empty",
			userID: 10,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"alias"})

				m.ExpectQuery(query).
					WithArgs(10).
					WillReturnRows(rows)
			},
			want:    []string{},
			wantErr: nil,
		},
		{
			name:   "collect_rows_error",
			userID: 10,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"alias"}).
					AddRow(123)

				m.ExpectQuery(query).
					WithArgs(10).
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
		{
			name:   "query_error",
			userID: 10,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(query).
					WithArgs(10).
					WillReturnError(errors.New("db error"))
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

			repo := NewUserRepo(mock)

			got, err := repo.GetRoommateFormTags(context.Background(), test.userID)
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

func TestUpdateRoommateForm(t *testing.T) {
	formQuery := regexp.QuoteMeta(`
		UPDATE roommate_forms
		SET gender = $1,
			birthday = $2,
			description = $3,
			updated_at = NOW()
		WHERE user_id = $4
		RETURNING id
	`)

	deleteTagsQuery := regexp.QuoteMeta(`
		DELETE FROM roommate_form_tags
		WHERE roommate_form_id = $1
	`)

	tagQuery := regexp.QuoteMeta(`
		INSERT INTO roommate_form_tags (roommate_form_id, roommate_tag_id)
		SELECT $1, rt.id
		FROM roommate_tags rt
		WHERE rt.alias = $2
	`)

	input := dto.CreateRoommateFormRequest{
		UserID:      10,
		Gender:      "female",
		Birthday:    "2001-02-03",
		Description: "updated roommate",
		Tags:        []string{"no_smoking", "no_pets"},
	}

	tests := []struct {
		name      string
		data      dto.CreateRoommateFormRequest
		setupMock func(m pgxmock.PgxPoolIface)
		wantErr   error
	}{
		{
			name: "OK",
			data: input,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id"}).AddRow(1)

				m.ExpectQuery(formQuery).
					WithArgs(input.Gender, input.Birthday, input.Description, input.UserID).
					WillReturnRows(rows)

				m.ExpectExec(deleteTagsQuery).
					WithArgs(1).
					WillReturnResult(pgxmock.NewResult("DELETE", 2))

				m.ExpectExec(tagQuery).
					WithArgs(1, "no_smoking").
					WillReturnResult(pgxmock.NewResult("INSERT", 1))

				m.ExpectExec(tagQuery).
					WithArgs(1, "no_pets").
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			wantErr: nil,
		},
		{
			name: "OK_empty_tags",
			data: dto.CreateRoommateFormRequest{
				UserID:      10,
				Gender:      "male",
				Birthday:    "2001-02-03",
				Description: "updated roommate",
				Tags:        []string{},
			},
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id"}).AddRow(1)

				m.ExpectQuery(formQuery).
					WithArgs("male", "2001-02-03", "updated roommate", 10).
					WillReturnRows(rows)

				m.ExpectExec(deleteTagsQuery).
					WithArgs(1).
					WillReturnResult(pgxmock.NewResult("DELETE", 0))
			},
			wantErr: nil,
		},
		{
			name: "form_not_found",
			data: input,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id"})

				m.ExpectQuery(formQuery).
					WithArgs(input.Gender, input.Birthday, input.Description, input.UserID).
					WillReturnRows(rows)
			},
			wantErr: entity.NotFoundError,
		},
		{
			name: "form_service_error",
			data: input,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(formQuery).
					WithArgs(input.Gender, input.Birthday, input.Description, input.UserID).
					WillReturnError(errors.New("db error"))
			},
			wantErr: entity.ServiceError,
		},
		{
			name: "delete_tags_error",
			data: input,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id"}).AddRow(1)

				m.ExpectQuery(formQuery).
					WithArgs(input.Gender, input.Birthday, input.Description, input.UserID).
					WillReturnRows(rows)

				m.ExpectExec(deleteTagsQuery).
					WithArgs(1).
					WillReturnError(errors.New("db error"))
			},
			wantErr: entity.ServiceError,
		},
		{
			name: "tag_insert_error",
			data: input,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id"}).AddRow(1)

				m.ExpectQuery(formQuery).
					WithArgs(input.Gender, input.Birthday, input.Description, input.UserID).
					WillReturnRows(rows)

				m.ExpectExec(deleteTagsQuery).
					WithArgs(1).
					WillReturnResult(pgxmock.NewResult("DELETE", 2))

				m.ExpectExec(tagQuery).
					WithArgs(1, "no_smoking").
					WillReturnError(errors.New("db error"))
			},
			wantErr: entity.ServiceError,
		},
		{
			name: "tag_not_found",
			data: input,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id"}).AddRow(1)

				m.ExpectQuery(formQuery).
					WithArgs(input.Gender, input.Birthday, input.Description, input.UserID).
					WillReturnRows(rows)

				m.ExpectExec(deleteTagsQuery).
					WithArgs(1).
					WillReturnResult(pgxmock.NewResult("DELETE", 2))

				m.ExpectExec(tagQuery).
					WithArgs(1, "no_smoking").
					WillReturnResult(pgxmock.NewResult("INSERT", 0))
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

			repo := NewUserRepo(mock)

			err = repo.UpdateRoommateForm(context.Background(), test.data)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetIncomingRoommateMatches(t *testing.T) {
	query := regexp.QuoteMeta(`
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
	`)

	userID := 10

	expectedUsers := []dto.RoommateUserDTO{
		{
			ID:          42,
			FirstName:   "John",
			LastName:    "Doe",
			AvatarURL:   lo.ToPtr("https://example.com/avatar.jpg"),
			PosterAlias: lo.ToPtr("flat-1"),
		},
		{
			ID:          43,
			FirstName:   "Jane",
			LastName:    "Smith",
			AvatarURL:   nil,
			PosterAlias: lo.ToPtr("flat-2"),
		},
	}

	tests := []struct {
		name      string
		userID    int
		setupMock func(m pgxmock.PgxPoolIface)
		want      []dto.RoommateUserDTO
		wantErr   error
	}{
		{
			name:   "OK",
			userID: userID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "first_name", "last_name", "avatar_url", "poster_alias",
				}).AddRow(
					42,
					"John",
					"Doe",
					lo.ToPtr("https://example.com/avatar.jpg"),
					lo.ToPtr("flat-1"),
				).AddRow(
					43,
					"Jane",
					"Smith",
					nil,
					lo.ToPtr("flat-2"),
				)

				m.ExpectQuery(query).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			want:    expectedUsers,
			wantErr: nil,
		},
		{
			name:   "empty",
			userID: userID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "first_name", "last_name", "avatar_url", "poster_alias",
				})

				m.ExpectQuery(query).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			want:    []dto.RoommateUserDTO{},
			wantErr: nil,
		},
		{
			name:   "collect_rows_error",
			userID: userID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "first_name", "last_name", "avatar_url", "poster_alias",
				}).AddRow(
					"bad_id",
					"John",
					"Doe",
					nil,
					lo.ToPtr("flat-1"),
				)

				m.ExpectQuery(query).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
		{
			name:   "query_error",
			userID: userID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(query).
					WithArgs(userID).
					WillReturnError(errors.New("db error"))
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

			repo := NewUserRepo(mock)

			got, err := repo.GetIncomingRoommateMatches(context.Background(), test.userID)
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

func TestGetMatchedRoommateMatches(t *testing.T) {
	query := regexp.QuoteMeta(`
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
	`)

	userID := 10

	expectedUsers := []dto.RoommateUserDTO{
		{
			ID:          42,
			FirstName:   "John",
			LastName:    "Doe",
			AvatarURL:   lo.ToPtr("https://example.com/avatar.jpg"),
			PosterAlias: lo.ToPtr("flat-1"),
		},
		{
			ID:          43,
			FirstName:   "Jane",
			LastName:    "Smith",
			AvatarURL:   nil,
			PosterAlias: lo.ToPtr("flat-2"),
		},
	}

	tests := []struct {
		name      string
		userID    int
		setupMock func(m pgxmock.PgxPoolIface)
		want      []dto.RoommateUserDTO
		wantErr   error
	}{
		{
			name:   "OK",
			userID: userID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "first_name", "last_name", "avatar_url", "poster_alias",
				}).AddRow(
					42,
					"John",
					"Doe",
					lo.ToPtr("https://example.com/avatar.jpg"),
					lo.ToPtr("flat-1"),
				).AddRow(
					43,
					"Jane",
					"Smith",
					nil,
					lo.ToPtr("flat-2"),
				)

				m.ExpectQuery(query).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			want:    expectedUsers,
			wantErr: nil,
		},
		{
			name:   "empty",
			userID: userID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "first_name", "last_name", "avatar_url", "poster_alias",
				})

				m.ExpectQuery(query).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			want:    []dto.RoommateUserDTO{},
			wantErr: nil,
		},
		{
			name:   "collect_rows_error",
			userID: userID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "first_name", "last_name", "avatar_url", "poster_alias",
				}).AddRow(
					"bad_id",
					"John",
					"Doe",
					nil,
					lo.ToPtr("flat-1"),
				)

				m.ExpectQuery(query).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: entity.ServiceError,
		},
		{
			name:   "query_error",
			userID: userID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectQuery(query).
					WithArgs(userID).
					WillReturnError(errors.New("db error"))
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

			repo := NewUserRepo(mock)

			got, err := repo.GetMatchedRoommateMatches(context.Background(), test.userID)
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

func TestDeletePosterRoommatesByUserID(t *testing.T) {
	query := regexp.QuoteMeta(`
		DELETE FROM poster_roommates
		WHERE user_id = $1
	`)

	userID := 10

	tests := []struct {
		name      string
		userID    int
		setupMock func(m pgxmock.PgxPoolIface)
		wantErr   error
	}{
		{
			name:   "OK",
			userID: userID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(query).
					WithArgs(userID).
					WillReturnResult(pgxmock.NewResult("DELETE", 2))
			},
			wantErr: nil,
		},
		{
			name:   "OK_empty",
			userID: userID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(query).
					WithArgs(userID).
					WillReturnResult(pgxmock.NewResult("DELETE", 0))
			},
			wantErr: nil,
		},
		{
			name:   "service_error",
			userID: userID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(query).
					WithArgs(userID).
					WillReturnError(errors.New("db error"))
			},
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewUserRepo(mock)

			err = repo.DeletePosterRoommatesByUserID(context.Background(), test.userID)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDeleteRoommateForm(t *testing.T) {
	deleteTagsQuery := regexp.QuoteMeta(`
		DELETE FROM roommate_form_tags rft
		USING roommate_forms rf
		WHERE rft.roommate_form_id = rf.id
		  AND rf.user_id = $1
	`)

	deleteFormQuery := regexp.QuoteMeta(`
		DELETE FROM roommate_forms
		WHERE user_id = $1
	`)

	userID := 10

	tests := []struct {
		name      string
		userID    int
		setupMock func(m pgxmock.PgxPoolIface)
		wantErr   error
	}{
		{
			name:   "OK",
			userID: userID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(deleteTagsQuery).
					WithArgs(userID).
					WillReturnResult(pgxmock.NewResult("DELETE", 2))

				m.ExpectExec(deleteFormQuery).
					WithArgs(userID).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
			wantErr: nil,
		},
		{
			name:   "OK_without_tags",
			userID: userID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(deleteTagsQuery).
					WithArgs(userID).
					WillReturnResult(pgxmock.NewResult("DELETE", 0))

				m.ExpectExec(deleteFormQuery).
					WithArgs(userID).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
			wantErr: nil,
		},
		{
			name:   "not_found",
			userID: userID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(deleteTagsQuery).
					WithArgs(userID).
					WillReturnResult(pgxmock.NewResult("DELETE", 0))

				m.ExpectExec(deleteFormQuery).
					WithArgs(userID).
					WillReturnResult(pgxmock.NewResult("DELETE", 0))
			},
			wantErr: entity.NotFoundError,
		},
		{
			name:   "delete_tags_error",
			userID: userID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(deleteTagsQuery).
					WithArgs(userID).
					WillReturnError(errors.New("db error"))
			},
			wantErr: entity.ServiceError,
		},
		{
			name:   "delete_form_error",
			userID: userID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(deleteTagsQuery).
					WithArgs(userID).
					WillReturnResult(pgxmock.NewResult("DELETE", 2))

				m.ExpectExec(deleteFormQuery).
					WithArgs(userID).
					WillReturnError(errors.New("db error"))
			},
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewUserRepo(mock)

			err = repo.DeleteRoommateForm(context.Background(), test.userID)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDeleteRoommateMatch(t *testing.T) {
	query := regexp.QuoteMeta(`
		DELETE FROM roommate_matches
		WHERE (from_user_id = $1 AND to_user_id = $2)
		   OR (from_user_id = $2 AND to_user_id = $1)
	`)

	fromUserID := 10
	toUserID := 42

	tests := []struct {
		name       string
		fromUserID int
		toUserID   int
		setupMock  func(m pgxmock.PgxPoolIface)
		wantErr    error
	}{
		{
			name:       "OK_delete_both",
			fromUserID: fromUserID,
			toUserID:   toUserID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(query).
					WithArgs(fromUserID, toUserID).
					WillReturnResult(pgxmock.NewResult("DELETE", 2))
			},
			wantErr: nil,
		},
		{
			name:       "OK_delete_one",
			fromUserID: fromUserID,
			toUserID:   toUserID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(query).
					WithArgs(fromUserID, toUserID).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
			wantErr: nil,
		},
		{
			name:       "not_found",
			fromUserID: fromUserID,
			toUserID:   toUserID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(query).
					WithArgs(fromUserID, toUserID).
					WillReturnResult(pgxmock.NewResult("DELETE", 0))
			},
			wantErr: entity.NotFoundError,
		},
		{
			name:       "service_error",
			fromUserID: fromUserID,
			toUserID:   toUserID,
			setupMock: func(m pgxmock.PgxPoolIface) {
				m.ExpectExec(query).
					WithArgs(fromUserID, toUserID).
					WillReturnError(errors.New("db error"))
			},
			wantErr: entity.ServiceError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			test.setupMock(mock)

			repo := NewUserRepo(mock)

			err = repo.DeleteRoommateMatch(context.Background(), test.fromUserID, test.toUserID)
			if test.wantErr != nil {
				require.ErrorIs(t, err, test.wantErr)
				return
			}

			require.NoError(t, err)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
