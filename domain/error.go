package domain

import "errors"

var (
	ErrInternalServerError = errors.New("internal Server Error")
	ErrLoginExist          = errors.New("login is busy")
	ErrUserNotFound        = errors.New("user not found")
	ErrBadData             = errors.New("check data error")
	ErrUserIDAbsent        = errors.New("user id absent")
	ErrDataVersionAbsent   = errors.New("data version absent")
	ErrDataOutdated        = errors.New("data outdated")
	ErrDataInsert          = errors.New("data insert error")
	ErrDataUpdate          = errors.New("data update error")
)
