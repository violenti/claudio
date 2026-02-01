package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/violenti/claudio/internal/ai"
	"github.com/violenti/claudio/internal/ui"
)

func main() {

	motors := []ai.Provider{
		ai.OpenAI{},
		ai.Claude{},
	}

	m := ui.InitialModel(motors)

	// Bubble Tea
	p := tea.NewProgram(m)

	finalModel, err := p.Run()
	if err != nil {
		color.Red("Error : %v\n", err)
		os.Exit(1)
	}

	// 4. Lógica post-selección:
	// Convertimos el resultado (que es una interfaz tea.Model) a nuestro struct ui.Model
	if mResult, ok := finalModel.(ui.Model); ok {

		// Verificamos si el usuario seleccionó algo antes de salir
		if mResult.Selected != nil {
			color.Yellow("\nYou select: %s\n", mResult.Selected.Name())

			// Aquí es donde llamas a la lógica de la IA
			question, err := mResult.Selected.Question("Hi, how can you help me?")
			if err != nil {
				color.Red("Error querying AI: %v\n", err)
				return
			}
			color.Yellow("Question:", question)
		} else {
			color.Red("\nProgram completed without selection.")
		}
	}
}
