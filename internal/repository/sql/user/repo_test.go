package user

import (
	"context"
	"regexp"
	"testing"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestGetUserByEmail(t *testing.T) {
	query := regexp.QuoteMeta(`SELECT id, email, salt, hashed_password, provider FROM users WHERE email=$1`)
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
					"id", "email", "salt", "hashed_password", "provider",
				}).AddRow(1, "test@mail.ru", lo.ToPtr("salt123"), lo.ToPtr("hash123"), nil)

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
					"id", "email", "salt", "hashed_password", "provider",
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
					"id", "email", "salt", "hashed_password", "provider",
				}).AddRow("bad_id", "test@mail.ru", lo.ToPtr("salt123"), lo.ToPtr("hash123"), nil)

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
	query := regexp.QuoteMeta(`SELECT u.id, u.email, p.first_name, p.last_name, p.avatar_url , p.phone
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
			RETURNING id, email, hashed_password, salt, provider`)

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
					"id", "email", "hashed_password", "salt", "provider",
				}).AddRow(1, "test@mail.ru", lo.ToPtr("hash123"), lo.ToPtr("salt123"), nil)

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
					"id", "email", "hashed_password", "salt", "provider",
				}).AddRow("bad_id", "test@mail.ru", lo.ToPtr("hash123"), lo.ToPtr("salt123"), nil)

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
	userQuery := regexp.QuoteMeta(`INSERT INTO users (email, provider, provider_id, profile_id) VALUES ($1, $2, $3, $4) 
			RETURNING id, email, hashed_password, salt, provider`)

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
					"id", "email", "hashed_password", "salt", "provider",
				}).AddRow(1, "test@mail.ru", nil, nil, lo.ToPtr("provider"))

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
					"id", "email", "hashed_password", "salt", "provider",
				}).AddRow("bad_id", "test@mail.ru", nil, nil, lo.ToPtr("provider"))

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
	query := regexp.QuoteMeta(`SELECT id, email, salt, hashed_password, provider FROM users WHERE email=$1 AND provider=$2`)
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
					"id", "email", "salt", "hashed_password", "provider",
				}).AddRow(1, "test@mail.ru", nil, nil, lo.ToPtr(provider))

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
					"id", "email", "salt", "hashed_password", "provider",
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
					"id", "email", "salt", "hashed_password", "provider",
				}).AddRow("bad_id", "test@mail.ru", nil, nil, lo.ToPtr(provider))

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
