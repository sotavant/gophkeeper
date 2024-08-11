package grpc

import (
	"context"
	"errors"
	clientDomain "gophkeeper/client/domain"
	"gophkeeper/internal"
	pb "gophkeeper/proto"
	"io"
	"os"
	"path/filepath"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DataClient struct {
	client pb.DataServiceClient
}

func NewDataClient(client pb.DataServiceClient) *DataClient {
	return &DataClient{
		client: client,
	}
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

	resp, err := c.client.SaveData(ctx, &pb.SaveDataRequest{
		Data: pbData,
	})

	if err != nil {
		if status.Code(err) == codes.Internal {
			internal.Logger.Errorw("error while saving data", "error", err)
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

func (c *DataClient) UploadFile(ctx context.Context, data *clientDomain.Data, encryptedFilePath string) error {
	buf := make([]byte, 1024)

	uploadedFile, err := os.Open(encryptedFilePath)
	if err != nil {
		internal.Logger.Errorw("error while opening encrypted file", "error", err)
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
			DataId:      dData.ID,
			DataVersion: dData.Version,
			FileName:    filepath.Base(uploadedFile.Name()),
			FileChunk:   chunk,
		})

		assert.NoError(t, err)
	}
	batchNum += 1

	_, err = stream.CloseAndRecv()
	assert.NoError(t, err)

	uploadedFilePath := filepath.Join("/tmp/uploaded", uploadedFile.GetSaveFileSubDir(dData), filepath.Base(tmpFile.Name()))
	if _, err = os.Stat(uploadedFilePath); err != nil {
		assert.NoError(t, err)
	}

	err = os.Remove(uploadedFilePath)
	assert.NoError(t, err)
}
