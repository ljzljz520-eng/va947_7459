package store

import (
	"encoding/json"
	"fmt"
)

func encode(value any) ([]byte, error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("encode value: %w", err)
	}
	return bytes, nil
}

func decode(data []byte, target any) error {
	if len(data) == 0 {
		return fmt.Errorf("empty stored value")
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("decode value: %w", err)
	}
	return nil
}

func key(value string) []byte {
	return []byte(value)
}

func keyForParts(parts ...string) []byte {
	encoded, _ := encode(parts)
	return encoded
}
