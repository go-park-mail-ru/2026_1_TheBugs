package dto

type ChatResult struct {
	Answer      string   `json:"answer"`
	Status      string   `json:"status"`
	MissingInfo []string `json:"missing_info,omitempty"`
	Reason      string   `json:"reason,omitempty"`
}
