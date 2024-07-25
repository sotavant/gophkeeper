package grpc

import (
	"context"
	"gophkeeper/internal"
	"gophkeeper/internal/server/repository/pgsql"
	"gophkeeper/internal/test"
	"gophkeeper/proto"
	"gophkeeper/user"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUserServer_Register(t *testing.T) {
	internal.InitLogger()
	ctx := context.Background()

	pool, err := test.InitConnection(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, pool, "no databases init")

	defer func(ctx context.Context, pool *pgxpool.Pool) {
		err = test.CleanData(ctx, pool)
		assert.NoError(t, err)
	}(ctx, pool)

	repo, err := pgsql.NewUserRepository(ctx, pool, test.TestUsersTable)
	assert.NoError(t, err)

	server := NewUserServer(user.NewService(repo))

	tests := []struct {
		name    string
		request proto.RegisterRequest
		wantErr codes.Code
	}{
		{
			name: "emptyUser",
			request: proto.RegisterRequest{
				User: nil,
			},
			wantErr: codes.InvalidArgument,
		},
		{
			name: "shortLogin",
			request: proto.RegisterRequest{
				User: &proto.User{
					Login:    "a",
					Password: "sdfffdd",
				},
			},
			wantErr: codes.InvalidArgument,
		},
		{
			name: "bigLogin",
			request: proto.RegisterRequest{
				User: &proto.User{
					Login: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
						"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
					Password: "sdfffdd",
				},
			},
			wantErr: codes.InvalidArgument,
		},
		{
			name: "shortPassword",
			request: proto.RegisterRequest{
				User: &proto.User{
					Login:    "aaa",
					Password: "sdff",
				},
			},
			wantErr: codes.InvalidArgument,
		},
		{
			name: "bigPassword",
			request: proto.RegisterRequest{
				User: &proto.User{
					Login:    "aaa",
					Password: "sdffsdffsdffsdff",
				},
			},
			wantErr: codes.InvalidArgument,
		},
		{
			name: "ok",
			request: proto.RegisterRequest{
				User: &proto.User{
					Login:    "aaa",
					Password: "aaaddddd",
				},
			},
			wantErr: codes.OK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := server.Register(ctx, &tt.request)
			if (err != nil) && tt.wantErr == codes.OK {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr != codes.OK {
				e, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, e.Code(), tt.wantErr)
				return
			}

			assert.NotEmpty(t, got)
		})
	}
}
