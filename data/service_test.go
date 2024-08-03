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

	_, err = pgsql.NewFileRepository(ctx, pool, test.FileTestTable)
	repo, err := pgsql.NewDataRepository(ctx, pool, test.DataTestTable, test.FileTestTable, test.UsersTestTable)
	assert.NoError(t, err)

	var versionFirst int64 = 1
	var versionSecond int64 = 2
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

	err = repo.Insert(ctx, testData)
	assert.NoError(t, err)

	err = repo.Insert(ctx, testData2)
	assert.NoError(t, err)

	service := NewService(repo)

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
				dd, err = repo.GetById(ctx, tt.data.ID, []string{})
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
