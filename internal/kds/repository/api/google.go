package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/leftovers-2025/kds_backend/internal/kds/datasource"
	"github.com/leftovers-2025/kds_backend/internal/kds/entity"
	"github.com/leftovers-2025/kds_backend/internal/kds/port"
)

const (
	GOOGLE_OAUTH_TOKEN_URL = "https://oauth2.googleapis.com/token"
	GOOGLE_USER_API_URL    = "https://www.googleapis.com/oauth2/v3/userinfo"
)

type ApiGoogleRepository struct {
	client *datasource.GoogleApiClient
}

func NewApiGoogleRepository(client *datasource.GoogleApiClient) port.GoogleRepository {
	if client == nil {
		panic("nil GoogleApiClient")
	}
	return &ApiGoogleRepository{
		client: client,
	}
}

type GoogleOauthResponse struct {
	AccessToken           string `json:"access_token"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn string `json:"refresh_token_expires_in"`
	Scope                 string `json:"scope"`
	TokenType             string `json:"token_type"`
	IdToken               string `json:"id_token"`
}

type GoogleUserResponse struct {
	Sub   string `json:"sub"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Oauthのコードの認証を行いユーザーを取得する
func (r *ApiGoogleRepository) CodeAuthorization(code string) (*entity.GoogleUser, error) {
	// リクエスト作成
	req, err := r.createOauthRequest(code)
	if err != nil {
		return nil, err
	}
	// クライアント用意
	client := &http.Client{}

	// Oauthリクエスト送信
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("oauth request failed")
	}

	// ボディを読む
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// JSON変換
	oauthResponse := GoogleOauthResponse{}
	err = json.Unmarshal(body, &oauthResponse)
	if err != nil {
		return nil, err
	}

	// ユーザー情報取得リクエスト作成
	req, err = r.createUserRequest(oauthResponse.AccessToken)
	if err != nil {
		return nil, err
	}

	// ユーザー情報取得リクエスト送信
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("user infomation request failed")
	}

	// ボディを読む
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// JSON変換
	userResponse := GoogleUserResponse{}
	err = json.Unmarshal(body, &userResponse)
	if err != nil {
		return nil, err
	}

	// Googleユーザー作成
	googleUser, err := entity.NewGoogleUser(userResponse.Sub, userResponse.Name, userResponse.Email)
	if err != nil {
		return nil, err
	}
	return googleUser, nil
}

// ユーザー情報を取得するリクエストを作成
func (r *ApiGoogleRepository) createUserRequest(accessToken string) (*http.Request, error) {
	req, err := http.NewRequest(
		"POST",
		GOOGLE_USER_API_URL,
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	return req, nil
}

// Oauthの認証を完了するリクエストを作成
func (r *ApiGoogleRepository) createOauthRequest(code string) (*http.Request, error) {
	// フォームデータ
	values := url.Values{}
	values.Set("code", code)
	values.Set("client_id", r.client.ClientId())
	values.Set("client_secret", r.client.ClientSecret())
	values.Set("redirect_uri", r.client.RedirectURI())
	values.Set("grant_type", "authorization_code")

	// リクエスト
	req, err := http.NewRequest(
		"POST",
		GOOGLE_OAUTH_TOKEN_URL,
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}
