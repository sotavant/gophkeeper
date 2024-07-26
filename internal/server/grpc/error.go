package grpc

import (
	"errors"
	"gophkeeper/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func getError(err error) error {
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrInternalServerError):
		return status.Error(codes.Internal, err.Error())
	case errors.Is(err, domain.ErrLoginExist):
		return status.Error(codes.AlreadyExists, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
