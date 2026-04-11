package entity

type Facility struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Alias string `db:"alias"`
}
