package pgsql

import (
	"context"
	"gophkeeper/server/domain"
	"strings"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const UsersTableName = "users"

// UserRepository структура для взаимодействия с таблицей пользвателей
type UserRepository struct {
	DBPoll    *pgxpool.Pool
	tableName string
}

func NewUserRepository(ctx context.Context, pool *pgxpool.Pool, tableName string) (*UserRepository, error) {
	err := createUsersTable(ctx, pool, tableName)
	if err != nil {
		return nil, err
	}

	return &UserRepository{
		DBPoll:    pool,
		tableName: tableName,
	}, nil
}

// GetByLogin Получить пользователя по ИД
func (u *UserRepository) GetByLogin(ctx context.Context, login string) (domain.User, error) {
	query := u.setUserTableName(`select id, login, password from #T# where login = $1`)

	return u.getOne(ctx, query, login)
}

// Store добавить нового пользователя
func (u *UserRepository) Store(ctx context.Context, user domain.User) (uint64, error) {
	var id uint64
	query := u.setUserTableName(`insert into #T# (login, password) values ($1, $2) returning id`)

	err := u.DBPoll.QueryRow(ctx, query, user.Login, user.Password).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (u *UserRepository) getOne(ctx context.Context, query string, args ...interface{}) (user domain.User, err error) {
	rows, err := u.DBPoll.Query(ctx, query, args...)
	if err != nil {
		return user, err
	}

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.User])
	if err != nil {
		return user, err
	}

	for _, user = range users {
		return user, nil
	}

	return
}

func createUsersTable(ctx context.Context, pool *pgxpool.Pool, tableName string) error {
	query := strings.ReplaceAll(`create table if not exists #T#
		(
			id    serial primary key,
			login  varchar not null,
			password varchar not null
		);`, "#T#", tableName)

	_, err := pool.Exec(ctx, query)

	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) setUserTableName(query string) string {
	return strings.ReplaceAll(query, "#T#", u.tableName)
}
