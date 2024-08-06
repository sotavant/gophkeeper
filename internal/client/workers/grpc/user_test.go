package grpc

import (
	"context"
	"fmt"
	"gophkeeper/internal"
	g "gophkeeper/internal/server/grpc"
	"gophkeeper/internal/server/repository/pgsql"
	"gophkeeper/internal/test"
	pb "gophkeeper/proto"
	"gophkeeper/user"
	"net"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

var lis *bufconn.Listener

const bufSize = 1024 * 1024

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestUserClient_Registration(t *testing.T) {
	internal.InitLogger()
	ctx := context.Background()

	pool, err := test.InitConnection(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, pool, "no databases init")

	defer func(ctx context.Context, pool *pgxpool.Pool) {
		err = test.CleanData(ctx, pool, []string{test.UsersTestTable})
		assert.NoError(t, err)
	}(ctx, pool)

	repo, err := pgsql.NewUserRepository(ctx, pool, test.UsersTestTable)
	assert.NoError(t, err)

	server := g.NewUserServer(user.NewService(repo))

	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, server)
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

	pbClient := pb.NewUserServiceClient(conn)
	client := NewUserClient(pbClient)

	type args struct {
		login    string
		password string
	}
	tests := []struct {
		name          string
		args          args
		wantErr       bool
		wantErrorCode codes.Code
	}{
		{
			name: "success",
			args: args{
				login:    "test",
				password: "testicaaa",
			},
			wantErr:       false,
			wantErrorCode: codes.OK,
		},
		{
			name: "short_login",
			args: args{
				login:    "t",
				password: "testic",
			},
			wantErr:       true,
			wantErrorCode: codes.InvalidArgument,
		},
		{
			name: "short_pass",
			args: args{
				login:    "tttt",
				password: "tes",
			},
			wantErr:       true,
			wantErrorCode: codes.InvalidArgument,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var token string
			token, err = client.Registration(tt.args.login, tt.args.password)
			if tt.wantErr {
				fmt.Println(err.Error())
				assert.Equal(t, tt.wantErrorCode, status.Code(err))
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}
}
