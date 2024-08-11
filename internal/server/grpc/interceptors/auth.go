// Package interceptors пакет для модификации данных запроса от клиента к серверу
package interceptors

import (
	"context"
	"gophkeeper/internal/server/auth"
	"gophkeeper/proto"
	domain2 "gophkeeper/server/domain"
	"gophkeeper/server/user"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	authenticatedMetaNotFound = "not found Authorization meta"
	wrongAuthenticatedMeta    = "wrong Authorization meta"
)

type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context получить context из потока данных
func (w *wrappedStream) Context() context.Context {
	return w.ctx
}

// SetContext записать context в поток данных
func (w *wrappedStream) SetContext(ctx context.Context) {
	w.ctx = ctx
}

func (w *wrappedStream) RecvMsg(m interface{}) error {
	return w.ServerStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m interface{}) error {
	return w.ServerStream.SendMsg(m)
}

// StreamContextWrapper обертка для потока данных, для управления context
type StreamContextWrapper interface {
	grpc.ServerStream
	SetContext(context.Context)
}

func newStreamContextWrapper(ss grpc.ServerStream) StreamContextWrapper {
	ctx := ss.Context()
	return &wrappedStream{
		ss,
		ctx,
	}
}

// Auth получение из запроса авторизационных данных и запись их в контекст
func Auth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	if info.FullMethod == proto.UserService_Register_FullMethodName || info.FullMethod == proto.UserService_Login_FullMethodName {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		vals := md[domain2.AuthorizationMetaKey]
		if len(vals) > 0 {
			val := vals[0]
			if !strings.Contains(val, domain2.TokenSubstr) {
				return nil, status.Errorf(codes.Unauthenticated, wrongAuthenticatedMeta)
			}

			var userID uint64
			token := strings.TrimSpace(strings.Replace(val, domain2.TokenSubstr, "", -1))
			userID, err = auth.GetUserID(token)
			if err != nil {
				return nil, status.Errorf(codes.Internal, domain2.ErrInternalServerError.Error())
			}

			if userID == 0 {
				return nil, status.Errorf(codes.Unauthenticated, wrongAuthenticatedMeta)
			}

			respCtx := context.WithValue(ctx, user.ContextUserIDKey{}, userID)

			return handler(respCtx, req)
		}
		return nil, status.Errorf(codes.Unauthenticated, authenticatedMetaNotFound)
	} else {
		return nil, status.Errorf(codes.Unauthenticated, authenticatedMetaNotFound)
	}
}

// StreamAuth получение из потокового запроса авторизационных данных и запись их в контекст
func StreamAuth(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	md, ok := metadata.FromIncomingContext(ss.Context())
	if ok {
		vals := md[domain2.AuthorizationMetaKey]
		if len(vals) > 0 {
			val := vals[0]
			if !strings.Contains(val, domain2.TokenSubstr) {
				return status.Errorf(codes.Unauthenticated, wrongAuthenticatedMeta)
			}

			var userID uint64
			token := strings.TrimSpace(strings.Replace(val, domain2.TokenSubstr, "", -1))
			userID, err := auth.GetUserID(token)
			if err != nil {
				return status.Errorf(codes.Internal, domain2.ErrInternalServerError.Error())
			}

			if userID == 0 {
				return status.Errorf(codes.Unauthenticated, wrongAuthenticatedMeta)
			}

			respCtx := context.WithValue(ss.Context(), user.ContextUserIDKey{}, userID)
			sw := newStreamContextWrapper(ss)
			sw.SetContext(respCtx)

			return handler(srv, sw)
		}
		return status.Errorf(codes.Unauthenticated, authenticatedMetaNotFound)
	} else {
		return status.Errorf(codes.Unauthenticated, authenticatedMetaNotFound)
	}
}
