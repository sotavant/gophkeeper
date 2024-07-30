package test

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

const UsersTestTable = "testUsers"
const DataTestTable = "testData"

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

func CleanData(ctx context.Context, pool *pgxpool.Pool) error {
	query := `drop table if exists $1, $2;`

	_, err := pool.Exec(ctx, query, UsersTestTable, DataTestTable)
	if err != nil {
		return err
	}

	return nil
}
