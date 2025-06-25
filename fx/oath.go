package fx

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Oauth struct {
	ClientID     string
	ClientSecret string
	TokenURL     string
	HTTPClient   *http.Client
	Ctx          context.Context

	accessToken string
	expiresAt   time.Time
	mu          sync.Mutex
}

func NewOauth(ctx context.Context, client *http.Client, clientID, clientSecret, tokenURL string) *Oauth {
	return &Oauth{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
		HTTPClient:   client,
		Ctx:          ctx,
	}
}

func (o *Oauth) GetAccessToken() (string, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if time.Now().Before(o.expiresAt) && o.accessToken != "" {
		return o.accessToken, nil
	}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	req, err := http.NewRequestWithContext(o.Ctx, "POST", o.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(o.ClientID, o.ClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := o.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token request failed: %s", resp.Status)
	}

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	o.accessToken = result.AccessToken
	o.expiresAt = time.Now().Add(time.Duration(result.ExpiresIn-60) * time.Second) // refresh 1 min early

	return o.accessToken, nil
}
