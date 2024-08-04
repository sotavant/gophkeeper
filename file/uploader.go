package file

import (
	"bytes"
	"gophkeeper/domain"
	"gophkeeper/internal"
	"os"
	"path/filepath"
)

type Uploader struct {
	FilePath   string
	buffer     *bytes.Buffer
	OutputFile *os.File
	SavePath   string
}

func NewUploader(savePath string) *Uploader {
	return &Uploader{
		SavePath: savePath,
	}
}

func (u *Uploader) SetFile(fileName, path string) error {
	savePath := filepath.Join(u.SavePath, path)
	err := os.MkdirAll(savePath, os.ModePerm)
	if err != nil {
		internal.Logger.Infow("err in create directory", "err", err)
		return domain.ErrInternalServerError
	}

	u.FilePath = filepath.Join(savePath, fileName)
	file, err := os.Create(u.FilePath)
	if err != nil {
		internal.Logger.Infow("err in create file", "err", err)
		return domain.ErrInternalServerError
	}

	u.OutputFile = file
	return nil
}

func (u *Uploader) Write(chunk []byte) error {
	if u.OutputFile == nil {
		return nil
	}

	_, err := u.OutputFile.Write(chunk)
	return err
}

func (u *Uploader) Close() error {
	if u.OutputFile != nil {
		return u.OutputFile.Close()
	}

	return nil
}
