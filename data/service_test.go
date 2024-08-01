package data

import (
	"context"
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
		err = test.CleanData(ctx, pool)
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

	repo, err := pgsql.NewDataRepository(ctx, pool, test.DataTestTable)
	assert.NoError(t, err)

	testData := &domain.Data{
		Name:    "testic",
		Login:   "test",
		Pass:    "test",
		Version: 1,
		UID:     userId,
	}

	testData2 := &domain.Data{
		Name:    "testic2",
		Login:   "test",
		Pass:    "test",
		Version: 1,
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
		name string
		data domain.Data
		want want
	}{
		{
			name: "success new data", // need check userId and version
			data: domain.Data{
				Login:   "test",
				Pass:    "test",
				CardNum: "test",
				Text:    "test",
				Name:    "test",
				Meta:    "test",
				UID:     userId,
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "name exist", // need check userId and version
			data: domain.Data{
				Login:   "test",
				Pass:    "test",
				CardNum: "test",
				Text:    "test",
				Name:    testData.Name,
				Meta:    "test",
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
				Text:    "success update",
				Version: testData.Version,
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "wrong update: name exist",
			data: domain.Data{
				ID:      testData.ID,
				Name:    testData2.Name,
				Version: testData2.Version,
			},
			want: want{
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				dataRepo: tt.fields.dataRepo,
			}
			if err := s.UpsertData(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UpsertData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
