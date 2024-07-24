package user

import (
	"context"
	"gophkeeper/domain"
	"gophkeeper/internal/server/repository/pgsql"
	"gophkeeper/internal/test"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestService_Register(t *testing.T) {
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

	service := NewService(repo)

	type fields struct {
		userRepo UserRepository
	}
	type args struct {
		ctx  context.Context
		user domain.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Service{
				userRepo: tt.fields.userRepo,
			}
			got, err := u.Register(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Register() got = %v, want %v", got, tt.want)
			}
		})
	}
}
