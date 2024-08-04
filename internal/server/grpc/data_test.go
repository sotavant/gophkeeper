package grpc

import (
	"context"
	"gophkeeper/data"
	"gophkeeper/domain"
	"gophkeeper/file"
	"gophkeeper/internal"
	"gophkeeper/internal/server/auth"
	"gophkeeper/internal/server/grpc/interceptors"
	"gophkeeper/internal/server/repository/pgsql"
	"gophkeeper/internal/test"
	pb "gophkeeper/proto"
	user2 "gophkeeper/user"
	"io"
	"net"
	"os"
	"path/filepath"
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

var lis *bufconn.Listener

const bufSize = 1024 * 1024

const (
	testUsersTable = "d_users"
	testFileTable  = "d_files"
	testDataTable  = "d_data"
)

func TestDataServer_SaveData(t *testing.T) {
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

	userRepo, err := pgsql.NewUserRepository(ctx, pool, testUsersTable)
	assert.NoError(t, err)

	user := &domain.User{
		Login:    "test",
		Password: "test",
	}

	userID, err := userRepo.Store(ctx, *user)
	assert.NoError(t, err)
	assert.NotZero(t, userID)

	fileRepo, err = pgsql.NewFileRepository(ctx, pool, testFileTable)
	repo, err := pgsql.NewDataRepository(ctx, pool, testDataTable, testFileTable, testUsersTable)
	assert.NoError(t, err)

	service := data.NewService(repo, fileRepo)

	var versionFirst uint64 = 1
	login := "test"
	name := "test"
	testData := &domain.Data{
		Name:    name,
		Login:   &login,
		Pass:    &login,
		Version: versionFirst,
		UID:     userID,
	}

	err = repo.Insert(ctx, testData)
	assert.NoError(t, err)

	server := NewDataServer(service, "/", file.NewService(fileRepo))

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

func TestDataServer_UploadFile(t *testing.T) {
	userTable := "p_users"
	fileTable := "p_files"
	dataTable := "p_data"

	ctx := context.Background()
	internal.InitLogger()
	pool, err := test.InitConnection(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, pool, "no databases init")

	defer func(ctx context.Context, pool *pgxpool.Pool) {
		err = test.CleanData(ctx, pool, []string{userTable, fileTable, dataTable})
		assert.NoError(t, err)
	}(ctx, pool)

	fileRepo, err := pgsql.NewFileRepository(ctx, pool, fileTable)
	assert.NoError(t, err)

	userRepo, err := pgsql.NewUserRepository(ctx, pool, userTable)
	assert.NoError(t, err)

	user := &domain.User{
		Login:    "test",
		Password: "test",
	}
	userID, err := userRepo.Store(ctx, *user)

	repo, err := pgsql.NewDataRepository(ctx, pool, dataTable, fileTable, userTable)
	assert.NoError(t, err)
	dData := domain.Data{
		Name:    "5",
		Version: 1,
		UID:     userID,
	}
	err = repo.Insert(ctx, &dData)
	assert.NoError(t, err)

	service := data.NewService(repo, fileRepo)
	server := NewDataServer(service, "/tmp/uploaded", file.NewService(fileRepo))

	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer(grpc.UnaryInterceptor(interceptors.Auth), grpc.StreamInterceptor(interceptors.StreamAuth))
	pb.RegisterDataServiceServer(s, server)
	go func() {
		err = s.Serve(lis)
		assert.NoError(t, err)
	}()

	token, err := auth.BuildJWTString(userID)
	assert.NoError(t, err)

	md := metadata.Pairs(domain.AuthorizationMetaKey, domain.TokenSubstr+" "+token)
	mCtx := metadata.NewOutgoingContext(ctx, md)
	ctxx := context.WithValue(mCtx, user2.ContextUserIDKey{}, userID)

	conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		assert.NoError(t, err)
	}(conn)

	client := pb.NewDataServiceClient(conn)
	stream, err := client.UploadFile(ctxx)
	assert.NoError(t, err)

	tmpFile, err := os.CreateTemp("/tmp", "test_upload")
	assert.NoError(t, err)

	n, err := tmpFile.Write([]byte("some super text"))
	assert.NoError(t, err)
	assert.NotZero(t, n)
	defer func(name string) {
		err = os.Remove(name)
		assert.NoError(t, err)
	}(tmpFile.Name())

	buf := make([]byte, 1024)
	batchNum := 1

	fileForSend, err := os.Open(tmpFile.Name())
	for {
		var num int
		num, err = fileForSend.Read(buf)
		if err == io.EOF {
			break
		}
		assert.NoError(t, err)
		chunk := buf[:num]
		err = stream.Send(&pb.UploadFileRequest{
			DataId:      dData.ID,
			DataVersion: dData.Version,
			FileName:    filepath.Base(fileForSend.Name()),
			FileChunk:   chunk,
		})

		assert.NoError(t, err)
	}
	batchNum += 1

	_, err = stream.CloseAndRecv()
	assert.NoError(t, err)

	uploadedFilePath := filepath.Join("/tmp/uploaded", file.GetSaveFileSubDir(dData), filepath.Base(tmpFile.Name()))
	if _, err = os.Stat(uploadedFilePath); err != nil {
		assert.NoError(t, err)
	}

	err = os.Remove(uploadedFilePath)
	assert.NoError(t, err)
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}
