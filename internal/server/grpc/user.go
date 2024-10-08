package grpc

import (
	"context"
	"gophkeeper/server/domain"
	"gophkeeper/server/user"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gophkeeper/internal"
)
import pb "gophkeeper/proto"

// UserServer структура обеспечивающая регистрация/авторизацию пользвателя
type UserServer struct {
	pb.UnimplementedUserServiceServer
	Service *user.Service
}

func NewUserServer(s *user.Service) *UserServer {
	return &UserServer{
		Service: s,
	}
}

// Register регистрация пользвателя
func (u *UserServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	ur := &userRequest{}
	if err := ur.Bind(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := u.Service.Register(ctx, ur.User)
	if err != nil {
		return nil, getError(err)
	}

	return &pb.RegisterResponse{
		Token: token,
		Error: "",
	}, nil
}

// Login авторизация пользвателя
func (u *UserServer) Login(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	ur := &userRequest{}
	if err := ur.Bind(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := u.Service.Login(ctx, ur.User)
	if err != nil {
		return nil, getError(err)
	}

	return &pb.RegisterResponse{
		Token: token,
		Error: "",
	}, nil
}

type userRequest struct {
	domain.User
}

// Bind отображения в данных пользователя из запроса в модель сервера
func (u *userRequest) Bind(req *pb.RegisterRequest) error {
	v, err := protovalidate.New()
	if err != nil {
		internal.Logger.Fatalw("failed to initialize validator", "err", err)
	}

	if err = v.Validate(req.User); err != nil {
		internal.Logger.Errorw("user validation error", "err", err)
		return err
	}

	u.Password = req.User.Password
	u.Login = req.User.Login

	return nil
}
