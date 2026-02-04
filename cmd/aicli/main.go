package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
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
		ai.MockIA{},
	}

	fd := int(os.Stdin.Fd())
	oldState, err := term.GetState(fd)
	if err != nil {
		panic(err)
	}

	// --- Bubble Tea (menú) ---
	p := tea.NewProgram(ui.InitialModel(motors))
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

	cmd := exec.Command(os.Args[0], "chat", strconv.Itoa(selectedIndex))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error restarting: %v\n", err)
	}
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
		ai.MockIA{},
	}

	if selectedIndex < 0 || selectedIndex >= len(motors) {
		fmt.Println("Error: invalid provider index")
		return
	}
	switch selectedIndex {
	case 0: // OpenAI
		key := getAPIKey("OPENAI_API_KEY", "OpenAI")
		motors[0] = ai.OpenAI{Token: key}
	case 1: // Claude
		key := getAPIKey("ANTHROPIC_API_KEY", "Claude")
		motors[1] = ai.Claude{ApiKey: key}

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

	color.Yellow("--- Talking with: %s ---", selectedProvider.Name())
	color.Cyan("(Type 'exit' to quit)")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\r")
		fmt.Print(color.GreenString("You: "))
		os.Stdout.Sync()

		fmt.Printf("")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			break
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		if strings.EqualFold(input, "exit") {
			color.Magenta("Bye!")
			break
		}

		response, err := selectedProvider.Question(input)
		if err != nil {
			color.Red("Error: %v", err)
			continue
		}

		fmt.Printf("\n%s %s\n", color.WhiteString("Claudio:"), response)
	}
}
