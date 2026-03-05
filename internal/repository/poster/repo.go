package poster

import (
	"database/sql"
	"fmt"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PosterRepo struct {
	db *sql.DB
}

func NewPosterRepo(dataSourceName string) (*PosterRepo, error) {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	return &PosterRepo{
		db: db,
	}, nil
}

func (r *PosterRepo) GetPosters(limit, offset int) ([]*entity.Poster, error) {
	rows, err := r.db.Query(
		`SELECT id, price, image_url, address, metro, area, floor, type
		FROM posters
		ORDER BY id
		LIMIT $1 OFFSET $2`,
		limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posters := make([]*entity.Poster, 0, limit)

	for rows.Next() {
		var p entity.Poster
		var img sql.NullString

		err := rows.Scan(&p.Id, &p.Price, &img, &p.Address, &p.Metro, &p.Area, &p.Floor, &p.Type)
		if err != nil {
			return nil, err
		}

		if img.Valid {
			p.ImgURL = &img.String
		} else {
			p.ImgURL = nil
		}
		posters = append(posters, &p)
	}

	return posters, rows.Err()
}

func (r *PosterRepo) CountPosters() int {
	var count int
	r.db.QueryRow("SELECT COUNT(*) FROM posters").Scan(&count)
	return count
}
