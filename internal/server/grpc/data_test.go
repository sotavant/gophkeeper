package grpc

import (
	"context"
	"gophkeeper/data"
	"gophkeeper/domain"
	"gophkeeper/internal"
	"gophkeeper/internal/server/repository/pgsql"
	"gophkeeper/internal/test"
	pb "gophkeeper/proto"
	user2 "gophkeeper/user"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	testUsersTable = "d_users"
	testFileTable  = "d_files"
	testDataTable  = "d_data"
)

func TestDataServer_SaveData(t *testing.T) {
	ctx := context.Background()
	internal.InitLogger()
	pool, err := test.InitConnection(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, pool, "no databases init")

	defer func(ctx context.Context, pool *pgxpool.Pool) {
		err = test.CleanData(ctx, pool, []string{testUsersTable, testDataTable, testFileTable})
		assert.NoError(t, err)
	}(ctx, pool)

	userRepo, err := pgsql.NewUserRepository(ctx, pool, testUsersTable)
	assert.NoError(t, err)

	user := &domain.User{
		Login:    "test",
		Password: "test",
	}

	userID, err := userRepo.Store(ctx, *user)
	assert.NoError(t, err)
	assert.NotZero(t, userID)

	_, err = pgsql.NewFileRepository(ctx, pool, testFileTable)
	repo, err := pgsql.NewDataRepository(ctx, pool, testDataTable, testFileTable, testUsersTable)
	assert.NoError(t, err)

	service := data.NewService(repo)

	var versionFirst int64 = 1
	login := "test"
	name := "test"
	testData := &domain.Data{
		Name:    name,
		Login:   &login,
		Pass:    &login,
		Version: versionFirst,
		UID:     userID,
	}

	err = service.DataRepo.Insert(ctx, testData)
	assert.NoError(t, err)

	server := NewDataServer(service)

	tests := []struct {
		name    string
		req     *pb.SaveDataRequest
		ctx     context.Context
		wantErr bool
		errCode codes.Code
	}{
		{
			name: "success insert",
			req: &pb.SaveDataRequest{
				Data: &pb.Data{
					Name:  name + "_",
					Login: "test",
				},
			},
			ctx:     context.WithValue(ctx, user2.ContextUserIDKey{}, userID),
			wantErr: false,
		},
		{
			name: "success update",
			req: &pb.SaveDataRequest{
				Data: &pb.Data{
					Id:      testData.ID,
					Name:    name,
					Login:   "test1",
					Version: testData.Version,
				},
			},
			ctx:     context.WithValue(ctx, user2.ContextUserIDKey{}, userID),
			wantErr: false,
		},
		{
			name: "with empty name",
			req: &pb.SaveDataRequest{
				Data: &pb.Data{
					Id:    testData.ID,
					Name:  "",
					Login: "test1",
				},
			},
			ctx:     context.WithValue(ctx, user2.ContextUserIDKey{}, userID),
			wantErr: true,
			errCode: codes.InvalidArgument,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err = server.SaveData(tt.ctx, tt.req)
			if tt.wantErr {
				assert.NotNil(t, err)
				e, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.errCode, e.Code())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
