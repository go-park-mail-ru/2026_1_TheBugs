package order

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	repository "github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/sql"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/ctxLogger"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type OrderRepo struct {
	pool repository.DB
}

func NewOrderRepo(pool repository.DB) *OrderRepo {
	return &OrderRepo{
		pool: pool,
	}
}

func (r *OrderRepo) Create(ctx context.Context, order *dto.Order) (int, error) {
	log := ctxLogger.GetLogger(ctx).WithField("op", "OrderRepo.Create")
	log.Info("start db query")

	query := `
		INSERT INTO handlings (user_id, category_id, description)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var handlingID int
	err := r.pool.QueryRow(ctx, query,
		order.UserID,
		order.CategoryID,
		order.Description,
	).Scan(&handlingID)

	if err != nil {
		return 0, repository.HandelPgErrors(err)
	}

	return handlingID, nil
}

func (r *OrderRepo) InsertPhotos(ctx context.Context, orderID int, photos []dto.PhotoInput) error {
	log := ctxLogger.GetLogger(ctx).WithField("op", "HandlingRepo.InsertPhotos")
	log.Info("start db query")

	if len(photos) == 0 {
		return nil
	}

	args := make([]any, 0, len(photos)*3)
	list := make([]string, 0, len(photos))

	for i, photo := range photos {
		base := i * 3
		list = append(list, fmt.Sprintf("($%d, $%d, $%d)", base+1, base+2, base+3))
		args = append(args, photo.Path, photo.Order, orderID)
	}

	query := `
		INSERT INTO handling_photos (img_url, sequence_order, handling_id)
		VALUES ` + strings.Join(list, ",")

	_, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return repository.HandelPgErrors(err)
	}

	return nil
}

func (r *OrderRepo) GetByUserID(ctx context.Context, userID int) ([]entity.Order, error) {
	log := ctxLogger.GetLogger(ctx).WithField("op", "OrderRepo.GetByUserID")
	log.Info("start db query")

	query := `
		SELECT 
			h.id,
			hc.name,
			h.status,
			h.created_at
		FROM handlings h
		JOIN handling_categories hc ON hc.id = h.category_id
		WHERE h.user_id = $1
		ORDER BY h.created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}
	defer rows.Close()

	orders := make([]entity.Order, 0)

	for rows.Next() {
		var o entity.Order

		err = rows.Scan(
			&o.ID,
			&o.CategoryName,
			&o.Status,
			&o.CreatedAt,
		)
		if err != nil {
			return nil, repository.HandelPgErrors(err)
		}

		orders = append(orders, o)
	}

	if err = rows.Err(); err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	return orders, nil
}

func (r *OrderRepo) GetAll(ctx context.Context) ([]entity.Order, error) {
	log := ctxLogger.GetLogger(ctx).WithField("op", "OrderRepo.GetAll")
	log.Info("start db query")

	query := `
		SELECT 
			h.id,
			hc.name,
			h.status,
			h.created_at
		FROM handlings h
		JOIN handling_categories hc ON hc.id = h.category_id
		ORDER BY h.created_at DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}
	defer rows.Close()

	orders := make([]entity.Order, 0)

	for rows.Next() {
		var o entity.Order

		err = rows.Scan(
			&o.ID,
			&o.CategoryName,
			&o.Status,
			&o.CreatedAt,
		)
		if err != nil {
			return nil, repository.HandelPgErrors(err)
		}

		orders = append(orders, o)
	}

	if err = rows.Err(); err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	return orders, nil
}
