package view

import (
	"fmt"
	"gophkeeper/client/user"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type AuthModel struct {
	focusIndex    int
	inputs        []textinput.Model
	msg           string
	isLoginAction bool
}

func initAuthModel(isLoginAction bool) AuthModel {
	m := AuthModel{
		inputs:        make([]textinput.Model, 2),
		isLoginAction: isLoginAction,
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = loginPlaceholder
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = PasswordPlaceholder
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.inputs[i] = t
	}

	return m
}

func (m AuthModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m AuthModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "ctrl+w":
			var cmd tea.Cmd
			rm := RootModel{}
			return rm, tea.Batch(cmd, rm.Init())

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				return m.registrateUser(msg)
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m AuthModel) registrateUser(teaMsg tea.Msg) (tea.Model, tea.Cmd) {
	var pass, login string
	var userModel UserModel

	for _, v := range m.inputs {
		switch v.Placeholder {
		case loginPlaceholder:
			login = v.Value()
		case PasswordPlaceholder:
			pass = v.Value()
		}
	}

	err := user.Auth(login, pass, m.isLoginAction)

	if err != nil {
		m.msg = getError(err)
		return m, m.updateInputs(teaMsg)
	}

	var cmd tea.Cmd

	return userModel, cmd
}

func (m *AuthModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m AuthModel) View() string {
	var b strings.Builder

	if len(m.msg) > 0 {
		b.WriteString(errorStyle.Render(m.msg) + "\n\n")
	}

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("'ctrl+w' to main window\n'ctrl-c' to quit"))

	return b.String()
}
