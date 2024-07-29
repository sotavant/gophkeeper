package data

import (
	"context"
	"gophkeeper/domain"
)

type Service struct {
	dataRepo DataRepository
}

type DataRepository interface {
}

func NewService(d DataRepository) Service {
	return Service{
		dataRepo: d,
	}
}

func (s Service) UpsertData(ctx context.Context, data domain.Data) (int64, error) {
	return 0, nil
}
