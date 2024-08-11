package view

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	ProgramInfoChoice  = 0
	RegistrationChoice = 1
	LoginChoice        = 2
)

var choices = map[int]string{
	ProgramInfoChoice:  "Get program info",
	RegistrationChoice: "Registration",
	LoginChoice:        "Login",
}

// RootModel стартовая модель
// позволяет просмотреть данные о программе, зарегистрировать или авторизоваться
type RootModel struct {
	cursor       int
	choice       string
	BuildVersion string
	BuildDate    string
	msg          string
}

func (m RootModel) Init() tea.Cmd {
	return nil
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "enter":
			//var cmd tea.Cmd
			// Send the choice on the channel and exit.
			return m.Do()
		case "down", "j":
			m.cursor++
			if m.cursor >= len(choices) {
				m.cursor = 0
			}
		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(choices) - 1
			}
		}
	}

	return m, nil
}

func (m RootModel) Do() (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.cursor {
	case ProgramInfoChoice:
		m.msg = fmt.Sprintf("Build date: %s; Build Version: %s", m.BuildDate, m.BuildVersion)
	case RegistrationChoice:
		return initAuthModel(false), cmd
	case LoginChoice:
		return initAuthModel(true), cmd
	}

	return m, tea.Batch(cmd, m.Init())
}

func (m RootModel) View() string {
	s := strings.Builder{}

	if len(m.msg) > 0 {
		s.WriteString(infoStyle.Render(m.msg) + "\n\n")
	}

	s.WriteString("What kind of Bubble Tea would you like to order?\n\n")

	for i := 0; i < len(choices); i++ {
		if m.cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(choices[i])
		s.WriteString("\n")
	}
	s.WriteString("\n(press q to quit)\n")

	return s.String()
}
