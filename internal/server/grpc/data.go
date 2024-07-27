package grpc

import (
	"context"
	"gophkeeper/data"
	pb "gophkeeper/proto"
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

}
