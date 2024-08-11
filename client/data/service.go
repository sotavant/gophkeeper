package data

import (
	"context"
	"gophkeeper/client/domain"
	domain2 "gophkeeper/domain"
	"gophkeeper/internal"
	"gophkeeper/internal/client"
	"gophkeeper/internal/client/workers/grpc/interceptors"
	"gophkeeper/internal/crypto"
	"os"
	"path/filepath"
)

func SaveData(data domain.Data) (uint64, uint64, error) {
	var encryptedFilePath string

	ctx := context.WithValue(context.Background(), interceptors.ContextUserTokenKey{}, client.AppInstance.User.Token)
	// hash data
	hashedData, err := encryptData(data)
	if err != nil {
		return 0, 0, domain.ErrEncryptData
	}

	// save data
	err = client.AppInstance.DataClient.SaveData(ctx, hashedData)

	if err != nil {
		return 0, 0, err
	}

	data.ID = hashedData.ID
	data.Version = hashedData.Version

	client.AppInstance.DecryptedData[data.ID] = data

	// upload file
	if data.FilePath != "" {
		encryptedFilePath, err = encryptFile(data.FilePath)
		err = client.AppInstance.DataClient.UploadFile(ctx, &data, encryptedFilePath, filepath.Base(data.FilePath))
		if err != nil {
			return hashedData.ID, hashedData.Version, err
		}

		data.FileName = filepath.Base(data.FilePath)
		data.Version = hashedData.Version

		client.AppInstance.DecryptedData[data.ID] = data
	}

	return data.ID, data.Version, err
}

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

func GetDataList() ([]domain2.DataName, error) {
	ctx := context.WithValue(context.Background(), interceptors.ContextUserTokenKey{}, client.AppInstance.User.Token)

	list, err := client.AppInstance.DataClient.GetList(ctx)
	if err != nil {
		return nil, err
	}

	return list, nil
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
