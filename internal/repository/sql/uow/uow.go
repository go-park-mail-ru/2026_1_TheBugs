package uow

import (
	"context"

	repository "github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/sql"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/sql/auth"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/sql/complex"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/sql/order"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/sql/poster"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/sql/user"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase"
)

type SQLStorage struct {
	db      repository.DB
	users   usecase.UserRepo
	auth    usecase.AuthRepo
	posters usecase.PosterRepo
	company usecase.UtilityCompanyRepo
	order   usecase.OrderRepo
}

func NewSQLStorage(db repository.DB) *SQLStorage {
	return &SQLStorage{
		db:      db,
		users:   user.NewUserRepo(db),
		auth:    auth.NewAuthRepo(db),
		posters: poster.NewPosterRepo(db),
		order:   order.NewOrderRepo(db),
		company: complex.NewUtilityCompanyRepo(db),
	}
}

func (s *SQLStorage) Users() usecase.UserRepo {
	return s.users
}

func (s *SQLStorage) UtilityCompany() usecase.UtilityCompanyRepo {
	return s.company
}

func (s *SQLStorage) Autho() usecase.AuthRepo {
	return s.auth
}

func (s *SQLStorage) Posters() usecase.PosterRepo {
	return s.posters
}

func (s *SQLStorage) Order() usecase.OrderRepo {
	return s.order
}

func (s *SQLStorage) Do(
	ctx context.Context,
	fn func(r usecase.UnitOfWork) error,
) error {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	txStore := NewSQLStorage(tx)

	if err := fn(txStore); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
