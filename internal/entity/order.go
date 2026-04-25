package entity

import "time"

type Order struct {
	ID           int
	CategoryName string
	Status       string
	CreatedAt    time.Time
}
