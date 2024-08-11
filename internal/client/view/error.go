package view

import (
	"errors"
	"gophkeeper/client/domain"
)

func getError(err error) string {
	switch {
	case errors.Is(err, domain.ErrRegisterDataLength):
		return errorStyle.Render(err.Error())
	default:
		return errorStyle.Render(err.Error())
	}
}
