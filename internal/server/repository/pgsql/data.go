package pgsql

import (
	"context"
	"gophkeeper/domain"
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

func (d *DataRepository) Insert(ctx context.Context, data *domain.Data) error {
	query := d.setTableName(`insert into #T# (name, uid, login, pass, text, card_num, meta, version) values ($1, $2, $3, $4, $5, $6, $7, $8) returning id`)

	err := d.DBPoll.QueryRow(ctx, query, data.Name, data.UID, data.Login, data.Pass, data.CardNum, data.Meta, data.Version).Scan(&data.ID)
	if err != nil {
		return err
	}

	return nil
}

func (d *DataRepository) Update(ctx context.Context, data domain.Data) error {
	query := d.setTableName(`update #T# set
		name = $1, 
		login = $2,
		pass = $3,
		text = $4,
		card_num = $5,
		meta = $6,
		version = $7
		where id = $8
	`)

	_, err := d.DBPoll.Exec(ctx, query, data.Name, data.Login, data.Pass, data.Text, data.CardNum, data.Meta, data.Version, data.ID)

	if err != nil {
		return err
	}

	return nil
}

func (d *DataRepository) setTableName(query string) string {
	return strings.ReplaceAll(query, "#T#", d.tableName)
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
