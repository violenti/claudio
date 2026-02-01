package ui

import (
	"fmt"
	"strings"
)

func (m Model) View() string {
	if m.quitting {
		return "Exit...\n"
	}

	var b strings.Builder

	banner := `
  ########################
  #   Claudio   #
  ########################`

	b.WriteString(banner + "\n\n")
	b.WriteString(titleStyle.Render("Chooise the ai:") + "\n\n")

	for i, p := range m.providers {
		if m.cursor == i {
			b.WriteString(selStyle.Render(fmt.Sprintf("> %s", p.Name())))
		} else {
			b.WriteString(itemStyle.Render(p.Name()))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n (Use the arrow keys to move, press enter to select)\n")

	return b.String()
}
