// Package jsonx unifies Parse behaviour across tool adapters:
// nil input yields the zero value, non-RawMessage yields an error
// instead of panicking on a blind type assertion.
package jsonx

import (
	"encoding/json"
	"fmt"
)

func Decode(in any, v any) error {
	if in == nil {
		return nil
	}
	raw, ok := in.(json.RawMessage)
	if !ok {
		return fmt.Errorf("expected json.RawMessage, got %T", in)
	}
	if len(raw) == 0 {
		return nil
	}
	return json.Unmarshal(raw, v)
}
