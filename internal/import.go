package internal

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func importStruct[T any](path string, structure *T) error {
	text, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	switch {
	case strings.HasSuffix(path, ".json"):
		return json.Unmarshal(text, structure)
	case strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml"):
		return yaml.Unmarshal(text, structure)
	default:
		return errors.New("unsupported file type")
	}
}
