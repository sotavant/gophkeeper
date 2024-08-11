package data

import (
	"context"
	"gophkeeper/client/domain"
	"gophkeeper/internal"
	"gophkeeper/internal/client"
	"gophkeeper/internal/crypto"
	"os"
	"path/filepath"
)

func SaveData(data domain.Data) (uint64, error) {
	var encryptedFilePath string

	ctx := context.Background()
	// hash data
	hashedData, err := encryptData(data)
	if err != nil {
		return 0, domain.ErrEncryptData
	}

	// save data
	err = client.AppInstance.DataClient.SaveData(ctx, hashedData)
	if err != nil {
		return 0, err
	}

	// upload file
	if data.FilePath != "" {
		encryptedFilePath, err = encryptFile(data.FilePath)

	}

	// set file name
}

func encryptData(data domain.Data) (*domain.Data, error) {
	var pass, login, cardNum, text, meta []byte
	var err error
	var hashedData *domain.Data

	if data.Pass != "" {
		pass, err = crypto.Encrypt(client.AppInstance.User.StorageKey, []byte(data.Pass))
		if err != nil {
			return nil, err
		}
	}

	if data.Login != "" {
		login, err = crypto.Encrypt(client.AppInstance.User.StorageKey, []byte(data.Login))
		if err != nil {
			return nil, err
		}
	}

	if data.CardNum != "" {
		cardNum, err = crypto.Encrypt(client.AppInstance.User.StorageKey, []byte(data.CardNum))
		if err != nil {
			return nil, err
		}
	}

	if data.Text != "" {
		text, err = crypto.Encrypt(client.AppInstance.User.StorageKey, []byte(data.Text))
		if err != nil {
			return nil, err
		}
	}

	if data.Meta != "" {
		meta, err = crypto.Encrypt(client.AppInstance.User.StorageKey, []byte(data.Meta))
		if err != nil {
			return nil, err
		}
	}

	hashedData = &domain.Data{
		Name:    data.Name,
		Pass:    string(pass),
		CardNum: string(cardNum),
		Text:    string(text),
		Login:   string(login),
		Meta:    string(meta),
	}

	return hashedData, nil
}

func encryptFile(filePath string) (string, error) {
	text, err := os.ReadFile(filePath)
	if err != nil {
		internal.Logger.Errorw("error reading file", "error", err)
		return "", domain.ErrReadingFile
	}

	cryptedText, err := crypto.Encrypt(client.AppInstance.User.StorageKey, text)
	if err != nil {
		internal.Logger.Errorw("error encrypting file", "error", err)
		return "", domain.ErrEncryptData
	}

	f, err := os.CreateTemp(filepath.FromSlash("/tmp"), client.AppInstance.User.Login)
	if err != nil {
		internal.Logger.Errorw("error creating temp file", "error", err)
		return "", domain.ErrEncryptData
	}

	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			internal.Logger.Fatalw("error creating temp file", "error", err)
		}
	}(f)

	_, err = f.Write(cryptedText)

	return f.Name(), nil
}
