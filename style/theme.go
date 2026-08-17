// Package style provides a shared Base16 color theme for bubbletea/lipgloss
// applications. The theme is loaded automatically on import from
// ~/.config/ilovetui/config.yaml (falling back to the embedded
// default config). Access colors and styles via the package-level variable S.
//
//	import "github.com/anotherhadi/ilovetui/style"
//
//	s := lipgloss.NewStyle().Foreground(style.S.Primary)
//	box := style.RenderWithTitle(style.S.PanelFocused, "Title", content, w, h)
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

// S is the active theme. It is populated automatically at import time and can
// be reloaded at any point by calling Init, InitFrom, or InitFromBytes.
var S Styles

func init() {
	path := DefaultConfigPath()
	if data, err := os.ReadFile(path); err == nil {
		if s, err := stylesFromBytes(data); err == nil {
			S = s
			return
		}
	}
	// Silent fallback: embedded default always works.
	s, _ := stylesFromBytes(DefaultConfig)
	S = s
}

// Init reloads S from the user config file, falling back to the embedded
// default if the file is missing. Returns an error only on parse failures.
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

// InitFrom reloads S from an explicit file path.
func InitFrom(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("style: read config: %w", err)
	}
	return InitFromBytes(data)
}

// InitFromBytes reloads S from raw YAML. Accepts hex strings with or without
// the leading '#'.
func InitFromBytes(data []byte) error {
	s, err := stylesFromBytes(data)
	if err != nil {
		return err
	}
	S = s
	return nil
}

// WriteDefaultConfig writes the embedded default config to path, creating
// parent directories as needed. No-op if the file already exists.
func WriteDefaultConfig(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("style: create config dir: %w", err)
	}
	if err := os.WriteFile(path, DefaultConfig, 0o600); err != nil {
		return fmt.Errorf("style: write config: %w", err)
	}
	return nil
}

// DefaultConfigPath returns the canonical user config path,
// respecting $XDG_CONFIG_HOME.
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
