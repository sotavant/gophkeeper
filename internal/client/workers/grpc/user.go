package grpc

import (
	"context"
	"gophkeeper/client/domain"
	pb "gophkeeper/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserClient struct {
	client pb.UserServiceClient
}

func NewUserClient(client pb.UserServiceClient) *UserClient {
	return &UserClient{
		client: client,
	}
}

func (c *UserClient) Registration(login, password string) (token string, err error) {
	var response *pb.RegisterResponse

	response, err = c.client.Register(context.Background(), &pb.RegisterRequest{
		User: &pb.User{
			Login:    login,
			Password: password,
		},
	})

	if err != nil {
		if status.Code(err) == codes.Internal {
			return "", domain.ErrRegisterRequest
		}

		return "", err
	}

	if len(response.Token) == 0 {
		return "", domain.ErrRegisterRequest
	}

	return response.Token, nil
}

func (c *UserClient) Login(login, password string) (token string, err error) {
	var response *pb.RegisterResponse

	response, err = c.client.Login(context.Background(), &pb.RegisterRequest{
		User: &pb.User{
			Login:    login,
			Password: password,
		},
	})

	if err != nil {
		if status.Code(err) == codes.Internal {
			return "", domain.ErrRegisterRequest
		}

		return "", err
	}

	if len(response.Token) == 0 {
		return "", domain.ErrRegisterRequest
	}

	return response.Token, nil
}
