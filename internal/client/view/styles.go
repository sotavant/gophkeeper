package view

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffa657")) // orange
	blueStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#57caff")) // orange
	actionsStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#16e931"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#ef1515"))
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#f4d934"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle
	noStyle      = lipgloss.NewStyle()
	helpStyle    = blurredStyle

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)
