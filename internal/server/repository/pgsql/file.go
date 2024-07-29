package pgsql

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

const FileTableName = "file"

type FileRepository struct {
	DBPoll    *pgxpool.Pool
	tableName string
}

func NewFileRepository(ctx context.Context, pool *pgxpool.Pool, tableName string) (*FileRepository, error) {
	err := createFileTable(ctx, pool, tableName)
	if err != nil {
		return nil, err
	}

	return &FileRepository{
		DBPoll:    pool,
		tableName: tableName,
	}, nil
}

func createFileTable(ctx context.Context, pool *pgxpool.Pool, tableName string) error {
	query := strings.ReplaceAll(`create table if not exists #T#
		(
			id    serial primary key,
			name varchar(255) not null,
			path varchar(255) not null
		);`, "#T#", tableName)

	_, err := pool.Exec(ctx, query)

	if err != nil {
		return err
	}
	return nil
}
