package ai

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeConfigFile(t *testing.T, dir string, content any) {
	t.Helper()
	configDir := filepath.Join(dir, ".claudio")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(content)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(configDir, "config.json"), data, 0644); err != nil {
		t.Fatal(err)
	}
}

func TestReadModels_ValidConfig(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	writeConfigFile(t, tmp, map[string]any{
		"models": map[string]any{
			"anthropic": map[string]string{
				"Claude Haiku": "claude-3-5-haiku-20241022",
				"Claude Opus":  "claude-3-opus-20240229",
			},
			"openai": map[string]string{
				"GPT-5 Mini": "gpt-5-mini",
			},
		},
	})

	config, err := ReadModels()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(config.Models) != 2 {
		t.Errorf("expected 2 providers, got %d", len(config.Models))
	}

	anthropic := config.Models["anthropic"]
	if anthropic["Claude Haiku"] != "claude-3-5-haiku-20241022" {
		t.Errorf("unexpected model ID: %q", anthropic["Claude Haiku"])
	}
	if anthropic["Claude Opus"] != "claude-3-opus-20240229" {
		t.Errorf("unexpected model ID: %q", anthropic["Claude Opus"])
	}

	openai := config.Models["openai"]
	if openai["GPT-5 Mini"] != "gpt-5-mini" {
		t.Errorf("unexpected model ID: %q", openai["GPT-5 Mini"])
	}
}

func TestReadModels_MissingFile(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	_, err := ReadModels()
	if err == nil {
		t.Error("expected error for missing config file, got nil")
	}
}

func TestReadModels_InvalidJSON(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	configDir := filepath.Join(tmp, ".claudio")
	os.MkdirAll(configDir, 0755)
	os.WriteFile(filepath.Join(configDir, "config.json"), []byte("{invalid json}"), 0644)

	_, err := ReadModels()
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestReadModels_EmptyModels(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	writeConfigFile(t, tmp, map[string]any{"models": map[string]any{}})

	config, err := ReadModels()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(config.Models) != 0 {
		t.Errorf("expected 0 providers, got %d", len(config.Models))
	}
}
