package pgsql

import (
	"context"
	"gophkeeper/server/domain"
	"strings"

	"github.com/jackc/pgx/v5"
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

func (f *FileRepository) Get(ctx context.Context, id uint64) (*domain.File, error) {
	query := f.setTableName(`select * from #T# where id = $1`)

	row, err := f.getOne(ctx, query, id)
	if err != nil {
		return nil, err
	}

	return &row, nil
}

func (f *FileRepository) Insert(ctx context.Context, file *domain.File) error {
	query := f.setTableName(`insert into #T# (name, path) values ($1, $2) returning id`)

	err := f.DBPoll.QueryRow(ctx, query, file.Name, file.Path).Scan(&file.ID)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileRepository) Update(ctx context.Context, file *domain.File) error {
	query := f.setTableName(`update #T# set
		name = $1, 
		path = $2
		where id = $3
	`)

	_, err := f.DBPoll.Exec(ctx, query, file.Name, file.Path, file.ID)

	if err != nil {
		return err
	}

	return nil
}

func (f *FileRepository) Delete(ctx context.Context, id uint64) error {
	query := f.setTableName(`delete from #T# where id = $1`)
	_, err := f.DBPoll.Exec(ctx, query, id)
	return err
}

func (f *FileRepository) setTableName(query string) string {
	return strings.ReplaceAll(query, "#T#", f.tableName)
}

func (f *FileRepository) getOne(ctx context.Context, query string, args ...interface{}) (file domain.File, err error) {
	rows, err := f.DBPoll.Query(ctx, query, args...)
	if err != nil {
		return
	}

	files, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.File])
	if err != nil {
		return
	}

	for _, file = range files {
		return
	}

	return
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
