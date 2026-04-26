package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	es "github.com/elastic/go-elasticsearch/v9"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
)

func ApplyRangeFilter[T []any](filters T, field string, max *int, min *int) T {
	if max != nil || min != nil {
		rangeFilter := make(map[string]int)
		if max != nil {
			rangeFilter["lte"] = *max
		}
		if min != nil {
			rangeFilter["gte"] = *min
		}

		filters = append(filters, RangeQuery{Range: map[string]map[string]int{
			field: rangeFilter,
		}})
	}
	return filters
}

type ESRepo struct {
	client *es.Client
}

func ParseInSliceStruct[T any](body io.Reader) ([]T, int, error) {
	raw, err := io.ReadAll(body)
	if err != nil {
		return nil, 0, fmt.Errorf("read body: %w", err)
	}

	var resp ESResponse[T]
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, 0, fmt.Errorf("unmarshal: %w", err)
	}

	res := make([]T, 0, resp.Res.Total.Value)
	for _, hit := range resp.Res.Hits {
		res = append(res, hit.Source)
	}

	return res, resp.Res.Total.Value, nil
}

func ApplyPostersFilters(filters dto.PostersFiltersDTO) SearchQuery {
	var should []any
	var notMust []any

	if filters.SearchQuery != nil {
		matchQuery := MultiMatchQuery{
			Query: *filters.SearchQuery,
			Type:  "best_fields",
			//Fuzziness: "AUTO",
			Fields: []string{
				"city^100",
				// "description^5",
				"address^20",
				// "flat_category^10",
				"station_name^10",
				"district^10",
				"company_name^15",
				//"facilities.name^10",
			},
		}
		should = append(should, map[string]any{"multi_match": matchQuery})
	} else {
		should = append(should, map[string]any{"match_all": map[string]any{}})
	}

	var filter []any
	if filters.UtilityCompany != nil {
		termQuery := TermQuery{
			Term: map[string]any{
				"company_alias": *filters.UtilityCompany,
			},
		}
		filter = append(filter, termQuery)
	}

	if filters.Category != nil {
		termQuery := TermQuery{
			Term: map[string]any{
				"category_alias": *filters.Category,
			},
		}
		filter = append(filter, termQuery)
	}
	if filters.RoomCount != nil {
		termQuery := TermQuery{
			Term: map[string]any{
				"room_count": *filters.RoomCount,
			},
		}
		filter = append(filter, termQuery)
	}

	if len(filters.Facilities) > 0 {
		termQuery := TermsQuery{
			Terms: map[string]any{
				"facilities.alias": filters.Facilities,
			},
		}
		filter = append(filter, termQuery)
	}

	if filters.IsNotFirstFloor {
		notMust = append(notMust, map[string]any{
			"term": map[string]any{"floor": 1},
		})
	}
	if filters.IsNotLastFloor {
		notMust = append(notMust, map[string]any{
			"script": map[string]any{
				"script": "doc['floor'].size() > 0 && doc['building_floor'].size() > 0 && doc['floor'].value == doc['building_floor'].value",
			},
		})
	}

	filter = ApplyRangeFilter(filter, "price", filters.MaxPrice, filters.MinPrice)
	filter = ApplyRangeFilter(filter, "area", filters.MaxSquare, filters.MinSquare)
	filter = ApplyRangeFilter(filter, "floor", filters.MaxFlatFloor, filters.MinFlatFloor)
	filter = ApplyRangeFilter(filter, "building_floor", filters.MaxBuildingFloor, filters.MinBuildingFloor)

	return SearchQuery{
		Size:           filters.Limit,
		From:           filters.Offset,
		TrackTotalHits: true,
		Query: Query{Bool: BoolQuery{
			Should:  should,
			Filter:  filter,
			MustNot: notMust,
		}},
	}
}

func NewESRepo(client *es.Client) *ESRepo {
	return &ESRepo{client: client}
}

func (r *ESRepo) SearchPosters(ctx context.Context, filters dto.PostersFiltersDTO) (*dto.PostersResponse, error) {
	searchQuery := ApplyPostersFilters(filters)

	queryBody, err := json.Marshal(searchQuery)
	if err != nil {
		return nil, fmt.Errorf("marshal query: %w", err)
	}

	fmt.Printf("ES Query: %s\n", string(queryBody))

	res, err := r.client.Search(
		r.client.Search.WithIndex("_all"),
		r.client.Search.WithBody(strings.NewReader(string(queryBody))),
	)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}
	defer res.Body.Close()

	posters, total, err := ParseInSliceStruct[entity.PosterFlat](res.Body)
	if err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	return &dto.PostersResponse{
		Posters: dto.PostersToPostersDTO(posters),
		Len:     total,
	}, nil
}

func (r *ESRepo) GetClustersByMapBounds(ctx context.Context, coords dto.MapBounds, filters dto.PostersFiltersDTO) ([]entity.ClusterPoint, error) {
	precision := (float64(coords.Zoom) * 0.6)

	// query := map[string]interface{}{
	// 	"size": 0,
	// 	"query": map[string]interface{}{
	// 		"bool": map[string]interface{}{
	// 			"filter": []interface{}{
	// 				map[string]interface{}{
	// 					"geo_bounding_box": map[string]interface{}{
	// 						"geo": map[string]interface{}{
	// 							"top_left": map[string]interface{}{
	// 								"lat": coords.BBox.NorthEast.Lat,
	// 								"lon": coords.BBox.SouthWest.Lon,
	// 							},
	// 							"bottom_right": map[string]interface{}{
	// 								"lat": coords.BBox.SouthWest.Lat,
	// 								"lon": coords.BBox.NorthEast.Lon,
	// 							},
	// 						},
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// 	"aggs": map[string]interface{}{
	// 		"by_grid": map[string]interface{}{
	// 			"geohash_grid": map[string]interface{}{
	// 				"field":     "geo",
	// 				"size":      50,
	// 				"precision": precision,
	// 			},
	// 			"aggs": map[string]interface{}{
	// 				"centroid": map[string]interface{}{"geo_centroid": map[string]string{"field": "geo"}},
	// 			},
	// 		},
	// 	},
	// }

	searchQuery := ApplyPostersFilters(filters)
	searchQuery.TrackTotalHits = false

	searchQuery.Sourse = map[string]any{
		"includes": []string{"buckets"},
	}

	searchQuery.Query.Bool.Filter = append(searchQuery.Query.Bool.Filter, map[string]any{
		"geo_bounding_box": map[string]any{
			"geo": map[string]any{
				"top_left": map[string]any{
					"lat": coords.BBox.NorthEast.Lat,
					"lon": coords.BBox.SouthWest.Lon,
				},
				"bottom_right": map[string]any{
					"lat": coords.BBox.SouthWest.Lat,
					"lon": coords.BBox.NorthEast.Lon,
				},
			},
		},
	})

	searchQuery.Aggs = map[string]any{
		"by_grid": map[string]any{
			"geohash_grid": map[string]any{
				"field":     "geo",
				"size":      50,
				"precision": precision,
			},
			"aggs": map[string]any{
				"centroid": map[string]any{"geo_centroid": map[string]string{"field": "geo"}},
			},
		},
	}

	queryBody, err := json.Marshal(searchQuery)
	if err != nil {
		return nil, fmt.Errorf("marshal query: %w", err)
	}
	fmt.Println(string(queryBody))

	res, err := r.client.Search(
		r.client.Search.WithIndex("posters"),
		r.client.Search.WithContext(ctx),
		r.client.Search.WithBody(strings.NewReader(string(queryBody))),
	)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}
	defer res.Body.Close()

	var esResp struct {
		Aggregations json.RawMessage `json:"aggregations"`
	}
	if err := json.NewDecoder(res.Body).Decode(&esResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return parseSimpleClusters(esResp.Aggregations)
}

func parseSimpleClusters(aggs json.RawMessage) ([]entity.ClusterPoint, error) {
	var aggsResp struct {
		ByGrid struct {
			Buckets []struct {
				Key      string `json:"key"`
				DocCount int64  `json:"doc_count"`
				Centroid struct {
					Location struct {
						Lat float64 `json:"lat"`
						Lon float64 `json:"lon"`
					} `json:"location"`
				} `json:"centroid"`
			} `json:"buckets"`
		} `json:"by_grid"`
	}

	if err := json.Unmarshal(aggs, &aggsResp); err != nil {
		return nil, fmt.Errorf("unmarshal aggs: %w", err)
	}

	clusters := make([]entity.ClusterPoint, 0, len(aggsResp.ByGrid.Buckets))

	for i, bucket := range aggsResp.ByGrid.Buckets {
		clusters = append(clusters, entity.ClusterPoint{
			ID:    int64(i + 1),
			Lat:   bucket.Centroid.Location.Lat,
			Lon:   bucket.Centroid.Location.Lon,
			Count: int64(bucket.DocCount),
		})
	}

	return clusters, nil
}

func (r *ESRepo) GetPostersByMapBounds(ctx context.Context, coords dto.MapBounds, filters dto.PostersFiltersDTO) ([]entity.AnyPoint, error) {

	searchQuery := ApplyPostersFilters(filters)
	searchQuery.TrackTotalHits = false

	searchQuery.Query.Bool.Filter = append(searchQuery.Query.Bool.Filter, map[string]any{
		"geo_bounding_box": map[string]any{
			"geo": map[string]any{
				"top_left": map[string]any{
					"lat": coords.BBox.NorthEast.Lat,
					"lon": coords.BBox.SouthWest.Lon,
				},
				"bottom_right": map[string]any{
					"lat": coords.BBox.SouthWest.Lat,
					"lon": coords.BBox.NorthEast.Lon,
				},
			},
		},
	})
	searchQuery.Aggs = map[string]any{
		"clusters": map[string]any{
			"geohash_grid": map[string]any{
				"field":     "geo",
				"size":      50,
				"precision": 9,
			},
			"aggs": map[string]any{
				"centroid": map[string]any{
					"geo_centroid": map[string]any{
						"field": "geo",
					},
				},
				"price_min": map[string]any{
					"min": map[string]any{"field": "price"},
				},
				"top_1": map[string]any{
					"top_hits": map[string]any{
						"_source": map[string]any{
							"includes": []string{"id", "price", "alias"},
						},
						"size": 1,
					},
				},
			},
		},
	}
	queryBody, err := json.Marshal(searchQuery)
	if err != nil {
		return nil, fmt.Errorf("marshal query: %w", err)
	}
	fmt.Println(queryBody)

	var resp struct {
		Aggregations struct {
			Clusters struct {
				Buckets []struct {
					DocCount int64 `json:"doc_count"`
					Centroid struct {
						Location struct{ Lat, Lon float64 } `json:"location"`
					} `json:"centroid"`
					PriceMin struct {
						Value float64 `json:"value"`
					} `json:"price_min"`
					TopHit struct {
						Hits struct {
							Hits []struct {
								Source struct {
									ID    int64   `json:"id"`
									Price float64 `json:"price"`
									Alias string  `json:"alias"`
								} `json:"_source"`
							} `json:"hits"`
						} `json:"hits"`
					} `json:"top_1"`
				} `json:"buckets"`
			} `json:"clusters"`
		} `json:"aggregations"`
	}

	res, err := r.client.Search(
		r.client.Search.WithIndex("posters"),
		r.client.Search.WithContext(ctx),
		r.client.Search.WithBody(bytes.NewReader(queryBody)),
	)
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	result := make([]entity.AnyPoint, 0, len(resp.Aggregations.Clusters.Buckets))
	for i, b := range resp.Aggregations.Clusters.Buckets {
		p := entity.AnyPoint{
			ID:       int64(i + 1),
			Lat:      b.Centroid.Location.Lat,
			Lon:      b.Centroid.Location.Lon,
			Count:    &b.DocCount,
			PriceMin: &b.PriceMin.Value,
			Group:    b.DocCount > 1,
		}
		if b.DocCount == 1 && len(b.TopHit.Hits.Hits) > 0 {
			p.Price = &b.TopHit.Hits.Hits[0].Source.Price
			p.Alias = &b.TopHit.Hits.Hits[0].Source.Alias
		}
		result = append(result, p)
	}

	return result, nil
}

func (r *ESRepo) DeletePoster(ctx context.Context, posterID int) error {
	var filter []any
	termQuery := TermQuery{
		Term: map[string]any{
			"id": posterID,
		},
	}
	filter = append(filter, termQuery)
	searchQuery := SearchQuery{
		Size:           1,
		TrackTotalHits: true,
		Query: Query{Bool: BoolQuery{
			Filter: filter,
		}},
	}
	queryBody, err := json.Marshal(searchQuery)
	if err != nil {
		return fmt.Errorf("marshal query: %w", err)
	}
	fmt.Println(queryBody)

	res, err := r.client.DeleteByQuery(
		[]string{"posters"},
		bytes.NewReader(queryBody),
		r.client.DeleteByQuery.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("delete by query request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		b, _ := io.ReadAll(res.Body)
		return fmt.Errorf("delete by query returned %s: %s", res.Status(), string(b))
	}

	return nil

}
