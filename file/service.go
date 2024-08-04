package file

import (
	"context"
	"gophkeeper/domain"
	"gophkeeper/internal"
	"os"
	"strconv"
)

type Service struct {
	repo FileRepository
}

type FileRepository interface {
	Insert(ctx context.Context, file *domain.File) error
	Update(ctx context.Context, file *domain.File) error
	Get(ctx context.Context, id uint64) (*domain.File, error)
}

func NewService(repo FileRepository) *Service {
	return &Service{repo: repo}
}

// if update need delete old file
func (s *Service) Save(ctx context.Context, file *domain.File) error {
	if file.ID == 0 {
		err := s.repo.Insert(ctx, file)
		if err != nil {
			internal.Logger.Infow("error while inserting file", "error", err)
			return domain.ErrDataInsert
		}

		return nil
	}

	dbFile, err := s.repo.Get(ctx, file.ID)
	if err != nil {
		internal.Logger.Infow("error while getting file", "error", err)
		return domain.ErrInternalServerError
	}

	if err = os.Remove(dbFile.Path); err != nil {
		internal.Logger.Infow("error while removing file", "error", err)
		return domain.ErrInternalServerError
	}

	if err = s.repo.Update(ctx, file); err != nil {
		internal.Logger.Infow("error while updating file", "error", err)
		return domain.ErrInternalServerError
	}

	return nil
}

func GetSaveFileSubDir(data domain.Data) string {
	return "/" + strconv.FormatUint(data.UID, 10) + "/" + strconv.FormatUint(data.ID, 10)
}
