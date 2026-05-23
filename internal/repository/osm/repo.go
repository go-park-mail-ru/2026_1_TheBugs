package osm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/ctxLogger"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/geo"
)

type Node[T any] struct {
	ID   int     `json:"id"`
	Lat  float32 `json:"lat"`
	Lon  float32 `json:"lon"`
	Tags T       `json:"tags"`
}

type OSMResponse[T any] struct {
	Elements []Node[T] `json:"elements"`
}
type MetroTags struct {
	Name  string `json:"name"`
	Color string `json:"lightblue"`
}
type OSMRepo struct {
}

func NewOSMRepo() *OSMRepo {
	return &OSMRepo{}
}

func (r *OSMRepo) GetMetroStationByRadius(ctx context.Context, buildingGeo dto.GeographyDTO, radius entity.Metre) ([]entity.MetroStation, error) {
	op := "GetMetroStationsByRadius"
	log := ctxLogger.GetLogger(ctx).WithField("op", op)

	query := fmt.Sprintf(`[out:json][timeout:25];
		node(around:%d,%f,%f)["station"="subway"];
		out body 1;`, radius, buildingGeo.Lat, buildingGeo.Lon)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://overpass-api.de/api/interpreter",
		strings.NewReader(query))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "DomDeli/1.0.0")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("yandex request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Infof("Yandex response: %d %s", resp.StatusCode, string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("yandex status %d: %s", resp.StatusCode, string(body))
	}
	var result OSMResponse[MetroTags]

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	stations := make([]entity.MetroStation, 0, len(result.Elements))

	for _, el := range result.Elements {
		stations = append(stations, entity.MetroStation{
			ID:          el.ID,
			StationName: el.Tags.Name,
			StationGEO: geo.GeographyPoint{
				Lat: float64(el.Lat),
				Lon: float64(el.Lon),
			},
		},
		)
	}
	log.Info(result)
	log.Info(stations)

	return stations, nil
}
