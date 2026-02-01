package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/violenti/claudio/internal/ai"
)

type Model struct {
	providers []ai.Provider
	cursor    int
	Selected  ai.Provider
	quitting  bool
}

func InitialModel(p []ai.Provider) Model {
	return Model{
		providers: p,
	}
}

var (
	titleStyle = lipgloss.NewStyle().MarginLeft(2).Foreground(lipgloss.Color("205")).Bold(true)
	itemStyle  = lipgloss.NewStyle().PaddingLeft(4)
	selStyle   = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("86")).Bold(true)
)
