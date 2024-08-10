package view

import (
	"fmt"
	"gophkeeper/client/domain"
	"gophkeeper/internal"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type field struct {
	key,
	name string
}

const (
	nameFieldKey    = "name"
	loginFieldKey   = "login"
	passFieldKey    = "pass"
	cardNumFieldKey = "card_num"
	fileFieldKey    = "file"
)

var dataFields = []field{
	{
		key:  nameFieldKey,
		name: nameFieldName,
	},
	{
		key:  loginFieldKey,
		name: loginFieldName,
	},
	{
		key:  passFieldKey,
		name: passFieldName,
	},
	{
		key:  cardNumFieldKey,
		name: cardNumFieldName,
	},
	{
		key:  fileFieldKey,
		name: fileFieldName,
	},
}

type DataFieldsModel struct {
	focusIndex int
	inputs     []textinput.Model
	errMsg     string
	msg        string
	data       domain.Data
}

func InitDataFieldsModel(data domain.Data) DataFieldsModel {
	m := DataFieldsModel{
		inputs: make([]textinput.Model, len(dataFields)),
		data:   data,
	}

	var t textinput.Model

	for i, n := range dataFields {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.Prompt = fmt.Sprintf("%-17s:  ", n.name)
		t.CharLimit = 32

		switch n.key {
		case nameFieldKey:
			t.Focus()
			t.SetValue(data.Name)
		case passFieldKey:
			t.Placeholder = PasswordPlaceholder
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
			t.SetValue(data.Pass)
		case cardNumFieldKey:
			t.CharLimit = 20
			t.Placeholder = "4505 **** **** 1234"
			t.Width = 30
			t.Validate = ccnValidator
			t.SetValue(data.CardNum)
		case fileFieldKey:
			t.CharLimit = 200
			t.SetValue(data.FilePath)
		case loginFieldKey:
			t.SetValue(data.Login)
		}

		m.inputs[i] = t
	}

	return m
}

func (m DataFieldsModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m DataFieldsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msgType := msg.(type) {
	case tea.KeyMsg:
		switch msgType.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "ctrl+w":
			var cmd tea.Cmd
			rm := RootModel{}
			return rm, tea.Batch(cmd, rm.Init())
		case "ctrl+u": // to user view
			var cmd tea.Cmd
			rm := UserModel{}
			return rm, tea.Batch(cmd, rm.Init())
		// save data
		case "ctrl+s":
			m.saveData()
		// to text view
		case "ctrl+t":
			dt := InitDataTextModel(m.getData())
			return dt, dt.Init()
		// to meta view
		case "ctrl+a":
			dt := initMetaModel(m.getData())
			return dt, dt.Init()
		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msgType.String()

			if s == "enter" {
				internal.Logger.Info("focusIndex", m.focusIndex, "inputsLen", m.getInputsCount())
			}
			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == m.getInputsCount() {
				var cmd tea.Cmd
				m.saveData()
				return m, cmd
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > m.getInputsCount() {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = m.getInputsCount()
			}

			cmds := make([]tea.Cmd, m.getInputsCount())
			state := m.setFocusedState()
			if state != nil {
				cmds[m.focusIndex] = m.setFocusedState()
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	return m.updateInputs(msg)
}

func (m DataFieldsModel) updateInputs(msg tea.Msg) (DataFieldsModel, tea.Cmd) {
	cmds := make([]tea.Cmd, m.getInputsCount())

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	counter := 0
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
		counter++
	}

	return m, tea.Batch(cmds...)
}

func (m DataFieldsModel) View() string {
	var b strings.Builder

	if len(m.errMsg) > 0 {
		b.WriteString(errorStyle.Render(m.errMsg) + "\n\n")
	}

	if len(m.msg) > 0 {
		b.WriteString(infoStyle.Render(m.msg) + "\n\n")
	}

	if m.data.ID != 0 {
		b.WriteString("Data ID: " + blueStyle.Render(strconv.FormatUint(m.data.ID, 10)) + "\n\n")
	}

	// simple input
	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		b.WriteRune('\n')
	}

	b.WriteRune('\n')

	button := &blurredButton
	if m.focusIndex == m.getInputsCount() {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "%s\n\n", *button)

	b.WriteString(actionsStyle.Render("'ctrl+t' to edit text window"))
	b.WriteRune('\n')
	b.WriteString(actionsStyle.Render("'ctrl+a' to edit meta window"))
	b.WriteRune('\n')
	b.WriteString(actionsStyle.Render("'ctrl+s' save data"))
	b.WriteRune('\n')
	b.WriteString(helpStyle.Render("'ctrl+w' to main window\n'ctrl-c' to quit"))

	return b.String()
}

func (m DataFieldsModel) getInputsCount() int {
	return len(m.inputs)
}

func (m DataFieldsModel) setFocusedState() tea.Cmd {
	var cmd tea.Cmd

	for i := 0; i < len(m.inputs); i++ {
		if m.focusIndex == i {
			cmd = m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
			m.inputs[i].PromptStyle = noStyle
			m.inputs[i].TextStyle = noStyle
		}
	}

	return cmd
}

func (m DataFieldsModel) getData() domain.Data {
	for i, v := range m.inputs {
		switch dataFields[i].key {
		case nameFieldKey:
			m.data.Name = v.Value()
		case loginFieldKey:
			m.data.Login = v.Value()
		case passFieldKey:
			m.data.Pass = v.Value()
		case cardNumFieldKey:
			m.data.CardNum = v.Value()
		case fileFieldKey:
			m.data.FilePath = strings.TrimSpace(v.Value())
		}
	}

	return m.data
}

func (m *DataFieldsModel) saveData() {
	id, err := saveData(m.getData())
	if err != nil {
		m.errMsg = err.Error()
	} else {
		m.data.ID = id
		m.msg = "data saved"
	}
}

// Validator functions to ensure valid input
func ccnValidator(s string) error {
	// Credit Card Number should a string less than 20 digits
	// It should include 16 integers and 3 spaces
	if len(s) > 16+3 {
		return fmt.Errorf("CCN is too long")
	}

	if len(s) == 0 || len(s)%5 != 0 && (s[len(s)-1] < '0' || s[len(s)-1] > '9') {
		return fmt.Errorf("CCN is invalid")
	}

	// The last digit should be a number unless it is a multiple of 4 in which
	// case it should be a space
	if len(s)%5 == 0 && s[len(s)-1] != ' ' {
		return fmt.Errorf("CCN must separate groups with spaces")
	}

	// The remaining digits should be integers
	c := strings.ReplaceAll(s, " ", "")
	_, err := strconv.ParseInt(c, 10, 64)

	return err
}
