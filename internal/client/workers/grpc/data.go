package grpc

import (
	"context"
	"errors"
	clientDomain "gophkeeper/client/domain"
	domain2 "gophkeeper/domain"
	"gophkeeper/internal"
	pb "gophkeeper/proto"
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
	}

	return data, nil
}

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

	internal.Logger.Info(pbData)
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

	internal.Logger.Infow("saved data", "id", resp.GetDataId(), "version", resp.GetDataVersion())
	data.ID = resp.GetDataId()
	data.Version = resp.GetDataVersion()

	return nil
}

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

	return nil
}
