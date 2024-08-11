package data

import (
	"context"
	"gophkeeper/internal"
	domain2 "gophkeeper/server/domain"
	"gophkeeper/server/file"
	"path/filepath"
)

type Service struct {
	DataRepo Repository
	FileRepo FileRepository
}

type Repository interface {
	Insert(ctx context.Context, data *domain2.Data) error
	Update(ctx context.Context, data domain2.Data) error
	Get(ctx context.Context, id uint64) (*domain2.Data, error)
	GetByUser(ctx context.Context, id uint64, uid uint64) (*domain2.Data, error)
	GetByNameAndUserID(ctx context.Context, uid uint64, name string) (uint64, error)
	SetFile(ctx context.Context, data domain2.Data) error
	GetList(ctx context.Context, uid uint64) ([]domain2.DataName, error)
	Delete(ctx context.Context, id uint64) error
}

type FileRepository interface {
	Get(ctx context.Context, id uint64) (*domain2.File, error)
}

func NewService(d Repository, fileRepo FileRepository) *Service {
	return &Service{
		DataRepo: d,
		FileRepo: fileRepo,
	}
}

func (s Service) UpsertData(ctx context.Context, data *domain2.Data) error {
	if data.ID == 0 {
		uniq, err := s.checkName(ctx, data, nil)
		if err != nil {
			internal.Logger.Errorw("error while checking name", "err", err)
			return domain2.ErrCheckDataName
		}

		if !uniq {
			return domain2.ErrDataNameNotUniq
		}

		data.Version = getVersion()

		err = s.DataRepo.Insert(ctx, data)
		if err != nil {
			internal.Logger.Errorw("error while inserting data", "err", err)
			return domain2.ErrDataInsert
		}
	} else {
		oldRow, err := s.DataRepo.Get(ctx, data.ID)
		if err != nil {
			internal.Logger.Errorw("error while fetching data", "id", data.ID, "err", err)
			return domain2.ErrInternalServerError
		}

		err = s.updateVersion(oldRow, data)
		if err != nil {
			return err
		}

		uniq, err := s.checkName(ctx, data, oldRow)
		if err != nil {
			internal.Logger.Errorw("error while checking name", "err", err)
			return domain2.ErrCheckDataName
		}

		if !uniq {
			return domain2.ErrDataNameNotUniq
		}

		err = s.DataRepo.Update(ctx, *data)
		if err != nil {
			internal.Logger.Errorw("error while updating data", "id", data.ID, "err", err)
			return domain2.ErrDataUpdate
		}
	}

	return nil
}

// check data isset and for correct user
// if isset fileId check it equal to database.data.fileId
// if isset fileId check exist
// check data version
func (s Service) CheckUploadFileData(ctx context.Context, data domain2.Data) error {
	d, err := s.DataRepo.Get(ctx, data.ID)
	if err != nil {
		internal.Logger.Errorw("error while fetching data", "id", data.ID, "err", err)
		return domain2.ErrInternalServerError
	}

	if d == nil || d.UID != data.UID {
		return domain2.ErrDataNotFound
	}

	if data.FileID != nil && *data.FileID != 0 {
		if d.FileID == nil || *d.FileID != *data.FileID {
			return domain2.ErrBadFileID
		}

		if _, err = s.FileRepo.Get(ctx, *data.FileID); err != nil {
			internal.Logger.Errorw("error while fetching file", "id", *data.FileID, "err", err)
			return domain2.ErrInternalServerError
		}
	}

	return nil
}

// if data.fileId - remove old file and update row
// if new file - save file, and save data
func (s Service) SaveDataFile(ctx context.Context, data *domain2.Data, filePath string, f file.Service) error {
	dFile := domain2.File{
		Name: filepath.Base(filePath),
		Path: filePath,
		ID:   *data.FileID,
	}

	if err := f.Save(ctx, &dFile); err != nil {
		return err
	}

	data.FileID = &dFile.ID

	oldRow, err := s.DataRepo.Get(ctx, data.ID)
	if err != nil {
		internal.Logger.Errorw("error while fetching data", "id", data.ID, "err", err)
		return domain2.ErrInternalServerError
	}

	err = s.updateVersion(oldRow, data)
	if err != nil {
		return err
	}

	err = s.DataRepo.SetFile(ctx, *data)
	if err != nil {
		internal.Logger.Errorw("error while updating data", "id", data.ID, "err", err)
		return domain2.ErrInternalServerError
	}

	return nil
}

func (s Service) GetList(ctx context.Context, uid uint64) (list []domain2.DataName, err error) {
	list, err = s.DataRepo.GetList(ctx, uid)
	if err != nil {
		internal.Logger.Errorw("error while fetching data", "uid", uid, "err", err)
		return list, domain2.ErrInternalServerError
	}

	return
}

func (s Service) Get(ctx context.Context, dataID uint64, uid uint64) (data *domain2.Data, err error) {
	data, err = s.DataRepo.GetByUser(ctx, dataID, uid)
	if err != nil {
		internal.Logger.Errorw("error while fetching data", "id", dataID, "err", err)
		return data, domain2.ErrInternalServerError
	}

	if data == nil {
		return data, domain2.ErrDataNotFound
	}

	return
}

func (s Service) Delete(ctx context.Context, dataID, uid uint64, fs file.Service) error {
	data, err := s.DataRepo.GetByUser(ctx, dataID, uid)
	if err != nil {
		internal.Logger.Errorw("error while fetching data", "id", dataID, "err", err)
		return domain2.ErrInternalServerError
	}

	if data == nil {
		return domain2.ErrDataNotFound
	}

	err = s.DataRepo.Delete(ctx, dataID)
	if err != nil {
		internal.Logger.Errorw("error while deleting data", "id", dataID, "err", err)
		return domain2.ErrInternalServerError
	}

	if data.FileID != nil && *data.FileID != 0 {
		err = fs.Delete(ctx, *data.FileID)
		if err != nil {
			internal.Logger.Errorw("error while deleting file", "id", *data.FileID, "err", err)
			return domain2.ErrInternalServerError
		}
	}

	return nil
}

func (s Service) updateVersion(oldRow *domain2.Data, newRow *domain2.Data) error {
	if newRow.Version == 0 {
		return domain2.ErrDataVersionAbsent
	}

	if oldRow.Version != newRow.Version {
		return domain2.ErrDataOutdated
	}

	newRow.Version = getVersion()

	return nil
}

// new: check by name and userId
// update: get old row, if name changed find row with same name and userId
func (s Service) checkName(ctx context.Context, data *domain2.Data, oldData *domain2.Data) (uniq bool, err error) {
	var id uint64

	if data.ID != 0 && data.Name == oldData.Name {
		uniq = true
		return
	}

	id, err = s.DataRepo.GetByNameAndUserID(ctx, data.UID, data.Name)
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
