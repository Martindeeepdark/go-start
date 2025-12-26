package cache

import (
	"encoding/json"
)

// Marshal serializes a value to JSON string for storage in cache
func Marshal(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Unmarshal deserializes a JSON string from cache to a value
func Unmarshal(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}
