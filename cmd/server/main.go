package main

import (
	"context"
	"gophkeeper/data"
	"gophkeeper/file"
	"gophkeeper/internal"
	"gophkeeper/internal/crypto"
	"gophkeeper/internal/server"
	grpc2 "gophkeeper/internal/server/grpc"
	"gophkeeper/internal/server/repository/pgsql"
	pb "gophkeeper/proto"
	"gophkeeper/user"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	var listen net.Listener

	ctx := context.Background()
	app, err := server.InitApp(ctx)

	if err != nil {
		panic(err)
	}

	listen, err = net.Listen("tcp", app.Address)
	if err != nil {
		internal.Logger.Fatalw("failed to listen", "err", err)
	}

	s := initGRPCServer(ctx, app)

	jobsDone := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sigint
		s.GracefulStop()
		close(jobsDone)
		internal.Logger.Infow("shutdown complete")
	}()

	go func() {
		if err = s.Serve(listen); err != nil {
			internal.Logger.Fatalw("failed to grpc serve", "err", err)
		}
	}()

	<-jobsDone
}

func initGRPCServer(ctx context.Context, app *server.App) *grpc.Server {
	var err error
	var ch *crypto.Cipher
	var interceptors []grpc.UnaryServerInterceptor

	ch, err = crypto.NewCipher()
	if err != nil {
		internal.Logger.Fatalw("error initializing cipher", "err", err)
	}

	s := grpc.NewServer(grpc.Creds(ch.GetServerGRPCTransportCreds()), grpc.ChainUnaryInterceptor(interceptors...))

	userRepo, err := pgsql.NewUserRepository(ctx, app.DBPool, pgsql.UsersTableName)
	if err != nil {
		panic(err)
	}

	fileRepo, err := pgsql.NewFileRepository(ctx, app.DBPool, pgsql.FileTableName)
	if err != nil {
		panic(err)
	}

	dataRepo, err := pgsql.NewDataRepository(ctx, app.DBPool, pgsql.DataTableName, pgsql.UsersTableName, pgsql.FileTableName)
	if err != nil {
		panic(err)
	}

	userService := user.NewService(userRepo)
	fileService := file.NewService(fileRepo)
	dataService := data.NewService(dataRepo, fileRepo)

	pb.RegisterUserServiceServer(s, grpc2.NewUserServer(userService))
	pb.RegisterDataServiceServer(s, grpc2.NewDataServer(dataService, app.FilesSavePath, fileService))

	return s
}
