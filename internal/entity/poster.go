package entity

type Poster struct {
	Id      int     `db:"id"`
	Price   float64 `db:"price"`
	ImgURL  *string `db:"avatar_url"`
	Address string  `db:"address"`
	Metro   *string `db:"station_name"`
	Area    float64 `db:"area"`
	Floor   int     `db:"floor"`
}
