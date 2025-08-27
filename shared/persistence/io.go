//go:build !js

package persistence

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

func SaveState(pathOrKey string) error {
	data, err := json.MarshalIndent(State, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}
	return os.WriteFile(pathOrKey, data, 0644)
}

func LoadState(pathOrKey string) error {
	data, err := os.ReadFile(pathOrKey)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrStateNotFound
		}
		return fmt.Errorf("failed to read save file: %w", err)
	}
	return json.Unmarshal(data, &State)
}
