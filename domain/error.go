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
	ErrCheckDataName       = errors.New("error in check data name")
	ErrDataNameNotUniq     = errors.New("data name not uniq")
	ErrDataNotFound        = errors.New("data not found")
	ErrBadFileID           = errors.New("bad file id")
	ErrFileNotFound        = errors.New("file not found")
)
