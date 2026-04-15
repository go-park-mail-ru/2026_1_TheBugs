package entity

import "github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/geo"

type MetroStation struct {
	ID          int                `db:"id"`
	StationName string             `db:"station_name"`
	StationGEO  geo.GeographyPoint `db:"metro_geo"`
}
