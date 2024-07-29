package pgsql

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

const DataTableName = "data"

type DataRepository struct {
	DBPoll    *pgxpool.Pool
	tableName string
}

func NewDataRepository(ctx context.Context, pool *pgxpool.Pool, tableName string) (*DataRepository, error) {
	err := createDataTable(ctx, pool, tableName)
	if err != nil {
		return nil, err
	}

	return &DataRepository{
		DBPoll:    pool,
		tableName: tableName,
	}, nil
}

func createDataTable(ctx context.Context, pool *pgxpool.Pool, tableName string) error {
	query := strings.ReplaceAll(`create table if not exists #T#
		(
			id    serial primary key,
			name varchar(255) not null,
			uid      integer not null
        		constraint user___fk
            		references users,
			file   integer not null
			    constraint user___fk_file
			    references file,
    		login    varchar,
    		pass     varchar,
    		text     text,
    		card_num varchar,
    		meta     varchar,
    		version integer not null
		);`, "#T#", tableName)

	_, err := pool.Exec(ctx, query)

	if err != nil {
		return err
	}
	return nil
}
