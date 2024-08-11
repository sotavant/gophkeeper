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
	ErrGetDataList            = errors.New("error in get data list")
	ErrGetData                = errors.New("error in get data")
	ErrDataNotFound           = errors.New("data not found")
	ErrCreationFileSaveDir    = errors.New("error in creation save dir")
	ErrDownloadFile           = errors.New("error in download file request")
	ErrDeleteData             = errors.New("error in delete request")
)
