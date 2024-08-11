// Package domain в данном пакете представлены модели данных
package domain

type Data struct {
	ID,
	FileID,
	Version uint64
	Name,
	Pass,
	CardNum,
	Text,
	FilePath,
	FileName,
	Login,
	Meta string
}
