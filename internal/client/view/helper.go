package view

import (
	"errors"
	"fmt"
	"gophkeeper/client/domain"
	"os"
	"strings"
	"unicode/utf8"
)

func showData(d domain.Data) string {
	var res string

	if d.ID != 0 {
		res += fmt.Sprintf("%-17s:  %d\n", "ID", d.ID)
	}
	res += fmt.Sprintf("%-17s:  %s\n", nameFieldName, d.Name)
	res += fmt.Sprintf("%-17s:  %s\n", loginFieldName, d.Login)
	res += fmt.Sprintf("%-17s:  %s\n", passFieldName, strings.Repeat("*", utf8.RuneCountInString(d.Pass)))
	res += fmt.Sprintf("%-17s:  %s\n", cardNumFieldName, d.CardNum)
	res += fmt.Sprintf("%-17s:  %s\n", fileFieldName, d.FilePath)

	return res
}

func saveData(d domain.Data) (uint64, error) {
	if d.Name == "" {
		return 0, errors.New("name is required")
	}

	if d.FilePath != "" {
		if _, err := os.Stat(d.FilePath); errors.Is(err, os.ErrNotExist) {
			return 0, errors.New("file does not exist")
		}
	}

}
