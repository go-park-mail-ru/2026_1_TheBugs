package entity

import (
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/geo"
)

type PosterFlat struct {
	ID           int     `db:"id" json:"id"`
	Price        float64 `db:"price" json:"price"`
	ImgURL       *string `db:"avatar_url" json:"avatar_url"`
	Address      string  `db:"address" json:"address"`
	Metro        *string `db:"station_name" json:"station_name"`
	Area         float64 `db:"area" json:"area"`
	Alias        string  `db:"alias" json:"alias"`
	Floor        *int    `db:"floor" json:"floor"`
	FlatCategory *string `db:"flat_category" json:"flat_category"`
}

type Poster struct {
	ID            int     `db:"id"`
	Alias         string  `db:"alias"`
	Address       string  `db:"address"`
	Area          float64 `db:"area"`
	Price         float64 `db:"price"`
	AvatarURl     *string `db:"avatar_url"`
	CategoryName  string  `db:"category_name"`
	CategoryAlias string  `db:"category_alias"`
}

type AnyPoint struct {
	ID       int64   `db:"id" json:"id"`
	Lat      float64 `db:"lat" json:"lat"`
	Lon      float64 `db:"lon" json:"lon"`
	Count    *int64  `db:"count" json:"count,omitempty"`
	PriceMin *int64  `db:"price_min" json:"priceMin,omitempty"`
	PriceMax *int64  `db:"price_max" json:"priceMax,omitempty"`
	Price    *int64  `db:"price" json:"price,omitempty"`
	Alias    *string `db:"alias" json:"alias,omitempty"`
	Cluster  bool    `db:"cluster" json:"cluster"`
}

type PosterPoint struct {
	ID      int64   `db:"id" json:"id"`
	Lat     float64 `db:"lat" json:"lat"`
	Lon     float64 `db:"lon" json:"lon"`
	Price   int64   `db:"price" json:"price"`
	Alias   string  `db:"alias" json:"alias"`
	Cluster bool    `db:"cluster" json:"cluster"`
}

type ClusterPoint struct {
	ID       int64   `db:"id" json:"id"`
	Lat      float64 `db:"lat" json:"lat"`
	Lon      float64 `db:"lon" json:"lon"`
	Count    int64   `db:"count" json:"count"`
	PriceMin int64   `db:"price_min" json:"priceMin"`
	PriceMax int64   `db:"price_max" json:"priceMax"`
}

type PosterImage struct {
	ImgURL string `db:"img_url"`
	Order  int    `db:"sequence_order"`
}

type PosterById struct {
	ID            int     `db:"id"`
	Alias         string  `db:"alias"`
	Price         float64 `db:"price"`
	Category      string  `db:"category_name"`
	CategoryAlias string  `db:"category_alias"`
	Description   string  `db:"description"`

	Area       float64 `db:"area"`
	PropertyID int     `db:"property_id"`

	Geo        geo.GeographyPoint  `db:"building_geo"`
	Address    string              `db:"address"`
	District   *string             `db:"district"`
	Metro      *string             `db:"station_name"`
	MetroGeo   *geo.GeographyPoint `db:"metro_geo"`
	City       string              `db:"city_name"`
	FloorCount int                 `db:"floor_count"`

	Images     []PosterImage `db:"-"`
	Facilities []Facility    `db:"-"`

	SellerFirstName string  `db:"first_name"`
	SellerLastName  string  `db:"last_name"`
	SellerAvatarURL *string `db:"seller_avatar_url"`
	Phone           string  `db:"phone"`

	CompanyName      *string `db:"company_name"`
	CompanyAvatarURL *string `db:"company_avatar_url"`
	CompanyAlias     *string `db:"company_alias"`
	CompanyID        *int    `db:"company_id"`
}
