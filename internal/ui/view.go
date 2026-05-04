package ui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

func (m Model) View() string {
	if m.quitting {
		return "Exit...\n"
	}

	var b strings.Builder

	banner := `
############################################################
#   _____ _        _    _   _ ____ ___ ___                 #
#  / ____| |      / \  | | | |  _ \_ _/ _ \                #
# | |    | |     / _ \ | | | | | | | | | | |               #
# | |____| |___ / ___ \| |_| | |_| | | |_| |               #
#  \_____|_____/_/   \_\\___/|____/___\___/                #
#                                                          #
############################################################

                                         .;+xx:
                                         ;xxxxx:
                                          ..;:x;
                                         ::;::+:;.
                            ..          ::::+X+;xx;
                            . ..   ..   :::;x+;xxxx;
                              ..... .       .xxxxxxx+
                              ..  :..       ;xxxxxxxx+:
                              .. ....      +xxxxxxxxx+. :
                                .. ..    ....:;::;..      ...
                                 .  .. .. :                  ..
                                 .    :  .                      :
                                 .       .                        .
                                 .       .            .....       .xx;
                                 ..      .               : ..     ;xxx+
                                  .:.:....                :.     ++:xxx:
                                        ..                .    :x++.:xx;
                                        ..               :   .xxxxx+.xx:
                                        .               .    .::xxxx.+x:
                                        .               ..     ;.+xx.+x:
                                        .                   .. ++:x+.+x.
                                        ..                     ;x:x+.+x
                                         .                     ;x.x;.++
                                          :                    :x ++ ;+
                                           ..                ..    +.
                                             ... .::....    ..
                                               ::;      .;:;.
                                                ;;.      .:;.
                                                :::       :::
                                       :;;;;;:;;::::      .:;.
                                       :;:::::::.         ::::::;:::.
                                                            .:::;;;::;
                                                              ..:;::..
  ########################`

	// Apply yellow color to the banner
	yellowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("226")) // Bright yellow
	b.WriteString(yellowStyle.Render(banner) + "\n\n")

	switch m.state {
	case stateProviders:
		b.WriteString(titleStyle.Render("Choose the ai:") + "\n\n")
		for i, p := range m.providers {
			if m.cursor == i {
				b.WriteString(selStyle.Render(fmt.Sprintf("> %s", p.Name())))
			} else {
				b.WriteString(itemStyle.Render(p.Name()))
			}
			b.WriteString("\n")
		}
		b.WriteString("\n (Use the arrow keys to move, press enter to select)\n")

	case stateModels:
		b.WriteString(titleStyle.Render(fmt.Sprintf("Choose the model for %s:", m.pendingProvider.Name())) + "\n\n")
		for i, name := range m.modelNames {
			if m.cursor == i {
				b.WriteString(selStyle.Render(fmt.Sprintf("> %s", name)))
			} else {
				b.WriteString(itemStyle.Render(name))
			}
			b.WriteString("\n")
		}
		b.WriteString("\n (Use the arrow keys to move, press enter to select, esc to go back)\n")
	}

	return b.String()
}
