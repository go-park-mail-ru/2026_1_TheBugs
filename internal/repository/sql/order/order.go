package order

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	repository "github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/sql"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/ctxLogger"
	"github.com/jackc/pgx/v5"
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

func (r *OrderRepo) GetByID(ctx context.Context, orderID int) (*entity.OrderFull, error) {
	log := ctxLogger.GetLogger(ctx).WithField("op", "OrderRepo.GetByID")
	log.Info("start db query")

	query := `
		SELECT 
			h.id,
			h.user_id,
			hc.name,
			h.status,
			h.description,
			h.created_at
		FROM handlings h
		JOIN handling_categories hc ON hc.id = h.category_id
		WHERE h.id = $1
	`

	var order entity.OrderFull

	err := r.pool.QueryRow(ctx, query, orderID).Scan(
		&order.ID,
		&order.UserID,
		&order.CategoryName,
		&order.Status,
		&order.Description,
		&order.CreatedAt,
	)
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	return &order, nil
}

func (r *OrderRepo) GetOrderImages(ctx context.Context, id int) ([]entity.OrderPhoto, error) {
	query := `
		SELECT hp.img_url, hp.sequence_order AS "order"
		FROM handling_photos hp
		WHERE hp.handling_id = $1
		ORDER BY hp.sequence_order
	`

	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	defer rows.Close()

	images, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.OrderPhoto])
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	return images, rows.Err()
}

func (r *OrderRepo) FinishOrder(ctx context.Context, orderID int, adminID int) error {
	query := `
		UPDATE handlings
		SET status = 'finished',
		    admin_id = $1,
		    updated_at = NOW()
		WHERE id = $2
	`

	_, err := r.pool.Exec(ctx, query, adminID, orderID)
	if err != nil {
		return repository.HandelPgErrors(err)
	}

	return nil
}
