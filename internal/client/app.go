package client

import (
	"crypto/sha1"
	"errors"
	"flag"
	"gophkeeper/client/domain"
	"gophkeeper/internal"
	g "gophkeeper/internal/client/workers/grpc"
	"gophkeeper/internal/client/workers/grpc/interceptors"
	"gophkeeper/internal/crypto"
	pb "gophkeeper/proto"
	"os"
	"path/filepath"

	"golang.org/x/crypto/pbkdf2"
	"google.golang.org/grpc"
)

type AppUser struct {
	Token,
	Login string
	StorageKey []byte
}

type App struct {
	UserClient    *g.UserClient
	DataClient    *g.DataClient
	User          AppUser
	DecryptedData map[uint64]domain.Data
	DataSavePath  string
}

var AppInstance *App

func InitApp() error {
	var err error

	internal.InitLogger()
	c := initConfig()
	if err = checkConfig(c); err != nil {
		return err
	}

	AppInstance = &App{
		DecryptedData: make(map[uint64]domain.Data),
		DataSavePath:  c.fileSavePath,
	}

	err = initGRPCUserClient(c)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) SetStorageKey(login, pass string) {
	a.User.StorageKey = pbkdf2.Key([]byte(pass), []byte(login), 4096, 32, sha1.New)
}

func initGRPCUserClient(cnf *config) error {
	ch, err := crypto.NewCipher(cnf.cryptoKeysPath)
	if err != nil {
		internal.Logger.Fatalw("failed to init crypto cipher", "error", err)
	}

	conn, err := grpc.NewClient(
		cnf.address, grpc.WithTransportCredentials(ch.GetClientGRPCTransportCreds()),
		grpc.WithUnaryInterceptor(interceptors.Auth),
		grpc.WithStreamInterceptor(interceptors.StreamAuth),
	)

	if err != nil {
		internal.Logger.Fatalw("failed to create grpc client", "error", err)
	}

	AppInstance.UserClient = g.NewUserClient(pb.NewUserServiceClient(conn))
	AppInstance.DataClient = g.NewDataClient(pb.NewDataServiceClient(conn))

	return nil
}

type config struct {
	address,
	fileSavePath,
	cryptoKeysPath string
}

const serverAddressVAr = "SERVER_ADDRESS"
const cryptoKeysPath = "CRYPTO_KEYS_PATH"

func initConfig() *config {
	c := new(config)

	flag.StringVar(&c.address, "a", "", "server address")
	flag.StringVar(&c.cryptoKeysPath, "c", "", "crypto keys path")
	flag.StringVar(&c.fileSavePath, "f", "", "save files path")

	flag.Parse()

	if envVar := os.Getenv(serverAddressVAr); envVar != "" {
		c.address = envVar
	}

	if envVar := os.Getenv(cryptoKeysPath); envVar != "" {
		c.cryptoKeysPath = envVar
	}

	if c.cryptoKeysPath != "" {
		c.cryptoKeysPath = filepath.FromSlash(c.cryptoKeysPath)
	}

	return c
}

func checkConfig(c *config) error {
	if c.address == "" || c.fileSavePath == "" || c.cryptoKeysPath == "" {
		return errors.New("please, check configs")
	}

	return nil
}
