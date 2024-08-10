package view

// View for add text data

import (
	"fmt"
	"gophkeeper/client/domain"
	"gophkeeper/internal"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type dataMetaModel struct {
	textarea textarea.Model
	err      error
	data     domain.Data
}

func initMetaModel(d domain.Data) dataMetaModel {
	ti := textarea.New()
	ti.Placeholder = "Meta name : Meta value"
	ti.SetValue(d.Text)
	ti.Focus()

	return dataMetaModel{
		data:     d,
		textarea: ti,
		err:      nil,
	}
}

func (m dataMetaModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m dataMetaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	var needUpdate = true

	switch msgType := msg.(type) {
	case tea.KeyMsg:
		if msgType.Type == tea.KeyRunes && (len(msgType.Runes) > 1 || msgType.String() == "alt+\\") {
			needUpdate = false
		}

		switch msgType.Type {
		case tea.KeyEsc:
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyCtrlD:
			dt := InitDataFieldsModel(m.getData())
			return dt, dt.Init()
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

func (m dataMetaModel) View() string {
	var b strings.Builder

	_, err := fmt.Fprintf(
		&b,
		"%s\n\n%s\n%s\n\n",
		showData(m.data),
		cursorStyle.Render(metaFieldName),
		m.textarea.View(),
	)

	b.WriteString(actionsStyle.Render("'ctrl+t' to text window"))
	b.WriteRune('\n')
	b.WriteString(actionsStyle.Render("'ctrl+d' to edit data window"))
	b.WriteRune('\n')
	b.WriteString(helpStyle.Render("'ctrl+w' to main window\n'ctrl-c' to quit"))

	if err != nil {
		internal.Logger.Fatalw("err while updating text", "err", err)
	}

	return b.String()
}

func (m dataMetaModel) getData() domain.Data {
	m.data.Meta = m.textarea.Value()
	return m.data
}
