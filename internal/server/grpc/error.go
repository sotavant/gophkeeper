package grpc

import (
	"errors"
	"gophkeeper/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func getError(err error) error {
	switch {
	case errors.Is(err, domain.ErrUserIDAbsent):
		return status.Error(codes.Unauthenticated, "user id absent")
	case errors.Is(err, domain.ErrBadData), errors.Is(err, domain.ErrDataVersionAbsent):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrInternalServerError), errors.Is(err, domain.ErrDataInsert), errors.Is(err, domain.ErrDataUpdate):
		return status.Error(codes.Internal, err.Error())
	case errors.Is(err, domain.ErrLoginExist):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrDataOutdated):
		return status.Error(codes.FailedPrecondition, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
