package spotify

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	spotifyAuthBaseUrl = "https://accounts.spotify.com"
)

// AuthorizationCodeGrantResponse represents the
// spotify response for obtaining an access / refresh token.
type AuthorizationCodeGrantResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// RefreshTokenResponse represents the spotify response
// for refreshing an access token.
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// AuthClient represent a spotify authentication client
// to authenticate the user using oauth.
type AuthClient interface {
	// BuildAuthUrl builds the oauth url with the
	// clientId, the redirect uri and the scopes.
	BuildAuthUrl() string

	//GetTokenPair requests a new token pair.
	GetTokenPair(code string) (*AuthorizationCodeGrantResponse, error)

	// RequestRefreshedToken refreshes an access token
	// using a refresh token.
	RequestRefreshedToken(refreshToken string) (string, error)
}

type authClient struct {
	clientId     string
	clientSecret string
	redirectUri  string
	scopes       []string

	httpClient *http.Client
}

// NewAuthClient creates a new spotify authentication
// client with a given configuration.
func NewAuthClient(
	clientId string,
	clientSecret string,
	redirectUri string,
	scopes []string,
) AuthClient {
	return &authClient{
		clientId:     clientId,
		clientSecret: clientSecret,
		redirectUri:  redirectUri,
		scopes:       scopes,

		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *authClient) BuildAuthUrl() string {
	q := url.Values{}
	q.Add("response_type", "code")
	q.Add("client_id", c.clientId)
	q.Add("redirect_uri", c.redirectUri)

	if len(c.scopes) > 0 {
		q.Add("scope", strings.Join(c.scopes, ","))
	}

	return fmt.Sprintf("%s/authorize?%s", spotifyAuthBaseUrl, q.Encode())
}

func (c *authClient) GetTokenPair(code string) (*AuthorizationCodeGrantResponse, error) {
	data := url.Values{}
	data.Add("grant_type", "authorization_code")
	data.Add("code", code)
	data.Add("client_id", c.clientId)
	data.Add("client_secret", c.clientSecret)
	data.Add("redirect_uri", c.redirectUri)

	url := fmt.Sprintf("%s/api/token", spotifyAuthBaseUrl)

	req, err := http.NewRequest(
		http.MethodPost,
		url,
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		fmt.Println(res.StatusCode)
		return nil, errors.New("request failed")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var tokenResponse AuthorizationCodeGrantResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return nil, err
	}

	return &tokenResponse, nil
}

func (c *authClient) RequestRefreshedToken(refreshToken string) (string, error) {
	data := url.Values{}
	data.Add("grant_type", "refresh_token")
	data.Add("refresh_token", refreshToken)
	data.Add("client_id", c.clientId)
	data.Add("client_secret", c.clientSecret)

	url := fmt.Sprintf("%s/api/token", spotifyAuthBaseUrl)

	req, err := http.NewRequest(
		http.MethodPost,
		url,
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		fmt.Println(res.StatusCode)
		return "", errors.New("request failed")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var tokenResponse RefreshTokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return "", err
	}

	return tokenResponse.AccessToken, nil
}
