package frame

import (
	"context"
	"encoding/json"
	"html/template"
	"io"
	"maps"
	"net/http"
)

type Component struct {
	HTML template.HTML  `json:"html"`
	CSS  template.CSS   `json:"css"`
	Map  map[string]any `json:"map,omitempty"`
	HTTP *http.Client
	CTX  context.Context
}

func NewComponent(html template.HTML, css template.CSS, http *http.Client, ctx context.Context) *Component {
	return &Component{
		HTML: html,
		CSS:  css,
		HTTP: http,
		CTX:  ctx,
		Map:  make(map[string]any),
	}
}

func (c *Component) Source(url string) error {
	req, err := http.NewRequestWithContext(c.CTX, "GET", url, nil)
	if err != nil {
		return err
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	maps.Copy(c.Map, result)
	return nil
}
