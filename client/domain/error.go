package domain

import "errors"

var (
	ErrRegisterDataLength     = errors.New("login or password too short")
	ErrSomethingWrong         = errors.New("something wrong")
	ErrRegisterRequest        = errors.New("error in register request")
	ErrSaveDataRequest        = errors.New("error in save data request")
	ErrSaveDataIdErrorRequest = errors.New("error in save data request")
	ErrEncryptData            = errors.New("error in encrypt data")
	ErrReadingFile            = errors.New("error in reading file")
	ErrUploadFile             = errors.New("error in upload file")
)
