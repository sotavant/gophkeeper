package pgsql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
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

func (d *DataRepository) GetByNameAndUserID(ctx context.Context, uid int64, name string) (int64, error) {
	query := d.setTableName(`select id from #T# where uid = $1 and name = $2`)

	data, err := d.getOne(ctx, query, uid, name)
	if err != nil {
		return 0, err
	}

	return data.ID, nil
}

func (d *DataRepository) GetById(ctx context.Context, id int64, fields []string) (*domain.Data, error) {
	selectFields := "*"
	if len(fields) > 0 {
		selectFields = strings.Join(fields, ",")
	}

	query := d.setTableName(fmt.Sprintf(`select %s from #T# where id = $1`, selectFields))

	row, err := d.getOne(ctx, query, id)
	if err != nil {
		return nil, err
	}

	return &row, nil
}

func (d *DataRepository) getOne(ctx context.Context, query string, args ...interface{}) (data domain.Data, err error) {
	rows, err := d.DBPoll.Query(ctx, query, args...)
	if err != nil {
		return data, err
	}

	datas, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Data])
	if err != nil {
		return data, err
	}

	for _, data = range datas {
		return data, nil
	}

	return
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
    		version integer not null,
    		constraint  name_unique UNIQUE (name, uid)
		);`, "#T#", tableName)

	_, err := pool.Exec(ctx, query)

	if err != nil {
		return err
	}
	return nil
}
