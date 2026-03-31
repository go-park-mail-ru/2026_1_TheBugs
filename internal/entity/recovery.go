package entity

type RecoverSession struct {
	Email    string `json:"email"`
	Code     string `json:"code"`
	Attempts int    `json:"attempts"`
	Verified bool   `json:"verified"`
}
