package style

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

//go:embed default.yaml
var DefaultConfig []byte

var S Styles

func init() {
	path := DefaultConfigPath()
	if data, err := os.ReadFile(path); err == nil {
		if s, err := stylesFromBytes(data); err == nil {
			S = s
			return
		}
	}

	s, _ := stylesFromBytes(DefaultConfig)
	S = s
}

func Init() error {
	path := DefaultConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		s, e := stylesFromBytes(DefaultConfig)
		if e != nil {
			return e
		}
		S = s
		return nil
	}
	return InitFromBytes(data)
}

func InitFrom(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("style: read config: %w", err)
	}
	return InitFromBytes(data)
}

func InitFromBytes(data []byte) error {
	s, err := stylesFromBytes(data)
	if err != nil {
		return err
	}
	S = s
	return nil
}

func DefaultConfigPath() string {
	return filepath.Join(configDir(), "ilovetui", "config.yaml")
}

func stylesFromBytes(data []byte) (Styles, error) {
	var base configYAML
	if err := yaml.Unmarshal(DefaultConfig, &base); err != nil {
		return Styles{}, fmt.Errorf("style: parse default config: %w", err)
	}
	var user configYAML
	if err := yaml.Unmarshal(data, &user); err != nil {
		return Styles{}, fmt.Errorf("style: parse config: %w", err)
	}
	merged := mergeConfig(base, user)
	return newStyles(merged.Colors, merged.NerdFonts, merged.Border), nil
}

func configDir() string {
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return dir
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config")
}
