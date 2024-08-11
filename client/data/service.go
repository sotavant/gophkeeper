// Package data пакет служит для взаимодействия cli приложение с сервером
package data

import (
	"context"
	"errors"
	"gophkeeper/client/domain"
	"gophkeeper/internal"
	"gophkeeper/internal/client"
	"gophkeeper/internal/client/workers/grpc/interceptors"
	"gophkeeper/internal/crypto"
	domain2 "gophkeeper/server/domain"
	"os"
	"path/filepath"
	"strconv"
)

// SaveData сохранение данных на сервере
// Если в данных имеется файл, то дополнительным запросом происходит его сохранение
// На сервер отправляются зашифрованне паролем пользователя данные
func SaveData(data domain.Data) (domain.Data, error) {
	var encryptedFilePath string

	ctx := context.WithValue(context.Background(), interceptors.ContextUserTokenKey{}, client.AppInstance.User.Token)
	// hash data
	hashedData, err := encryptData(data)
	if err != nil {
		return data, domain.ErrEncryptData
	}

	// save data
	err = client.AppInstance.DataClient.SaveData(ctx, hashedData)

	if err != nil {
		return data, err
	}

	data.ID = hashedData.ID
	data.Version = hashedData.Version

	client.AppInstance.DecryptedData[data.ID] = data

	// upload file
	if data.FilePath != "" {
		encryptedFilePath, err = encryptFile(data.FilePath)
		err = client.AppInstance.DataClient.UploadFile(ctx, &data, encryptedFilePath, filepath.Base(data.FilePath))
		if err != nil {
			return data, err
		}

		data.FileName = filepath.Base(data.FilePath)
		data.Version = hashedData.Version

		client.AppInstance.DecryptedData[data.ID] = data
	}

	return data, nil
}

// GetData получение данных с сервера
// после получения, данные раскодируются паролем пользователя
func GetData(id uint64) (*domain.Data, error) {
	var err error
	data, ok := client.AppInstance.DecryptedData[id]
	ctx := context.WithValue(context.Background(), interceptors.ContextUserTokenKey{}, client.AppInstance.User.Token)

	if !ok {
		var gotData, decrypted *domain.Data
		gotData, err = client.AppInstance.DataClient.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		if gotData == nil {
			return nil, domain.ErrDataNotFound
		}

		decrypted, err = decryptData(*gotData)
		data = *decrypted
	}

	client.AppInstance.DecryptedData[id] = data

	return &data, nil
}

// GetDataList получить список данных пользователя в кратком формате (ID, Name)
func GetDataList() ([]domain2.DataName, error) {
	ctx := context.WithValue(context.Background(), interceptors.ContextUserTokenKey{}, client.AppInstance.User.Token)

	list, err := client.AppInstance.DataClient.GetList(ctx)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// DownloadFile скачать файл пользователя с сервера
// после скачивания файл раскодируется
func DownloadFile(data domain.Data) (string, error) {
	ctx := context.WithValue(context.Background(), interceptors.ContextUserTokenKey{}, client.AppInstance.User.Token)

	dataSavePath := filepath.Join(client.AppInstance.DataSavePath, client.AppInstance.User.Login, strconv.FormatUint(data.ID, 10))
	tmpSavePath, err := os.MkdirTemp(filepath.FromSlash("/tmp"), client.AppInstance.User.Login)
	if err != nil {
		return "", errors.New("cannot create temporary directory")
	}

	tmpFilePath, err := client.AppInstance.DataClient.DownloadFile(ctx, data, tmpSavePath, data.FileName)
	if err != nil {
		return tmpFilePath, err
	}

	return decryptFile(tmpFilePath, filepath.Join(dataSavePath, data.FileName))
}

// DeleteData удалить данные
func DeleteData(id uint64) error {
	ctx := context.WithValue(context.Background(), interceptors.ContextUserTokenKey{}, client.AppInstance.User.Token)
	err := client.AppInstance.DataClient.DeleteData(ctx, id)

	delete(client.AppInstance.DecryptedData, id)

	if err != nil {
		return err
	}

	return nil
}

func encryptData(data domain.Data) (*domain.Data, error) {
	var pass, login, cardNum, text, meta string
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
		Version: data.Version,
		ID:      data.ID,
		Name:    data.Name,
		Pass:    pass,
		CardNum: cardNum,
		Text:    text,
		Login:   login,
		Meta:    meta,
	}

	return hashedData, nil
}

func decryptData(data domain.Data) (*domain.Data, error) {
	var pass, login, cardNum, text, meta string
	var err error
	var decryptedData *domain.Data

	if data.Pass != "" {
		pass, err = crypto.Decrypt(client.AppInstance.User.StorageKey, data.Pass)
		if err != nil {
			return nil, err
		}
	}

	if data.Login != "" {
		login, err = crypto.Decrypt(client.AppInstance.User.StorageKey, data.Login)
		if err != nil {
			return nil, err
		}
	}

	if data.CardNum != "" {
		cardNum, err = crypto.Decrypt(client.AppInstance.User.StorageKey, data.CardNum)
		if err != nil {
			return nil, err
		}
	}

	if data.Text != "" {
		text, err = crypto.Decrypt(client.AppInstance.User.StorageKey, data.Text)
		if err != nil {
			return nil, err
		}
	}

	if data.Meta != "" {
		meta, err = crypto.Decrypt(client.AppInstance.User.StorageKey, data.Meta)
		if err != nil {
			return nil, err
		}
	}

	decryptedData = &domain.Data{
		Version:  data.Version,
		ID:       data.ID,
		Name:     data.Name,
		Pass:     pass,
		CardNum:  cardNum,
		Text:     text,
		Login:    login,
		Meta:     meta,
		FileName: data.FileName,
		FileID:   data.FileID,
	}

	return decryptedData, nil
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

	_, err = f.Write([]byte(cryptedText))

	return f.Name(), nil
}

func decryptFile(inputFile, outputFile string) (string, error) {
	text, err := os.ReadFile(inputFile)
	if err != nil {
		internal.Logger.Errorw("error reading file", "error", err)
		return "", domain.ErrReadingFile
	}

	cryptedText, err := crypto.Decrypt(client.AppInstance.User.StorageKey, string(text))
	if err != nil {
		internal.Logger.Errorw("error encrypting file", "error", err)
		return "", domain.ErrEncryptData
	}

	err = os.MkdirAll(filepath.Dir(outputFile), 0755)
	if err != nil {
		internal.Logger.Errorw("error creating output directory", "error", err)
		return "", domain.ErrEncryptData
	}

	f, err := os.Create(outputFile)
	if err != nil {
		internal.Logger.Errorw("error open output file", "error", err)
		return "", domain.ErrEncryptData
	}

	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			internal.Logger.Fatalw("error creating temp file", "error", err)
		}
	}(f)

	_, err = f.Write([]byte(cryptedText))

	return f.Name(), nil
}
