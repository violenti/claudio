package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/violenti/claudio/internal/ai"
)

type mockProvider struct {
	name string
}

func (m mockProvider) Name() string                    { return m.name }
func (m mockProvider) Question(string) (string, error) { return "", nil }

func providers(names ...string) []ai.Provider {
	ps := make([]ai.Provider, len(names))
	for i, n := range names {
		ps[i] = mockProvider{n}
	}
	return ps
}

func pressKey(m Model, key string) Model {
	var msg tea.KeyMsg
	switch key {
	case "enter":
		msg = tea.KeyMsg{Type: tea.KeyEnter}
	case "up":
		msg = tea.KeyMsg{Type: tea.KeyUp}
	case "down":
		msg = tea.KeyMsg{Type: tea.KeyDown}
	case "esc":
		msg = tea.KeyMsg{Type: tea.KeyEsc}
	default:
		msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
	}
	result, _ := m.Update(msg)
	return result.(Model)
}

func TestUpdate_SelectProvider_NoModels(t *testing.T) {
	m := InitialModel(providers("OpenAI", "Claude"), map[string]map[string]string{})

	m = pressKey(m, "enter")

	if m.Selected == nil {
		t.Fatal("expected Selected to be set")
	}
	if m.Selected.Name() != "OpenAI" {
		t.Errorf("expected OpenAI, got %q", m.Selected.Name())
	}
}

func TestUpdate_SelectProvider_WithModels_TransitionsToSubmenu(t *testing.T) {
	modelOpts := map[string]map[string]string{
		"Claude": {
			"Claude Opus":  "claude-3-opus-20240229",
			"Claude Haiku": "claude-3-5-haiku-20241022",
		},
	}
	m := InitialModel(providers("OpenAI", "Claude"), modelOpts)

	m = pressKey(m, "down") // move to Claude
	m = pressKey(m, "enter")

	if m.state != stateModels {
		t.Errorf("expected stateModels, got %v", m.state)
	}
	if m.pendingProvider == nil || m.pendingProvider.Name() != "Claude" {
		t.Error("pendingProvider should be Claude")
	}
	if len(m.modelNames) != 2 {
		t.Errorf("expected 2 model names, got %d", len(m.modelNames))
	}
	if m.cursor != 0 {
		t.Errorf("cursor should reset to 0, got %d", m.cursor)
	}
}

func TestUpdate_SelectModel(t *testing.T) {
	modelOpts := map[string]map[string]string{
		"Claude": {
			"Claude Haiku": "claude-3-5-haiku-20241022",
			"Claude Opus":  "claude-3-opus-20240229",
		},
	}
	m := InitialModel(providers("Claude"), modelOpts)

	m = pressKey(m, "enter") // enter submenu
	if m.state != stateModels {
		t.Fatal("expected to be in stateModels")
	}

	firstModelName := m.modelNames[0]
	m = pressKey(m, "enter") // select first model

	if m.Selected == nil {
		t.Fatal("expected Selected to be set after model selection")
	}
	if m.Selected.Name() != "Claude" {
		t.Errorf("expected Selected = Claude, got %q", m.Selected.Name())
	}
	if m.SelectedModel != modelOpts["Claude"][firstModelName] {
		t.Errorf("SelectedModel = %q, want %q", m.SelectedModel, modelOpts["Claude"][firstModelName])
	}
}

func TestUpdate_EscapeFromModels(t *testing.T) {
	modelOpts := map[string]map[string]string{
		"Claude": {"Claude Haiku": "claude-3-5-haiku-20241022"},
	}
	m := InitialModel(providers("OpenAI", "Claude"), modelOpts)

	m = pressKey(m, "down")
	m = pressKey(m, "enter")
	if m.state != stateModels {
		t.Fatal("expected stateModels")
	}

	m = pressKey(m, "esc")

	if m.state != stateProviders {
		t.Errorf("expected stateProviders after esc, got %v", m.state)
	}
	if m.cursor != 0 {
		t.Errorf("expected cursor reset to 0, got %d", m.cursor)
	}
}

func TestUpdate_Navigation_Providers(t *testing.T) {
	m := InitialModel(providers("A", "B", "C"), map[string]map[string]string{})

	if m.cursor != 0 {
		t.Errorf("initial cursor should be 0, got %d", m.cursor)
	}

	m = pressKey(m, "down")
	if m.cursor != 1 {
		t.Errorf("cursor should be 1, got %d", m.cursor)
	}

	m = pressKey(m, "down")
	m = pressKey(m, "down") // boundary
	if m.cursor != 2 {
		t.Errorf("cursor should stop at 2, got %d", m.cursor)
	}

	m = pressKey(m, "up")
	if m.cursor != 1 {
		t.Errorf("cursor should be 1 after up, got %d", m.cursor)
	}
}

func TestUpdate_Navigation_Models(t *testing.T) {
	modelOpts := map[string]map[string]string{
		"Claude": {"Model A": "a", "Model B": "b", "Model C": "c"},
	}
	m := InitialModel(providers("Claude"), modelOpts)
	m = pressKey(m, "enter") // enter submenu

	m = pressKey(m, "down")
	if m.cursor != 1 {
		t.Errorf("expected cursor 1, got %d", m.cursor)
	}

	m = pressKey(m, "down")
	m = pressKey(m, "down") // boundary
	if m.cursor != 2 {
		t.Errorf("cursor should stop at last model, got %d", m.cursor)
	}
}

func TestUpdate_ModelNames_AreSorted(t *testing.T) {
	modelOpts := map[string]map[string]string{
		"Claude": {"Zebra Model": "z", "Alpha Model": "a", "Middle": "m"},
	}
	m := InitialModel(providers("Claude"), modelOpts)
	m = pressKey(m, "enter")

	if m.modelNames[0] != "Alpha Model" {
		t.Errorf("first model should be 'Alpha Model', got %q", m.modelNames[0])
	}
	if m.modelNames[2] != "Zebra Model" {
		t.Errorf("last model should be 'Zebra Model', got %q", m.modelNames[2])
	}
}
