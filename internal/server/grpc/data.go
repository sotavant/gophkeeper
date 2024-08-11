package grpc

import (
	"context"
	"gophkeeper/data"
	"gophkeeper/domain"
	file2 "gophkeeper/file"
	"gophkeeper/internal"
	pb "gophkeeper/proto"
	"gophkeeper/user"
	"io"
	"os"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type DataServer struct {
	pb.UnimplementedDataServiceServer
	Service       data.Service
	FileService   file2.Service
	filesSavePath string
}

func NewDataServer(s *data.Service, filesSavePath string, f *file2.Service) *DataServer {
	return &DataServer{
		Service:       *s,
		filesSavePath: filesSavePath,
		FileService:   *f,
	}
}

func (s *DataServer) SaveData(ctx context.Context, req *pb.SaveDataRequest) (*pb.SaveDataResponse, error) {
	ur := &dataRequest{&domain.Data{}}
	if err := ur.Bind(ctx, req); err != nil {
		return nil, getError(err)
	}

	err := s.Service.UpsertData(ctx, ur.Data)
	if err != nil {
		return nil, getError(err)
	}

	return &pb.SaveDataResponse{
		DataId:      ur.ID,
		DataVersion: ur.Version,
	}, nil
}

func (s *DataServer) GetData(ctx context.Context, req *pb.GetDataRequest) (*pb.GetDataResponse, error) {
	var err error
	var dbFile *domain.File

	ctxUID := ctx.Value(user.ContextUserIDKey{}).(uint64)
	if ctxUID == 0 {
		return nil, domain.ErrUserIDAbsent
	}

	v, err := protovalidate.New()
	if err != nil {
		internal.Logger.Fatalw("failed to initialize validator", "err", err)
	}

	if err = v.Validate(req); err != nil {
		internal.Logger.Errorw("upload request validation error", "err", err)
		return nil, domain.ErrBadData
	}

	d, err := s.Service.Get(ctx, req.GetId(), ctxUID)
	if err != nil {
		return nil, getError(err)
	}

	if d.FileID != nil {
		dbFile, err = s.FileService.Get(ctx, *d.FileID)
		if err != nil {
			return nil, getError(err)
		}
	}

	return getDataResponse(*d, dbFile), nil
}

func (s *DataServer) DeleteData(ctx context.Context, req *pb.DeleteDataRequest) (*emptypb.Empty, error) {
	var err error

	ctxUID := ctx.Value(user.ContextUserIDKey{}).(uint64)
	if ctxUID == 0 {
		return nil, domain.ErrUserIDAbsent
	}

	v, err := protovalidate.New()
	if err != nil {
		internal.Logger.Fatalw("failed to initialize validator", "err", err)
	}

	if err = v.Validate(req); err != nil {
		internal.Logger.Errorw("upload request validation error", "err", err)
		return nil, domain.ErrBadData
	}

	err = s.Service.Delete(ctx, req.GetId(), ctxUID, s.FileService)

	return nil, nil
}

func (s *DataServer) GetDataList(ctx context.Context, empty *emptypb.Empty) (*pb.DataListResponse, error) {
	ctxUID := ctx.Value(user.ContextUserIDKey{}).(uint64)
	if ctxUID == 0 {
		return nil, domain.ErrUserIDAbsent
	}

	list, err := s.Service.GetList(ctx, ctxUID)
	if err != nil {
		return nil, getError(err)
	}

	if len(list) == 0 {
		return &pb.DataListResponse{}, nil
	}

	return getDataListResponse(list), nil
}

// test: file with same name
// validate data
// check row with id exist
// check version
// uploader -> fileService for save file -> dataService for add file
func (s *DataServer) UploadFile(stream pb.DataService_UploadFileServer) error {
	validated := false
	ur := &dataRequest{&domain.Data{}}

	var fileSize uint32 = 0
	file := file2.NewUploader(s.filesSavePath)

	defer func() {
		if err := file.Close(); err != nil {
			internal.Logger.Fatalw("error closing file", "error", err)
		}
	}()

	for {
		req, err := stream.Recv()
		if err == io.EOF || req == nil {
			break
		}

		if err != nil {
			return status.Errorf(codes.Internal, "error receiving data: %v", err)
		}

		if !validated {
			if err = ur.BindUploadFile(stream.Context(), req); err != nil {
				return getError(err)
			}

			if err = s.Service.CheckUploadFileData(stream.Context(), *ur.Data); err != nil {
				return getError(err)
			}

			validated = true
		}

		if file.FilePath == "" {
			var dir string
			dir = file2.GetSaveFileSubDir(*ur.Data)
			if err = file.SetFile(req.GetFileName(), dir); err != nil {
				return getError(err)
			}
		}

		chunk := req.GetFileChunk()
		fileSize += uint32(len(chunk))
		if err = file.Write(chunk); err != nil {
			internal.Logger.Infow("error in write chunk", "err", err)
			return status.Errorf(codes.Internal, "error in write file chunk")
		}
	}

	if fileSize == 0 {
		return status.Errorf(codes.InvalidArgument, "file size is zero")
	}

	if err := s.Service.SaveDataFile(stream.Context(), ur.Data, file.FilePath, s.FileService); err != nil {
		return getError(err)
	}

	return stream.SendAndClose(&pb.FileUploadResponse{
		FileId:      *ur.Data.FileID,
		DataVersion: ur.Data.Version,
		Size:        fileSize,
	})
}

func (s *DataServer) DownloadFile(req *pb.DownloadFileRequest, stream pb.DataService_DownloadFileServer) error {
	dr := &DownloadFileRequest{}
	if err := dr.BindDownloadFileRequest(stream.Context(), req); err != nil {
		return getError(err)
	}

	dbData, err := s.Service.Get(stream.Context(), dr.DataID, dr.UID)
	if err != nil {
		return getError(err)
	}

	if dbData.FileID == nil || *dbData.FileID != dr.FileID {
		return getError(domain.ErrFileNotFound)
	}

	file, err := s.FileService.Get(stream.Context(), dr.FileID)
	if err != nil {
		return getError(err)
	}

	if file == nil {
		return getError(domain.ErrFileNotFound)
	}

	bufferSize := 1024 * 1024
	var bytesRead int

	osFile, err := os.Open(file.Path)
	if err != nil {
		internal.Logger.Errorw("error opening file", "err", err)
		return status.Error(codes.Internal, "error opening file")
	}

	defer func(osFile *os.File) {
		err = osFile.Close()
		if err != nil {
			internal.Logger.Errorw("error closing file", "err", err)
		}
	}(osFile)

	buff := make([]byte, bufferSize)
	for {
		bytesRead, err = osFile.Read(buff)
		if err != nil {
			if err == io.EOF {
				break
			}

			internal.Logger.Errorw("error reading file", "err", err)
			return status.Error(codes.Internal, "error reading file")
		}

		resp := &pb.DownloadFileResponse{
			FileChunk: buff[:bytesRead],
		}

		if err = stream.Send(resp); err != nil {
			internal.Logger.Errorw("error sending response", "err", err)
			return status.Error(codes.Internal, "error sending response")
		}
	}

	return nil
}

type dataRequest struct {
	*domain.Data
}

func (d *dataRequest) BindUploadFile(ctx context.Context, req *pb.UploadFileRequest) error {
	ctxUID := ctx.Value(user.ContextUserIDKey{}).(uint64)
	if ctxUID == 0 {
		return domain.ErrUserIDAbsent
	}

	v, err := protovalidate.New()
	if err != nil {
		internal.Logger.Fatalw("failed to initialize validator", "err", err)
	}

	if err = v.Validate(req); err != nil {
		internal.Logger.Errorw("upload request validation error", "err", err)
		return domain.ErrBadData
	}

	fileId := req.GetFileId()
	d.ID = req.GetDataId()
	d.Version = req.GetDataVersion()
	d.FileID = &fileId
	d.UID = ctxUID

	return nil
}

func (d *dataRequest) Bind(ctx context.Context, req *pb.SaveDataRequest) error {
	ctxUID := ctx.Value(user.ContextUserIDKey{}).(uint64)
	if ctxUID == 0 {
		return domain.ErrUserIDAbsent
	}

	v, err := protovalidate.New()
	if err != nil {
		internal.Logger.Fatalw("failed to initialize validator", "err", err)
	}

	if err = v.Validate(req.Data); err != nil {
		internal.Logger.Errorw("user validation error", "err", err)
		return domain.ErrBadData
	}

	reqData := req.GetData()
	login := reqData.GetLogin()
	pass := reqData.GetPass()
	text := reqData.GetText()
	meta := reqData.GetMeta()
	cardNum := reqData.GetCardNum()

	d.Data.ID = reqData.GetId()
	d.Version = reqData.GetVersion()
	d.Name = reqData.GetName()
	d.Login = &login
	d.Pass = &pass
	d.Text = &text
	d.Meta = &meta
	d.CardNum = &cardNum
	d.UID = ctxUID

	return nil
}

func getDataResponse(data domain.Data, file *domain.File) *pb.GetDataResponse {
	fileName := ""
	if file != nil {
		fileName = file.Name
	}

	respData := &pb.Data{
		Id:       data.ID,
		Name:     data.Name,
		FileName: fileName,
		Version:  data.Version,
	}

	if data.Login != nil {
		respData.Login = *data.Login
	}

	if data.Pass != nil {
		respData.Pass = *data.Pass
	}

	if data.Text != nil {
		respData.Text = *data.Text
	}

	if data.CardNum != nil {
		respData.CardNum = *data.CardNum
	}

	if data.Meta != nil {
		respData.Meta = *data.Meta
	}

	return &pb.GetDataResponse{Data: respData}
}

func getDataListResponse(data []domain.DataName) *pb.DataListResponse {
	dataList := make([]*pb.DataList, len(data))

	for i, d := range data {
		dataList[i] = &pb.DataList{
			Name: d.Name,
			Id:   d.ID,
		}
	}

	return &pb.DataListResponse{
		DataList: dataList,
	}
}

type DownloadFileRequest struct {
	DataID uint64
	UID    uint64
	FileID uint64
}

func (d *DownloadFileRequest) BindDownloadFileRequest(ctx context.Context, req *pb.DownloadFileRequest) error {
	ctxUID := ctx.Value(user.ContextUserIDKey{}).(uint64)
	if ctxUID == 0 {
		return domain.ErrUserIDAbsent
	}

	v, err := protovalidate.New()
	if err != nil {
		internal.Logger.Fatalw("failed to initialize validator", "err", err)
	}

	if err = v.Validate(req); err != nil {
		internal.Logger.Errorw("user validation error", "err", err)
		return domain.ErrBadData
	}

	d.DataID = req.GetDataID()
	d.UID = ctxUID
	d.FileID = req.GetFileID()

	return nil
}
