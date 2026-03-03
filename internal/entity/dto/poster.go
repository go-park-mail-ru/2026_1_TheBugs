package dto

type PosterDTO struct {
	Price   float64 `json:"price"`
	ImgURL  *string `json:"img_url"`
	Address string  `json:"address"`
	Metro   string  `json:"metro"`
	Area    float64 `json:"area"`
	Floor   int     `json:"floor"`
	Type    string  `json:"type"`
}
