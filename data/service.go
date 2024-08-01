package data

import (
	"context"
	"gophkeeper/domain"
	"gophkeeper/internal"
)

type Service struct {
	dataRepo Repository
}

type Repository interface {
	Insert(ctx context.Context, data *domain.Data) error
	Update(ctx context.Context, data domain.Data) error
	GetById(ctx context.Context, id int64, selectFields []string) (*domain.Data, error)
	GetByNameAndUserID(ctx context.Context, uid int64, name string) (int64, error)
}

func NewService(d Repository) *Service {
	return &Service{
		dataRepo: d,
	}
}

func (s Service) UpsertData(ctx context.Context, data *domain.Data) error {
	if data.ID == nil {
		uniq, err := s.checkName(ctx, data, nil)
		if err != nil {
			internal.Logger.Errorw("error while checking name", "err", err)
			return domain.ErrCheckDataName
		}

		if !uniq {
			return domain.ErrDataNameNotUniq
		}

		version := getVersion()
		data.Version = &version

		err = s.dataRepo.Insert(ctx, data)
		if err != nil {
			return domain.ErrDataInsert
		}
	} else {
		oldRow, err := s.dataRepo.GetById(ctx, *data.ID, []string{"version", "name"})
		err = s.updateVersion(oldRow, data)
		if err != nil {
			return err
		}

		uniq, err := s.checkName(ctx, data, oldRow)
		if err != nil {
			internal.Logger.Errorw("error while checking name", "err", err)
			return domain.ErrCheckDataName
		}

		if !uniq {
			return domain.ErrDataNameNotUniq
		}

		err = s.dataRepo.Update(ctx, *data)
		if err != nil {
			return domain.ErrDataUpdate
		}
	}

	return nil
}

func (s Service) updateVersion(oldRow *domain.Data, newRow *domain.Data) error {
	if newRow.Version == nil {
		return domain.ErrDataVersionAbsent
	}

	if *oldRow.Version != *newRow.Version {
		return domain.ErrDataOutdated
	}

	version := getVersion()
	newRow.Version = &version

	return nil
}

// new: check by name and userId
// update: get old row, if name changed find row with same name and userId
func (s Service) checkName(ctx context.Context, data *domain.Data, oldData *domain.Data) (uniq bool, err error) {
	var id int64

	if data.ID != nil && data.Name == oldData.Name {
		uniq = true
		return
	}

	id, err = s.dataRepo.GetByNameAndUserID(ctx, *data.UID, data.Name)
	if err != nil {
		return
	}

	if id != 0 {
		uniq = false
	}

	uniq = true

	return
}
