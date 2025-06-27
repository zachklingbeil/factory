package json

import (
	"encoding/json"
	"fmt"
)

// Print value as indented JSON to the standard output or logs error when value cannot be marshaled.
func (j *JSON) Print(value any) {
	json, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling Frame to JSON: %v\n", err)
		return
	}
	fmt.Println(string(json))
}

// Simplify processes a slice of objects ([]any), flattens each object, and removes empty fields.
func (j *JSON) Simplify(input []any, prefix string) []any {
	result := make([]any, 0, len(input))

	for _, item := range input {
		if obj, ok := item.(map[string]any); ok {
			if flatMap := j.flattenMap(obj, prefix); len(flatMap) > 0 {
				result = append(result, flatMap)
			}
		}
	}
	return result
}

// flattenMap recursively flattens a map with the given prefix
func (j *JSON) flattenMap(m map[string]any, prefix string) map[string]any {
	flatMap := make(map[string]any)
	j.flatten(m, prefix, flatMap)
	return flatMap
}

// flatten recursively processes map entries and populates the flat map
func (j *JSON) flatten(m map[string]any, prefix string, result map[string]any) {
	for key, value := range m {
		newKey := j.buildKey(prefix, key)

		switch v := value.(type) {
		case map[string]any:
			j.flatten(v, newKey, result)
		case []any:
			j.flattenArray(v, newKey, result)
		default:
			if !j.isEmpty(v) {
				result[newKey] = v
			}
		}
	}
}

// flattenArray processes array values and flattens nested structures
func (j *JSON) flattenArray(arr []any, prefix string, result map[string]any) {
	for i, item := range arr {
		arrayKey := fmt.Sprintf("%s[%d]", prefix, i)

		if nestedMap, ok := item.(map[string]any); ok {
			j.flatten(nestedMap, arrayKey, result)
		} else if !j.isEmpty(item) {
			result[arrayKey] = item
		}
	}
}

// buildKey constructs the flattened key with proper prefix handling
func (j *JSON) buildKey(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + "." + key
}

// isEmpty checks if a value is considered empty
func (j *JSON) isEmpty(v any) bool {
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
