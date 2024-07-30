package data

import (
	"context"
	"gophkeeper/domain"
)

type Service struct {
	dataRepo DataRepository
}

type DataRepository interface {
	Insert(ctx context.Context, data *domain.Data) error
	Update(ctx context.Context, data domain.Data) error
	GetVersion(ctx context.Context, id int64) (int64, error)
}

func NewService(d DataRepository) Service {
	return Service{
		dataRepo: d,
	}
}

func (s Service) UpsertData(ctx context.Context, data *domain.Data) error {
	if data.ID == 0 {
		data.Version = getVersion()
		err := s.dataRepo.Insert(ctx, data)
		if err != nil {
			return domain.ErrDataInsert
		}
	} else {
		err := s.updateVersion(ctx, data)
		if err != nil {
			return err
		}

		err = s.dataRepo.Update(ctx, *data)
		if err != nil {
			return domain.ErrDataUpdate
		}
	}

	return nil
}

func (s Service) updateVersion(ctx context.Context, data *domain.Data) error {
	if data.Version == 0 {
		return domain.ErrDataVersionAbsent
	}

	currentVersion, err := s.dataRepo.GetVersion(ctx, data.ID)
	if err != nil {
		return domain.ErrInternalServerError
	}

	if currentVersion != data.Version {
		return domain.ErrDataOutdated
	}

	data.Version = getVersion()

	return nil
}
