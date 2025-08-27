//go:build js && wasm

package persistence

import (
	"encoding/json"
	"fmt"
	"syscall/js"
)

func SaveState(pathOrKey string) error {
	data, err := json.Marshal(State)
	if err != nil {
		return fmt.Errorf("failed to marshal state for localstorage: %w", err)
	}
	localStorage := js.Global().Get("localStorage")
	localStorage.Call("setItem", pathOrKey, string(data))
	return nil
}

func LoadState(pathOrKey string) error {
	localStorage := js.Global().Get("localStorage")
	value := localStorage.Call("getItem", pathOrKey)

	if value.IsNull() || value.IsUndefined() {
		return ErrStateNotFound
	}

	data := []byte(value.String())
	return json.Unmarshal(data, &State)
}
