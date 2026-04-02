package entity

type City struct {
	ID   int    `db:"id"`
	Name string `db:"city_name"`
}
