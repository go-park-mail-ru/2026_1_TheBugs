package entity

type Flat struct {
	PropertyID   int    `db:"property_id"`
	FlatCategory string `db:"flat_category"`
	Number       int    `db:"number"`
	Floor        int    `db:"floor"`
	RoomCount    int    `db:"room_count"`
}
