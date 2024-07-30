package grpc

import (
	"context"
	"gophkeeper/data"
	"gophkeeper/domain"
	"gophkeeper/internal"
	pb "gophkeeper/proto"
	"gophkeeper/user"

	"github.com/bufbuild/protovalidate-go"
)

type DataServer struct {
	pb.UnimplementedDataServiceServer
	Service data.Service
}

func NewDataServer(s data.Service) DataServer {
	return DataServer{
		Service: s,
	}
}

func (s *DataServer) SaveData(ctx context.Context, req *pb.SaveDataRequest) (*pb.SaveDataResponse, error) {
	ur := &dataRequest{}
	if err := ur.Bind(req); err != nil {
		return nil, getError(err)
	}

	ur.Data.UID = ctx.Value(user.ContextUserIDKey{}).(int64)
	if ur.Data.UID == 0 {
		return nil, getError(domain.ErrUserIDAbsent)
	}

	id, err := s.Service.UpsertData(ctx, ur.Data)
	if err != nil {
		return nil, getError(err)
	}

	return &pb.SaveDataResponse{
		DataId: id,
	}, nil
}

type dataRequest struct {
	domain.Data
}

func (d *dataRequest) Bind(req *pb.SaveDataRequest) error {
	v, err := protovalidate.New()
	if err != nil {
		internal.Logger.Fatalw("failed to initialize validator", "err", err)
	}

	if err = v.Validate(req.Data); err != nil {
		internal.Logger.Errorw("user validation error", "err", err)
		return domain.ErrBadData
	}

	reqData := req.GetData()
	d.ID = reqData.GetId()
	d.Login = reqData.GetLogin()
	d.Pass = reqData.GetPass()
	d.Text = reqData.GetText()
	d.Meta = reqData.GetMeta()
	d.CardNum = reqData.GetCardNum()

	return nil
}
