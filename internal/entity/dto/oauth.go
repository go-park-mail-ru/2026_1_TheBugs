package dto

type OAuthUserCred struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	UserID       int    `json:"user_id"`
	State        string `json:"state"`
	Scope        string `json:"scope"`
}

type OAuthCodeFlow struct {
	Code         string `json:"code"`
	DeviceID     string `json:"device_id"`
	State        string `json:"state"`
	CodeVerifier string `json:"code_verifier"`
}
