package geo

import (
	"fmt"
	"strconv"
	"strings"
)

type GeographyPoint struct {
	Lat float64
	Lon float64
}

func (g *GeographyPoint) Scan(value interface{}) error {
	switch t := value.(type) {
	case string:
		s := t // Work on copy to avoid bugs
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
		g.Lon = lon // Note: Lon first in POINT(lon lat)
		g.Lat = lat
		return nil
	default:
		return fmt.Errorf("unsupported type for GeographyPoint: %T", t)
	}
}

func (g GeographyPoint) Value() (string, error) {
	return fmt.Sprintf("ST_GeogFromText('SRID=4326;POINT(%f %f)')", g.Lon, g.Lat), nil
}
