package fx

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2/clientcredentials"
)

// NewOAuthClient returns an authenticated HTTP client (machine-to-machine, no user interaction)
func NewOAuthClient(ctx context.Context, clientID, clientSecret, tokenURL string, scopes []string) (*http.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	clientConfig := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
		Scopes:       scopes,
	}

	// Get token and create HTTP client
	client := clientConfig.Client(ctx)
	if client == nil {
		return nil, fmt.Errorf("failed to create OAuth client")
	}

	// Test the client by making a token request to validate credentials
	token, err := clientConfig.Token(ctx)
	if err != nil {
		return nil, fmt.Errorf("OAuth client credentials flow failed: %w", err)
	}
	fmt.Printf("✓ OAuth client authenticated successfully (token expires: %v)\n", token.Expiry)
	return client, nil
}
