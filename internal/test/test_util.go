// Package test вспомогательный пакет для тестирования
package test

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

const UsersTestTable = "test_users"
const DataTestTable = "test_data"
const FileTestTable = "test_file"

func InitConnection(ctx context.Context) (*pgxpool.Pool, error) {
	dns := os.Getenv("TEST_DATABASE_DSN")
	if dns == "" {
		return nil, nil
	}

	dbConn, err := pgxpool.New(ctx, dns)

	if err != nil {
		return nil, err
	}

	err = dbConn.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}

func CleanData(ctx context.Context, pool *pgxpool.Pool, tables []string) error {
	tt := strings.Join(tables, ",")
	query := fmt.Sprintf(`drop table if exists %s;`, tt)

	_, err := pool.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
