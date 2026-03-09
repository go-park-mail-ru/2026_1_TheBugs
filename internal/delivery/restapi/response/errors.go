package response

type ErrorResponse struct {
	Error   string  `json:"error"`
	Details *string `json:"details,omitempty"`
}

type ValidationErrorResponse struct {
	Error string `json:"error"`
	Field string `json:"field"`
}
