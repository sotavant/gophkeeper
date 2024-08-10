package main

// A simple example that shows how to retrieve a value from a Bubble Tea
// program after the Bubble Tea has exited.

import (
	"fmt"
	"gophkeeper/client/domain"
	"gophkeeper/internal"
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
	internal.InitLogger()
	/*	err := client.InitApp()
		if err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}*/

	/*if _, err = tea.NewProgram(view.RootModel{BuildDate: buildDate, BuildVersion: buildVersion}).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}*/
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	var data domain.Data
	if _, err := tea.NewProgram(view.InitDataFieldsModel(data)).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
