package io

import (
	"fmt"
	"net/http"
)

// Header option builders
func WithHeader(key, value string) GetOption {
	return func(req *http.Request) {
		if value != "" {
			req.Header.Set(key, value)
		}
	}
}

func WithAPIKey(apiKey string) GetOption {
	return WithHeader("X-API-Key", apiKey)
}

func WithBearerToken(token string) GetOption {
	return func(req *http.Request) {
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
	}
}

func WithBasicAuth(username, password string) GetOption {
	return func(req *http.Request) {
		if username != "" || password != "" {
			req.SetBasicAuth(username, password)
		}
	}
}

// Query parameter option builders
func WithQuery(key string, values ...string) GetOption {
	return func(req *http.Request) {
		if len(values) == 0 {
			return
		}

		query := req.URL.Query()
		if len(values) == 1 {
			query.Set(key, values[0])
		} else {
			for _, value := range values {
				query.Add(key, value)
			}
		}
		req.URL.RawQuery = query.Encode()
	}
}

func WithQueryMap(params map[string]string) GetOption {
	return func(req *http.Request) {
		if len(params) == 0 {
			return
		}

		query := req.URL.Query()
		for key, value := range params {
			query.Set(key, value)
		}
		req.URL.RawQuery = query.Encode()
	}
}

// Simplify processes any input, flattens it, and removes empty fields.
func (j *Json) Simplify(input any) any {
	result := make(map[string]any)

	// Use a stack to avoid recursion
	type stackItem struct {
		value  any
		prefix string
	}

	stack := []stackItem{{value: input, prefix: ""}}

	for len(stack) > 0 {
		// Pop from stack
		item := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		switch v := item.value.(type) {
		case map[string]any:
			for key, value := range v {
				newKey := buildKey(item.prefix, key)
				stack = append(stack, stackItem{value: value, prefix: newKey})
			}
		case []any:
			for i, value := range v {
				arrayKey := buildArrayKey(item.prefix, i)
				stack = append(stack, stackItem{value: value, prefix: arrayKey})
			}
		default:
			if !isEmpty(v) {
				result[item.prefix] = v
			}
		}
	}

	if len(result) == 0 {
		return nil
	}

	return any(result)
}

func buildKey(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + "." + key
}

func buildArrayKey(prefix string, index int) string {
	arrayKey := fmt.Sprintf("[%d]", index)
	if prefix == "" {
		return arrayKey
	}
	return prefix + arrayKey
}

// isEmpty checks if a value is considered empty
func isEmpty(v any) bool {
	if v == nil {
		return true
	}

	switch val := v.(type) {
	case string:
		return val == ""
	case []any:
		return len(val) == 0
	case map[string]any:
		return len(val) == 0
	default:
		return false
	}
}
