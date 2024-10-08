package view

// View for add text data

import (
	"fmt"
	"gophkeeper/client/data"
	"gophkeeper/client/domain"
	"gophkeeper/internal"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type dataTextModel struct {
	textarea textarea.Model
	err      error
	data     domain.Data
	msg      string
}

// InitDataTextModel модель для добавления текстовых данных
func InitDataTextModel(d domain.Data) dataTextModel {
	ti := textarea.New()
	ti.Placeholder = "Add your text..."
	ti.SetValue(d.Text)
	ti.Focus()

	return dataTextModel{
		data:     d,
		textarea: ti,
		err:      nil,
	}
}

func (m dataTextModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m dataTextModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	var needUpdate = true

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		if msgType.Type == tea.KeyRunes && (len(msgType.Runes) > 1 || msgType.String() == "alt+\\" || msgType.String() == "alt+]") {
			needUpdate = false
		}

		internal.Logger.Info("type", msgType.Type, msgType.Type == tea.KeyCtrlM)
		internal.Logger.Info("string", msgType.String())
		switch msgType.Type {
		case tea.KeyEsc:
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		// to data list
		case tea.KeyCtrlL:
			dt := InitDataListModel()
			return dt, dt.Init()
		// to meta view
		case tea.KeyCtrlA:
			dt := initMetaModel(m.getData())
			return dt, dt.Init()
		case tea.KeyCtrlD:
			dt := InitDataFieldsModel(m.getData())
			return dt, dt.Init()
		// save data
		case tea.KeyCtrlS:
			m.saveData()
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msgType
		return m, nil
	}

	if needUpdate {
		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m dataTextModel) View() string {
	var b strings.Builder

	if m.err != nil {
		b.WriteString(errorStyle.Render(m.err.Error()) + "\n\n")
	}

	if len(m.msg) > 0 {
		b.WriteString(infoStyle.Render(m.msg) + "\n\n")
	}

	_, err := fmt.Fprintf(
		&b,
		"%s\n\n%s\n%s\n\n",
		showData(m.data),
		cursorStyle.Render(textFieldName),
		m.textarea.View(),
	)

	b.WriteString(actionsStyle.Render("'ctrl+a' to meta window"))
	b.WriteRune('\n')
	b.WriteString(actionsStyle.Render("'ctrl+d' to edit data window"))
	b.WriteRune('\n')
	b.WriteString(actionsStyle.Render("'ctrl+s' save data"))
	b.WriteRune('\n')
	b.WriteString(actionsStyle.Render("'ctrl+l' to data list"))
	b.WriteRune('\n')
	b.WriteString(helpStyle.Render("'ctrl+w' to main window\n'ctrl-c' to quit"))

	if err != nil {
		internal.Logger.Fatalw("err while updating text", "err", err)
	}

	return b.String()
}

func (m *dataTextModel) saveData() {
	gotData, err := data.SaveData(m.getData())

	if err != nil {
		m.err = err
	} else {
		m.data.ID = gotData.ID
		m.data.Version = gotData.Version
		m.data.FileID = gotData.FileID
		m.msg = "data saved"
	}
}

func (m dataTextModel) getData() domain.Data {
	m.data.Text = m.textarea.Value()
	return m.data
}
