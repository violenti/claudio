package ui

import (
	"sort"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "esc":
			if m.state == stateModels {
				m.state = stateProviders
				m.cursor = 0
			}

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			switch m.state {
			case stateProviders:
				if m.cursor < len(m.providers)-1 {
					m.cursor++
				}
			case stateModels:
				if m.cursor < len(m.modelNames)-1 {
					m.cursor++
				}
			}

		case "enter":
			switch m.state {
			case stateProviders:
				if len(m.providers) == 0 {
					return m, nil
				}
				p := m.providers[m.cursor]
				models, hasModels := m.modelOptions[p.Name()]
				if hasModels && len(models) > 0 {
					m.pendingProvider = p
					m.state = stateModels
					names := make([]string, 0, len(models))
					for name := range models {
						names = append(names, name)
					}
					sort.Strings(names)
					m.modelNames = names
					m.cursor = 0
				} else {
					m.Selected = p
					return m, tea.Quit
				}

			case stateModels:
				if len(m.modelNames) == 0 {
					return m, nil
				}
				modelName := m.modelNames[m.cursor]
				m.SelectedModel = m.modelOptions[m.pendingProvider.Name()][modelName]
				m.Selected = m.pendingProvider
				return m, tea.Quit
			}
		}
	}

	return m, nil
}
