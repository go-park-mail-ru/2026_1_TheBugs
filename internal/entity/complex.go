package entity

import "github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/geo"

type UtilityCompany struct {
	ID          int                `db:"id"`
	Phone       string             `db:"phone"`
	GEO         geo.GeographyPoint `db:"geo"`
	Address     string             `db:"address"`
	Alias       string             `db:"alias"`
	AvatarURL   *string            `db:"avatar_url"`
	CompanyName string             `db:"company_name"`
	Description string             `db:"description"`
}

type UtilityCompanyCard struct {
	ID          int     `db:"id"`
	Alias       string  `db:"alias"`
	AvatarURL   *string `db:"avatar_url"`
	CompanyName string  `db:"company_name"`
}

type UtilityCompanyPhoto struct {
	ID               *int    `db:"id"`
	UtilityCompanyID *int    `db:"utility_company_id"`
	ImgURL           *string `db:"img_url"`
	Order            *int    `db:"sequence_order"`
}

type Developer struct {
	ID            int     `db:"id"`
	AvatarURL     *string `db:"avatar_url"`
	DeveloperName string  `db:"developer_name"`
}
