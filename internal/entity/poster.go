package entity

type Poster struct {
	Id      int      `db:"id"`
	Price   float64  `db:"price"`
	ImgURL  *string  `db:"image_url"`
	Address string   `db:"address"`
	Metro   *string  `db:"metro"`
	Area    float64  `db:"area"`
	Floor   int      `db:"floor"`
	Beds    *int     `db:"beds"`
	Rating  *float64 `db:"rating"`
	Type    string   `db:"type"`
}
