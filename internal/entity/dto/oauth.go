package dto

type OAuthCodeFlow struct {
	Code         string  `json:"code"`
	DeviceID     *string `json:"device_id"`
	State        *string `json:"state"`
	CodeVerifier *string `json:"code_verifier"`
}
