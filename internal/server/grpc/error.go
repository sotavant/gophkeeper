package grpc

import (
	"errors"
	"gophkeeper/server/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func getError(err error) error {
	switch {
	case errors.Is(err, domain.ErrUserIDAbsent):
		return status.Error(codes.Unauthenticated, "user id absent")
	case
		errors.Is(err, domain.ErrBadData),
		errors.Is(err, domain.ErrDataVersionAbsent),
		errors.Is(err, domain.ErrDataNameNotUniq),
		errors.Is(err, domain.ErrBadFileID):
		return status.Error(codes.InvalidArgument, err.Error())
	case
		errors.Is(err, domain.ErrUserNotFound),
		errors.Is(err, domain.ErrDataNotFound),
		errors.Is(err, domain.ErrFileNotFound):
		return status.Error(codes.NotFound, err.Error())
	case
		errors.Is(err, domain.ErrInternalServerError),
		errors.Is(err, domain.ErrDataInsert),
		errors.Is(err, domain.ErrDataUpdate),
		errors.Is(err, domain.ErrCheckDataName):
		return status.Error(codes.Internal, err.Error())
	case errors.Is(err, domain.ErrLoginExist):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrDataOutdated):
		return status.Error(codes.FailedPrecondition, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
