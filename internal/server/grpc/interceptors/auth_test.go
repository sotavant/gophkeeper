package interceptors

import (
	"context"
	"gophkeeper/domain"
	"gophkeeper/internal"
	"gophkeeper/internal/server/auth"
	pb "gophkeeper/proto"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/emptypb"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

type TestServer struct {
	pb.UnimplementedTestServiceServer
}

func (t *TestServer) Test(ctx context.Context, res *emptypb.Empty) (*pb.TestResponse, error) {
	return nil, nil
}

func TestAuth(t *testing.T) {
	internal.InitLogger()
	var ctx context.Context
	var userID uint64 = 1

	token, err := auth.BuildJWTString(userID)
	assert.NoError(t, err)

	tests := []struct {
		name       string
		token      string
		wantStatus codes.Code
	}{
		{
			"good token",
			domain.TokenSubstr + " " + token,
			codes.OK,
		},
		{
			"without token",
			"",
			codes.Unauthenticated,
		},
		{
			"bad token",
			domain.TokenSubstr + " sdfsdf",
			codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lis = bufconn.Listen(bufSize)
			s := grpc.NewServer(grpc.UnaryInterceptor(Auth))
			pb.RegisterTestServiceServer(s, &TestServer{})
			go func() {
				err = s.Serve(lis)
				assert.NoError(t, err)
			}()

			var conn *grpc.ClientConn
			conn, err = grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
			assert.NoError(t, err)
			defer func(conn *grpc.ClientConn) {
				err = conn.Close()
				assert.NoError(t, err)
			}(conn)

			if tt.token != "" {
				md := metadata.Pairs(domain.AuthorizationMetaKey, tt.token)
				ctx = metadata.NewOutgoingContext(context.Background(), md)
			} else {
				ctx = context.Background()
			}

			client := pb.NewTestServiceClient(conn)
			_, err = client.Test(ctx, nil)
			if tt.wantStatus == codes.OK {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, status.Code(err), tt.wantStatus)
			}
		})
	}
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}
