package client

import (
	"errors"
	"flag"
	"gophkeeper/internal"
	g "gophkeeper/internal/client/workers/grpc"
	"gophkeeper/internal/crypto"
	pb "gophkeeper/proto"
	"os"
	"path/filepath"

	"google.golang.org/grpc"
)

type AppUser struct {
	Token,
	Login string
	StorageKey []byte
}

type App struct {
	UserClient *g.UserClient
	User       AppUser
}

var AppInstance *App

func InitApp() error {
	var err error

	internal.InitLogger()
	c := initConfig()
	if err = checkConfig(c); err != nil {
		return err
	}

	AppInstance = &App{}
	err = initGRPCUserClient(c)
	if err != nil {
		return err
	}

	return nil
}

func initGRPCUserClient(cnf *config) error {
	ch, err := crypto.NewCipher(cnf.cryptoKeysPath)
	if err != nil {
		internal.Logger.Fatalw("failed to init crypto cipher", "error", err)
	}

	conn, err := grpc.NewClient(cnf.address, grpc.WithTransportCredentials(ch.GetClientGRPCTransportCreds()))
	if err != nil {
		internal.Logger.Fatalw("failed to create grpc client", "error", err)
	}

	c := pb.NewUserServiceClient(conn)
	AppInstance.UserClient = g.NewUserClient(c)

	return nil
}

type config struct {
	address,
	cryptoKeysPath string
}

const serverAddressVAr = "SERVER_ADDRESS"
const cryptoKeysPath = "CRYPTO_KEYS_PATH"

func initConfig() *config {
	c := new(config)

	flag.StringVar(&c.address, "a", "", "server address")
	flag.StringVar(&c.cryptoKeysPath, "c", "", "crypto keys path")

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
	if c.address == "" {
		return errors.New("please, check configs")
	}

	return nil
}
