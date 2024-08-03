package data

import (
	"context"
	"fmt"
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
	if data.ID == 0 {
		uniq, err := s.checkName(ctx, data, nil)
		if err != nil {
			internal.Logger.Errorw("error while checking name", "err", err)
			return domain.ErrCheckDataName
		}

		if !uniq {
			return domain.ErrDataNameNotUniq
		}

		data.Version = getVersion()

		err = s.dataRepo.Insert(ctx, data)
		if err != nil {
			fmt.Println(data)
			internal.Logger.Errorw("error while inserting data", "err", err)
			return domain.ErrDataInsert
		}
	} else {
		oldRow, err := s.dataRepo.GetById(ctx, data.ID, []string{})
		if err != nil {
			internal.Logger.Errorw("error while fetching data", "id", data.ID, "err", err)
			return domain.ErrInternalServerError
		}

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
			fmt.Println(data.UID, data.Name)
			internal.Logger.Errorw("error while updating data", "id", data.ID, "err", err)
			return domain.ErrDataUpdate
		}
	}

	return nil
}

func (s Service) updateVersion(oldRow *domain.Data, newRow *domain.Data) error {
	if newRow.Version == 0 {
		return domain.ErrDataVersionAbsent
	}

	if oldRow.Version != newRow.Version {
		return domain.ErrDataOutdated
	}

	newRow.Version = getVersion()

	return nil
}

// new: check by name and userId
// update: get old row, if name changed find row with same name and userId
func (s Service) checkName(ctx context.Context, data *domain.Data, oldData *domain.Data) (uniq bool, err error) {
	var id int64

	if data.ID != 0 && data.Name == oldData.Name {
		uniq = true
		return
	}

	id, err = s.dataRepo.GetByNameAndUserID(ctx, data.UID, data.Name)
	if err != nil {
		return
	}

	if id != 0 {
		uniq = false
		return
	}

	uniq = true

	return
}
