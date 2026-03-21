package complex

import (
	"context"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/repository"
	"github.com/jackc/pgx/v5"
)

type UtilityCompanyPgRepo struct {
	pool repository.DB
}

func NewUtilityCompanyRepo(pool repository.DB) *UtilityCompanyPgRepo {
	return &UtilityCompanyPgRepo{pool: pool}
}

func (r UtilityCompanyPgRepo) GetUtilityCompanyByID(ctx context.Context, id int) (*entity.UtilityCompany, error) {

	sql := `SELECT id, phone, ST_AsText(geo), address, avatar_url FROM utility_complex WHERE id=$1`

	row, err := r.pool.Query(ctx, sql, id)

	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	defer row.Close()

	complex, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[entity.UtilityCompany])
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}
	return &complex, nil

}

func (r UtilityCompanyPgRepo) GetUtilityCompanyByAlias(ctx context.Context, alias string) (*dto.UtilityCompanyDTO, error) {
	sql := `
		SELECT uc.id, uc.phone, uc.company_name, uc.description, ST_AsText(uc.geo) AS geo, uc.address, uc.avatar_url, uc.alias,
		       up.id as photo_id, up.utility_company_id, up.img_url, up.sequence_order, d.id, d.developer_name, d.avatar_url
		FROM utility_companies uc
		LEFT JOIN developers d ON d.id = uc.developer_id
		LEFT JOIN utility_companies_photos up ON uc.id = up.utility_company_id
		WHERE uc.alias = $1
	`

	row, err := r.pool.Query(ctx, sql, alias)
	if err != nil {
		return nil, repository.HandelPgErrors(err)
	}

	defer row.Close()
	var u entity.UtilityCompany
	var photos []entity.UtilityCompanyPhoto
	var d entity.Developer
	found := false
	for row.Next() {
		found = true
		var p entity.UtilityCompanyPhoto
		err := row.Scan(&u.ID, &u.Phone, &u.CompanyName, &u.Description, &u.GEO, &u.Address, &u.AvatarURL, &u.Alias, &p.ID, &p.UtilityCompanyID, &p.ImgURL, &p.Order, &d.ID, &d.DeveloperName, &d.AvatarURL)
		if err != nil {
			return nil, repository.HandelPgErrors(err)
		}
		if p.ID != nil {
			photos = append(photos, p)
		}
	}
	if !found {
		return nil, repository.HandelPgErrors(pgx.ErrNoRows)
	}
	return dto.ToUtilityCompanyDTO(&u, photos, &d), nil

}
