package fx

import (
	"fmt"
)

type Flows interface {
	AuthHeader(token string) map[string]string
	Authorization() string
	Authentication() string
	Enrollment() string
	Recovery() string
	Settings() string
	Sessions() string
	Logout() string
}

func NewAuthentikFlows(baseURL, token string) Flows {
	return &flows{baseURL: baseURL, token: token}
}

type flows struct {
	baseURL string
	token   string
}

// AuthHeader returns the Authorization header as a map for use in HTTP requests.
func (f *flows) AuthHeader(token string) map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + token,
	}
}

// Authentik implements the flows interface.
func (f *flows) Authorization() string {
	return fmt.Sprintf("%s/application/o/authorize/", f.baseURL)
}
func (f *flows) Authentication() string {
	return fmt.Sprintf("%s/application/o/flow/authentication/", f.baseURL)
}
func (f *flows) Enrollment() string {
	return fmt.Sprintf("%s/application/o/flow/enrollment/", f.baseURL)
}
func (f *flows) Recovery() string {
	return fmt.Sprintf("%s/application/o/flow/recovery/", f.baseURL)
}
func (f *flows) Settings() string {
	return fmt.Sprintf("%s/application/o/flow/settings/", f.baseURL)
}
func (f *flows) Sessions() string {
	return fmt.Sprintf("%s/application/o/flow/sessions/", f.baseURL)
}
func (f *flows) Logout() string {
	return fmt.Sprintf("%s/application/o/flow/logout/", f.baseURL)
}
