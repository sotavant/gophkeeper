package view

import (
	"errors"
	"fmt"
	"gophkeeper/client/data"
	"gophkeeper/client/domain"
	"gophkeeper/client/user"
	domain2 "gophkeeper/domain"
	"gophkeeper/internal/client"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type DataListModel struct {
	cursor   int
	choice   string
	msg      string
	dataList []domain2.DataName
	errMsg   string
}

func InitDataListModel() DataListModel {
	var errmsg string

	dataList, err := data.GetDataList()
	if err != nil {
		errmsg = err.Error()
	}

	m := DataListModel{
		dataList: dataList,
		errMsg:   errmsg,
	}

	return m
}

func (m DataListModel) Init() tea.Cmd {
	return nil
}

func (m DataListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.cursor >= len(m.dataList) {
				m.cursor = 0
			}
		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.dataList) - 1
			}
		}
	}

	return m, nil
}

func (m DataListModel) Do() (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	dataId := m.dataList[m.cursor].ID
	data, err := data.GetData(dataId)
	if err != nil {
		m.errMsg = err.Error()
		return m, tea.Batch(cmd, m.Init())
	}

	if data == nil {
		m.errMsg = domain.ErrDataNotFound.Error()
		return m, tea.Batch(cmd, m.Init())
	}

	dt := InitDataFieldsModel(*data)

	return dt, tea.Batch(cmd, dt.Init())
}

func (m DataListModel) View() string {
	s := strings.Builder{}

	if len(client.AppInstance.User.Token) == 0 {
		getError(errors.New("sorry( auth data is empty"))
	}

	if len(m.msg) > 0 {
		s.WriteString(infoStyle.Render(m.msg) + "\n\n")
	}

	s.WriteString(strings.TrimSpace(actionsStyle.Render("Choose data and press 'enter' for go to view/edit\n\n")))

	for i := 0; i < len(m.dataList); i++ {
		if m.cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(fmt.Sprintf("dataID: %d, dataName: %s\n", m.dataList[i].ID, m.dataList[i].Name))
	}

	s.WriteString(helpStyle.Render("\n\n'ctrl+w' to main window"))
	s.WriteString("\n(press q to quit)\n")

	return s.String()
}
