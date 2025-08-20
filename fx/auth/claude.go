package auth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/oauth2/clientcredentials"
)

type Client struct {
	oauthClient *http.Client
	baseURL     string
	services    map[string]string
	mu          sync.RWMutex
}

func NewClient(baseURL, clientID, clientSecret string) *Client {
	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     baseURL + "/application/o/token/",
	}

	return &Client{
		oauthClient: config.Client(context.Background()),
		baseURL:     baseURL,
		services:    make(map[string]string),
	}
}

func (c *Client) Initialize(ctx context.Context) error {
	var response struct {
		Results []struct {
			Slug string `json:"slug"`
		} `json:"results"`
	}

	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/v3/providers/oauth2/", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.oauthClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	for _, provider := range response.Results {
		c.services[provider.Slug] = provider.Slug
	}
	return nil
}

func (c *Client) Get(ctx context.Context, provider, endpoint string) (any, error) {
	c.mu.RLock()
	baseURL := c.services[provider]
	c.mu.RUnlock()

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.oauthClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result any
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}

func (c *Client) Post(ctx context.Context, provider, endpoint string, body any) (any, error) {
	c.mu.RLock()
	baseURL := c.services[provider]
	c.mu.RUnlock()

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = strings.NewReader(string(jsonBody))
	}

	req, err := http.NewRequestWithContext(ctx, "POST", baseURL+endpoint, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.oauthClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result any
	err = json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}

func (c *Client) ListProviders() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var providers []string
	for provider := range c.services {
		providers = append(providers, provider)
	}
	return providers
}
