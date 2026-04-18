package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadYAML(path string, out any) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config %s: %w", path, err)
	}
	if err := yaml.Unmarshal(content, out); err != nil {
		return fmt.Errorf("unmarshal config %s: %w", path, err)
	}
	return nil
}
