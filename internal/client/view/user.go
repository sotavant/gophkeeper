package view

import (
	"errors"
	"fmt"
	"gophkeeper/client/domain"
	"gophkeeper/client/user"
	"gophkeeper/internal/client"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	DataListChoice = 0
	AddDataChoice  = 1
)

var userModelChoices = map[int]string{
	DataListChoice: "Get data list",
	AddDataChoice:  "Add data",
}

// UserModel модель для авторизованного пользователя
// позволяет перейти к списку данных, либо в форму добавления
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
		case "ctrl+w":
			var cmd tea.Cmd
			rm := RootModel{}

			user.ResetUser()

			return rm, tea.Batch(cmd, rm.Init())
		case "enter":
			//var cmd tea.Cmd
			// Send the choice on the channel and exit.
			return m.Do()
		case "down", "j":
			m.cursor++
			if m.cursor >= len(userModelChoices) {
				m.cursor = 0
			}
		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(userModelChoices) - 1
			}
		}
	}

	return m, nil
}

func (m UserModel) Do() (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.cursor {
	case DataListChoice:
		return InitDataListModel(), cmd
	case AddDataChoice:
		var data domain.Data
		return InitDataFieldsModel(data), cmd
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

	s.WriteString("What you want to do?\n\n")

	for i := 0; i < len(userModelChoices); i++ {
		if m.cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(userModelChoices[i])
		s.WriteString("\n")
	}

	s.WriteString(helpStyle.Render("\n\n'ctrl+w' to main window"))
	s.WriteString("\n(press q to quit)\n")

	return s.String()
}
