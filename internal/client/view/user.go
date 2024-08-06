package view

import (
	"errors"
	"fmt"
	"gophkeeper/internal/client"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type UserModel struct {
	cursor int
	choice string
	msg    string
}

func (m UserModel) Init() tea.Cmd {
	return nil
}

func (m UserModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m UserModel) Do() (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.cursor {
	case ProgramInfoChoice:
		m.msg = fmt.Sprintf("Build date: %s; Build Version: %s")
	case RegistrationChoice:
		return initialRegistrationModel(), cmd
	}

	return m, tea.Batch(cmd, m.Init())
}

func (m UserModel) View() string {
	s := strings.Builder{}

	if len(client.AppInstance.User.Token) == 0 {
		getError(errors.New("sorry( auth data is empty"))
	} else {
		s.WriteString(infoStyle.Render(fmt.Sprintf("Hi!, %s", client.AppInstance.User.Login)) + "\n\n")
	}

	if len(m.msg) > 0 {
		s.WriteString(infoStyle.Render(m.msg) + "\n\n")
	}

	s.WriteString("What you want?\n\n")

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
