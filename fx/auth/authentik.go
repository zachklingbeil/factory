package auth

import (
	"context"
	"encoding/json"
	"fmt"

	"goauthentik.io/api/v3"
)

type Auth struct {
	*api.APIClient // Authentik API client for management
}

// NewAuth creates a new Authentik API client with the given baseURL and apikey.
func NewAuth(baseURL, apiKey string) *Auth {
	cfg := api.NewConfiguration()
	cfg.Host = baseURL
	cfg.Scheme = "https"
	cfg.DefaultHeader = map[string]string{
		"Authorization": "Bearer " + apiKey,
	}
	client := api.NewAPIClient(cfg)
	return &Auth{
		APIClient: client,
	}
}

// TestConnection fetches and prints the current user info to test the connection.
func (a *Auth) WhoAmI() error {
	user, _, err := a.CoreApi.CoreUsersMeRetrieve(context.TODO()).Execute()
	if err != nil {
		return err
	}

	out, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(out))
	return nil
}

// ListUsers lists all users in Authentik.
func (a *Auth) ListUsers() error {
	users, _, err := a.CoreApi.CoreUsersList(context.TODO()).Execute()
	if err != nil {
		return err
	}
	out, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

// ListApplications lists all applications in Authentik.
func (a *Auth) ListApplications() error {
	apps, _, err := a.CoreApi.CoreApplicationsList(context.TODO()).Execute()
	if err != nil {
		return err
	}
	out, err := json.MarshalIndent(apps, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}
