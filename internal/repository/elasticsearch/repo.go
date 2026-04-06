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
	var must []interface{}

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
		must = append(must, map[string]interface{}{"multi_match": matchQuery})
	} else {
		must = append(must, map[string]interface{}{"match_all": map[string]interface{}{}})
	}

	var filter []interface{}
	if filters.UtilityCompany != nil {
		termQuery := TermQuery{
			Term: map[string]string{
				"company_alias": *filters.UtilityCompany,
			},
		}
		filter = append(filter, termQuery)
	}

	searchQuery := SearchQuery{
		Size:           filters.Limit,
		From:           filters.Offset,
		TrackTotalHits: true,
		Query: Query{Bool: BoolQuery{
			Must:   must,
			Filter: filter,
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
