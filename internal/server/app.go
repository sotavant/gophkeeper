package server

import (
	"context"
	"errors"
	"flag"
	"gophkeeper/internal"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	runAddressVar  = "RUN_ADDRESS"
	databaseURIVar = "DATABASE_URI"
	saveFilesPath  = "FILES_SAVE_PATH"
)

type App struct {
	Address       string
	DBPool        *pgxpool.Pool
	FilesSavePath string
}

type config struct {
	runAddress,
	databaseURI,
	saveFilePath string
}

func InitApp(ctx context.Context) (*App, error) {
	internal.InitLogger()
	c := initConfig()
	if err := checkConfig(c); err != nil {
		return nil, err
	}

	dbPool, err := initDB(ctx, c.databaseURI)
	if err != nil {
		return nil, err
	}

	return &App{
		Address:       c.runAddress,
		DBPool:        dbPool,
		FilesSavePath: c.saveFilePath,
	}, nil
}

func initConfig() *config {
	c := new(config)

	flag.StringVar(&c.runAddress, "a", "", "server address")
	flag.StringVar(&c.databaseURI, "d", "", "database uri")
	flag.StringVar(&c.saveFilePath, "f", "", "save files path")

	flag.Parse()

	if envVar := os.Getenv(runAddressVar); envVar != "" {
		c.runAddress = envVar
	}

	if envVar := os.Getenv(databaseURIVar); envVar != "" {
		c.databaseURI = envVar
	}

	if envVar := os.Getenv(saveFilesPath); envVar != "" {
		c.databaseURI = envVar
	}

	return c
}

func initDB(ctx context.Context, DSN string) (*pgxpool.Pool, error) {
	dbConf, err := pgxpool.ParseConfig(DSN)
	if err != nil {
		return nil, err
	}

	dbPool, err := pgxpool.NewWithConfig(ctx, dbConf)
	if err != nil {
		return nil, err
	}

	err = dbPool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return dbPool, nil
}

func checkConfig(c *config) error {
	if c.runAddress == "" || c.databaseURI == "" || c.saveFilePath == "" {
		return errors.New("please, check configs")
	}

	return nil
}
