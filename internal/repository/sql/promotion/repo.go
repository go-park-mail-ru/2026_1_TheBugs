package promotion

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	repository "github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository/sql"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/jackc/pgx/v5"
)

type PromotionRepo struct {
	pool repository.DB
}

func NewPromotionRepo(pool repository.DB) *PromotionRepo {
	return &PromotionRepo{
		pool: pool,
	}
}

func (r *PromotionRepo) Create(ctx context.Context, data dto.CreatePromotionDTO) (int, error) {
	query := `
		INSERT INTO posters_promotions (poster_id, promotion_id, ends_at, payment_id, amount_paid, user_id)
		SELECT $1, prom.id, $3, $4, $5, $6
		FROM promotions AS prom
		WHERE prom.code = $2
		RETURNING posters_promotions.id;
	`

	var promotionsID int
	err := r.pool.QueryRow(ctx, query,
		data.PosterID, data.PromotionCode, data.EndsAt, data.PaymentID, data.AmountPaid, data.UserID,
	).Scan(&promotionsID)

	if err != nil {
		return 0, repository.HandelPgErrors(err)
	}

	return promotionsID, nil
}

func (r *PromotionRepo) GetByCode(ctx context.Context, code string) (*entity.Promotion, error) {
	sql := `
		SELECT p.id, p.code, p.name,
               p.description, p.duration_days, p.price
        FROM promotions p
		WHERE p.code = $1;
	`
	rows, err := r.pool.Query(ctx, sql, code)
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}
	prom, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Promotion])
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}
	return &prom, nil
}

func (r *PromotionRepo) GetByPaymentID(ctx context.Context, paymentID string) (*entity.PosterPromotion, error) {
	sql := `
		SELECT p.id, p.poster_id, p.promotion_id,
               p.status, p.payment_id, p.user_id
        FROM posters_promotions p
		WHERE p.payment_id = $1;
	`
	rows, err := r.pool.Query(ctx, sql, paymentID)
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}
	prom, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.PosterPromotion])
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}
	return &prom, nil
}
func (r *PromotionRepo) GetActiveByPosterID(ctx context.Context, posterID int) (*entity.PosterPromotion, error) {
	sql := `
		SELECT p.id, p.poster_id, p.promotion_id,
               p.status, p.payment_id, p.user_id
        FROM posters_promotions p
		WHERE p.poster_id = $1 AND status = 'active'
		LIMIT 1;
	`
	rows, err := r.pool.Query(ctx, sql, posterID)
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}
	prom, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.PosterPromotion])
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}
	return &prom, nil
}
func (r *PromotionRepo) UpdateStatus(ctx context.Context, paymentID string, status string) error {
	sql := `
		UPDATE posters_promotions SET status = $1 WHERE payment_id = $2;
	`
	_, err := r.pool.Exec(ctx, sql, status, paymentID)
	if err != nil {
		return repository.HandelPgErrors(err)
	}

	return nil
}
func (r *PromotionRepo) Activate(ctx context.Context, paymentID string, startAt time.Time) error {
	sql := `
		UPDATE posters_promotions SET started_at = $1, status = 'active' WHERE payment_id = $2;
	`
	_, err := r.pool.Exec(ctx, sql, startAt, paymentID)
	if err != nil {
		return repository.HandelPgErrors(err)
	}

	return nil
}

func (r *PromotionRepo) GetByUserID(ctx context.Context, userID int) ([]dto.UserPromotionDTO, error) {
	sql := `
		SELECT p.poster_id, p.promotion_id,
               p.status, p.ends_at
        FROM posters_promotions p
		WHERE p.user_id = $1 AND status = 'active';
	`
	rows, err := r.pool.Query(ctx, sql, userID)
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}
	proms, err := pgx.CollectRows(rows, pgx.RowToStructByName[dto.UserPromotionDTO])
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}
	return proms, nil
}
