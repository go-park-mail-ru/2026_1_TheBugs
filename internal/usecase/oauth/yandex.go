package oauth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/ctxLogger"
)

const oAuthYandexURI = "https://oauth.yandex.ru/token"
const yandexUserInfoURI = "https://login.yandex.ru/info?format=json"

type YandexCred struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type YandexPublicUserInfo struct {
	ID           string `json:"id"`
	Login        string `json:"login"`
	ClientID     string `json:"client_id"`
	Psuid        string `json:"psuid"`
	DefaultEmail string `json:"default_email"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	DisplayName  string `json:"display_name"`
	RealName     string `json:"real_name"`
	DefaultPhone struct {
		ID     int    `json:"id"`
		Number string `json:"number"`
	} `json:"default_phone"`
}

func ChangeYandexCodeToAccessToken(ctx context.Context, flow dto.OAuthCodeFlow) (*YandexCred, error) {
	op := "ChangeYandexCodeToAccessToken"
	log := ctxLogger.GetLogger(ctx).WithField("op", op)

	creds := config.Config.OAuth.YandexClientID + ":" + config.Config.OAuth.YandexClientSecret
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(creds))

	reqBody := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {config.Config.OAuth.YandexClientID},
		"redirect_uri":  {config.Config.OAuth.YandexRedirectURI},
		"code_verifier": {*flow.CodeVerifier},
		"code":          {flow.Code},
		"state":         {*flow.State},
	}

	log.Infof("Yandex token request body: %s", reqBody.Encode())
	log.Infof("Yandex Basic Auth: %s", basicAuth[:20]+"...")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		oAuthYandexURI,
		strings.NewReader(reqBody.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", basicAuth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
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

	var result struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token,omitempty"`
		Scope        string `json:"scope"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	cred := &YandexCred{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		Scope:        result.Scope,
		TokenType:    result.TokenType,
		ExpiresIn:    result.ExpiresIn,
	}

	return cred, nil
}

func GetYandexUserPublicInfo(ctx context.Context, accessToken string) (*YandexPublicUserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		yandexUserInfoURI,
		nil)
	if err != nil {
		return nil, fmt.Errorf("create yandex user info request: %w", err)
	}

	req.Header.Set("Authorization", "OAuth  "+accessToken)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("yandex user info request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("yandex user info status %d: %s", resp.StatusCode, string(body))
	}

	var userInfo YandexPublicUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("decode public_info response: %w", err)
	}

	return &userInfo, nil

}
