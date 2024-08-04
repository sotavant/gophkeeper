package interceptors

import (
	"context"
	"gophkeeper/domain"
	"gophkeeper/internal/server/auth"
	"gophkeeper/user"
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

func (w *wrappedStream) Context() context.Context {
	return w.ctx
}

func (w *wrappedStream) SetContext(ctx context.Context) {
	w.ctx = ctx
}

func (w *wrappedStream) RecvMsg(m interface{}) error {
	return w.ServerStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m interface{}) error {
	return w.ServerStream.SendMsg(m)
}

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

func Auth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		vals := md[domain.AuthorizationMetaKey]
		if len(vals) > 0 {
			val := vals[0]
			if !strings.Contains(val, domain.TokenSubstr) {
				return nil, status.Errorf(codes.Unauthenticated, wrongAuthenticatedMeta)
			}

			var userID uint64
			token := strings.TrimSpace(strings.Replace(val, domain.TokenSubstr, "", -1))
			userID, err = auth.GetUserID(token)
			if err != nil {
				return nil, status.Errorf(codes.Internal, domain.ErrInternalServerError.Error())
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

func StreamAuth(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	md, ok := metadata.FromIncomingContext(ss.Context())
	if ok {
		vals := md[domain.AuthorizationMetaKey]
		if len(vals) > 0 {
			val := vals[0]
			if !strings.Contains(val, domain.TokenSubstr) {
				return status.Errorf(codes.Unauthenticated, wrongAuthenticatedMeta)
			}

			var userID uint64
			token := strings.TrimSpace(strings.Replace(val, domain.TokenSubstr, "", -1))
			userID, err := auth.GetUserID(token)
			if err != nil {
				return status.Errorf(codes.Internal, domain.ErrInternalServerError.Error())
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
