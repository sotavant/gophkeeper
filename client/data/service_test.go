package data

import (
	"context"
	"gophkeeper/client/domain"
	"gophkeeper/data"
	serverDomain "gophkeeper/domain"
	"gophkeeper/file"
	"gophkeeper/internal"
	"gophkeeper/internal/client"
	g "gophkeeper/internal/client/workers/grpc"
	interceptors2 "gophkeeper/internal/client/workers/grpc/interceptors"
	"gophkeeper/internal/server/auth"
	grpc2 "gophkeeper/internal/server/grpc"
	"gophkeeper/internal/server/grpc/interceptors"
	"gophkeeper/internal/server/repository/pgsql"
	"gophkeeper/internal/test"
	pb "gophkeeper/proto"
	"net"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

var lis *bufconn.Listener

const bufSize = 1024 * 1024

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func Test_hashData(t *testing.T) {
	client.AppInstance = &client.App{}
	client.AppInstance.SetStorageKey("login", "pass")

	type args struct {
		data domain.Data
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test hash data",
			args: args{
				data: domain.Data{
					Name:  "name",
					Login: "login",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encryptData(tt.args.data)
			assert.NoError(t, err)
			assert.Equal(t, tt.args.data.Name, got.Name)
			assert.NotEqual(t, tt.args.data.Login, got.Login)
			assert.Equal(t, "", got.Pass)
		})
	}
}

func TestSaveData(t *testing.T) {
	testFile := createTestFile(t)
	ctx := context.Background()

	// prepare server side
	userTable := "p_users"
	fileTable := "p_files"
	dataTable := "p_data"

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

	// test user
	user := &serverDomain.User{
		Login:    "test",
		Password: "test",
	}
	userID, err := userRepo.Store(ctx, *user)
	user.ID = userID

	repo, err := pgsql.NewDataRepository(ctx, pool, dataTable, fileTable, userTable)
	assert.NoError(t, err)

	// server grpc server
	service := data.NewService(repo, fileRepo)
	server := grpc2.NewDataServer(service, "/tmp/uploaded", file.NewService(fileRepo))

	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer(grpc.UnaryInterceptor(interceptors.Auth), grpc.StreamInterceptor(interceptors.StreamAuth))
	pb.RegisterDataServiceServer(s, server)
	go func() {
		err = s.Serve(lis)
		assert.NoError(t, err)
	}()

	// client instance
	conn := initTestAppInstance(t, user)
	assert.NoError(t, err)
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		assert.NoError(t, err)
	}(conn)

	type args struct {
		data *domain.Data
	}
	tests := []struct {
		name         string
		args         args
		needFileName bool
	}{
		{
			name: "without file", // check id, version in appinstance.data
			args: args{
				data: &domain.Data{
					Name:  "name",
					Login: "login",
				},
			},
			needFileName: false,
		},
		{
			name: "with file",
			args: args{
				data: &domain.Data{
					Name:     "second name",
					Login:    "login",
					FilePath: testFile.Name(),
				},
			},
			needFileName: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var version uint64 = 0
			var gotData domain.Data
			gotData, err = SaveData(*tt.args.data)

			assert.NoError(t, err)
			assert.NotNil(t, client.AppInstance.DecryptedData[gotData.ID])
			assert.Greater(t, client.AppInstance.DecryptedData[gotData.ID].Version, version)
			assert.Equal(t, client.AppInstance.DecryptedData[gotData.ID].Version, gotData.Version)

			if tt.needFileName {
				assert.Greater(t, client.AppInstance.DecryptedData[gotData.ID].FileID, version)
				assert.NotEmpty(t, client.AppInstance.DecryptedData[gotData.ID].FilePath)
			}
		})
	}
}

func TestGet(t *testing.T) {
	ctx := context.Background()

	// prepare server side
	userTable := "d_users"
	fileTable := "d_files"
	dataTable := "d_data"

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

	// test user
	user := &serverDomain.User{
		Login:    "test",
		Password: "test",
	}
	userID, err := userRepo.Store(ctx, *user)
	user.ID = userID

	repo, err := pgsql.NewDataRepository(ctx, pool, dataTable, fileTable, userTable)
	assert.NoError(t, err)

	// server grpc server
	service := data.NewService(repo, fileRepo)
	server := grpc2.NewDataServer(service, "/tmp/uploaded", file.NewService(fileRepo))

	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer(grpc.UnaryInterceptor(interceptors.Auth), grpc.StreamInterceptor(interceptors.StreamAuth))
	pb.RegisterDataServiceServer(s, server)
	go func() {
		err = s.Serve(lis)
		assert.NoError(t, err)
	}()

	// client instance
	conn := initTestAppInstance(t, user)
	assert.NoError(t, err)
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		assert.NoError(t, err)
	}(conn)

	testData := domain.Data{
		Name: "test",
		Pass: "test",
	}
	testData, err = SaveData(testData)
	assert.NoError(t, err)

	type args struct {
		id uint64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "ok", // check id, version in appinstance.data
			args: args{
				id: testData.ID,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var compareVersion uint64 = 0
			var got *domain.Data
			got, err = GetData(tt.args.id)

			assert.NoError(t, err)

			assert.NotNil(t, client.AppInstance.DecryptedData[got.ID])
			assert.Greater(t, client.AppInstance.DecryptedData[got.ID].Version, compareVersion)
			assert.Equal(t, got.Name, testData.Name)

		})
	}
}

func TestGetDataList(t *testing.T) {
	ctx := context.Background()

	// prepare server side
	userTable := "d_users"
	fileTable := "d_files"
	dataTable := "d_data"

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

	// test user
	user := &serverDomain.User{
		Login:    "test",
		Password: "test",
	}
	userID, err := userRepo.Store(ctx, *user)
	user.ID = userID

	repo, err := pgsql.NewDataRepository(ctx, pool, dataTable, fileTable, userTable)
	assert.NoError(t, err)

	// server grpc server
	service := data.NewService(repo, fileRepo)
	server := grpc2.NewDataServer(service, "/tmp/uploaded", file.NewService(fileRepo))

	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer(grpc.UnaryInterceptor(interceptors.Auth), grpc.StreamInterceptor(interceptors.StreamAuth))
	pb.RegisterDataServiceServer(s, server)
	go func() {
		err = s.Serve(lis)
		assert.NoError(t, err)
	}()

	// client instance
	conn := initTestAppInstance(t, user)
	assert.NoError(t, err)
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		assert.NoError(t, err)
	}(conn)

	testData := domain.Data{
		Name: "test",
		Pass: "test",
	}
	testData, err = SaveData(testData)
	assert.NoError(t, err)

	tests := []struct {
		name string
	}{
		{
			name: "ok", // check id, version in appinstance.data
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []serverDomain.DataName
			got, err = GetDataList()
			assert.NoError(t, err)

			assert.Equal(t, 1, len(got))
			assert.Equal(t, got[0].Name, testData.Name)
			assert.Equal(t, got[0].ID, testData.ID)
		})
	}
}

func createTestFile(t *testing.T) os.File {
	tmpFile, err := os.CreateTemp("/tmp", "test_upload")
	assert.NoError(t, err)

	n, err := tmpFile.Write([]byte("some super text"))
	assert.NoError(t, err)
	assert.NotZero(t, n)

	return *tmpFile
}

func initTestAppInstance(t *testing.T, user *serverDomain.User) *grpc.ClientConn {
	client.AppInstance = &client.App{
		DecryptedData: make(map[uint64]domain.Data),
	}
	client.AppInstance.User.Login = user.Login
	client.AppInstance.SetStorageKey(user.Login, user.Password)

	token, err := auth.BuildJWTString(user.ID)
	assert.NoError(t, err)

	client.AppInstance.User.Token = token

	conn, err := grpc.NewClient(
		"passthrough://bufnet", grpc.WithContextDialer(bufDialer),
		grpc.WithUnaryInterceptor(interceptors2.Auth),
		grpc.WithStreamInterceptor(interceptors2.StreamAuth),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		internal.Logger.Fatalw("failed to create grpc client", "error", err)
	}

	client.AppInstance.UserClient = g.NewUserClient(pb.NewUserServiceClient(conn))
	client.AppInstance.DataClient = g.NewDataClient(pb.NewDataServiceClient(conn))

	return conn
}
