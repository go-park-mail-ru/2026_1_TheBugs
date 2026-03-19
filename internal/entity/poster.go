package entity

import "github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/geo"

type Poster struct {
	ID      int     `db:"id"`
	Price   float64 `db:"price"`
	ImgURL  *string `db:"avatar_url"`
	Address string  `db:"address"`
	Metro   *string `db:"station_name"`
	Area    float64 `db:"area"`
	Floor   int     `db:"floor"`
}

type PosterImage struct {
	ImgURL string `db:"img_url"`
	Order  int    `db:"sequence_order"`
}

type PosterById struct {
	ID          int     `db:"id"`
	Alias       string  `db:"alias"`
	Price       float64 `db:"price"`
	Category    string  `db:"category"`
	Description string  `db:"description"`

	Area       float64 `db:"area"`
	PropertyID int     `db:"property_id"`

	Geo        geo.GeographyPoint  `db:"building_geo"`
	Address    string              `db:"address"`
	District   *string             `db:"district"`
	Metro      *string             `db:"station_name"`
	MetroGeo   *geo.GeographyPoint `db:"metro_geo"`
	City       string              `db:"city_name"`
	FloorCount int                 `db:"floor_count"`

	Images []PosterImage `db:"-"`

	SellerFirstName string `db:"first_name"`
	SellerLastName  string `db:"last_name"`
	Phone           string `db:"phone"`

	CompanyName *string `db:"company_name"`
	LogoURL     *string `db:"company_avatar_url"`
}

type Flat struct {
	PropertyID   int    `db:"property_id"`
	FlatCategory string `db:"flat_category"`
	Number       int    `db:"number"`
	Floor        int    `db:"floor"`
}
