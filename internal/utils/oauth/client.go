package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2026_1_TheBugs/config"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity/dto"
)

const redirectURI = "https://dom-deli.ru/oauth/vk"
const oAuthURI = "https://id.vk.ru/oauth2/auth"
const publicInfoURI = "https://id.vk.ru/oauth2/public_info"

func ChangeCodeToAccessToken(ctx context.Context, flow dto.OAuthCodeFlow) (*dto.OAuthUserCred, error) {
	reqBody := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {config.Config.OAuth.VKClientID},
		"redirect_uri":  {redirectURI},
		"code_verifier": {flow.CodeVerifier},
		"code":          {flow.Code},
		"device_id":     {flow.DeviceID},
		"state":         {flow.State},
	}
	log.Println(reqBody)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		oAuthURI,
		strings.NewReader(reqBody.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("vk request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("vk status %d: %s", resp.StatusCode, string(body))
	}
	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		IDToken      string `json:"id_token"`
		UserID       int    `json:"user_id"`
		ExpiresIn    int    `json:"expires_in"`
		Scope        string `json:"scope"`
		State        string `json:"state"`
		Error        string `json:"error,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if result.Error != "" {
		return nil, fmt.Errorf("vk error: %s", result.Error)
	}

	cred := &dto.OAuthUserCred{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		IDToken:      result.IDToken,
		UserID:       result.UserID,
		ExpiresIn:    result.ExpiresIn,
		State:        result.State,
	}

	return cred, nil
}

type VKPublicUserInfo struct {
	User struct {
		UserID    string `json:"user_id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Phone     string `json:"phone"`
		Avatar    string `json:"avatar"`
		Email     string `json:"email"`
	} `json:"user"`
}

func GetUserPublicInfo(ctx context.Context, idToken string) (*VKPublicUserInfo, error) {
	reqBody := url.Values{
		"client_id": {config.Config.OAuth.VKClientID},
		"id_token":  {idToken},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		publicInfoURI,
		strings.NewReader(reqBody.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create public_info request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("vk public_info request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("vk public_info status %d: %s", resp.StatusCode, string(body))
	}

	var userInfo VKPublicUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("decode public_info response: %w", err)
	}

	return &userInfo, nil
}
