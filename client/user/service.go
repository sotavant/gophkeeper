package user

import (
	"gophkeeper/client/domain"
	"gophkeeper/internal/client"
	"unicode/utf8"
)

const (
	PasswordMinLen = 6
	LoginMinLen    = 2
)

type Service struct{}

func Auth(login, pass string, isLogin bool) error {
	var token string

	err := validateRegisterCredential(login, pass)
	if err != nil {
		return err
	}

	if isLogin {
		token, err = client.AppInstance.UserClient.Login(login, pass)
	} else {
		token, err = client.AppInstance.UserClient.Registration(login, pass)
	}
	if err != nil {
		return err
	}

	client.AppInstance.User.Token = token
	client.AppInstance.User.Login = login
	client.AppInstance.SetStorageKey(login, pass)

	return nil
}

func ResetUser() {
	client.AppInstance.User.Login = ""
	client.AppInstance.User.Token = ""
}

func validateRegisterCredential(login, pass string) error {
	if utf8.RuneCountInString(login) < LoginMinLen || utf8.RuneCountInString(pass) < PasswordMinLen {
		return domain.ErrRegisterDataLength
	}

	return nil
}
