package domain

import "errors"

var (
	ErrRegisterDataLength     = errors.New("login or password too short")
	ErrSomethingWrong         = errors.New("something wrong")
	ErrRegisterRequest        = errors.New("err in register request")
	ErrSaveDataRequest        = errors.New("err in save data request")
	ErrSaveDataIdErrorRequest = errors.New("err in save data request")
)
