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

func Auth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		vals := md[domain.AuthorizationMetaKey]
		if len(vals) > 0 {
			val := vals[0]
			if !strings.Contains(val, domain.TokenSubstr) {
				return nil, status.Errorf(codes.Unauthenticated, wrongAuthenticatedMeta)
			}

			var userID int64
			token := strings.TrimSpace(strings.Replace(val, domain.TokenSubstr, "", -1))
			userID, err = auth.GetUserID(token)
			if err != nil {
				return nil, status.Errorf(codes.Internal, domain.ErrInternalServerError.Error())
			}

			if userID == -1 {
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
