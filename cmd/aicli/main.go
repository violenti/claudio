package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/violenti/claudio/internal/ai"
	"github.com/violenti/claudio/internal/ui"
	"golang.org/x/term"
)

func getAPIKey(envVar string, providerNmae string) string {
	key := os.Getenv(envVar)
	if key == "" {
		color.Red("%s not set. Enter API Key:", envVar)
		bytesKey, _ := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		key = string(bytesKey)
	}
	return key
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "chat" {
		chatMode()
		return
	}

	motors := []ai.Provider{
		ai.OpenAI{},
		ai.Claude{},
		ai.Ollama{},
	}

	modelOptions := map[string]map[string]string{}
	config, err := ai.ReadModels()
	if err == nil {
		providerKeys := map[string]string{
			ai.Claude{}.Name(): "anthropic",
			ai.OpenAI{}.Name(): "openai",
			ai.Ollama{}.Name(): "ollama",
		}
		for _, p := range motors {
			if key, ok := providerKeys[p.Name()]; ok {
				if models, exists := config.Models[key]; exists {
					modelOptions[p.Name()] = models
				}
			}
		}
	}

	fd := int(os.Stdin.Fd())
	oldState, err := term.GetState(fd)
	if err != nil {
		panic(err)
	}

	// --- Bubble Tea (menú) ---
	p := tea.NewProgram(ui.InitialModel(motors, modelOptions))
	finalModel, err := p.Run()

	_ = term.Restore(fd, oldState)

	fmt.Print("\033c")
	fmt.Print("\033[?1049l")
	fmt.Print("\033[?2004l")
	fmt.Print("\033[?1000l")
	fmt.Print("\033[?1006l")
	fmt.Print("\033[?1015l")
	fmt.Print("\033[0m")
	fmt.Print("\033[?25h")

	_ = os.Stdout.Sync()
	_ = os.Stderr.Sync()

	_ = term.Restore(int(os.Stdin.Fd()), oldState)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	mResult, ok := finalModel.(ui.Model)
	if !ok || mResult.Selected == nil {
		fmt.Println("No selection made.")
		return
	}

	selectedIndex := -1
	for i, provider := range motors {
		if provider.Name() == mResult.Selected.Name() {
			selectedIndex = i
			break
		}
	}
	_ = term.Restore(fd, oldState)

	cmd := exec.Command(os.Args[0], "chat", strconv.Itoa(selectedIndex), mResult.SelectedModel)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error restarting: %v\n", err)
	}
}

// wrapText wraps text to fit within the specified width
func wrapText(text string, width int) []string {
	var lines []string
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{text}
	}

	currentLine := ""
	for _, word := range words {
		if len(currentLine)+len(word)+1 > width {
			if currentLine != "" {
				lines = append(lines, currentLine)
				currentLine = word
			} else {
				// Word is longer than width, split it
				lines = append(lines, word[:width])
				currentLine = word[width:]
			}
		} else {
			if currentLine != "" {
				currentLine += " " + word
			} else {
				currentLine = word
			}
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}
	return lines
}

func chatMode() {
	if len(os.Args) < 3 {
		fmt.Println("Error: missing provider index")
		return
	}

	selectedIndex, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("Error parsing provider index: %v\n", err)
		return
	}

	motors := []ai.Provider{
		ai.OpenAI{},
		ai.Claude{},
		ai.Ollama{},
	}

	if selectedIndex < 0 || selectedIndex >= len(motors) {
		fmt.Println("Error: invalid provider index")
		return
	}

	selectedModel := ""
	if len(os.Args) > 3 {
		selectedModel = os.Args[3]
	}

	switch selectedIndex {
	case 0: // OpenAI
		key := getAPIKey("OPENAI_API_KEY", "OpenAI")
		motors[0] = ai.OpenAI{Token: key, Model: selectedModel}
	case 1: // Claude
		key := getAPIKey("ANTHROPIC_API_KEY", "Claude")
		motors[1] = ai.Claude{ApiKey: key, Model: selectedModel}
	case 2: // Ollama
		motors[2] = ai.Ollama{}
	}

	selectedProvider := motors[selectedIndex]

	fmt.Print("\033c")

	fd := int(os.Stdin.Fd())
	if term.IsTerminal(fd) {
		oldState, err := term.MakeRaw(fd)
		if err == nil {
			_ = term.Restore(fd, oldState)
		}
	}

	// Draw header
	fmt.Println()
	fmt.Printf("  %s\n", color.YellowString("Talking with: %s", selectedProvider.Name()))
	fmt.Printf("  %s\n", color.CyanString("Commands: 'exit' to quit, 'menu' to return, Ctrl+C twice for menu"))
	fmt.Println(strings.Repeat("─", 65))
	fmt.Println()

	// Configure readline for better input handling
	rl, err := readline.NewEx(&readline.Config{
		Prompt:                 "", // We'll set this dynamically
		HistoryFile:            "/tmp/claudio_history.txt",
		AutoComplete:           nil,
		InterruptPrompt:        "^C",
		EOFPrompt:              "exit",
		HistorySearchFold:      true,
		DisableAutoSaveHistory: false,
		VimMode:                false,
	})
	if err != nil {
		fmt.Printf("Error initializing readline: %v\n", err)
		return
	}
	defer rl.Close()

	returnToMenu := false
	ctrlCCount := 0

	for {
		// Static prompt with animated arrow
		rl.SetPrompt(color.GreenString("You ▸ "))

		// Use readline for input with cursor movement support
		input, err := rl.Readline()

		if err != nil {
			if err == readline.ErrInterrupt {
				ctrlCCount++
				if ctrlCCount == 1 {
					// First Ctrl+C
					fmt.Printf("%s\n", color.YellowString("  Press Ctrl+C again to return to menu"))
					continue
				} else {
					// Second Ctrl+C
					returnToMenu = true
					fmt.Printf("\n%s\n", color.YellowString("  Returning to menu..."))
					break
				}
			} else {
				break
			}
		}

		// Reset Ctrl+C counter on successful input
		ctrlCCount = 0

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		if strings.EqualFold(input, "exit") {
			fmt.Printf("\n%s\n", color.MagentaString("  Bye!"))
			return
		}

		if strings.EqualFold(input, "menu") {
			returnToMenu = true
			fmt.Printf("\n%s\n", color.YellowString("  Returning to menu..."))
			break
		}

		// Show thinking animation
		fmt.Println()
		stopSpinner := make(chan bool, 1)
		go func() {
			spinners := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
			i := 0
			for {
				select {
				case <-stopSpinner:
					fmt.Print("\r") // Clear the spinner line
					return
				default:
					fmt.Printf("\r%s %s", color.MagentaString("Claudio thinking"), color.YellowString(spinners[i%len(spinners)]))
					time.Sleep(100 * time.Millisecond)
					i++
				}
			}
		}()

		response, err := selectedProvider.Question(input)
		stopSpinner <- true

		if err != nil {
			fmt.Printf("%s\n\n", color.RedString("  Error: %v", err))
			continue
		}

		// Display response with simple formatting
		fmt.Printf("\r%s\n", color.MagentaString("Claudio:"))
		lines := wrapText(response, 60)
		for _, line := range lines {
			fmt.Printf("  %s\n", line)
		}
		fmt.Println()
		fmt.Println(strings.Repeat("─", 65))
		fmt.Println()
	}

	// If user wants to return to menu, restart the program without chat mode
	if returnToMenu {
		cmd := exec.Command(os.Args[0])
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			fmt.Printf("Error returning to menu: %v\n", err)
		}
	}
}
