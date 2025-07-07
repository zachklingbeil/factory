package fx

import (
	"encoding/json"
	"fmt"
)

// Print value as indented JSON to the standard output or logs error when value cannot be marshaled.
func (a *API) Print(value any) {
	json, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling Frame to JSON: %v\n", err)
		return
	}
	fmt.Println(string(json))
}

// Simplify processes a slice of objects ([]any), flattens each object, and removes empty fields.
func (a *API) Simplify(input []any, prefix string) []any {
	result := make([]any, 0, len(input))

	for _, item := range input {
		if obj, ok := item.(map[string]any); ok {
			if flatMap := a.flattenMap(obj, prefix); len(flatMap) > 0 {
				result = append(result, flatMap)
			}
		}
	}
	return result
}

// flattenMap recursively flattens a map with the given prefix
func (a *API) flattenMap(m map[string]any, prefix string) map[string]any {
	flatMap := make(map[string]any)
	a.flatten(m, prefix, flatMap)
	return flatMap
}

// flatten recursively processes map entries and populates the flat map
func (a *API) flatten(m map[string]any, prefix string, result map[string]any) {
	for key, value := range m {
		newKey := a.buildKey(prefix, key)

		switch v := value.(type) {
		case map[string]any:
			a.flatten(v, newKey, result)
		case []any:
			a.flattenArray(v, newKey, result)
		default:
			if !a.isEmpty(v) {
				result[newKey] = v
			}
		}
	}
}

// flattenArray processes array values and flattens nested structures
func (a *API) flattenArray(arr []any, prefix string, result map[string]any) {
	for i, item := range arr {
		arrayKey := fmt.Sprintf("%s[%d]", prefix, i)

		if nestedMap, ok := item.(map[string]any); ok {
			a.flatten(nestedMap, arrayKey, result)
		} else if !a.isEmpty(item) {
			result[arrayKey] = item
		}
	}
}

// buildKey constructs the flattened key with proper prefix handling
func (a *API) buildKey(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + "." + key
}

// isEmpty checks if a value is considered empty
func (a *API) isEmpty(v any) bool {
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
