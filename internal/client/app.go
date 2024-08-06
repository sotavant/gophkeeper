package client

import (
	"errors"
	"flag"
	"gophkeeper/internal"
	g "gophkeeper/internal/client/workers/grpc"
	"gophkeeper/internal/crypto"
	pb "gophkeeper/proto"
	"os"

	"google.golang.org/grpc"
)

type AppUser struct {
	Token,
	Login string
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
	err = initGRPCUserClient(c.address)
	if err != nil {
		return err
	}

	return nil
}

func initGRPCUserClient(serverAddr string) error {
	ch, err := crypto.NewCipher()
	if err != nil {
		internal.Logger.Fatalw("failed to init crypto cipher", "error", err)
	}

	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(ch.GetClientGRPCTransportCreds()))
	if err != nil {
		internal.Logger.Fatalw("failed to create grpc client", "error", err)
	}

	c := pb.NewUserServiceClient(conn)
	AppInstance.UserClient = g.NewUserClient(c)

	return nil
}

type config struct {
	address string
}

const serverAddressVAr = "SERVER_ADDRESS"

func initConfig() *config {
	c := new(config)

	flag.StringVar(&c.address, "a", "", "server address")

	flag.Parse()

	if envVar := os.Getenv(serverAddressVAr); envVar != "" {
		c.address = envVar
	}

	return c
}

func checkConfig(c *config) error {
	if c.address == "" {
		return errors.New("please, check configs")
	}

	return nil
}
