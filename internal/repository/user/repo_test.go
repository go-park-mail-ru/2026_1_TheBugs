package user

import (
	"context"
	"regexp"
	"testing"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"
)

func TestGetUserByEmail(t *testing.T) {
	query := regexp.QuoteMeta(`SELECT id, email, salt, hashed_password FROM users WHERE email=$1`)
	salt := "salt123"
	hash := "hash123"
	expectedUser := &entity.User{
		ID:             1,
		Email:          "test@mail.ru",
		Salt:           &salt,
		HashedPassword: &hash,
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
					"id", "email", "salt", "hashed_password",
				}).AddRow(1, "test@mail.ru", "salt123", "hash123")

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
					"id", "email", "salt", "hashed_password",
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
					"id", "email", "salt", "hashed_password",
				}).AddRow("bad_id", "test@mail.ru", "salt123", "hash123")

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

			got, err := repo.GetUserByEmail(context.Background(), test.email)
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
	profileQuery := regexp.QuoteMeta(`INSERT INTO profile (phone, first_name, last_name) VALUES ($1, $2, $3) RETURNING id`)
	userQuery := regexp.QuoteMeta(`INSERT INTO users (email, hashed_password, salt, profile_id) VALUES ($1, $2, $3, $4) 
			RETURNING id, email, hashed_password, salt`)

	hashedPwd := "hash123"
	salt := "salt123"
	phone := "123456"
	firstName := "John"
	lastName := "Doe"

	inputDTO := dto.CreateUserDTO{
		Email:          "test@mail.ru",
		HashedPassword: &hashedPwd,
		Salt:           &salt,
		Phone:          phone,
		FirstName:      firstName,
		LastName:       lastName,
	}
	salt = "salt123"
	hash := "hash123"
	expectedUser := &entity.User{
		ID:             1,
		Email:          "test@mail.ru",
		Salt:           &salt,
		HashedPassword: &hash,
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
					"id", "email", "hashed_password", "salt",
				}).AddRow(1, "test@mail.ru", "hash123", "salt123")

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
					"id", "email", "hashed_password", "salt",
				}).AddRow("bad_id", "test@mail.ru", "hash123", "salt123")

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

			got, err := repo.CreateUser(context.Background(), test.dto)
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
