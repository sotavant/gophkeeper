package interceptors

import (
	"context"
	"errors"
	"gophkeeper/domain"
	"gophkeeper/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ContextUserTokenKey struct{}

func Auth(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	var err error

	ctx, err = setAuthMeta(ctx, method)
	if err != nil {
		return err
	}

	return invoker(ctx, method, req, reply, cc, opts...)
}

func StreamAuth(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (stream grpc.ClientStream, err error) {
	ctx, err = setAuthMeta(ctx, method)
	if err != nil {
		return
	}

	return streamer(ctx, desc, cc, method, opts...)
}

func setAuthMeta(ctx context.Context, method string) (context.Context, error) {
	if method == proto.UserService_Register_FullMethodName || method == proto.UserService_Login_FullMethodName {
		return ctx, nil
	}

	token := ctx.Value(ContextUserTokenKey{}).(string)
	if token == "" {
		return ctx, errors.New("invalid token")
	}

	md := metadata.Pairs(domain.AuthorizationMetaKey, domain.TokenSubstr+" "+token)
	ctx = metadata.NewOutgoingContext(ctx, md)
	return ctx, nil
}
