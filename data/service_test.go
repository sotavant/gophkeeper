package data

import (
	"context"
	"errors"
	"gophkeeper/domain"
	"gophkeeper/internal"
	"gophkeeper/internal/server/repository/pgsql"
	"gophkeeper/internal/test"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestService_UpsertData(t *testing.T) {
	ctx := context.Background()
	internal.InitLogger()
	pool, err := test.InitConnection(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, pool, "no databases init")

	defer func(ctx context.Context, pool *pgxpool.Pool) {
		err = test.CleanData(ctx, pool, []string{test.UsersTestTable, test.DataTestTable, test.FileTestTable})
		assert.NoError(t, err)
	}(ctx, pool)

	userRepo, err := pgsql.NewUserRepository(ctx, pool, test.UsersTestTable)
	assert.NoError(t, err)

	user := &domain.User{
		Login:    "test",
		Password: "test",
	}

	userId, err := userRepo.Store(ctx, *user)
	assert.NoError(t, err)
	assert.NotZero(t, userId)

	service := GetTestService(ctx, t, pool)

	var versionFirst uint64 = 1
	var versionSecond uint64 = 2
	login := "test"
	successTextData := "success update"
	testData := &domain.Data{
		Name:    "testic",
		Login:   &login,
		Pass:    &login,
		Version: versionFirst,
		UID:     userId,
	}

	testData2 := &domain.Data{
		Name:    "testic2",
		Login:   &login,
		Pass:    &login,
		Version: versionSecond,
		UID:     userId,
	}

	err = service.DataRepo.Insert(ctx, testData)
	assert.NoError(t, err)

	err = service.DataRepo.Insert(ctx, testData2)
	assert.NoError(t, err)

	type want struct {
		err error
	}

	tests := []struct {
		name          string
		data          domain.Data
		want          want
		updateVersion bool
	}{
		{
			name: "success new data", // need check userId and version
			data: domain.Data{
				Login:   &login,
				Pass:    &login,
				CardNum: &login,
				Text:    &login,
				Name:    "test",
				Meta:    &login,
				UID:     userId,
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "name exist", // need check userId and version
			data: domain.Data{
				Login:   &login,
				Pass:    &login,
				CardNum: &login,
				Text:    &login,
				Name:    testData.Name,
				Meta:    &login,
				UID:     userId,
			},
			want: want{
				err: domain.ErrDataNameNotUniq,
			},
		},
		{
			name: "success update",
			data: domain.Data{
				ID:      testData.ID,
				Text:    &successTextData,
				Version: testData.Version,
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "wrong update name exist",
			data: domain.Data{
				ID:      testData.ID,
				UID:     testData.UID,
				Name:    testData2.Name,
				Version: testData.Version,
			},
			want: want{
				err: domain.ErrDataNameNotUniq,
			},
			updateVersion: true,
		},
		{
			name: "wrong update version absent",
			data: domain.Data{
				ID:   testData.ID,
				Name: testData.Name,
			},
			want: want{
				err: domain.ErrDataVersionAbsent,
			},
		},
		{
			name: "wrong update bad version",
			data: domain.Data{
				ID:      testData.ID,
				Name:    testData.Name,
				Version: versionSecond,
			},
			want: want{
				err: domain.ErrDataOutdated,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.updateVersion {
				var dd *domain.Data
				dd, err = service.DataRepo.Get(ctx, tt.data.ID)
				assert.NoError(t, err)
				tt.data.Version = dd.Version
			}

			err = service.UpsertData(ctx, &tt.data)
			if tt.want.err == nil {
				assert.NoError(t, err)
			}

			if !errors.Is(err, tt.want.err) {
				t.Errorf("UpsertData() error = %v, wantErr %v", err, tt.want.err)
			}
		})
	}
}

func GetTestService(ctx context.Context, t *testing.T, pool *pgxpool.Pool) *Service {
	var err error

	fileRepo, err := pgsql.NewFileRepository(ctx, pool, test.FileTestTable)
	repo, err := pgsql.NewDataRepository(ctx, pool, test.DataTestTable, test.FileTestTable, test.UsersTestTable)
	assert.NoError(t, err)
	return NewService(repo, fileRepo)
}

func TestService_CheckUploadFileData(t *testing.T) {
	ctx := context.Background()
	internal.InitLogger()
	pool, err := test.InitConnection(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, pool, "no databases init")

	defer func(ctx context.Context, pool *pgxpool.Pool) {
		err = test.CleanData(ctx, pool, []string{"c_users", "c_data", "c_files"})
		assert.NoError(t, err)
	}(ctx, pool)

	userRepo, err := pgsql.NewUserRepository(ctx, pool, "c_users")
	assert.NoError(t, err)

	user := &domain.User{
		Login:    "test",
		Password: "test",
	}
	user1 := &domain.User{
		Login:    "test1",
		Password: "test",
	}

	userId, err := userRepo.Store(ctx, *user)
	assert.NoError(t, err)
	assert.NotZero(t, userId)

	userId1, err := userRepo.Store(ctx, *user1)
	assert.NoError(t, err)
	assert.NotZero(t, userId1)

	fileRepo, err := pgsql.NewFileRepository(ctx, pool, "c_files")
	assert.NoError(t, err)
	file := domain.File{
		Name: "pup",
		Path: "/dfff/sdd",
	}
	err = fileRepo.Insert(ctx, &file)
	assert.NoError(t, err)

	repo, err := pgsql.NewDataRepository(ctx, pool, "c_data", "c_files", "c_users")
	assert.NoError(t, err)
	service := NewService(repo, fileRepo)

	var fileId uint64 = 1
	var fileId1 uint64 = 2

	tests := []struct {
		name       string
		insertData *domain.Data
		newData    domain.Data
		file       domain.File
		wantErr    error
	}{
		{
			name: "wrong_data_id",
			insertData: &domain.Data{
				Name:    "1",
				Version: 1,
				UID:     userId,
				FileID:  nil,
			},
			newData: domain.Data{
				UID:    userId,
				FileID: &fileId,
			},
			wantErr: domain.ErrDataNotFound,
		},
		{
			name: "bad uid",
			insertData: &domain.Data{
				Name:    "2",
				Version: 1,
				UID:     userId,
				FileID:  nil,
			},
			newData: domain.Data{
				UID:    userId1,
				FileID: &fileId,
			},
			wantErr: domain.ErrDataNotFound,
		},
		{
			name: "empty db.fileId",
			insertData: &domain.Data{
				Name:    "3",
				Version: 1,
				UID:     userId,
				FileID:  nil,
			},
			newData: domain.Data{
				UID:    userId,
				FileID: &fileId,
			},
			wantErr: domain.ErrBadFileID,
		},
		{
			name: "empty data.fileId",
			insertData: &domain.Data{
				Name:    "4",
				Version: 1,
				UID:     userId,
				FileID:  nil,
			},
			newData: domain.Data{
				UID:    userId,
				FileID: nil,
			},
			wantErr: nil,
		},
		{
			name: "not equal data.fileId and db.dataId",
			insertData: &domain.Data{
				Name:    "5",
				Version: 1,
				UID:     userId,
				FileID:  &fileId,
			},
			newData: domain.Data{
				UID:    userId,
				FileID: &fileId1,
			},
			wantErr: domain.ErrBadFileID,
		},
		{
			name: "bd.fileId not null, data.fileId not null",
			insertData: &domain.Data{
				Name:    "6",
				Version: 1,
				UID:     userId,
				FileID:  &file.ID,
			},
			newData: domain.Data{
				UID:    userId,
				FileID: &file.ID,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = repo.Insert(ctx, tt.insertData)
			assert.NoError(t, err)
			tt.newData.ID = tt.insertData.ID

			err = service.CheckUploadFileData(ctx, tt.newData)
			if tt.wantErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
