package grpc

import (
	"context"
	"gophkeeper/data"
	"gophkeeper/domain"
	"gophkeeper/file"
	"gophkeeper/internal"
	"gophkeeper/internal/server/auth"
	g "gophkeeper/internal/server/grpc"
	"gophkeeper/internal/server/grpc/interceptors"
	"gophkeeper/internal/server/repository/pgsql"
	"gophkeeper/internal/test"
	pb "gophkeeper/proto"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func TestDataClient_SaveDate(t *testing.T) {
	internal.InitLogger()
	ctx := context.Background()
	existingUserLogin := "kaka"

	pool, err := test.InitConnection(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, pool, "no databases init")

	defer func(ctx context.Context, pool *pgxpool.Pool) {
		//err = test.CleanData(ctx, pool, []string{test.DataTestTable, test.FileTestTable, test.UsersTestTable})
		assert.NoError(t, err)
	}(ctx, pool)

	// test user
	repo, err := pgsql.NewUserRepository(ctx, pool, test.UsersTestTable)
	assert.NoError(t, err)

	userId, err := repo.Store(ctx, domain.User{
		Login:    existingUserLogin,
		Password: "kakadud",
	})

	assert.NoError(t, err)
	assert.NotEqual(t, 0, userId)

	// test file repo
	fileRepo, err := pgsql.NewFileRepository(ctx, pool, test.FileTestTable)
	assert.NoError(t, err)

	// test data
	dataRepo, err := pgsql.NewDataRepository(ctx, pool, test.DataTestTable, test.FileTestTable, test.UsersTestTable)
	assert.NoError(t, err)

	testData := &domain.Data{
		Name:    "test",
		Version: 1,
		UID:     userId,
	}

	err = dataRepo.Insert(ctx, testData)
	assert.NoError(t, err)

	server := g.NewDataServer(data.NewService(dataRepo, fileRepo), "/tmp", file.NewService(fileRepo))

	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer(grpc.UnaryInterceptor(interceptors.Auth))
	pb.RegisterDataServiceServer(s, server)
	go func() {
		err = s.Serve(lis)
		assert.NoError(t, err)
	}()

	// client
	conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		assert.NoError(t, err)
	}(conn)

	token, err := auth.BuildJWTString(userId)
	assert.NoError(t, err)
	md := metadata.Pairs(domain.AuthorizationMetaKey, domain.TokenSubstr+" "+token)
	mCtx := metadata.NewOutgoingContext(ctx, md)

	pbClient := pb.NewDataServiceClient(conn)
	client := NewDataClient(pbClient)

	tests := []struct {
		name          string
		data          *domain.Data
		wantErr       bool
		wantErrorCode codes.Code
	}{
		{
			name: "success",
			data: &domain.Data{
				Name:    "success",
				Version: 1,
			},
			wantErr:       false,
			wantErrorCode: codes.OK,
		},
		{
			name: "name exist",
			data: &domain.Data{
				Name:    testData.Name,
				Version: 1,
			},
			wantErr:       true,
			wantErrorCode: codes.InvalidArgument,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = client.SaveData(mCtx, tt.data)
			if tt.wantErr {
				assert.Equal(t, tt.wantErrorCode, status.Code(err))
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, tt.data.ID)
			}
		})
	}
}
