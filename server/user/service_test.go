package user

import (
	"context"
	"gophkeeper/internal"
	"gophkeeper/internal/server/repository/pgsql"
	"gophkeeper/internal/test"
	domain2 "gophkeeper/server/domain"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestService_Register(t *testing.T) {
	testUserTable := "reg_test_users"
	ctx := context.Background()
	internal.InitLogger()

	pool, err := test.InitConnection(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, pool, "no databases init")

	defer func(ctx context.Context, pool *pgxpool.Pool) {
		err = test.CleanData(ctx, pool, []string{testUserTable})
		assert.NoError(t, err)
	}(ctx, pool)

	repo, err := pgsql.NewUserRepository(ctx, pool, testUserTable)
	assert.NoError(t, err)

	service := NewService(repo)

	user := &domain2.User{
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
		user       domain2.User
		wantErr    bool
		wantResult wantResult
	}{
		{
			name:    "login busy",
			user:    *user,
			wantErr: true,
			wantResult: wantResult{
				err: domain2.ErrLoginExist,
			},
		},
		{
			name: "correct",
			user: domain2.User{
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

func TestService_Auth(t *testing.T) {
	tableName := "test_users_auth"
	internal.InitLogger()
	ctx := context.Background()

	pool, err := test.InitConnection(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, pool, "no databases init")

	defer func(ctx context.Context, pool *pgxpool.Pool) {
		err = test.CleanData(ctx, pool, []string{tableName})
		assert.NoError(t, err)
	}(ctx, pool)

	repo, err := pgsql.NewUserRepository(ctx, pool, tableName)
	assert.NoError(t, err)

	service := NewService(repo)

	user := &domain2.User{
		Login:    "test",
		Password: "testtest",
	}

	hashedPass, err := HashPassword(user.Password)
	assert.NoError(t, err)
	storedUser := domain2.User{
		Login:    user.Login,
		Password: hashedPass,
	}

	id, err := repo.Store(ctx, storedUser)
	assert.NoError(t, err)
	assert.NotZero(t, id)

	type wantResult struct {
		err error
	}

	tests := []struct {
		name       string
		user       domain2.User
		wantErr    bool
		wantResult wantResult
	}{
		{
			name:    "correct",
			user:    *user,
			wantErr: false,
		},
		{
			name: "badLogin",
			user: domain2.User{
				Login:    "test1",
				Password: "testtest",
			},
			wantErr: true,
			wantResult: wantResult{
				err: domain2.ErrUserNotFound,
			},
		},
		{
			name: "badPass",
			user: domain2.User{
				Login:    "test",
				Password: "testtesttt",
			},
			wantErr: true,
			wantResult: wantResult{
				err: domain2.ErrUserNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.Login(ctx, tt.user)
			if tt.wantErr {
				assert.Equal(t, tt.wantResult.err, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, got)
		})
	}
}
