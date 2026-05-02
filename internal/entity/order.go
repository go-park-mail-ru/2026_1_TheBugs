package entity

import "time"

type Order struct {
	ID           int
	CategoryName string
	Status       string
	CreatedAt    time.Time
}

type OrderFull struct {
	ID           int
	UserID       int
	CategoryName string
	Status       string
	Description  string
	CreatedAt    time.Time
	Photos       []OrderPhoto
}

type OrderPhoto struct {
	ImgURL string
	Order  int
}
