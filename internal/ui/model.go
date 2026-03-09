package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/violenti/claudio/internal/ai"
)

type state int

const (
	stateProviders state = iota
	stateModels
)

type Model struct {
	providers       []ai.Provider
	cursor          int
	Selected        ai.Provider
	SelectedModel   string
	quitting        bool
	state           state
	modelOptions    map[string]map[string]string // provider name → {display name → model ID}
	modelNames      []string                     // current submenu model display names
	pendingProvider ai.Provider
}

func InitialModel(p []ai.Provider, modelOptions map[string]map[string]string) Model {
	return Model{
		providers:    p,
		modelOptions: modelOptions,
	}
}

var (
	titleStyle = lipgloss.NewStyle().
			MarginLeft(2).
			Foreground(lipgloss.Color("205")).
			Bold(true)

	itemStyle = lipgloss.NewStyle().
			PaddingLeft(4)

	selStyle = lipgloss.NewStyle().
			PaddingLeft(2).
			Foreground(lipgloss.Color("86")).
			Bold(true)
)
