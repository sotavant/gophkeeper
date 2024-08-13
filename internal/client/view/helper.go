package view

import (
	"fmt"
	"gophkeeper/client/domain"
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
	res += fmt.Sprintf("%-17s:  %s\n", fileNameFieldName, d.FileName)

	return res
}
