package main

import (
	"fmt"
	"gophkeeper/internal"
	"gophkeeper/internal/client"
	"gophkeeper/internal/client/view"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// Build info.
// Need define throw ldflags:
//
//	go build -ldflags "-X 'main.buildDate=$(date +'%Y/%m/%d')' -X 'main.buildVersion=$(git rev-parse --short HEAD)' -X 'main.buildCryptoKeysPath=path' -X 'main.buildSaveFilePath=path'"
var (
	buildVersion,
	buildDate,
	buildCryptoKeysPath,
	buildSaveFilePath string
)

func main() {
	internal.InitLogger()
	err := client.InitApp(buildCryptoKeysPath, buildSaveFilePath)
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if len(os.Getenv("DEBUG")) > 0 {
		var f *os.File
		f, err = tea.LogToFile("debug.log", "debug")
		if err != nil {
			internal.Logger.Fatalw("Error logging to file", "error", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	if _, err = tea.NewProgram(view.RootModel{BuildDate: buildDate, BuildVersion: buildVersion}).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
