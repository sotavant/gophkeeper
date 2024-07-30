package data

import (
	"context"
	"gophkeeper/domain"
	"gophkeeper/internal"
	"gophkeeper/internal/server/repository/pgsql"
	"gophkeeper/internal/test"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestService_UpsertData(t *testing.T) {
	ctx := context.Background()
	internal.InitLogger()

	pool, err := test.InitConnection(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, pool, "no databases init")

	defer func(ctx context.Context, pool *pgxpool.Pool) {
		err = test.CleanData(ctx, pool)
		assert.NoError(t, err)
	}(ctx, pool)

	repo, err := pgsql.NewDataRepository(ctx, pool, test.DataTestTable)
	assert.NoError(t, err)

	service := NewService(repo)

	type want struct {
		wantErr bool
		errCode int
	}

	tests := []struct {
		name string
		data domain.Data
		want want
	}{
		{
			name: "new data", // need check id and version
		},
		{
			name: "new data with absent uid", // need check id and version
		},
		{
			name: "new data with bad uid", // need check id and version
		},
		{
			name: "success update data", // with same version
		},
		{
			name: "update data with outdated version", // with same version
		},
		{
			name: "update data with absent version", // with same version
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				dataRepo: tt.fields.dataRepo,
			}
			if err := s.UpsertData(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UpsertData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
