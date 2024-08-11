package file

import (
	"context"
	"gophkeeper/internal"
	"gophkeeper/internal/server/repository/pgsql"
	"gophkeeper/internal/test"
	"gophkeeper/server/domain"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestService_Save(t *testing.T) {
	testUsersTable := "u_s"
	testFileTable := "f_s"
	testDataTable := "d_s"

	var fileRepo *pgsql.FileRepository
	ctx := context.Background()
	internal.InitLogger()
	pool, err := test.InitConnection(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, pool, "no databases init")

	defer func(ctx context.Context, pool *pgxpool.Pool) {
		err = test.CleanData(ctx, pool, []string{testUsersTable, testDataTable, testFileTable})
		assert.NoError(t, err)
	}(ctx, pool)

	tmpFile, err := os.CreateTemp("/tmp", "test")
	assert.NoError(t, err)

	file := &domain.File{
		Name: tmpFile.Name(),
		Path: tmpFile.Name(),
	}

	fileRepo, err = pgsql.NewFileRepository(ctx, pool, testFileTable)
	assert.NoError(t, err)

	err = fileRepo.Insert(ctx, file)
	assert.NoError(t, err)

	service := NewService(fileRepo)

	tests := []struct {
		name    string
		file    *domain.File
		wantErr error
	}{
		{
			name: "new file",
			file: &domain.File{
				Name: "ddd",
				Path: "sss",
			},
			wantErr: nil,
		},
		{
			name: "update file",
			file: &domain.File{
				ID:   file.ID,
				Name: "ddd",
				Path: "sss",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = service.Save(ctx, tt.file)
			assert.NoError(t, err)
		})
	}
}
