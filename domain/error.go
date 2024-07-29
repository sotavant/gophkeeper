package domain

import "errors"

var (
	ErrInternalServerError = errors.New("internal Server Error")
	ErrLoginExist          = errors.New("login is busy")
	ErrUserNotFound        = errors.New("user not found")
	ErrBadData             = errors.New("check data error")
)
