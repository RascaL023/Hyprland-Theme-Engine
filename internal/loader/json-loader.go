package loader

import (
	"encoding/json"
	"os"
)

func LoadJSON[T any](path string) (*T, error) {
	data, err := os.ReadFile(path);
	if err != nil {
		return nil, err;
	}

	var t T;
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, err;
	}

	return &t, nil;
}

