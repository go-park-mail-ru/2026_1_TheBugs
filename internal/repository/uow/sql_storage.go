package uow

import (
	"context"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/poster"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/user"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase"
)

type SQLStorage struct {
	pool repository.DB
}

func (s *SQLStorage) Users() usecase.UserRepo {
	return user.NewUserRepo(s.pool)
}

//	func (s *SQLStorage) Autho() usecase.AuthRepo {
//		return auth.NewAuthRepo(s.pool)
//	}
func (s *SQLStorage) Poster() usecase.PosterRepo {
	return poster.NewPosterRepo(s.pool)
}

func (s *SQLStorage) Do(
	ctx context.Context,
	fn func(r usecase.Repositories) error,
) error {

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	txStore := &SQLStorage{
		pool: s.pool,
	}

	if err := fn(txStore); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
