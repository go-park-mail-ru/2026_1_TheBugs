package geo

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

type GeographyPoint struct {
	Lat float64
	Lon float64
}

func (g *GeographyPoint) Scan(value interface{}) error {
	t, ok := value.(string)
	if !ok {
		return fmt.Errorf("unsupported type for GeographyPoint: %T", t)
	}
	s := t
	s = strings.TrimPrefix(s, "POINT(")
	s = strings.TrimSuffix(s, ")")
	parts := strings.Split(s, " ")
	if len(parts) != 2 {
		return fmt.Errorf("invalid geography point: %s", t)
	}
	lon, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return fmt.Errorf("parse lon %s: %w", parts[0], err)
	}
	lat, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return fmt.Errorf("parse lat %s: %w", parts[1], err)
	}
	g.Lon = lon
	g.Lat = lat
	return nil
}

func (g GeographyPoint) Value() (driver.Value, error) {
	return fmt.Sprintf("SRID=4326;POINT(%f %f)", g.Lon, g.Lat), nil
}
