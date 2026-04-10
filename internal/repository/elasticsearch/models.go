package elasticsearch

type ESResponse[T any] struct {
	Res struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		Hits []struct {
			ID     string `json:"_id"`
			Source T      `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}
type Query struct {
	Bool BoolQuery `json:"bool"`
}
type SearchQuery struct {
	Size           int   `json:"size"`
	From           int   `json:"from"`
	TrackTotalHits any   `json:"track_total_hits"`
	Query          Query `json:"query"`
}

type BoolQuery struct {
	Must    []any `json:"must"`
	Filter  []any `json:"filter,omitempty"`
	MustNot []any `json:"must_not,omitempty"`
}

type MultiMatchQuery struct {
	Query     string   `json:"query"`
	Type      string   `json:"type"`
	Fuzziness string   `json:"fuzziness"`
	Fields    []string `json:"fields"`
}

type TermQuery struct {
	Term map[string]any `json:"term"`
}

type TermsQuery struct {
	Terms map[string]any `json:"terms"`
}

type ScriptQuery struct {
	Script map[string]string `json:"script"`
}

type RangeQuery struct {
	Range map[string]map[string]int `json:"range"`
}
