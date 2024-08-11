package file

import (
	"context"
	"gophkeeper/internal"
	domain2 "gophkeeper/server/domain"
	"os"
	"strconv"
)

type Service struct {
	repo FileRepository
}

type FileRepository interface {
	Insert(ctx context.Context, file *domain2.File) error
	Update(ctx context.Context, file *domain2.File) error
	Get(ctx context.Context, id uint64) (*domain2.File, error)
	Delete(ctx context.Context, id uint64) error
}

func NewService(repo FileRepository) *Service {
	return &Service{repo: repo}
}

// if update need delete old file
func (s *Service) Save(ctx context.Context, file *domain2.File) error {
	if file.ID == 0 {
		err := s.repo.Insert(ctx, file)
		if err != nil {
			internal.Logger.Infow("error while inserting file", "error", err)
			return domain2.ErrDataInsert
		}

		return nil
	}

	dbFile, err := s.repo.Get(ctx, file.ID)
	if err != nil {
		internal.Logger.Infow("error while getting file", "error", err)
		return domain2.ErrInternalServerError
	}

	if err = os.Remove(dbFile.Path); err != nil {
		internal.Logger.Infow("error while removing file", "error", err)
		return domain2.ErrInternalServerError
	}

	if err = s.repo.Update(ctx, file); err != nil {
		internal.Logger.Infow("error while updating file", "error", err)
		return domain2.ErrInternalServerError
	}

	return nil
}

func (s *Service) Get(ctx context.Context, id uint64) (*domain2.File, error) {
	file, err := s.repo.Get(ctx, id)
	if err != nil {
		internal.Logger.Infow("error while getting file", "error", err)
		return nil, domain2.ErrInternalServerError
	}

	if file == nil {
		return nil, domain2.ErrFileNotFound
	}

	return file, nil
}

func (s *Service) Delete(ctx context.Context, id uint64) error {
	file, err := s.repo.Get(ctx, id)
	if err != nil {
		internal.Logger.Infow("error while getting file", "error", err)
		return domain2.ErrInternalServerError
	}

	if file == nil {
		return domain2.ErrFileNotFound
	}

	if err = os.Remove(file.Path); err != nil {
		internal.Logger.Infow("error while removing file", "error", err)
		return domain2.ErrInternalServerError
	}

	if err = s.repo.Delete(ctx, file.ID); err != nil {
		internal.Logger.Infow("error while deleting file", "error", err)
		return domain2.ErrInternalServerError
	}

	return nil
}

func GetSaveFileSubDir(data domain2.Data) string {
	return "/" + strconv.FormatUint(data.UID, 10) + "/" + strconv.FormatUint(data.ID, 10)
}
