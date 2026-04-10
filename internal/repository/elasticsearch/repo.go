package elasticsearch

import (
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

func NewESRepo(client *es.Client) *ESRepo {
	return &ESRepo{client: client}
}

func (r *ESRepo) SearchPosters(ctx context.Context, filters dto.PostersFiltersDTO) (*dto.PostersResponse, error) {
	var should []any
	var notMust []any

	if filters.SearchQuery != nil {
		matchQuery := MultiMatchQuery{
			Query:     *filters.SearchQuery,
			Type:      "best_fields",
			Fuzziness: "AUTO",
			Fields: []string{
				"city^100",
				"description^5",
				"address^20",
				"flat_category^10",
				"station_name^10",
				"district^10",
				"company_name^15",
			},
		}
		should = append(should, map[string]any{"multi_match": matchQuery})
		should = append(should, map[string]interface{}{
			"nested": map[string]interface{}{
				"path": "facilities",
				"query": map[string]interface{}{
					"match": map[string]interface{}{
						"facilities.name": map[string]interface{}{
							"query": *filters.SearchQuery,
							"boost": float64(10),
						},
					},
				},
			},
		})
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

	searchQuery := SearchQuery{
		Size:           filters.Limit,
		From:           filters.Offset,
		TrackTotalHits: true,
		Query: Query{Bool: BoolQuery{
			Should:  should,
			Filter:  filter,
			MustNot: notMust,
		}},
	}

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
