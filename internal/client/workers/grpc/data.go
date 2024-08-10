package grpc

import (
	"context"
	"errors"
	clientDomain "gophkeeper/client/domain"
	"gophkeeper/domain"
	"gophkeeper/internal"
	pb "gophkeeper/proto"

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

func (c *DataClient) SaveDate(ctx context.Context, data *domain.Data) error {
	resp, err := c.client.SaveData(ctx, &pb.SaveDataRequest{
		Data: &pb.Data{
			Id:      data.ID,
			Name:    data.Name,
			Login:   *data.Login,
			Pass:    *data.Pass,
			Text:    *data.Text,
			CardNum: *data.CardNum,
			Meta:    *data.Meta,
			Version: data.Version,
		},
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
	return nil
}
