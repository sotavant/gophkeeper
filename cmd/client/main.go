package main

// A simple example that shows how to retrieve a value from a Bubble Tea
// program after the Bubble Tea has exited.

import (
	"fmt"
	"gophkeeper/internal/client/view"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// Build info.
// Need define throw ldflags:
//
//	go build -ldflags "-X 'main.buildDate=$(date +'%Y/%m/%d')' -X 'main.buildVersion=$(git rev-parse --short HEAD)'"
var (
	buildVersion string
	buildDate    string
)

func main() {
	if _, err := tea.NewProgram(view.RootModel{BuildDate: buildDate, BuildVersion: buildVersion}).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
