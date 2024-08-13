// Package grpc пакет для взаимодействия с сервером
package grpc

import (
	"context"
	"errors"
	clientDomain "gophkeeper/client/domain"
	"gophkeeper/internal"
	pb "gophkeeper/proto"
	domain2 "gophkeeper/server/domain"
	file2 "gophkeeper/server/file"
	"io"
	"os"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type DataClient struct {
	client pb.DataServiceClient
}

func NewDataClient(client pb.DataServiceClient) *DataClient {
	return &DataClient{
		client: client,
	}
}

// Get получение данных по ID
func (c *DataClient) Get(ctx context.Context, id uint64) (*clientDomain.Data, error) {
	req := &pb.GetDataRequest{Id: id}

	resp, err := c.client.GetData(ctx, req)

	if err != nil {
		if status.Code(err) == codes.Internal {
			internal.Logger.Errorw("error while get data", "error", err)
			return nil, clientDomain.ErrGetData
		}

		return nil, err
	}

	respData := resp.GetData()
	data := &clientDomain.Data{
		ID:       id,
		Version:  respData.GetVersion(),
		Name:     respData.GetName(),
		Pass:     respData.GetPass(),
		CardNum:  respData.GetCardNum(),
		Text:     respData.GetText(),
		FileName: respData.GetFileName(),
		Login:    respData.GetLogin(),
		Meta:     respData.GetMeta(),
		FileID:   respData.GetFileID(),
	}

	return data, nil
}

// GetList получения списка данных пользователя
func (c *DataClient) GetList(ctx context.Context) ([]domain2.DataName, error) {
	resp, err := c.client.GetDataList(ctx, &emptypb.Empty{})

	if err != nil {
		if status.Code(err) == codes.Internal {
			internal.Logger.Errorw("error while get data list", "error", err)
			return nil, clientDomain.ErrGetDataList
		}

		return nil, err
	}

	var dataList []domain2.DataName

	for _, data := range resp.GetDataList() {
		dd := domain2.DataName{
			Name: data.GetName(),
			ID:   data.GetId(),
		}

		dataList = append(dataList, dd)
	}

	return dataList, nil
}

// SaveData сохранение данных
func (c *DataClient) SaveData(ctx context.Context, data *clientDomain.Data) error {
	pbData := &pb.Data{
		Id:      data.ID,
		Name:    data.Name,
		Version: data.Version,
		Login:   data.Login,
		Pass:    data.Pass,
		Text:    data.Text,
		CardNum: data.CardNum,
		Meta:    data.Meta,
	}

	resp, err := c.client.SaveData(ctx, &pb.SaveDataRequest{
		Data: pbData,
	})

	if err != nil {
		if status.Code(err) == codes.Internal {
			internal.Logger.Errorw("error while send data ", "error", err)
			return clientDomain.ErrSaveDataRequest
		}

		return err
	}

	if resp.GetDataId() == 0 {
		return errors.New("data ID absent in response")
	}

	data.ID = resp.GetDataId()
	data.Version = resp.GetDataVersion()

	return nil
}

// UploadFile загрузка файла на сервер
func (c *DataClient) UploadFile(ctx context.Context, data *clientDomain.Data, encryptedFilePath, fileName string) error {
	var resp *pb.FileUploadResponse
	buf := make([]byte, 1024)

	uploadedFile, err := os.Open(encryptedFilePath)
	if err != nil {
		internal.Logger.Errorw("error while opening encrypted file", "error", err)
		return clientDomain.ErrUploadFile
	}

	stream, err := c.client.UploadFile(ctx)
	if err != nil {
		internal.Logger.Errorw("error while get stream", "error", err)
		return clientDomain.ErrUploadFile
	}

	for {
		var num int
		num, err = uploadedFile.Read(buf)
		if err == io.EOF {
			break
		}

		if err != nil {
			internal.Logger.Errorw("error while read encrypted file", "error", err)
			return clientDomain.ErrUploadFile
		}

		chunk := buf[:num]
		err = stream.Send(&pb.UploadFileRequest{
			DataId:      data.ID,
			DataVersion: data.Version,
			FileName:    fileName,
			FileChunk:   chunk,
		})

		if err != nil {
			internal.Logger.Errorw("error while send file stream", "error", err)
			return clientDomain.ErrUploadFile
		}
	}

	resp, err = stream.CloseAndRecv()

	if err != nil {
		internal.Logger.Errorw("error while receive file upload response", "error", err)
		return clientDomain.ErrUploadFile
	}

	data.Version = resp.GetDataVersion()
	data.FileID = resp.GetFileId()

	return nil
}

// DownloadFile скачать файл
func (c *DataClient) DownloadFile(ctx context.Context, data clientDomain.Data, filePath, fileName string) (string, error) {
	var rr *pb.DownloadFileResponse

	request := &pb.DownloadFileRequest{
		DataID: data.ID,
		FileID: data.FileID,
	}

	fileStreamResponse, err := c.client.DownloadFile(ctx, request)
	if err != nil {
		if status.Code(err) == codes.Internal {
			internal.Logger.Errorw("error while download file", "error", err)
			return "", clientDomain.ErrDownloadFile
		}
	}

	file := file2.NewUploader(filePath)
	defer func(file *file2.Uploader) {
		err = file.Close()
		if err != nil {
			internal.Logger.Fatalw("error while close file", "error", err)
		}
	}(file)

	var fileSize uint32 = 0

	for {
		rr, err = fileStreamResponse.Recv()
		if err == io.EOF || rr == nil {
			break
		}

		if err != nil {
			internal.Logger.Errorw("error while receive file download response", "error", err)
			return "", clientDomain.ErrDownloadFile
		}

		if file.FilePath == "" {
			if err = file.SetFile(fileName, ""); err != nil {
				internal.Logger.Errorw("error while create downloaded file", "error", err)
				return "", clientDomain.ErrDownloadFile
			}
		}

		chunk := rr.GetFileChunk()
		fileSize += uint32(len(chunk))
		if err = file.Write(chunk); err != nil {
			internal.Logger.Errorw("error while write downloaded file", "error", err)
			return "", clientDomain.ErrDownloadFile
		}
	}

	if fileSize == 0 {
		internal.Logger.Errorw("empty file")
		return "", clientDomain.ErrDownloadFile
	}

	return file.FilePath, nil
}

// DeleteData удалить данные
func (c *DataClient) DeleteData(ctx context.Context, id uint64) error {
	resp := &pb.DeleteDataRequest{Id: id}

	_, err := c.client.DeleteData(ctx, resp)
	if err != nil {
		if status.Code(err) == codes.Internal {
			internal.Logger.Errorw("error while delete data", "error", err)
			return clientDomain.ErrDeleteData
		}
	}

	return nil
}
