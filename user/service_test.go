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

	user := &domain.User{
		Login:    "test",
		Password: "test",
	}

	id, err := repo.Store(ctx, *user)
	assert.NoError(t, err)
	assert.NotZero(t, id)

	type wantResult struct {
		err error
	}

	tests := []struct {
		name       string
		user       domain.User
		wantErr    bool
		wantResult wantResult
	}{
		{
			name:    "login busy",
			user:    *user,
			wantErr: true,
			wantResult: wantResult{
				err: domain.ErrLoginExist,
			},
		},
		{
			name: "correct",
			user: domain.User{
				Login:    "test1",
				Password: "test1",
			},
			wantErr: false,
			wantResult: wantResult{
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.Register(ctx, tt.user)
			if tt.wantErr {
				assert.Equal(t, tt.wantResult.err, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, got)
		})
	}
}
