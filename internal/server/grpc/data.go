package grpc

import (
	"context"
	"gophkeeper/data"
	"gophkeeper/domain"
	"gophkeeper/internal"
	pb "gophkeeper/proto"
	"gophkeeper/user"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/protobuf/types/known/emptypb"
)

type DataServer struct {
	pb.UnimplementedDataServiceServer
	Service data.Service
}

func NewDataServer(s *data.Service) *DataServer {
	return &DataServer{
		Service: *s,
	}
}

func (s *DataServer) SaveData(ctx context.Context, req *pb.SaveDataRequest) (*pb.SaveDataResponse, error) {
	ur := &dataRequest{}
	if err := ur.Bind(req); err != nil {
		return nil, getError(err)
	}

	ctxUID := ctx.Value(user.ContextUserIDKey{}).(int64)
	ur.Data.UID = &ctxUID
	if *ur.Data.UID == 0 {
		return nil, getError(domain.ErrUserIDAbsent)
	}

	err := s.Service.UpsertData(ctx, &ur.Data)
	if err != nil {
		return nil, getError(err)
	}

	return &pb.SaveDataResponse{
		DataId: *ur.ID,
	}, nil
}

func (s *DataServer) GetData(ctx context.Context, req *pb.GetDataRequest) (*pb.GetDataResponse, error) {
	return nil, nil
}

func (s *DataServer) DeleteData(ctx context.Context, req *pb.DeleteDataRequest) (*pb.DeleteDataResponse, error) {
	return nil, nil
}

func (s *DataServer) GetDataList(ctx context.Context, empty *emptypb.Empty) (*pb.DataListResponse, error) {
	return nil, nil
}

func (s *DataServer) UploadFile(stream pb.DataService_UploadFileServer) error {
	return nil
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
	reqDataId := reqData.GetId()
	d.ID = &reqDataId
	d.Login = reqData.GetLogin()
	d.Pass = reqData.GetPass()
	d.Text = reqData.GetText()
	d.Meta = reqData.GetMeta()
	d.CardNum = reqData.GetCardNum()

	return nil
}
